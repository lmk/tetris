package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// sendInTheRoom 방에 있는 사용자에게 메시지를 보낸다.
func (wss *WebsocketServer) sendInTheRoom(roomId int, message *Message) {
	for _, client := range wss.rooms[roomId] {
		client.send <- message
	}
}

// newJoinRoom 방 생성, 입장 처리
func (wss *WebsocketServer) newJoinRoom(message *Message) {
	roomId := 0
	if message.Action == "new-room" {
		// 새로운 방을 생성한다.
		roomId = len(wss.rooms) + 1
		wss.rooms[roomId] = make(map[string]*Client)
	} else {
		roomId, _ = strconv.Atoi(message.Msg)
	}

	// 방으로 이동
	client := wss.rooms[WAITITNG_ROOM][message.Sender]
	client.roomId = roomId

	wss.rooms[roomId][message.Sender] = client

	// 대기실에서 삭제
	delete(wss.rooms[WAITITNG_ROOM], message.Sender)

	// 방에 입장한 사용자에게 보내기
	message.Action = "join-room"
	message.RoomList = make([]RoomInfo, 0)
	message.RoomList = append(message.RoomList, RoomInfo{roomId, []string{message.Sender}})
	wss.sendInTheRoom(roomId, message)
}

// leaveRoom 방 나가기 처리
func (wss *WebsocketServer) leaveRoom(message *Message) {
	roomId, _ := strconv.Atoi(message.Msg)

	// 대기실로 이동
	client := wss.rooms[roomId][message.Sender]
	client.roomId = WAITITNG_ROOM
	wss.rooms[WAITITNG_ROOM][message.Sender] = client

	// 방에서 나가기
	delete(wss.rooms[roomId], message.Sender)
	if len(wss.rooms[roomId]) == 0 {
		delete(wss.rooms, roomId)
	} else {
		// 방에 입장한 사용자에게 보내기
		message.Action = "leave-room"

		wss.sendInTheRoom(roomId, message)
	}

	// 대기실 입장 메시지 보내기
	message.Action = "join-room"
	wss.sendInTheRoom(WAITITNG_ROOM, message)
}

// listRoom 방 목록 보기 처리
func (wss *WebsocketServer) listRoom(message *Message) {
	message.RoomList = make([]RoomInfo, 0)

	for roomId, members := range wss.rooms {

		if roomId == WAITITNG_ROOM {
			continue
		}

		roomInfo := RoomInfo{roomId, make([]string, 0)}

		for nick := range members {
			roomInfo.Nicks = append(roomInfo.Nicks, nick)
		}

		message.RoomList = append(message.RoomList, roomInfo)
	}

	// 요청한 사용자에게 보내기
	message.Client.send <- message
}

// newBlock 블록 생성 처리
func (wss *WebsocketServer) newBlock(message *Message) {
	// 생성할 블록의 개수를 가져온다.
	count, _ := strconv.Atoi(message.Msg)
	if count <= 0 {
		Warning.Println("reuqest new-block count is 0")
		return
	}

	// 0~7 사이의 난수를 생성한다.
	nums := make([]int, 0)
	for i := 0; i < count; i++ {
		nums = append(nums, rand.Int()%8)
	}

	// 방에 입장한 사용자에게 보내기
	message.Action = "new-block"
	message.Msg = strings.Join(strings.Fields(fmt.Sprint(nums)), ",")
	wss.sendInTheRoom(message.Client.roomId, message)
}

// HandleMessage 메시지 핸들러
func (wss *WebsocketServer) HandleMessage(message *Message) {
	Trace.Println("HandleMessage:", message)
	switch message.Action {
	case "new-room", "join-room":
		go wss.newJoinRoom(message)

	case "leave-room":
		go wss.leaveRoom(message)

	case "list-room":
		go wss.listRoom(message)

	case "start-game", "over-game":
		// 같은 방 사용자에게 by pass
		go wss.sendInTheRoom(message.Client.roomId, message)

	case "new-block":
		go wss.newBlock(message)

	default:
		Warning.Println("Unknown Action:", message)
	}
}
