package main

func (wss *WebsocketServer) HandleMessage(message *Message) {
	switch message.Action {
	case "new-room":
		roomId := len(wss.rooms) + 1

		Info.Printf("new-room: %v: %v", roomId, message.Sender)

		// 새로만든 룸으로 이동
		client := wss.rooms[WAITITNG_ROOM][message.Sender]
		client.roomId = roomId

		wss.rooms[roomId] = make(map[string]*Client)
		wss.rooms[roomId][message.Sender] = client

		// 대기실에서 삭제
		delete(wss.rooms[WAITITNG_ROOM], message.Sender)

		// 방에 입장한 사용자에게 보내기
		message.Action = "join-room"
		message.RoomList = make([]RoomInfo, 0)
		message.RoomList = append(message.RoomList, RoomInfo{roomId, []string{message.Sender}})

		wss.rooms[roomId][message.Sender].send <- message

	case "list-room":

		Info.Println("list-room:", message)

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
		if _, ok := wss.rooms[WAITITNG_ROOM][message.Sender]; ok {
			wss.rooms[WAITITNG_ROOM][message.Sender].send <- message
		}

	default:
		Warning.Println("Unknown Action:", message)
	}
}
