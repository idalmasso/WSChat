package models

//ChatMessage is the type that will be passed betweeb websockets
type ChatMessage struct{
	Message string `json:"message"`
	User string `json:"user"`
}

