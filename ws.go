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

type Message struct {
	Action string `json:"action"`
	Sender string `json:"sender"`
	Msg    string `json:"msg,omitempty"`
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

	Info.Printf("Run Websocket Server")

	for {

		select {

		case client := <-wsServer.register:
			Info.Printf("Run Websocket Server register")
			connections := wsServer.clients[client.roomId]
			if connections == nil {
				connections = make(map[*Client]bool)
				wsServer.clients[client.roomId] = connections
			}
			wsServer.clients[client.roomId][client] = true
			Info.Printf("\t %s: %d", client.roomId, len(wsServer.clients[client.roomId]))

		case client := <-wsServer.unregister:
			Info.Printf("Run Websocket Server unregister")
			if _, ok := wsServer.clients[client.roomId]; ok {
				delete(wsServer.clients[client.roomId], client)
				close(client.send)
			}
			Info.Printf("\t %s: %d", client.roomId, len(wsServer.clients[client.roomId]))

		case message := <-wsServer.broadcast:
			wsServer.HandleMessage(message)
		}
	}
}

func (wsServer *WebsocketServer) HandleMessage(message *Message) {

	Info.Printf("HandleMessage: %v", message)

	switch message.Action {
	case "new-room":
		roomId := len(wsServer.clients) + 1
		Info.Printf("new-room: %d", roomId)

		// 새로만든 룸으로 이동

	case "list-room":
		Info.Printf("list-room: ")

	default:
		Warning.Println("Unknown Action:", message)
	}
}

func serveWs(ctx *gin.Context, roomId string, wsServer *WebsocketServer) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		Error.Fatalln(err)
	}
	defer conn.Close()
	/*
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				break
			}
			fmt.Printf("recv:%s", message)
			err = conn.WriteMessage(mt, message)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		}
	*/
	client := NewClient(roomId, conn, wsServer)

	Info.Println("New Client: ", roomId)

	wsServer.register <- client

	// sleep 1 초
	//time.Sleep(time.Second * 1)

	go client.Write()
	go client.Read()

}
