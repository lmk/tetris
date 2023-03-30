package main

import (
	"fmt"
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
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	Action string `json:"action"`
	Sender string `json:"sender"`
	Msg    string `json:"msg,omitempty"`
}

func NewWebsocketServer() *WebsocketServer {

	// 웹 소켓 서버를 생성한다.
	return &WebsocketServer{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (wsServer *WebsocketServer) Run(addr string) {

	for {

		select {

		case client := <-wsServer.register:
			wsServer.clients[client] = true

		case client := <-wsServer.unregister:
			if _, ok := wsServer.clients[client]; ok {
				delete(wsServer.clients, client)
				close(client.send)
			}

		case message := <-wsServer.broadcast:
			wsServer.HandleMessage(message)
		}
	}
}

func (wsServer *WebsocketServer) HandleMessage(message *Message) {

	switch message.Action {
	case "join":
		for client := range wsServer.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(wsServer.clients, client)
			}
		}

	case "message":
		for client := range wsServer.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(wsServer.clients, client)
			}
		}

	default:
		fmt.Println("Unknown Action:", message)
	}
}

func serveWs(ctx *gin.Context, roomId string, wsServer *WebsocketServer) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}

	client := NewClient(roomId, conn, wsServer)
	wsServer.register <- client
	go client.Write()
	go client.Read()

}
