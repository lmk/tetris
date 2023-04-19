package main

type Message struct {
	// Action:
	//    "new-room", "join-room", "leave-room", "list-room",
	//    "start-game", "over-game",
	//    "block-rotate", "block-left", "block-right", "block-down", "block-drop"
	//    "gift-full-blocks", "sync-game"
	Action string `json:"action"`
	Sender string `json:"sender"` // user-nick or game-id
	roomId int    `json:"-"`
	Data   string `json:"data,omitempty"`

	NextBlockIndexs []int `json:"nextBlocks,omitempty"`

	// for list-room
	RoomList []RoomInfo `json:"roomList,omitempty"`

	// for gift-full-blocks
	Cells [][]int `json:"cells,omitempty"`

	// for over-game, sync-game
	Score int `json:"score,omitempty"`

	// sync-game
	CurrentBlock Block `json:"block,omitempty"`
}

func (g *Game) SendGameOver() {

	g.managerCh <- &Message{
		Action:       "over-game",
		Sender:       g.owner,
		Cells:        g.cell,
		CurrentBlock: g.currentBlock,
		Score:        g.score,
	}
}

func (g *Game) SendStartGame() {

	g.managerCh <- &Message{
		Action:          "start-game",
		Sender:          g.owner,
		NextBlockIndexs: g.nextBlockIndexs,
		Cells:           g.cell,
		CurrentBlock:    g.currentBlock,
	}
}

func (g *Game) SendSyncGame() {
	g.managerCh <- &Message{
		Action:          "sync-game",
		Sender:          g.owner,
		NextBlockIndexs: g.nextBlockIndexs,
		Cells:           g.cell,
		CurrentBlock:    g.currentBlock,
		Score:           g.score,
	}
}

func (g *Game) SendGiftFullBlocks(gift [][]int) {

	g.managerCh <- &Message{
		Action: "gift-full-blocks",
		Sender: g.owner,
		Cells:  gift,
	}
}
