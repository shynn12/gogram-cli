package models

import "time"

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`
}

type UserDTO struct {
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`
}

type Message struct {
	ID     int       `json:"message_id"`
	UserID int       `json:"user_id"`
	Body   string    `json:"body"`
	Time   time.Time `json:"time"`
}

type MessageDTO struct {
	UserID int       `json:"user_id"`
	Body   string    `json:"body"`
	ChatID int       `json:"chat_id"`
	Time   time.Time `json:"time"`
}

type Messages struct {
	Msgs []MessageDTO `json:"messages"`
}

type Chat struct {
	ID   int    `json:"chat_id"`
	Name string `json:"chat_name"`
}

type Chats struct {
	Chats []Chat `json:"chats"`
}

type Error struct {
	Text string `json:"err_text"`
}
