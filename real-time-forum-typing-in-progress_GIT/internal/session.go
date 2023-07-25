package session

import (
	"time"

	"forum/utils"

	"github.com/google/uuid"
)

var sessions = make(map[string]*utils.Session)

func GetUsernameBySession(sessionID string) string {
	for _, session := range sessions {
		if session.ID == sessionID {
			return session.Username
		}
	}
	return "not found"
}

func NewSessionID() string {
	return uuid.New().String()
}

// FindSession returns the session for a given session ID
func FindSessionByUserId(userID int) *utils.Session {
	for _, session := range sessions {
		if session.UserID == userID {
			return session
		}
	}
	return nil
}
func FindSessionByUsername(username string) string {
	for _, session := range sessions {
		if session.Username == username {
			return session.ID
		}
	}
	return ""
}


func NewSession(userID int, username string) *utils.Session {
	// Find and delete any existing sessions for this user

	for id, session := range sessions {
		if session.UserID == userID {
			delete(sessions, id)
		}
	}

	// Create a new session for the user
	session := &utils.Session{
		ID:           NewSessionID(),
		UserID:       userID,
		Username:     username,
		LastActivity: time.Now(),
	}
	// Add to the map of sessions
	sessions[session.ID] = session

	return session
}

func Logout(sessionID string) {
	delete(sessions, sessionID)
}
