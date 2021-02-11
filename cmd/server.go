package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/idalmasso/WSChat/pkg/services"
)

var addr = flag.String("addr", "0.0.0.0:2000", "http service address")

func main() {
	flag.Parse()

	http.HandleFunc("/chat", handler)
	log.Printf("Listening at %s", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query()["id"][0]
	log.Println("UserJoined :", user)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
	
	u:=services.User{Username: user, Conn: conn}
	services.AddUser(&u)
	// Start listening to messages from user
	go u.ReceiveMessages()

}
