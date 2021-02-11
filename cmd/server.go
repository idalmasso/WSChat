package main

import (
	"flag"
	"log"
	"net/http"

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



func handler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query()["id"][0]
	log.Println("UserJoined :", user)
  u, err := services.AddUser(user, w,r); 
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Start listening to messages from user
	go u.ReceiveMessages()

}
