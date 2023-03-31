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

	// client.socket.SetReadLimit(512)
	// client.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	// client.socket.SetPongHandler(func(string) error { client.socket.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		var msg Message

		mt, message, err := client.socket.ReadMessage()
		if err != nil {
			Error.Println("Read Error:", err)
			break
		}

		msg.Msg = fmt.Sprintf("mt: %d, msg:%s", mt, string(message))

		Info.Println("READ ", msg.Msg)

		// err := client.socket.ReadJSON(&msg)
		// if err != nil {
		// 	fmt.Println("Read Error:", err)
		// 	break
		// }

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

			Info.Println("WRITE ", message.Msg)

			err := client.socket.WriteJSON(message)
			if err != nil {
				Error.Println("Write Error:", err)
				return
			}
		}
	}
}
