package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool)
var defaultPort = 6969

func main() {
	http.HandleFunc("/", handleWebSocket)
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		port = defaultPort
	}
	fmt.Printf("Server is listening on :%d!\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	clients[conn] = true
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		for client := range clients {
			if client != conn {
				if err := client.WriteMessage(messageType, p); err != nil {
					fmt.Println(err)
					delete(clients, client)
					client.Close()
				}
			}
		}
	}
}
