package main

type Message struct {
	Action   string     `json:"action"`
	Sender   string     `json:"sender"`
	Msg      string     `json:"msg,omitempty"`
	RoomList []RoomInfo `json:"roomList,omitempty"`
	Client   *Client    `json:"-"`
}
