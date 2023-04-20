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
	if _, ok := wss.rooms[roomId].Clients[nick]; ok {
		return true
	}
	return false
}

// sendInTheRoom 방에 있는 사용자에게 메시지를 보낸다.
func (wss *WebsocketServer) sendInTheRoom(roomId int, msg *Message) {
	for _, client := range wss.rooms[roomId].Clients {
		client.send <- msg
	}
}

// sendInTheRoomExceptSender Sender를 제외한 방에 있는 사용자에게 메시지를 보낸다.
func (wss *WebsocketServer) sendInTheRoomExceptSender(roomId int, msg *Message) {
	for _, client := range wss.rooms[roomId].Clients {
		if client.Nick != msg.Sender {
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
	for _, room := range wss.rooms {
		for nick, client := range room.Clients {
			if nick == msg.Data {
				client.send <- &Message{Action: "error", Data: "nick duplicate"}
				return
			}
		}
	}

	// 닉네임 변경
	client := wss.rooms[msg.RoomId].Clients[msg.Sender]
	client.Nick = msg.Data
	wss.rooms[msg.RoomId].Clients[client.Nick] = client
	delete(wss.rooms[msg.RoomId].Clients, msg.Sender)

	// 방에 입장한 사용자에게 보내기
	wss.sendInTheRoom(msg.RoomId, msg)
}

// newJoinRoom 방 생성, 입장 처리
func (wss *WebsocketServer) newJoinRoom(msg *Message) {
	roomId := -1

	player, ok := Manager.players[msg.Sender]
	if !ok {
		Error.Println(msg.Sender, " player not found")
		return
	}

	oldRoomId := player.RoomId

	if msg.Action == "new-room" {
		// 새로운 방을 생성한다.
		roomId = wss.getFreeRoomID()
		if roomId == -1 {
			Error.Println("room id not found")
			return
		}
		wss.rooms[roomId] = NewRoomInfo(roomId, player.Client, "room "+strconv.Itoa(roomId))
	} else {
		// 방에 입장한다.
		roomId, _ = strconv.Atoi(msg.Data)
		wss.rooms[roomId].Clients[msg.Sender] = player.Client
	}

	player.RoomId = roomId

	// 개임을 생성한다.
	Manager.NewGame(roomId, player.Client)

	// 이전방에서 삭제
	delete(wss.rooms[oldRoomId].Clients, msg.Sender)
	if len(wss.rooms[msg.RoomId].Clients) == 0 {
		delete(wss.rooms, oldRoomId)
	}

	// 방에 입장한 사용자에게 보내기
	msg.Action = "join-room"
	msg.RoomId = roomId
	msg.RoomList = append(msg.RoomList, *wss.rooms[roomId])

	wss.sendInTheRoom(roomId, msg)
}

func (wss *WebsocketServer) OutRoom(roomId int, nick string) {

	if _, ok := wss.rooms[roomId]; !ok {
		Error.Println("room not found, roomID:", roomId)
		return
	}

	delete(wss.rooms[roomId].Clients, nick)
	if roomId != WAITITNG_ROOM && len(wss.rooms[roomId].Clients) == 0 {
		// delete room
		delete(wss.rooms, roomId)
	} else {
		// change owner
		if wss.rooms[roomId].Owner == nick {
			for _, client := range wss.rooms[roomId].Clients {
				wss.rooms[roomId].Owner = client.Nick
				break
			}
		}

		// 방에 입장한 사용자에게 보내기
		msg := &Message{
			Action: "leave-room",
			RoomId: roomId,
			Sender: nick}

		wss.sendInTheRoom(msg.RoomId, msg)
	}
}

// leaveRoom 방 나가기 처리
func (wss *WebsocketServer) leaveRoom(msg *Message) {

	client := Manager.getClient(msg.Sender)
	if client == nil {
		Error.Println(msg.Sender, " player not found, roomID:", msg.RoomId)
		return
	}

	// 대기실로 이동
	wss.rooms[WAITITNG_ROOM].Clients[msg.Sender] = client

	// 게임 중이면 종료
	if client.game.IsPlaying() {
		client.game.Stop()
	}

	// 방에서 나가기
	wss.OutRoom(msg.RoomId, msg.Sender)

	// 대기실 입장 메시지 보내기
	msg.Action = "join-room"
	msg.RoomId = WAITITNG_ROOM
	msg.RoomList = append(msg.RoomList, *wss.rooms[WAITITNG_ROOM])

	wss.sendInTheRoom(WAITITNG_ROOM, msg)
}

// listRoom 방 목록 보기 처리
func (wss *WebsocketServer) listRoom(msg *Message) {

	msg.RoomList = make([]RoomInfo, 0, len(wss.rooms))

	for _, roomInfo := range wss.rooms {
		msg.RoomList = append(msg.RoomList, *roomInfo)
	}

	// 요청한 사용자에게 보내기
	wss.rooms[msg.RoomId].Clients[msg.Sender].send <- msg
}

func (wss *WebsocketServer) startGame(msg *Message) {
	for _, client := range wss.rooms[msg.RoomId].Clients {
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
	Trace.Println("HandleMessage:", msg.Action, msg.Sender)

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
		wss.sendInTheRoom(msg.RoomId, msg)

	case "gift-full-blocks":
		wss.sendInTheRoomExceptSender(msg.RoomId, msg)

	case "start-game":
		wss.startGame(msg)

	case "block-drop", "block-rotate", "block-left", "block-right", "block-down":
		wss.actionGame(msg)

	default:
		Warning.Println("Unknown Action:", msg)
	}
}
