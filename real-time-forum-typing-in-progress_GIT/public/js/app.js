
// créer une nouvelle connexion WebSocket
var socket = new WebSocket("ws://localhost:8080/ws");

// private message
let startIndex = {};
const messagePerPage = 10;
let allPrivateMessages = [];

let currentUser = "";

var chatBox = document.getElementById("chat_box");
var chat = document.getElementById("chat_box");

const login_notif = document.getElementById("login_notif");
const reg_notif = document.getElementById("reg_notif");
const forgot_password = document.getElementById("forgot_password");


var loginBox = document.querySelector(".login_form_container");
var registerBox = document.querySelector(".register_form_container")
var createPost = document.querySelector(".create-post")
var peopleOnline = document.querySelector(".people-online");
var signupSwitch = document.getElementById("register-return");
var loginSwitch = document.getElementById("login-return");
var privateMsg = document.getElementById("private-msg");

var sentMessage = document.getElementsByClassName("sent_message");
var receivedMessage = document.getElementsByClassName("received_message");

var posts = document.getElementById("posts");
var comments = document.getElementById("comments");

let typingTimeout;

var sdrUUID = "";

// A la connexion, on affiche un message dans la console
socket.onopen = (e) => {
    console.log("Connection established with the server");
}

// A la réception d'un message, on l'affiche dans la console
socket.onmessage = (e) => {
    var message = JSON.parse(e.data);

    if (message.Type == "User registered") {
        console.log("register success");
        reg_notif.style.display = "block";
        reg_notif.innerHTML = "Register success";
        reg_notif.style.color = "green";
        //hideRegisterFormOnLogin();

    } if (message.Type == "User not registered") {
        console.log("register failed");
        reg_notif.style.display = "block";
        reg_notif.innerHTML = "Register failed";
        reg_notif.style.color = "red";

    } if (message.Type == "login_success") {
        // create a cookie named session with the value of the session id
        document.cookie = "session=" + message.Data.session; "max-age=3600";
        login_notif.style.display = "block";
        login_notif.innerHTML = "Login success";
        login_notif.style.color = "green";
        enterForum();
        welcomeUser(message.Data);
        sendConfirmation();

        //hideRegisterFormOnLogin();
    }
    if (message.Type == "login_failed") {
        console.log("login failed");
        login_notif.style.display = "block";
        login_notif.innerHTML = "Login failed";
        login_notif.style.color = "red";

    } if (message.Type == "connected_users") {
        //console.log("connected users : " + message.Data);
        checkCurrentUser(message.Data);
    } if (message.Type == "Sortedusers") {
        DisplayConnectedUsers(message.Data);

    } if (message.Type == "PrvMsg") {
        RemoveMessages(); // maybe essential ? 
        allPrivateMessages = message.Data; //.sort((a, b) => new Date(a.Date) -  new Date(b.Date));
        //console.log("type startindex", startIndex)
        startIndex[currentUser] =
            allPrivateMessages.length > messagePerPage
                ? allPrivateMessages.length - messagePerPage
                : 0;
        DisplayMessages(allPrivateMessages, currentUser);
    } if (message.Type == "Posts") {
        DisplayPosts(message.Data);
    } if (message.Type == "AllPosts") {
        DisplayPosts(message.Data);
    } if (message.Type == "Comments") {
        DisplayComments(message.Data);
    } if (message.Type == "typingNotification") {
        console.log("typing");
        typingIndicator();
    }
}

/* ------------------------------------ Register ---------------------------------------- */

