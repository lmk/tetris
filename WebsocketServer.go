package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketServer struct {
	rooms      map[int]*RoomInfo // roomId, roomInfo
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewWebsocketServer() *WebsocketServer {

	// 웹 소켓 서버를 생성한다.
	return &WebsocketServer{
		rooms:      make(map[int]*RoomInfo),
		broadcast:  make(chan *Message, MAX_CHAN),
		register:   make(chan *Client, MAX_CHAN),
		unregister: make(chan *Client, MAX_CHAN),
	}
}

func (wss *WebsocketServer) getFreeRoomID() int {

	for i := 1; i < 100; i++ {
		if _, ok := wss.rooms[i]; !ok {
			return i
		}
	}

	return -1
}

func (wss *WebsocketServer) report() {

	report := "REPORT:"

	for roomId, info := range wss.rooms {
		//report += fmt.Sprintf("[%v:%d:%s],", roomId, len(info.Clients), info.Owner)
		report += fmt.Sprintf("[%v:%s:[", roomId, info.Owner)
		for nick := range info.Clients {
			report += fmt.Sprintf("%s,", nick)
		}
		report = report[:len(report)-1]
		report += "]]"
	}

	Info.Println(report)
}

func (wss *WebsocketServer) Run() {

	Info.Println("Websocket Server is running...")

	for {

		select {

		case client := <-wss.register:
			room := wss.rooms[WAITITNG_ROOM]
			if room == nil {
				wss.rooms[WAITITNG_ROOM] = NewRoomInfo(WAITITNG_ROOM, client, "Watting Room")
			} else {
				room.Clients[client.Nick] = client
			}

			Manager.Register(WAITITNG_ROOM, client)

			// send random client nick to client
			client.send <- &Message{Action: "new-nick", Data: client.Nick}

			Info.Println("\t register ", client.Nick, client.socket.RemoteAddr().String())

		case client := <-wss.unregister:
			roomId := Manager.getRoomId(client.Nick)
			if roomId == -1 {
				Error.Println("unregister error: unknown user ", client.Nick, client.socket.RemoteAddr().String())
				close(client.send)
				continue
			}

			if !wss.OutRoom(roomId, client.Nick) {
				Error.Println("unregister error: unknown room ", roomId, client.Nick, client.socket.RemoteAddr().String())
			}

			Manager.Unregister(client)

			wss.RefreshWaitingRoom()

			close(client.send)

			Info.Printf("\t unregister %v:%v, %v", roomId, client.Nick, client.socket.RemoteAddr().String())

		case message := <-wss.broadcast:
			wss.HandleMessage(message)

		case <-time.After(time.Millisecond * time.Duration(10000)):
			wss.report()
		}
	}
}

func serveWs(ctx *gin.Context, roomId int, wsServer *WebsocketServer) {

	Info.Println("serveWs", roomId)

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Error.Println(err)
	}

	client := NewClient(conn, wsServer)

	wsServer.register <- client

	go client.Write()
	go client.Read()
}
