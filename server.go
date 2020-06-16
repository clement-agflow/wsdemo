package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade the http connection to use the websocket protocol
	upgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Fail to upgrade http connection:", err)
	}
	defer conn.Close()

	// send message every second
	go func() {
		cnt := 1
		ticker := time.NewTicker(750 * time.Millisecond)

		for range ticker.C {
			msg := fmt.Sprintf("Toc %d", cnt)
			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Fatal("Fail to write:", err)
			}
			cnt++
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Fail to read message:", err)
		}
		log.Printf("received from client: %s", message)
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("listening on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
