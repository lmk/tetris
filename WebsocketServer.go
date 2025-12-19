package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	MAX_ROOMS = 100 // maximum number of rooms
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// 브라우저가 아닌 클라이언트는 Origin 헤더가 없을 수 있음
			return true
		}

		// AllowedOrigins 확인
		for _, allowed := range conf.AllowedOrigins {
			// "*"는 모든 출처 허용 (개발 환경용)
			if allowed == "*" {
				return true
			}

			// 정확한 출처 매칭
			if origin == allowed {
				return true
			}

			// 포트 포함 매칭 (예: http://localhost:8090)
			if origin == allowed+":"+fmt.Sprintf("%d", conf.Port) {
				return true
			}
		}

		Warning.Printf("CORS: Rejected origin: %s", origin)
		return false
	},
}

type WebsocketServer struct {
	rooms      map[int]*RoomInfo // roomId, roomInfo
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex // protects rooms map
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
	wss.mu.RLock()
	defer wss.mu.RUnlock()

	for i := 1; i < MAX_ROOMS; i++ {
		if _, ok := wss.rooms[i]; !ok {
			return i
		}
	}

	return -1
}

func (wss *WebsocketServer) Report() {

	for range time.Tick(time.Second * 10) {

		wss.mu.RLock()
		roomCount := len(wss.rooms)
		waitingRoom := wss.rooms[WAITITNG_ROOM]
		wss.mu.RUnlock()

		if roomCount == 0 || (roomCount == 1 && waitingRoom != nil && len(waitingRoom.Clients) == 0) {
			continue
		}

		var report strings.Builder
		wss.mu.RLock()
		for roomId, info := range wss.rooms {
			fmt.Fprintf(&report, "[%v:%s:[", roomId, info.Owner)
			first := true
			for nick := range info.Clients {
				if !first {
					report.WriteString(",")
				}
				report.WriteString(nick)
				first = false
			}
			report.WriteString("]]")
		}
		wss.mu.RUnlock()

		Info.Println("REPORT:" + report.String())
	}
}

func (wss *WebsocketServer) Run() {

	Info.Println("Websocket Server is running...")

	for {

		select {

		case client := <-wss.register:
			wss.mu.Lock()
			room := wss.rooms[WAITITNG_ROOM]
			if room == nil {
				wss.rooms[WAITITNG_ROOM] = NewRoomInfo(WAITITNG_ROOM, client, "Watting Room")
			} else {
				room.Clients[client.Nick] = client
			}
			wss.mu.Unlock()

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
