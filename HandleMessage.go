package main

import (
	"encoding/json"
	"strconv"
)

func (wss *WebsocketServer) isVaildRoomId(roomId int) bool {
	if _, ok := wss.rooms[roomId]; ok {
		return true
	}
	return false
}

func (wss *WebsocketServer) isVaildNick(roomId int, nick string) bool {
	if _, ok := wss.rooms[roomId][nick]; ok {
		return true
	}
	return false
}

// sendInTheRoom 방에 있는 사용자에게 메시지를 보낸다.
func (wss *WebsocketServer) sendInTheRoom(roomId int, message *Message) {
	for _, client := range wss.rooms[roomId] {
		client.send <- message
	}
}

// setNick 닉네임을 변경한다.
func (wss *WebsocketServer) setNick(message *Message) {
	// 기존과 같은 닉네임이면 변경하지 않는다.
	if message.Sender == message.Data {
		return
	}

	// nick 중복 체크
	for _, nicks := range wss.rooms {
		for nick := range nicks {
			if nick == message.Data {
				wss.rooms[message.roomId][message.Sender].send <- &Message{Action: "error", Data: "nick duplicate"}
				return
			}
		}
	}

	// 닉네임 변경
	client := wss.rooms[message.roomId][message.Sender]
	client.nick = message.Data
	wss.rooms[message.roomId][client.nick] = client
	delete(wss.rooms[message.roomId], message.Sender)

	// 방에 입장한 사용자에게 보내기
	wss.sendInTheRoom(message.roomId, message)
}

// newJoinRoom 방 생성, 입장 처리
func (wss *WebsocketServer) newJoinRoom(message *Message) {
	roomId := -1
	if message.Action == "new-room" {
		// 새로운 방을 생성한다.
		roomId = len(wss.rooms) + 1
		wss.rooms[roomId] = make(map[string]*Client)
	} else {
		roomId, _ = strconv.Atoi(message.Data)
	}

	// 방으로 이동
	client := wss.rooms[WAITITNG_ROOM][message.Sender]
	client.roomId = roomId

	// 개임을 생성한다.
	client.game = NewGame(client)

	// next block indexs를 보낸다.
	jsonNexts, _ := json.Marshal(client.game.nextBlockIndexs)
	client.send <- &Message{Action: "next-block", Data: string(jsonNexts)}

	wss.rooms[roomId][message.Sender] = client

	// 대기실에서 삭제
	delete(wss.rooms[WAITITNG_ROOM], message.Sender)

	// 방에 입장한 사용자에게 보내기
	message.Action = "join-room"
	wss.sendInTheRoom(roomId, message)
}

// leaveRoom 방 나가기 처리
func (wss *WebsocketServer) leaveRoom(message *Message) {
	roomId, _ := strconv.Atoi(message.Data)

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
	wss.rooms[message.roomId][message.Sender].send <- message
}

// // newBlock 블록 생성 처리
// func (wss *WebsocketServer) newBlock(message *Message) {
// 	// 생성할 블록의 개수를 가져온다.
// 	count, _ := strconv.Atoi(message.Msg)
// 	if count <= 0 {
// 		Warning.Println("reuqest new-block count is 0")
// 		return
// 	}

// 	// 0~7 사이의 난수를 생성한다.
// 	nums := make([]int, 0)
// 	for i := 0; i < count; i++ {
// 		nums = append(nums, rand.Int()%8)
// 	}

// 	// 방에 입장한 사용자에게 보내기
// 	message.Action = "new-block"
// 	message.Msg = strings.Join(strings.Fields(fmt.Sprint(nums)), ",")
// 	wss.sendInTheRoom(message.roomId, message)
// }

func (wss *WebsocketServer) game(message *Message) {

	switch message.Action {
	case "start-game":
		for _, client := range wss.rooms[message.roomId] {
			go client.game.Run()

			client.send <- message
		}

	case "next-block":
		//

	case "over-game":
		wss.sendInTheRoom(message.roomId, message)
	}
}

// HandleMessage 메시지 핸들러
func (wss *WebsocketServer) HandleMessage(message *Message) {
	Trace.Println("HandleMessage:", message)

	switch message.Action {
	case "set-nick":
		go wss.setNick(message)

	case "new-room", "join-room":
		go wss.newJoinRoom(message)

	case "leave-room":
		go wss.leaveRoom(message)

	case "list-room":
		go wss.listRoom(message)

	case "start-game", "over-game", "next-block":
		go wss.game(message)

	default:
		Warning.Println("Unknown Action:", message)
	}
}
