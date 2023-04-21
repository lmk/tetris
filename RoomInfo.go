package main

const (
	WAITITNG_ROOM = 0 // 대기실
)

type RoomInfo struct {
	RoomId  int                `json:"roomId"`
	Owner   string             `json:"owner"`
	Clients map[string]*Client `json:"nicks"`
	Title   string             `json:"title"`
	State   string             `json:"state"`
}

type GetState func() string

func NewRoomInfo(roomId int, client *Client, title string) *RoomInfo {

	return &RoomInfo{
		RoomId:  roomId,
		Owner:   client.Nick,
		Clients: map[string]*Client{client.Nick: client},
		Title:   title,
		State:   "ready",
	}
}

func (room *RoomInfo) GetState() string {

	ready, playing, over := 0, 0, 0
	for _, client := range room.Clients {
		if client.Game == nil {
			continue
		}

		switch client.Game.State {
		case "ready":
			ready++
		case "playing":
			playing++
		case "over":
			over++
		}
	}

	if playing > 0 {
		return "playing"
	}

	return "ready"
}

func appendRoomInfo(RoomList []RoomInfo, room *RoomInfo) []RoomInfo {

	room.State = room.GetState()

	return append(RoomList, *room)
}