// Envoie les données du formulaire au serveur
const submitFormRegister = () => {
    var form = document.getElementById("registerForm");
    // On récupère les données du formulaire
    var nickname = document.getElementById("R-nickname").value;
    var age = document.getElementById("R-age").value;
    var gender = document.getElementById("R-gender").value;
    var first_name = document.getElementById("R-first-name").value;
    var last_name = document.getElementById("R-last-name").value;
    var email = document.getElementById("R-email").value;
    var password = document.getElementById("R-password").value;

    // On transforme les données en JSON pour les envoyer au serveur afin que ce soit plus simple à traiter
    var data = ({
        "nickname": nickname,
        "age": age,
        "gender": gender,
        "first_name": first_name,
        "last_name": last_name,
        "email": email,
        "password": password
    });
    var type = "register"
    var msg = JSON.stringify({
        type: type,
        data: data
    })
    socket.send(msg);
    form.reset();
}

/* ------------------------------------ Login ---------------------------------------- */

const submitFormLogin = () => {
    var form = document.getElementById("loginForm");
    var identifier = document.getElementById("identifier").value;
    var password = document.getElementById("password").value;
    var data = ({
        "identifier": identifier,
        "password": password
    });
    var type = "login"
    var msg = JSON.stringify({
        type: type,
        data: data
    })
    socket.send(msg);
    form.reset();
}

window.addEventListener("beforeunload", () => {
    document.cookie = "session= ; max-age=0 ; SameSite=None; Secure";
})

const forgotPassword = () => {
    forgot_password.style.display = "block";
    forgot_password.innerHTML = "Tant pis pour toi";
    forgot_password.style.color = "red";
}

const welcomeUser = (data) => {
    var welcome = document.getElementById("welcome");
    welcome.innerHTML = "Welcome " + data.username;

}

const sendConfirmation = () => {
    var type = "getPosts";
    var msg = JSON.stringify({
        type: type,
    })
    socket.send(msg);
}

/* ------------------------------------ Logout ---------------------------------------- */

const logout = () => {
    var type = "logout";
    var msg = JSON.stringify({
        type: type,
        data: ""
    })
    socket.send(msg);
    document.cookie = "session= ; max-age=0 ; SameSite=None; Secure";
    location.reload();
}

/* ------------------------------------ Chat ---------------------------------------- */

const checkCurrentUser = (data) => {
    // delay the sending of the message to avoid sending it before the client is ready to receive it.

    setTimeout(function () {
        var session = getCookieValue();

        var Data = ({
            "session": session,
            "users": data
        });
        var type = "checkCurrentUser";
        var msg = JSON.stringify({
            type: type,
            data: Data
        })

        socket.send(msg);

    }, 666);
}



const DisplayConnectedUsers = (data) => {

    var connectedUsers = document.getElementById("connected-users");
    connectedUsers.innerHTML = "";
    if (data == null) {
        data = [];
    }
    for (var i = 0; i < data.length; i++) {
        var user = data[i];
        var userDiv = document.createElement("div");

        var userImg = document.createElement("i");
        userImg.className = "fa fa-user";
        var userName = document.createElement("p");
        userName.className = "user";
        userName.innerHTML = user;

        userDiv.appendChild(userImg);
        userDiv.appendChild(userName);

        connectedUsers.appendChild(userDiv);
    }
    var users = document.getElementsByClassName("user");
    for (var i = 0; i < users.length; i++) {
        users[i].addEventListener("click", (e) => {

            var username = e.target.innerText;
            //console.log(e.currentTarget)
            createChatWindow(username);
        });
    }
}

const createChatWindow = (username) => {
    var chatbox = document.getElementById("chat_window");
    var title = document.getElementById("chat_user");
    title.textContent = username;

    RemoveMessages();
    var data = ({
        "senderUUID": getCookieValue(),
        "receiverNickname": username
    });
    socket.send(JSON.stringify({
        type: "getPrvMsg",
        data: data
    }));
    //console.log("create chat window data", data)

    chatbox.style.display = "flex";
    chatbox.style.zIndex = "100";

    startIndex[username] = 0;
    currentUser = username;
    //console.log("create chat window startIndex", startIndex)

}

const closeChatWindow = () => {
    var chatbox = document.getElementById("chat_window");
    chatbox.style.display = "none";
    chatbox.style.zIndex = "-100";
}

