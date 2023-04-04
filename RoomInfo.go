package main

type RoomInfo struct {
	RoomId int      `json:"roomId"`
	Nicks  []string `json:"nicks"`
}

const (
	WAITITNG_ROOM = 0 // 대기실
)
