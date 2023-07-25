package websocket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	session "forum/internal"
	"forum/internal/database"
	"forum/utils"

	"github.com/gorilla/websocket"
)

var (
	connexions = make(map[*websocket.Conn]string)
	UsersUUID  = make(map[string]string)
	db, _      = database.InitDB()
	upgrader   = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	defer conn.Close() // close the connection when the function returns
	for {
		message, err := ReadMsg(w, r, conn)
		if err != nil {
			fmt.Println(err)
			//fmt.Println("User:" + conn.LocalAddr().String() + "disconnected")
			break
		}
		TypeOfMsg(conn, message)
	}
}

func ReadMsg(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) ([]byte, error) {
	_, message, err := conn.ReadMessage() // we read the message
	if err != nil {                       // if there is an error
		RemoveConnectedUser(conn)
		// log.Println(err) // we log the error
	}
	return message, err // we return the message
}

/* ------------------------------------ WS Receiver ---------------------------------------- */

func TypeOfMsg(conn *websocket.Conn, message []byte) {
	var data map[string]interface{}

	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println(err)
	}
	switch data["type"] {
	case "register":
		registerData := data["data"].(map[string]interface{})
		confirm := database.InsertUser(db, registerData)
		fmt.Println(confirm)
		sendConfirmation(conn, confirm)

	case "login":
		loginData := data["data"].(map[string]interface{})
		userSession, username := database.LoginUser(db, loginData)
		sendUserSession(conn, userSession, username)
	case "getPosts":
		DisplayAllPosts(db, conn)

	case "logout":
		RemoveConnectedUser(conn)
		//fmt.Println("User:" + conn.LocalAddr().String() + "disconnected")

	case "new_post":
		postData := data["data"].(map[string]interface{})
		storeAndDisplayPosts(db, postData)

	case "get_comments":
		commentsData := data["data"].(map[string]interface{})
		// fmt.Println(commentsData)
		log.Println("commentsData : ", commentsData)
		DisplayCommentsByPost(db, conn, commentsData)

	case "new_comment":
		commentData := data["data"].(map[string]interface{})
		// log.Println("commentData : ", commentData)
		storeAndDisplayComments(db, commentData)

	case "getPrvMsg":
		getPrvMsgData := data["data"].(map[string]interface{})
		// fmt.Println(getPrvMsgData)
		previousMessages(db, getPrvMsgData)

	case "PrivateMessage":
		messageData := data["data"].(map[string]interface{})
		//log.Println(messageData)
		storeAndDisplayMessage(db, messageData)
		UpdateConnectedUsers(connexions)
	case "checkCurrentUser":
		donn := data["data"].(map[string]interface{})
		users, current := checkCurrentUser(db, donn)
		SortUsersByLastMessage(db, conn, users, current)
	case "is_typing":
		messageData := data["data"].(map[string]interface{})
		TypingNotification(conn, messageData)
	}
}

/* ------------------------------------ Register ---------------------------------------- */

func sendConfirmation(conn *websocket.Conn, confirm string) {
	// send confirmation to the client
	confirmation := struct {
		Type string
	}{
		Type: confirm,
	}
	confirmationJSON, _ := json.Marshal(confirmation)
	conn.WriteMessage(1, confirmationJSON)
	fmt.Println("Confirmation sent to client")
}

/* ------------------------------------ Login ---------------------------------------- */

func sendUserSession(conn *websocket.Conn, userSession string, username string) {
	var confirmation string
	if userSession == "" {
		confirmation = "login_failed"
	} else if userSession != "" {
		confirmation = "login_success"
		AddConnectedUser(conn, username)
	}
	data := map[string]string{
		"username": username,
		"session":  userSession,
	}

	LoggedMsg := struct {
		Type string
		Data map[string]string
	}{
		Type: confirmation,
		Data: data,
	}
	LoggedMsgJson, err := json.Marshal(LoggedMsg)
	if err != nil {
		log.Println(err)
	}
	conn.WriteMessage(1, LoggedMsgJson)
}

/* ------------------------------------ Chat ---------------------------------------- */
func checkCurrentUser(db *sql.DB, donn map[string]interface{}) ([]string, string) {
	current := donn["session"].(string)
	current = session.GetUsernameBySession(current)
	users := donn["users"].([]interface{})
	// convert users to []string
	var usersString []string
	for _, user := range users {
		if user.(string) != current {
			usersString = append(usersString, user.(string))
		}
	}
	// fmt.Println(usersString)
	return usersString, current
}