/* ------------------------------------ Private Messages ---------------------------------------- */

/* 
const sendMessage = (username, message) => {
    console.log("Message envoyé à " + username + " : " + message);
}
 */
const submitFormMessage = () => {
    var senderSession = getCookieValue();// get the cookie value to get the session id of the sender
    var receiver = document.getElementById("chat_user").textContent;// get the receiver from the chat window
    var message = document.getElementById("messageContent").value;// get the message from the form
    if (message != "") {
    var data = ({
        "senderUUID": senderSession,
        "receiverNickname": receiver,
        "messageContent": message
    });
    var type = "PrivateMessage"
    var msg = JSON.stringify({
        type: type,
        data: data
    })
    //console.log(msg);
    socket.send(msg);
    //RemoveMessages();
    //vider le champ de texte
    document.getElementById("messageContent").value = "";
    }
}

const RemoveMessages = () => {
    //remove all children of the chat div
    var chat = document.getElementById("chat_box");
    chat.innerHTML = "";
}

const DisplayMessages = (data, username) => {
    if (sentMessage != null && receivedMessage != null) {
        for (var i = 0; i < sentMessage.length; i++) {
            sentMessage[i].remove();
        }
        for (var i = 0; i < receivedMessage.length; i++) {
            receivedMessage[i].remove();
        }
    }
    // Sender; Receiver; Content; Date
    //for (var i = data.length - 1 ; i > startIndex[username]; i--) {
    for (var i = startIndex[username]; i < data.length; i++) {
        //console.log("display message startindex", startIndex)
        var message = data[i];
        if (message.SenderUUID == getCookieValue()) {
            var sentMessage = document.createElement("div");
            sentMessage.className = "sent_message";
            var messageContent = document.createElement("p");
            messageContent.innerHTML = message.Content;

            if (messageContent.innerHTML.length > 20) {
                messageContent.style.wordWrap = "break-word";
                messageContent.style.wordBreak = "break-all";
                messageContent.style.height = "auto";
            }

            var messageDate = document.createElement("p");
            // show the date of the message in the format hh:mm

            dt = new Date(message.Date);
            messageDate.innerHTML = dt.getHours() + ":" + dt.getMinutes() + " " + dt.getDate() + "/" + dt.getMonth() + "/" + dt.getFullYear();

            messageContent.className = "message_content";

            messageDate.className = "message_date";

            sentMessage.appendChild(messageContent);
            sentMessage.appendChild(messageDate);

            chat.appendChild(sentMessage);

            if (messageContent.offsetHeight > 100) {
                messageContent.style.height = "auto";
                //messageContent.style.overflow = "auto";
            }

        } else {
            var receivedMessage = document.createElement("div");
            receivedMessage.className = "received_message";
            var messageContent = document.createElement("p");
            messageContent.innerHTML = message.Content;

            if (messageContent.innerHTML.length > 20) {
                messageContent.style.wordWrap = "break-word";
                messageContent.style.wordBreak = "break-all";
                messageContent.style.height = "auto";
            }

            var messageDate = document.createElement("p");
            dt = new Date(message.Date);
            messageDate.innerHTML = dt.getHours() + ":" + dt.getMinutes() + " " + dt.getDate() + "/" + dt.getMonth() + "/" + dt.getFullYear();

            messageContent.className = "message_content";
            messageDate.className = "message_date";

            receivedMessage.appendChild(messageContent);
            receivedMessage.appendChild(messageDate);

            chat.appendChild(receivedMessage);

            if (messageContent.offsetHeight > 100) {
                messageContent.style.height = "auto";
                //messageContent.style.overflow = "auto";
            }
        }
    }
    scrollToBottom();
}

