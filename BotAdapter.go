package main

import (
	"fmt"
	"strings"
	"time"

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

	botAdapter := &BotAdapter{
		toBot:   make(chan *Message, MAX_CHAN),
		fromBot: make(chan *Message, MAX_CHAN),
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

// newNick :
func (ba *BotAdapter) newNick(msg *Message) {
	ba.nick = msg.Data

	// nick을 bot으로 변경한다.
	go func() {
		<-time.NewTimer(500 * time.Millisecond).C
		ba.socket.WriteJSON(&Message{
			Action: "set-nick",
			Data:   strings.Replace(ba.nick, "user", "bot", 1),
			Sender: ba.nick,
		})
	}()

	// room에 join 한다.
	go func() {
		<-time.NewTimer(1 * time.Second).C
		ba.socket.WriteJSON(&Message{
			Action: "join-room",
			Data:   fmt.Sprintf("%d", ba.roomId),
			Sender: ba.nick,
		})
	}()

}

func (ba *BotAdapter) startGame(msg *Message) {
	ba.toBot <- msg
}

// Read : server -> bot
func (ba *BotAdapter) Read() {
	defer func() {
		close(ba.done)
		if r := recover(); r != nil {
			Error.Println("Write panic:", r)
		}
	}()
	for {

		var msg Message
		err := ba.socket.ReadJSON(&msg)
		if err != nil {
			Error.Println("websocket handelr read:", err)
			ba.socket.Close()
			return
		}

		switch msg.Action {
		case "new-nick":
			ba.newNick(&msg)

		case "set-nick":
			if msg.Sender == ba.nick {
				ba.nick = msg.Data
			}

		case "start-game":
			ba.startGame(&msg)

		case "over-game":
			if msg.Sender == ba.nick {
				ba.toBot <- &msg
			}

		default:
			ba.toBot <- &msg
		}
	}
}

// Write : bot -> server
func (ba *BotAdapter) Write() {
	defer func() {
		if r := recover(); r != nil {
			Error.Println("Write panic:", r)
		}
	}()

	for {
		select {
		case <-ba.done:
			return
		case msg, ok := <-ba.fromBot:
			if !ok {
				Error.Println("websocket handelr write:", ba.nick, "channel closed")
				return
			}

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
