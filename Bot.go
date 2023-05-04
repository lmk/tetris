package main

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

	bot BotEngine

	state string
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
	}

	switch bi.level {
	case LEVEL_BIGINER:
		bi.bot = NewBotBigenner()
	}

	return bi
}

func (b *Bot) startGame(msg *Message) {
	b.Cell = msg.Cells
	b.NextBlockIndexs = msg.BlockIndexs

	b.state = "start"
}

func (b *Bot) run() {
	for b.state != "end" {
		msg := <-b.botAdapter.toBot
		switch msg.Action {
		case "start-game":
			b.startGame(msg)

		case "gift-full-blocks":

		}
	}
}

func (b *Bot) end() {
	b.state = "end"
}
