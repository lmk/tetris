package main

import "time"

const (
	LEVEL_BIGINER       = "beginer"
	LEVEL_FAST_FINGER   = "fast-finger"
	LEVEL_ATTACKER2     = "attacker2"
	LEVEL_ATTACKER3     = "attacker3"
	LEVEL_ATTACKER_HALF = "attacker-half"
)

type Bot struct {
	level      string
	botAdapter *BotAdapter // to bot

	Cell            [][]int
	NextBlockIndexs []int
	x, y            int // for current block
	CurrentBlock    *Block

	botEngine BotEngine

	state string

	timer *time.Timer
}

type BotEngine interface {
	getNextAction() string
	getCycle() int
}

func NewBot(level string, ba *BotAdapter) *Bot {

	bi := &Bot{
		level:      level,
		botAdapter: ba,
		state:      "ready",
		x:          BOARD_CENTER,
		y:          0,
	}

	switch bi.level {
	case LEVEL_BIGINER:
		bi.botEngine = NewBotBigenner(bi)
	}

	return bi
}

func (b *Bot) startGame(msg *Message) {
	b.Cell = msg.Cells
	b.NextBlockIndexs = msg.BlockIndexs

	b.state = "start"

	b.timer = time.NewTimer(time.Duration(b.botEngine.getCycle()) * time.Millisecond)

	b.CurrentBlock = msg.CurrentBlock
}

func (b *Bot) end() {
	b.state = "end"
}

func (b *Bot) syncGame(msg *Message) {
	if b.state != "start" || msg.Sender != b.botAdapter.nick {
		return
	}

	b.Cell = msg.Cells
	b.NextBlockIndexs = msg.BlockIndexs
	b.CurrentBlock = msg.CurrentBlock
}

func (b *Bot) run() {
	for b.state != "end" {
		select {
		case msg := <-b.botAdapter.toBot:
			switch msg.Action {
			case "start-game":
				b.startGame(msg)

			case "sync-game":
				b.syncGame(msg)

			case "gift-full-blocks":
				// 처리가 필요할까?

			case "over-game":
				b.end()
			}

		case <-b.timer.C:
			b.botAdapter.fromBot <- &Message{
				Action: b.botEngine.getNextAction(),
				Sender: b.botAdapter.nick,
			}
		}

	}
}
