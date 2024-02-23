package main

import (
	"encoding/json" // Added for JSON encoding
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{}

func main() {

	initDB()

	r := mux.NewRouter()
	Update()
	r.HandleFunc("/ws", handleWebSocket)

	http.ListenAndServe(":8080", r)
}

// Database instance
var db *gorm.DB

// Database setup
func initDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/websocket?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err) // Handle errors gracefully in production
	}
}

// update for every 3 second
func Update() {

	user := User{ID: 1, Username: "alex"}
	db.Create(user)
	fmt.Println("is work ")

}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func RepoQuery() ([]User, error) {
	user := []User{}
	err := db.Find(&user).Error // Use `&users` to pass a pointer for modification
	return user, err
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		data, err := RepoQuery()
		if err != nil {
			fmt.Println(err)
			break
		}
		active := false
		ok := len(data)
		if len(data) <= ok {
			Update(bool)

		}
		// Marshal data to JSON for WebSocket transmission
		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		err = conn.WriteMessage(websocket.TextMessage, bytes)

		if err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(3 * time.Second)
	}
}