func SortUsersByLastMessage(db *sql.DB, conn *websocket.Conn, users []string, current string) {
	// get the last message of each user
	var lastMessages []utils.Message
	var noMessages []string
	// fmt.Println("-------------------")
	for _, user := range users {
		lastMessage, err := database.GetLastMessage(db, user, current)
		/* fmt.Println(lastMessage)
		fmt.Println(err) */
		if err != nil {
			noMessages = append(noMessages, user)
		} else {
			lastMessages = append(lastMessages, lastMessage)
		}
	}

	// fmt.Println("Current user : ", current)
	//	fmt.Println("lastMessages : ", lastMessages)

	// sort noMessages by alphabetical order
	for i := 0; i < len(noMessages); i++ {
		for j := i + 1; j < len(noMessages); j++ {
			if noMessages[i] > noMessages[j] {
				noMessages[i], noMessages[j] = noMessages[j], noMessages[i]
			}
		}
	}
	// sort lastMessages by date
	for i := 0; i < len(lastMessages); i++ {
		for j := i + 1; j < len(lastMessages); j++ {
			if lastMessages[i].Date < lastMessages[j].Date {
				lastMessages[i], lastMessages[j] = lastMessages[j], lastMessages[i]
			}
		}
	}
	// convert lastMessages to []string with only the username
	var lastMsgString []string
	for i := 0; i < len(lastMessages); i++ {
		if lastMessages[i].Sender == current {
			lastMsgString = append(lastMsgString, lastMessages[i].Receiver)
		} else {
			lastMsgString = append(lastMsgString, lastMessages[i].Sender)
		}
	}
	// fmt.Println("lastMsgString : ", lastMsgString)
	// fmt.Println("noMessages : ", noMessages)
	// merge the two slices
	var sortedUsers []string
	for i := 0; i < len(lastMsgString); i++ {
		sortedUsers = append(sortedUsers, lastMsgString[i])
	}
	for i := 0; i < len(noMessages); i++ {
		sortedUsers = append(sortedUsers, noMessages[i])
	}

	// if there is 2 identical usernames, remove one
	for i := 0; i < len(sortedUsers); i++ {
		for j := i + 1; j < len(sortedUsers); j++ {
			if sortedUsers[i] == sortedUsers[j] {
				sortedUsers = append(sortedUsers[:j], sortedUsers[j+1:]...)
			}
		}
	}

	// fmt.Println("sortedUsers : ", sortedUsers)
	// send the sortedUsers to the client
	// send the sortedUsers to the client
	Users := struct {
		Type string
		Data []string
	}{
		Type: "Sortedusers",
		Data: sortedUsers,
	}
	UsersJson, err := json.Marshal(Users)
	if err != nil {
		log.Println(err)
	}
	conn.WriteMessage(1, UsersJson)
}

func AddConnectedUser(conn *websocket.Conn, username string) {
	connexions[conn] = username
	UpdateConnectedUsers(connexions)
}

func RemoveConnectedUser(conn *websocket.Conn) {
	delete(connexions, conn)
	UpdateConnectedUsers(connexions)
}

// update the list of connected users for the client.
func UpdateConnectedUsers(connexions map[*websocket.Conn]string) {
	var users []string
	for _, user := range connexions {
		users = append(users, user)
	}

	// users = SortUsers(users)

	// fmt.Println(users)
	ConnectedUsers := struct {
		Type string
		Data []string
	}{
		Type: "connected_users",
		Data: users,
	}

	ConnectedUsersJson, err := json.Marshal(ConnectedUsers)
	if err != nil {
		log.Println(err)
	}

	for conn := range connexions {
		conn.WriteMessage(1, ConnectedUsersJson)
	}
}

func SortUsers(users []string) []string {
	sort.Strings(users)
	return users
}

/* ------------------------------------ Private Messages ---------------------------------------- */

func previousMessages(db *sql.DB, getPrvMsgData map[string]interface{}) {
	// fmt.Println(getPrvMsgData)

	senderUUID := getPrvMsgData["senderUUID"].(string)
	sender := session.GetUsernameBySession(senderUUID)
	senderAdress := GetConnByNickname(sender)

	receiver := getPrvMsgData["receiverNickname"].(string)
	// log.Printf("receiver : '%s'", receiver)
	receiverAdress := GetConnByNickname(receiver)

	messagesArray := database.DisplayMessages(db, sender, receiver)

	// put the messages in a MessageArray struct
	MsgArray := *new(utils.MessageArray)
	MsgArray.Type = "PrvMsg"
	MsgArray.Data = messagesArray

	// convert the struct to json
	MsgArrayJson, err := json.Marshal(MsgArray)
	if err != nil {
		log.Println(err)
	}

	// send the json to the clients
	senderAdress.WriteMessage(1, MsgArrayJson)
	receiverAdress.WriteMessage(1, MsgArrayJson)
}

func storeAndDisplayMessage(db *sql.DB, messageData map[string]interface{}) {
	senderUuid := messageData["senderUUID"].(string)
	sender := session.GetUsernameBySession(senderUuid)
	senderAdress := GetConnByNickname(sender)

	receiver := messageData["receiverNickname"].(string)
	receiverAdress := GetConnByNickname(receiver)

	message := messageData["messageContent"].(string)

	database.InsertPrivateMessage(db, sender, receiver, message, time.Now())

	messagesArray := database.DisplayMessages(db, sender, receiver)

	// put the messages in a MessageArray struct
	MsgArray := *new(utils.MessageArray)
	MsgArray.Type = "PrvMsg"
	MsgArray.Data = messagesArray

	// convert the struct to json
	MsgArrayJson, err := json.Marshal(MsgArray)
	if err != nil {
		log.Println(err)
	}

	// send the json to the sender
	senderAdress.WriteMessage(1, MsgArrayJson)
	// send the json to the receiver
	receiverAdress.WriteMessage(1, MsgArrayJson)
}

