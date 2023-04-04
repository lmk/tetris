package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	roomId int
	nick   string
	socket *websocket.Conn
	send   chan *Message
	ws     *WebsocketServer
}

func NewClient(roomId int, conn *websocket.Conn, ws *WebsocketServer) *Client {

	// 웹 소켓 클라이언트를 생성한다.
	return &Client{
		roomId: roomId,
		nick:   getRandomNick(),
		socket: conn,
		ws:     ws,
		send:   make(chan *Message),
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomNick() string {

	n := rand.Intn(9999)

	return fmt.Sprintf("user%d", n)
}

func (client *Client) Read() {

	defer func() {
		client.ws.unregister <- client
		client.socket.Close()
	}()

	// client.socket.SetReadLimit(512)
	// client.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	// client.socket.SetPongHandler(func(string) error { client.socket.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		var msg Message

		err := client.socket.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read Error:", err)
			break
		}

		Info.Printf("[READ] %v", msg)

		client.ws.broadcast <- &msg
	}
}

func (client *Client) Write() {

	defer func() {
		client.socket.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			Info.Printf("[WRITE] %v %v", ok, message)
			if !ok {
				// The ws closed the channel.
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := client.socket.WriteJSON(message)
				if err != nil {
					Error.Println("Write Error:", err)
					return
				}
			}
		}
	}
}