chatBox.addEventListener('scroll', function () {
    if (this.scrollTop === 0) {
        // L'utilisateur a atteint le haut de la fenêtre et il y a des messages précédents à charger

        if (startIndex[currentUser] - messagePerPage < 0) {
            startIndex[currentUser] = 0;
        } else {
            startIndex[currentUser] = startIndex[currentUser] - messagePerPage;
        }

        // Appel à une fonction pour charger les messages précédents depuis le serveur
        // fetchPreviousMessages();
        // Ou si vous avez déjà les messages précédents en mémoire, vous pouvez simplement appeler DisplayMessages avec les messages appropriés
        //console.log("fetching previous messages")
        let scrollHeight = this.scrollHeight;

        RemoveMessages();
        DisplayMessages(allPrivateMessages, currentUser);
        this.scrollTop = this.scrollHeight - scrollHeight;
    }
});

const scrollToBottom = () => {
    var chat = document.getElementById("chat_box");
    chat.scrollTop = chat.scrollHeight;
}

// Typing indicator


const typingEvent = () => {
    //var input = document.querySelectorAll("messageContent");
    /* document.addEventListener("keypress", function () {
        console.log("keypress");
    }); */

        //console.log("typing");
    var senderSession = getCookieValue();// get the cookie value to get the session id of the sender
    var receiver = document.getElementById("chat_user").textContent;// get the receiver from the chat window
    var data = ({
        "senderUUID": senderSession,
        "receiverNickname": receiver,
        
    });
    var type = "is_typing"
    var msg = JSON.stringify({
        type: type,
        data: data
    })
    socket.send(msg);
}

const typingIndicator = () => {
    //get the name of the user who is typing
    var typingUser = document.getElementById("chat_user").textContent;

    var typingIndicator = document.getElementById("typing_indicator");
        typingIndicator.innerHTML = typingUser + " is typing...";    
        typingIndicator.style.display = "block";

    clearTimeout(typingTimeout);

    
    typingTimeout = setTimeout(() => {
        typingIndicator.style.display = "none";
    }, 2250);

    // if the message has been sent, remove the typing indicator
    var chat = document.getElementById("chat_box");
    chat.addEventListener("scroll", () => {
        typingIndicator.style.display = "none";
    });
}

document.addEventListener("keydown", () => {
    clearTimeout(typingTimeout);
});


/* ------------------------------------ Post ---------------------------------------- */

const closePostWindow = () => {
    var postWindow = document.getElementById("new_post");
    postWindow.style.display = "none";
    postWindow.style.zIndex = "-100";
}


const submitFormPost = () => {
    // get the cookie value to get the session id of the sender
    var senderSession = getCookieValue();

    var form = document.getElementById("post_form");

    if (form) {
        var title = document.getElementById("new_post_title").value;
        var content = document.getElementById("new_post_content").value;
        var data = ({
            "title": title,
            "content": content,
            "author": senderSession
        });
        var type = "new_post"
        var msg = JSON.stringify({
            type: type,
            data: data
        })
        socket.send(msg);
        form.reset();

        //displayPost(data);
    } else {
        console.log("form not found")
    }
}

