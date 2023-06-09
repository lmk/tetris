package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Nick   string           `json:"-"`
	socket *websocket.Conn  `json:"-"`
	send   chan *Message    `json:"-"` // to websocket client
	ws     *WebsocketServer `json:"-"`
	Game   *Game            `json:"game,omitempty"`
}

func NewClient(conn *websocket.Conn, ws *WebsocketServer) *Client {

	// 웹 소켓 클라이언트를 생성한다.
	return &Client{
		Nick:   getRandomNick(),
		socket: conn,
		ws:     ws,
		send:   make(chan *Message, MAX_CHAN),
		Game:   nil,
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
		if r := recover(); r != nil {
			Error.Println("Write panic:", r)
		}
	}()

	// client.socket.SetReadLimit(512)
	// client.socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	// client.socket.SetPongHandler(func(string) error { client.socket.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		var msg Message

		err := client.socket.ReadJSON(&msg)
		if err != nil {
			Trace.Println("Read Error", client.Nick, err)
			break
		}

		msg.RoomId = Manager.getRoomId(msg.Sender)

		if !client.ws.isVaildRoomId(msg.RoomId) || !client.ws.isVaildNick(msg.RoomId, msg.Sender) {
			Error.Println("Invalid RoomId or Nick:", msg, client.socket.RemoteAddr().String())
			client.send <- &Message{Action: "error", Data: "Invalid RoomId or Nick"}
			break
		}

		//Info.Printf("[READ] %v", msg)

		client.ws.broadcast <- &msg
	}

	Trace.Println("Read end:", client.Nick)
}

func (client *Client) Write() {

	defer func() {
		client.socket.Close()
		if r := recover(); r != nil {
			Error.Println("Write panic", client.Nick, r)
		}
	}()

	for {
		message, ok := <-client.send
		//Info.Printf("[WRITE] %v %v", ok, message)
		if !ok {
			// The ws closed the channel.
			client.socket.WriteMessage(websocket.CloseMessage, []byte{})
			break
		} else {
			err := client.socket.WriteJSON(message)
			if err != nil {
				Trace.Println("Write Error", client.Nick, err)
				break
			}
		}
	}

	Trace.Println("Write end:", client.Nick)
}
