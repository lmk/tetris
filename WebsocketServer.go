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
	clients    map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewWebsocketServer() *WebsocketServer {

	// 웹 소켓 서버를 생성한다.
	return &WebsocketServer{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (wsServer *WebsocketServer) Run() {
	for {

		select {

		case client := <-wsServer.register:
			connections := wsServer.clients[client.roomId]
			if connections == nil {
				connections = make(map[*Client]bool)
				wsServer.clients[client.roomId] = connections
			}
			wsServer.clients[client.roomId][client] = true
			Info.Printf("\t register %s: %d, %v", client.roomId, len(wsServer.clients[client.roomId]), client)

		case client := <-wsServer.unregister:
			if _, ok := wsServer.clients[client.roomId]; ok {
				delete(wsServer.clients[client.roomId], client)
				close(client.send)
			}
			Info.Printf("\t unregister %s: %d", client.roomId, len(wsServer.clients[client.roomId]))

		case message := <-wsServer.broadcast:
			wsServer.HandleMessage(message)
		}
	}
}

func (wsServer *WebsocketServer) HandleMessage(message *Message) {
	switch message.Action {
	case "new-room":
		roomId := len(wsServer.clients) + 1
		Info.Printf("new-room: %d", roomId)

		// 새로만든 룸으로 이동

	case "list-room":

		message.RoomInfo = make([]RoomInfo, 0)

		for roomId, client := range wsServer.clients {
			if roomId == "list" { // list는 제외
				continue
			}

			roomInfo := RoomInfo{roomId, make([]string, 0)}

			for c := range client {
				roomInfo.Nicks = append(roomInfo.Nicks, c.nick)
			}

			message.RoomInfo = append(message.RoomInfo, roomInfo)
		}

		// 요청한 사용자에게 보내기
		for client := range wsServer.clients["list"] {
			Trace.Printf("CHECK nick %v, %v", client.nick, message.Sender)
			if client.nick == message.Sender {
				Info.Printf("send list-room: %v", message)
				client.send <- message
			}
		}

		Info.Printf("list-room: %v", message)

	default:
		Warning.Println("Unknown Action:", message)
	}
}

func serveWs(ctx *gin.Context, roomId string, wsServer *WebsocketServer) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Error.Fatalln(err)
	}
	client := NewClient(roomId, conn, wsServer)

	wsServer.register <- client

	go client.Write()
	go client.Read()
}
