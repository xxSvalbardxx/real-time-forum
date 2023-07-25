package utils

import "time"

type Session struct {
	ID           string
	UserID       int
	Username     string
	LastActivity time.Time
}

type Message struct {
	ID         int
	SenderUUID string
	Sender     string
	Receiver   string
	Content    string
	Date       string
}

type MessageArray struct {
	Type string
	Data []Message
}

type Post struct {
	ID      int
	Title   string
	Content string
	Author  string
}

type PostArray struct {
	Type string
	Data []Post
}

type Comment struct {
	ID      int
	Content string
	Author  string
	PostID  int
}

type CommentArray struct {
	Type string
	Data []Comment
}
