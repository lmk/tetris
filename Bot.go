package main

import "time"

const (
	LEVEL_BIGINER       = "beginner"
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
}

type BotEngine interface {
	init()
	getNextAction() string
	getCycle() int
}

func NewBot(level string, ba *BotAdapter) *Bot {

	bot := &Bot{
		level:      level,
		botAdapter: ba,
		state:      "ready",
		x:          BOARD_CENTER,
		y:          0,
	}

	switch bot.level {
	case LEVEL_BIGINER:
		bot.botEngine = NewBotBigenner(bot)
	}

	go bot.runMsgHandler()

	return bot
}

func (b *Bot) startGame(msg *Message) {
	b.Cell = msg.Cells
	b.NextBlockIndexs = msg.BlockIndexs

	b.state = "start"
	b.CurrentBlock = msg.CurrentBlock

	b.botEngine.init()

	go b.runActionLoop()
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

func (b *Bot) leaveRoom(msg *Message) bool {

	b.end()

	if msg.RoomList[0].Owner == b.botAdapter.nick {
		b.botAdapter.fromBot <- &Message{
			Action: "leave-room",
			RoomId: msg.RoomId,
			Sender: b.botAdapter.nick,
		}

		return true
	}

	return false
}

func (b *Bot) runMsgHandler() {

	for {

		msg, ok := <-b.botAdapter.toBot
		if !ok {
			Trace.Println("runMsgHandler", b.botAdapter.nick, "channel closed")
			break
		}
		Trace.Println("runMsgHandler", b.botAdapter.nick, msg.Action)

		switch msg.Action {
		case "start-game":
			b.startGame(msg)

		case "sync-game":
			b.syncGame(msg)

		case "gift-full-blocks":
			// 처리가 필요할까?

		case "over-game", "end-game":
			b.end()

		case "leave-room":
			if b.leaveRoom(msg) {
				close(b.botAdapter.toBot)
			}
		}
	}

	Trace.Println("runMsgHandler", b.botAdapter.nick, "end")
}

func (b *Bot) runActionLoop() {
	for b.state != "end" {

		time.Sleep(time.Millisecond * time.Duration(b.botEngine.getCycle()))

		b.botAdapter.fromBot <- &Message{
			Action: b.botEngine.getNextAction(),
			Sender: b.botAdapter.nick,
		}
	}

	Trace.Println("runActionLoop", b.botAdapter.nick, "end")
}
