package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/idalmasso/WSChat/pkg/models"
)

//User is the type that contains data for a single user. Fill it and send it to Add user to be added to the chat list
type User struct{
	Username string 
  Conn *websocket.Conn
	sendMutex sync.Mutex
}
//My little in memory db of users in chat
var chatUsers map[string]*User = make(map[string]*User)

var mutexUsers sync.Mutex

//AddUser adds a user to the chat. 
func AddUser(user *User) error{
	mutexUsers.Lock()
	defer mutexUsers.Unlock()
	if _, ok := chatUsers[user.Username]; ok{
		return fmt.Errorf("User already exists with that username")
	}
	chatUsers[user.Username]=user
	return nil
}
//RemoveUser removes a user from the chat.
func RemoveUser(username string) {
	mutexUsers.Lock()
	defer mutexUsers.Unlock()
	if _, ok := chatUsers[username]; ok{
		delete(chatUsers,username)
	}
}
//ReceiveMessages is a routine that receiv messages from a single user (todo=> send them to all other users)
func (u *User) ReceiveMessages(){
	for{
		var chatMessage models.ChatMessage 
		if err :=u.Conn.ReadJSON(&chatMessage); err!=nil{
			fmt.Println("ERROR "+u.Username, "Cannot decode the chat message", err.Error())
			u.Close()
			RemoveUser(u.Username)
			return
		}
		if chatMessage.User!= "" && chatMessage.User!=u.Username{
			fmt.Println("ERROR", "User sent message for another user")
			return
		}
		chatMessage.User=u.Username
		fmt.Println("Received message "+chatMessage.Message+" from user "+chatMessage.User)
		dispatchMessage(chatMessage)
	}
}
//Send sends a chat message to the user
func (u *User) Send(message models.ChatMessage) error {
	u.sendMutex.Lock()
	defer u.sendMutex.Unlock()

	// log.Println("Sending message to", p.name, ":", string(msg))
	if err:=u.Conn.WriteJSON(message); err!=nil{
	  return fmt.Errorf("ERROR for %s\n%s\n%s",u.Username, "Cannot encode the chat message", err.Error())
	}
	return nil
}

// Close connection of the user
func (u *User) Close() {
	time.Sleep(1 * time.Second)
	u.Conn.Close()
}
	
func dispatchMessage(chatMessage models.ChatMessage){
	for _,user:=range(chatUsers){
		if err:= user.Send(chatMessage); err!=nil{
			fmt.Println(err.Error())
		
		}
	}
}

