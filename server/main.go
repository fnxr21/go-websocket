package main

import (
	"fmt"
	"log"
	"net/http"
``
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	setupRoutes()
	fmt.Println("GO WEBSOCKET")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndPoint)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Home Page")

}

func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true

	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Successfully Connected . . .")

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
