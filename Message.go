package main

// Message is a struct for message
// websocket -> server : Action, Sender, Data
// server -> websocket : All fields
type Message struct {
	// Action:
	//    "new-room", "join-room", "leave-room", "list-room",
	//    "start-game", "over-game",
	//    "block-rotate", "block-left", "block-right", "block-down", "block-drop"
	//    "gift-full-blocks", "sync-game"
	Action string `json:"action"`
	Sender string `json:"sender"` // user-nick or game-id
	Data   string `json:"data,omitempty"`

	RoomId          int   `json:"roomId,omitempty"`
	NextBlockIndexs []int `json:"nextBlocks,omitempty"`

	// for list-room
	RoomList []RoomInfo `json:"roomList,omitempty"`

	// for gift-full-blocks
	Cells [][]int `json:"cells,omitempty"`

	// for over-game, sync-game
	Score int `json:"score,omitempty"`

	// sync-game
	CurrentBlock *Block `json:"block,omitempty"`
}
