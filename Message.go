package main

import "encoding/json"

type Message struct {
	// Action:
	//    "new-room", "join-room", "leave-room", "list-room",
	//    "start-game", "over-game",
	//    "next-block",
	//    "block-rotate", "block-left", "block-right", "block-down", "block-drop"
	//    "gift-full-blocks", "sync-game"
	Action string `json:"action"`
	Sender string `json:"sender"` // user-nick or game-id
	roomId int    `json:"-"`
	Data   string `json:"data,omitempty"`

	// for list-room
	RoomList []RoomInfo `json:"roomList,omitempty"`

	// for gift-full-blocks
	Cells [][]int `json:"cells,omitempty"`

	// for over-game, sync-game
	Score int `json:"score"`

	// sync-game
	CurrentBlock Block `json:"block,omitempty"`
}

func (g *Game) SendGameOver() {

	g.managerCh <- &Message{
		Action: "over-game",
		Sender: g.owner,
		Cells:  g.cell,
		Score:  g.score,
	}
}

func (g *Game) SendNextBlocks() {
	jsonNexts, _ := json.Marshal(g.nextBlockIndexs)
	g.managerCh <- &Message{
		Action: "next-block",
		Sender: g.owner,
		Data:   string(jsonNexts),
	}
}

func (g *Game) SendSyncGame() {

	g.managerCh <- &Message{
		Action:       "sync-game",
		Sender:       g.owner,
		Data:         "",
		Cells:        g.cell,
		CurrentBlock: g.currentBlock,
		Score:        g.score,
	}
}

func (g *Game) SendGiftFullBlocks(gift [][]int) {

	g.managerCh <- &Message{
		Action: "gift-full-blocks",
		Sender: g.owner,
		Cells:  gift,
	}
}
