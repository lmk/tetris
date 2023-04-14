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
	Sender string `json:"sender"`
	roomId int    `json:"-"`
	Data   string `json:"data,omitempty"`

	// for list-room
	RoomList []RoomInfo `json:"roomList,omitempty"`

	// for gift-full-blocks
	Cells          [][]int `json:"cells,omitempty"`
	CellsRowIndexs []int   `json:"-"`

	// for over-game, sync-game
	Score int `json:"score"`

	// sync-game
	CurrentBlock Block `json:"block,omitempty"`
}

func (g *Game) SendGameOver() {

	for _, client := range g.client.ws.rooms[g.client.roomId] {
		client.send <- &Message{
			Action: "over-game",
			Sender: g.client.nick,
			roomId: g.client.roomId,
			Cells:  g.cell,
			Score:  g.score,
		}
	}
}

func (g *Game) SendNextBlocks() {
	jsonNexts, _ := json.Marshal(g.nextBlockIndexs)
	g.client.send <- &Message{
		Action: "next-block",
		Sender: g.client.nick,
		roomId: g.client.roomId,
		Data:   string(jsonNexts),
	}
}

func (g *Game) SendSyncGame() {

	for _, client := range g.client.ws.rooms[g.client.roomId] {
		client.send <- &Message{
			Action:       "sync-game",
			Sender:       g.client.nick,
			roomId:       g.client.roomId,
			Data:         "",
			Cells:        g.cell,
			CurrentBlock: g.currentBlock,
			Score:        g.score,
		}
	}
}

func (g *Game) SendGiftFullBlocks(gift [][]int, fullBlocksRowIndexs []int) {

	for _, client := range g.client.ws.rooms[g.client.roomId] {
		if g.client.nick == client.nick {
			continue
		}

		client.send <- &Message{
			Action:         "gift-full-blocks",
			Sender:         g.client.nick,
			roomId:         g.client.roomId,
			Cells:          gift,
			CellsRowIndexs: fullBlocksRowIndexs,
		}
	}
}
