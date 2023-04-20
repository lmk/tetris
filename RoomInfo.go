package main

type RoomInfo struct {
	RoomId  int                `json:"roomId"`
	Owner   string             `json:"owner"`
	Clients map[string]*Client `json:"nicks"`
	Title   string             `json:"title"`
}

const (
	WAITITNG_ROOM = 0 // 대기실
)

func NewRoomInfo(roomId int, client *Client, title string) *RoomInfo {

	return &RoomInfo{
		RoomId:  roomId,
		Owner:   client.Nick,
		Clients: map[string]*Client{client.Nick: client},
		Title:   title,
	}
}
