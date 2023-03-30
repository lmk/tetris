package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	roomId string
	nick   string
	socket *websocket.Conn
	send   chan *Message
	ws     *WebsocketServer
}

func NewClient(roomId string, conn *websocket.Conn, ws *WebsocketServer) *Client {

	// 웹 소켓 클라이언트를 생성한다.
	return &Client{
		roomId: roomId,
		nick:   "",
		socket: conn,
		ws:     ws,
		send:   make(chan *Message),
	}
}

func (client *Client) Read() {

	defer func() {
		client.ws.unregister <- client
		client.socket.Close()
	}()

	for {
		var msg Message
		err := client.socket.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read Error:", err)
			break
		}

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
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := client.socket.WriteJSON(message)
			if err != nil {
				fmt.Println("Write Error:", err)
				return
			}
		}
	}
}
