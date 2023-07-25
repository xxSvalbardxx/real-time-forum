package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	session "forum/internal"
	"forum/utils"

	_ "github.com/mattn/go-sqlite3"
)

// var db, _ = InitDB()

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

/* --------------------------------------- User ------------------------------------------- */

func CreateUserTable(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname TEXT NOT NULL,
		age INTEGER NOT NULL,
		gender TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s", err, sqlStmt)
		return
	}
}

func InsertUser(db *sql.DB, Data map[string]interface{}) string {
	nickname := Data["nickname"]
	age := Data["age"]
	gender := Data["gender"]
	firstName := Data["first_name"]
	lastName := Data["last_name"]
	email := Data["email"]
	password := Data["password"]
	var confirmation string
	_, err := db.Exec("INSERT INTO users (nickname, age, gender, first_name, last_name, email, password) VALUES(?,?,?,?,?,?,?)", nickname, age, gender, firstName, lastName, email, password)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error inserting new user into the users table")
		confirmation = "User not registered"
		return confirmation
	}
	confirmation = "User registered"
	return confirmation
}

func LoginUser(db *sql.DB, Data map[string]interface{}) (string, string) {
	indentifier := Data["identifier"]
	password := Data["password"]

	var dbId int
	var dbNickname string
	var dbEmail string
	var dbPassword string

	// Check if the identifier is corresponding to an email or a nickname in the database AND if the password is correct
	err := db.QueryRow("SELECT id, nickname, email, password FROM users WHERE (nickname = ? AND password = ?) OR (email = ? AND password = ?)", indentifier, password, indentifier, password).Scan(&dbId, &dbNickname, &dbEmail, &dbPassword)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error logging in user")
		return "", ""
	}

	if sessionD := session.FindSessionByUserId(dbId); sessionD != nil {
		session.Logout(sessionD.ID)
	}
	sessionD := session.NewSession(dbId, dbNickname)

	return sessionD.ID, sessionD.Username
}

/* --------------------------------- Private message ------------------------------------ */

func CreatePrivateMessageTable(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS private_messages (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		sender TEXT NOT NULL,
		receiver TEXT NOT NULL,
		content TEXT NOT NULL,
		date TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s", err, sqlStmt)
		return
	}
}

func InsertPrivateMessage(db *sql.DB, sender string, receiver string, content string, date time.Time) {
	_, err := db.Exec("INSERT INTO private_messages (sender, receiver, content, date) VALUES (?, ?, ?, ?)", sender, receiver, content, date)
	if err != nil {
		log.Fatal(err)
	}
}

func DisplayMessages(db *sql.DB, sender string, receiver string) []utils.Message {
	messages := []utils.Message{}

	// Search for all messages by sender and receiver
	rows, err := db.Query("SELECT * FROM private_messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)", sender, receiver, receiver, sender)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// For each message, append it to the messages slice
	for rows.Next() {
		var id int
		var sender string
		var receiver string
		var content string
		var date string

		err = rows.Scan(&id, &sender, &receiver, &content, &date)
		if err != nil {
			log.Fatal(err)
		}
		senderUUID := session.FindSessionByUsername(sender)
		// reformat date string in time.Time
		/* date = strings.Replace(date, "T", " ", 1)
		date = strings.Replace(date, "Z", "", 1) */

		messages = append(messages, utils.Message{ID: id, SenderUUID: senderUUID, Sender: sender, Receiver: receiver, Content: content, Date: date})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return messages

	/*
		// Display all messages
		for _, message := range messages {
			fmt.Println(message)
		}
	*/
}

/* ---------------------------------------- Post ------------------------------------------- */

func CreatePostTable(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		author TEXT
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s", err, sqlStmt)
		return
	}
}

func InsertPost(db *sql.DB, title string, content string, author string) {
	_, err := db.Exec("INSERT INTO posts (title, content, author) VALUES (?, ?, ?)", title, content, author)
	if err != nil {
		log.Fatal(err)
	}
}

func DisplayPosts(db *sql.DB) []utils.Post {
	posts := []utils.Post{}

	// Search for all posts
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// For each post, append it to the posts slice
	for rows.Next() {
		var id int
		var title string
		var content string
		var author string

		err = rows.Scan(&id, &title, &content, &author)
		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, utils.Post{ID: id, Title: title, Content: content, Author: author})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return posts
}

/* ---------------------------------------- Comment ------------------------------------------- */

func CreateCommentTable(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		author TEXT,
		post_id INTEGER
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s", err, sqlStmt)
		return
	}
}

func InsertComment(db *sql.DB, content string, author string, post_id int) {
	_, err := db.Exec("INSERT INTO comments (content, author, post_id) VALUES (?, ?, ?)", content, author, post_id)
	if err != nil {
		log.Fatal(err)
	}
}

func GetCommentsByPost(db *sql.DB, post_id int) []utils.Comment {
	comments := []utils.Comment{}

	// Search for all comments by post_id
	rows, err := db.Query("SELECT * FROM comments WHERE post_id = ?", post_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// For each comment, append it to the comments slice
	for rows.Next() {
		var id int
		var content string
		var author string
		var post_id int

		err = rows.Scan(&id, &content, &author, &post_id)
		if err != nil {
			log.Fatal(err)
		}

		comments = append(comments, utils.Comment{ID: id, Content: content, Author: author, PostID: post_id})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return comments
}

/* -------------------------------------- Chat ---------------------------------------- */

func GetLastMessage(db *sql.DB, user string, current string) (utils.Message, error) {
	query := "SELECT * FROM private_messages WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?) ORDER BY date DESC LIMIT 1" // Get the last message sent by the user or the current user
	row := db.QueryRow(query, user, current, current, user)
	message := utils.Message{}
	err := row.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.Date)
	if err != nil {
		fmt.Println(err)
		return message, err

	}
	fmt.Println(message)
	return message, nil
}
