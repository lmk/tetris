package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type BotAdapter struct {
	toBot   chan *Message
	fromBot chan *Message
	done    chan struct{}
	socket  *websocket.Conn
	roomId  int
	nick    string // bot nickname
}

func NewBotAdapter(roomId int) *BotAdapter {

	// websocket 서버에 연결한다.
	wsURI := fmt.Sprintf("ws://%s:%d/ws", conf.Domain, conf.Port)
	conn, _, err := websocket.DefaultDialer.Dial(wsURI, nil)
	if err != nil {
		Error.Println("connect:", err)
		return nil
	}
	defer conn.Close()

	botAdapter := &BotAdapter{
		toBot:   make(chan *Message),
		fromBot: make(chan *Message),
		done:    make(chan struct{}),
		roomId:  roomId,
		socket:  conn,
	}

	// server -> bot
	go botAdapter.Read()

	// bot -> server
	go botAdapter.Write()

	return botAdapter
}

// newNick : nick을 받으면 방에 입장한다.
func (ba *BotAdapter) newNick(msg *Message) {
	ba.nick = msg.Data
	ba.socket.WriteJSON(&Message{
		Action: "join-room",
		Data:   fmt.Sprintf("%d", ba.roomId),
		Sender: ba.nick,
	})
}

// Read : server -> bot
func (ba *BotAdapter) Read() {
	defer close(ba.done)
	for {
		var msg Message
		err := ba.socket.ReadJSON(&msg)
		if err != nil {
			Error.Println("websocket handelr read:", err)
			return
		}

		switch msg.Action {
		case "new-nick":
			ba.newNick(&msg)

		case "start-game":
			ba.toBot <- &msg
		}
	}
}

// Write : bot -> server
func (ba *BotAdapter) Write() {
	for {
		select {
		case <-ba.done:
			return
		case msg := <-ba.fromBot:
			err := ba.socket.WriteJSON(msg)
			if err != nil {
				Error.Println("websocket handelr write:", err)
				return
			}
			// case <-interrupt:
			// 	Info.Println("interrupt")

			// 	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			// 	if err != nil {
			// 		Error.Println("write close:", err)
			// 		return
			// 	}
			// 	select {
			// 	case <-done:
			// 	}
			// 	return
		}

	}
}
