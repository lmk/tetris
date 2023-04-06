package main

import (
	"net/http"

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
	rooms      map[int]map[string]*Client // clients list (roomId, nick, client)
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewWebsocketServer() *WebsocketServer {

	// 웹 소켓 서버를 생성한다.
	return &WebsocketServer{
		rooms:      make(map[int]map[string]*Client),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (wss *WebsocketServer) Run() {
	for {

		select {

		case client := <-wss.register:
			connections := wss.rooms[client.roomId]
			if connections == nil {
				connections = make(map[string]*Client)
				wss.rooms[client.roomId] = connections
			}
			wss.rooms[client.roomId][client.nick] = client
			client.send <- &Message{Action: "new-nick", Msg: client.nick}
			Info.Printf("\t register %v:%v, %v", client.roomId, client.nick, len(wss.rooms[client.roomId]))

		case client := <-wss.unregister:
			if _, ok := wss.rooms[client.roomId]; ok {
				delete(wss.rooms[client.roomId], client.nick)
				if client.roomId != WAITITNG_ROOM && len(wss.rooms[client.roomId]) == 0 {
					delete(wss.rooms, client.roomId)
				}
				close(client.send)
			}
			Info.Printf("\t unregister %v:%v, %v", client.roomId, client.nick, len(wss.rooms[client.roomId]))

		case message := <-wss.broadcast:
			wss.HandleMessage(message)
		}
	}
}

func serveWs(ctx *gin.Context, roomId int, wsServer *WebsocketServer) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Error.Fatalln(err)
	}
	client := NewClient(roomId, conn, wsServer)

	wsServer.register <- client

	go client.Write()
	go client.Read()
}
