package main

import (
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
func (wss *WebsocketServer) sendInTheRoom(roomId int, msg *Message) {
	for _, client := range wss.rooms[roomId] {
		client.send <- msg
	}
}

// sendInTheRoomExceptSender Sender를 제외한 방에 있는 사용자에게 메시지를 보낸다.
func (wss *WebsocketServer) sendInTheRoomExceptSender(roomId int, msg *Message) {
	for _, client := range wss.rooms[roomId] {
		if client.nick != msg.Sender {
			client.send <- msg
		}
	}
}

// setNick 닉네임을 변경한다.
func (wss *WebsocketServer) setNick(msg *Message) {
	// 기존과 같은 닉네임이면 변경하지 않는다.
	if msg.Sender == msg.Data {
		return
	}

	// nick 중복 체크
	for _, nicks := range wss.rooms {
		for nick := range nicks {
			if nick == msg.Data {
				wss.rooms[msg.roomId][msg.Sender].send <- &Message{Action: "error", Data: "nick duplicate"}
				return
			}
		}
	}

	// 닉네임 변경
	client := wss.rooms[msg.roomId][msg.Sender]
	client.nick = msg.Data
	wss.rooms[msg.roomId][client.nick] = client
	delete(wss.rooms[msg.roomId], msg.Sender)

	// 방에 입장한 사용자에게 보내기
	wss.sendInTheRoom(msg.roomId, msg)
}

// newJoinRoom 방 생성, 입장 처리
func (wss *WebsocketServer) newJoinRoom(msg *Message) {
	roomId := -1
	if msg.Action == "new-room" {
		// 새로운 방을 생성한다.
		roomId = len(wss.rooms) + 1
		wss.rooms[roomId] = make(map[string]*Client)
	} else {
		roomId, _ = strconv.Atoi(msg.Data)
	}

	// 방으로 이동
	client := wss.rooms[WAITITNG_ROOM][msg.Sender]
	client.roomId = roomId

	// 개임을 생성한다.
	Manager.NewGame(roomId, client)

	wss.rooms[roomId][msg.Sender] = client

	// 대기실에서 삭제
	delete(wss.rooms[WAITITNG_ROOM], msg.Sender)

	// 방에 입장한 사용자에게 보내기
	msg.Action = "join-room"
	wss.sendInTheRoom(roomId, msg)
}

// leaveRoom 방 나가기 처리
func (wss *WebsocketServer) leaveRoom(msg *Message) {

	// 대기실로 이동
	client := wss.rooms[msg.roomId][msg.Sender]
	client.roomId = WAITITNG_ROOM
	wss.rooms[WAITITNG_ROOM][msg.Sender] = client

	if client.game.IsPlaying() {
		client.game.Stop()
	}

	// 방에서 나가기
	delete(wss.rooms[msg.roomId], msg.Sender)
	if len(wss.rooms[msg.roomId]) == 0 {
		delete(wss.rooms, msg.roomId)
	} else {
		// 방에 입장한 사용자에게 보내기
		msg.Action = "leave-room"

		wss.sendInTheRoom(msg.roomId, msg)
	}

	// 대기실 입장 메시지 보내기
	msg.Action = "join-room"
	wss.sendInTheRoom(WAITITNG_ROOM, msg)
}

// listRoom 방 목록 보기 처리
func (wss *WebsocketServer) listRoom(msg *Message) {
	msg.RoomList = make([]RoomInfo, 0)

	for roomId, members := range wss.rooms {

		if roomId == WAITITNG_ROOM {
			continue
		}

		roomInfo := RoomInfo{roomId, make([]string, 0)}

		for nick := range members {
			roomInfo.Nicks = append(roomInfo.Nicks, nick)
		}

		msg.RoomList = append(msg.RoomList, roomInfo)
	}

	// 요청한 사용자에게 보내기
	wss.rooms[msg.roomId][msg.Sender].send <- msg
}

func (wss *WebsocketServer) startGame(msg *Message) {
	for _, client := range wss.rooms[msg.roomId] {
		client.game.Start()
	}
}

func (wss *WebsocketServer) actionGame(msg *Message) {
	game := Manager.getGame(msg.Sender)
	if game == nil {
		Warning.Println("Unknown player:", msg.Sender)
		return
	}

	game.Action(msg)
}

// HandleMessage websocket clinet -> server의 메시지를 처리한다.
func (wss *WebsocketServer) HandleMessage(msg *Message) {
	Trace.Println("HandleMessage:", msg)

	switch msg.Action {
	case "set-nick":
		wss.setNick(msg)

	case "new-room", "join-room":
		wss.newJoinRoom(msg)

	case "leave-room":
		wss.leaveRoom(msg)

	case "list-room":
		wss.listRoom(msg)

	case "over-game", "sync-game":
		wss.sendInTheRoom(msg.roomId, msg)

	case "gift-full-blocks":
		wss.sendInTheRoomExceptSender(msg.roomId, msg)

	case "start-game":
		wss.startGame(msg)

	case "block-drop", "block-rotate", "block-left", "block-right", "block-down":
		wss.actionGame(msg)

	default:
		Warning.Println("Unknown Action:", msg)
	}
}
