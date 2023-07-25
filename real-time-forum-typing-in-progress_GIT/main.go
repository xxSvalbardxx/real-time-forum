package main

import (
	"fmt"
	"net/http"

	"forum/internal/database"
	websocket "forum/internal/ws"

	_ "github.com/mattn/go-sqlite3"
)

/*
func Suppr() {
    exec.Command("bash", "-c", "echo \"drop table posts; drop table users\" | sqlite3 foo.db").Run()
}
*/

// rewrite 

var db, _ = database.InitDB()

func main() {
	defer db.Close()
	database.CreateUserTable(db)
	database.CreatePostTable(db)
	database.CreateCommentTable(db)
	database.CreatePrivateMessageTable(db)

	http.Handle("/", http.FileServer(http.Dir("public")))
	// serve the images folder
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/ws", websocket.WebSocketHandler)
	
	http.HandleFunc("/exemple", Exemple)

	fmt.Println("Server is running on : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}


func Exemple(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Exemple")
}