func GetConnByNickname(nickname string) *websocket.Conn {
	for conn, user := range connexions {
		if user == nickname {
			return conn
		}
	}
	//fmt.Println(connexions)
	log.Println("Connexion introuvable")
	return nil
}

/*
func ReceivePrivateMessage(conn *websocket.Conn, db *sql.DB) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Erreur lors de la lecture du message privé")
			break
		}

		var privateMessageData map[string]interface{}
		err = json.Unmarshal(message, &privateMessageData)
		if err != nil {
			log.Println("Erreur de décodage du message JSON: ", err)
			continue
		}
		database.InsertPrivateMessage(db, privateMessageData)
	}
} */

// typing notification
func TypingNotification(conn *websocket.Conn, messageData map[string]interface{} ) {
	receiver := messageData["receiverNickname"].(string)
	receiverAdress := GetConnByNickname(receiver)

	MsgArray := *new(utils.MessageArray)
	MsgArray.Type = "typingNotification"

	// convert the struct to json
	MsgArrayJson, err := json.Marshal(MsgArray)
	if err != nil {
		log.Println(err)
	}
	
	// send the json to the receiver
	receiverAdress.WriteMessage(1, MsgArrayJson)
}
	

/* ------------------------------------ Post ---------------------------------------- */

func storeAndDisplayPosts(db *sql.DB, postData map[string]interface{}) {
	title := postData["title"].(string)
	content := postData["content"].(string)
	authorUUID := postData["author"].(string)
	author := session.GetUsernameBySession(authorUUID)

	if postData != nil {
		database.InsertPost(db, title, content, author)
	}

	PostsArray := database.DisplayPosts(db)

	// put the messages in a PostArray struct
	PostsArrayStruct := *new(utils.PostArray)
	PostsArrayStruct.Type = "Posts"
	PostsArrayStruct.Data = PostsArray

	// convert the struct to json
	PostsArrayJson, err := json.Marshal(PostsArrayStruct)
	if err != nil {
		log.Println(err)
	}

	// send the json to the clients
	for conn := range connexions {
		conn.WriteMessage(1, PostsArrayJson)
	}
}

func DisplayAllPosts(db *sql.DB, conn *websocket.Conn) {
	PostsArray := database.DisplayPosts(db)

	// put the messages in a PostArray struct
	PostsArrayStruct := *new(utils.PostArray)
	PostsArrayStruct.Type = "AllPosts"
	PostsArrayStruct.Data = PostsArray

	// convert the struct to json
	PostsArrayJson, err := json.Marshal(PostsArrayStruct)
	if err != nil {
		log.Println(err)
	}

	// send the json to the client
	conn.WriteMessage(1, PostsArrayJson)
}

/* ------------------------------------ Comments ---------------------------------------- */

func storeAndDisplayComments(db *sql.DB, commentData map[string]interface{}) {
	content := commentData["content"].(string)
	authorUUID := commentData["author"].(string)
	author := session.GetUsernameBySession(authorUUID)
	postIdString := commentData["postId"].(string)

	var err error
	postId, err := strconv.Atoi(postIdString)
	if err != nil {
		log.Println(err)
	}

	if commentData != nil {
		database.InsertComment(db, content, author, postId)
	}

	CommentsArray := database.GetCommentsByPost(db, postId)

	// put the messages in a PostArray struct
	CommentsArrayStruct := *new(utils.CommentArray)
	CommentsArrayStruct.Type = "Comments"
	CommentsArrayStruct.Data = CommentsArray

	// convert the struct to json
	CommentsArrayJson, err := json.Marshal(CommentsArrayStruct)
	if err != nil {
		log.Println(err)
	}

	// send the json to the clients
	for conn := range connexions {
		conn.WriteMessage(1, CommentsArrayJson)
	}
}

func DisplayCommentsByPost(db *sql.DB, conn *websocket.Conn, postData map[string]interface{}) {
	postIdString := postData["postId"].(string)

	var err error
	postId, err := strconv.Atoi(postIdString)
	if err != nil {
		log.Println(err)
	}

	CommentsArray := database.GetCommentsByPost(db, postId)

	// put the messages in a PostArray struct
	CommentsArrayStruct := *new(utils.CommentArray)
	CommentsArrayStruct.Type = "Comments"
	CommentsArrayStruct.Data = CommentsArray

	// convert the struct to json
	CommentsArrayJson, err := json.Marshal(CommentsArrayStruct)
	if err != nil {
		log.Println(err)
	}

	// send the json to the client
	conn.WriteMessage(1, CommentsArrayJson)
}