const DisplayPosts = (data) => {
    RemoveAllPosts();

    for (var i = 0; i < data.length; i++) {
        var post = data[i];
        var postDiv = document.createElement("div");
        postDiv.className = "post";
        // each post must have an id to be able to delete it
        postDiv.id = "post" + post.ID;

        var postTitle = document.createElement("h3");
        postTitle.innerHTML = post.Title;
        var postContent = document.createElement("p");
        postContent.innerHTML = post.Content;
        var postAuthor = document.createElement("p");
        postAuthor.innerHTML = post.Author;
        postAuthor.classList.add("post_author");
        var postComment = document.createElement("p");
        postComment.innerHTML = "Comments";
        postComment.classList.add("post_comment");

        postComment.addEventListener("click", (e) => {
            var postId = e.target.parentNode.id.replace("post", "");
            var commentWindow = document.createElement("div");
            commentWindow.id = "comments";
            var headDiv = document.createElement("div");
            headDiv.className = "comment_head";
            var commentTitle = document.createElement("h3");
            commentTitle.innerHTML = "Comments";
            var commentExit = document.createElement("i");
            commentExit.id = "comment_exit";
            commentExit.classList.add("fa");
            commentExit.classList.add("fa-times-circle");
            commentExit.addEventListener("click", () => {
                // delete the comment window
                commentWindow.remove();
            });
            var commentForm = document.createElement("form");
            commentForm.id = "comment_form";
            var commentInput = document.createElement("input");
            commentInput.type = "text";
            commentInput.id = "new_comment";
            commentInput.placeholder = "Write a comment...";
            var commentSubmit = document.createElement("input");
            commentSubmit.type = "submit";
            commentSubmit.value = "Send";
            var comments_list = document.createElement("div");
            comments_list.id = "comments_list";

            GetComments(postId);

            // let the user push enter to send the comment
            commentForm.addEventListener("submit", function (e) {

                e.preventDefault();

                var comment = document.getElementById("new_comment").value;
                var data = ({
                    "content": comment,
                    "author": getCookieValue(),
                    "postId": postId
                });
                var type = "new_comment"
                var msg = JSON.stringify({
                    type: type,
                    data: data
                })
                socket.send(msg);
                commentForm.reset();

            });

            headDiv.appendChild(commentTitle);
            headDiv.appendChild(commentExit);

            commentForm.appendChild(commentInput);
            commentForm.appendChild(commentSubmit);

            commentWindow.appendChild(headDiv);
            commentWindow.appendChild(commentForm);
            commentWindow.appendChild(comments_list);

            document.body.appendChild(commentWindow);
        });

        postDiv.appendChild(postTitle);
        postDiv.appendChild(postContent);
        postDiv.appendChild(postAuthor);
        postDiv.appendChild(postComment);

        posts.appendChild(postDiv);
    }
}

const RemoveAllPosts = () => {
    if (posts) {
        while (posts.firstChild) {
            posts.removeChild(posts.firstChild);
        }
    }
}

/* ------------------------------------ Comment ---------------------------------------- */

const GetComments = (postId) => {
    RemoveAllComments();

    var data = ({
        "postId": postId
    });
    var type = "get_comments"
    var msg = JSON.stringify({
        type: type,
        data: data
    })
    socket.send(msg);
}

const DisplayComments = (data) => {
    RemoveAllComments();
    var comments_list = document.getElementById("comments_list");
    for (var i = 0; i < data.length; i++) {
        var comment = data[i];
        var commentDiv = document.createElement("div");
        commentDiv.className = "comment";
        // each comment must have an id to be able to delete it
        commentDiv.id = "comment" + comment.Id;

        commentDiv.style = "border: 1px solid #ee00ff; margin: 10px; padding: 10px;";
        var commentContent = document.createElement("p");
        commentContent.innerHTML = comment.Content;
        commentContent.style = "margin: 15px;";
        var commentAuthor = document.createElement("p");
        commentAuthor.innerHTML = comment.Author;
        commentAuthor.style.textAlign = "right";
        commentAuthor.style.marginRight = "15px";
        commentAuthor.classList.add("comment_author");

        commentDiv.appendChild(commentContent);
        commentDiv.appendChild(commentAuthor);

        comments_list.appendChild(commentDiv);
    }
    document.getElementById("comments_list").scrollTop = document.getElementById("comments_list").scrollHeight;
}

const RemoveAllComments = () => {
    var comments_list = document.getElementById("comments_list");
    if (comments_list) {
        while (comments_list.firstChild) {
            comments_list.removeChild(comments_list.firstChild);
        }
    }
}




/* ------------------------------------ Cookies ---------------------------------------- */

const getCookieValue = () => {
    if (document.cookie.length != 0) {
        var array = document.cookie.split("=");
        var UUID = array[1];
        return UUID;
    }
    else {
        console.log("Cookie not available");
    }
}

// async function exemplJS()

const exemplJS = async () => {
    await fetch("http://localhost:8080/exemple")
}
