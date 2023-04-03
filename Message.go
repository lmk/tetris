package main

type Message struct {
	Action   string     `json:"action"`
	Sender   string     `json:"sender"`
	Msg      string     `json:"msg,omitempty"`
	RoomInfo []RoomInfo `json:"roomInfo,omitempty"`
}
