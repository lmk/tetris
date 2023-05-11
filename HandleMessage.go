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
		for nick := range room.Clients {
			if nick == msg.Data {
				Manager.players[msg.Sender].Client.send <- &Message{Action: "error", Data: "nick duplicate"}
				return
			}
		}
	}

	// 닉네임 변경
	client := wss.rooms[msg.RoomId].Clients[msg.Sender]
	client.Nick = msg.Data
	wss.rooms[msg.RoomId].Clients[client.Nick] = client
	delete(wss.rooms[msg.RoomId].Clients, msg.Sender)

	// 방장이면 owner도 변경
	if wss.rooms[msg.RoomId].Owner == msg.Sender {
		wss.rooms[msg.RoomId].Owner = msg.Data
	}

	// manager도 변경
	Manager.players[msg.Data] = Manager.players[msg.Sender]
	delete(Manager.players, msg.Sender)

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
	}

	// 개임을 생성한다.
	Manager.NewGame(roomId, player.Client)

	if !wss.OutRoom(oldRoomId, msg.Sender) {
		Error.Printf("[newJoinRoom] out room fail, roomID:%d, user:%s", oldRoomId, msg.Sender)
	}

	if !wss.InRoom(roomId, msg.Sender) {
		Error.Printf("[newJoinRoom] in room fail, roomID:%d, user:%s", roomId, msg.Sender)
	}

	wss.RefreshWaitingRoom()
}

func (wss *WebsocketServer) InRoom(roomId int, nick string) bool {
	player, ok := Manager.players[nick]
	if !ok {
		Error.Println(nick, " player not found")
		return false
	}

	player.RoomId = roomId
	wss.rooms[roomId].Clients[nick] = player.Client

	// 방에 입장한 사용자에게 보내기
	msg := &Message{
		Action: "join-room",
		RoomId: roomId,
		Sender: nick}
	msg.RoomList = appendRoomInfo(msg.RoomList, wss.rooms[roomId])

	wss.sendInTheRoom(roomId, msg)

	// 이미 play 중이면 옵져버 상태로 설정
	if wss.rooms[roomId].GetState() == "playing" {
		player.Client.Game.Ch <- &Message{Action: "observer"}
	}

	return true
}

func (wss *WebsocketServer) OutRoom(roomId int, nick string) bool {

	if _, ok := wss.rooms[roomId]; !ok {
		Error.Println("room not found, roomID:", roomId)
		return false
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

			if wss.rooms[roomId].Owner == nick && roomId == WAITITNG_ROOM {
				wss.rooms[roomId].Owner = ""
			}
		}

		// 방에 입장한 사용자에게 보내기
		msg := &Message{
			Action: "leave-room",
			RoomId: roomId,
			Sender: nick}
		msg.RoomList = appendRoomInfo(msg.RoomList, wss.rooms[roomId])

		wss.sendInTheRoom(msg.RoomId, msg)

		// 게임 오버 전파
		if roomId != WAITITNG_ROOM && wss.rooms[roomId].State == "playing" && Manager.getClient(nick).Game.IsPlaying() {
			msg.Action = "over-game"
			Manager.overGame(msg)
		}
	}

	return true
}

func (wss *WebsocketServer) RefreshWaitingRoom() {
	msg := &Message{
		Action: "list-room",
		RoomId: WAITITNG_ROOM,
		Sender: "server"}
	msg.RoomList = wss.getAllRoomInfo()
	wss.sendInTheRoom(WAITITNG_ROOM, msg)
}

// leaveRoom 방 나가기 처리
func (wss *WebsocketServer) leaveRoom(msg *Message) {

	client := Manager.getClient(msg.Sender)
	if client == nil {
		Error.Println(msg.Sender, " player not found, roomID:", msg.RoomId)
		return
	}

	// 게임 중이면 종료
	if client.Game.IsPlaying() {
		client.Game.Ch <- &Message{
			Action: "stop-game",
		}
	}

	// 방에서 나가기
	if !wss.OutRoom(msg.RoomId, msg.Sender) {
		Error.Printf("[leaveRoom] out room fail, roomID:%d, user:%s", msg.RoomId, msg.Sender)
	}

	// 대기실로 이동
	if !wss.InRoom(WAITITNG_ROOM, msg.Sender) {
		Error.Printf("[leaveRoom] in room fail, roomID:%d, user:%s", WAITITNG_ROOM, msg.Sender)
	}

	wss.RefreshWaitingRoom()
}

func (wss *WebsocketServer) getAllRoomInfo() []RoomInfo {
	roomList := make([]RoomInfo, 0, len(wss.rooms))

	for _, roomInfo := range wss.rooms {
		roomList = appendRoomInfo(roomList, roomInfo)
	}

	return roomList
}

// listRoom 방 목록 보기 처리
func (wss *WebsocketServer) listRoom(msg *Message) {

	msg.RoomList = wss.getAllRoomInfo()

	// 요청한 사용자에게 보내기
	wss.rooms[msg.RoomId].Clients[msg.Sender].send <- msg
}

func (wss *WebsocketServer) listRank(msg *Message) {
	count, err := strconv.Atoi(msg.Data)
	if err != nil {
		Error.Println("listRank count error:", err, msg.Data)
		count = 5
	}
	msg.RankList = Manager.getRankList(count)

	if msg.RankList != nil && len(msg.RankList) > 0 {
		wss.rooms[msg.RoomId].Clients[msg.Sender].send <- msg
	}
}

func (wss *WebsocketServer) startGame(msg *Message) {
	for _, client := range wss.rooms[msg.RoomId].Clients {
		client.Game.Start()
	}
}

// actionGame 게임 동작 처리: "block-drop", "block-rotate", "block-left", "block-right", "block-down"
func (wss *WebsocketServer) actionGame(msg *Message) {
	game := Manager.getGame(msg.Sender)
	if game == nil {
		Warning.Println("Unknown player:", msg.Sender)
		return
	}

	if !game.IsPlaying() {
		Warning.Println("Not playing:", msg.Sender)
		return
	}
	Trace.Println("actionGame s:", msg.Action, msg.Sender, len(game.Ch))
	game.Ch <- msg
	Trace.Println("actionGame e:", msg.Action, msg.Sender)
}

// addBot 봇 추가
func (wss *WebsocketServer) addBot(msg *Message) {
	// BotFather에게 봇 추가 요청
	BotFather.fromManager <- msg

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

	case "list-rank":
		wss.listRank(msg)

	case "over-game", "sync-game", "end-game":
		wss.sendInTheRoom(msg.RoomId, msg)

	case "gift-full-blocks":
		wss.sendInTheRoomExceptSender(msg.RoomId, msg)

	case "start-game":
		wss.startGame(msg)

	case "block-drop", "block-rotate", "block-left", "block-right", "block-down":
		wss.actionGame(msg)

	case "add-bot":
		wss.addBot(msg)

	default:
		Warning.Println("Unknown Action:", msg)
	}
}
