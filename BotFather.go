package main

type botFather struct {
	botList     map[*Bot]int  // bot list (bot, roomId)
	fromBot     chan *Message // Bot -> BotFather
	fromManager chan *Message // Manager -> BotFather
}

var BotFather *botFather

func init() {
	BotFather = newBotFather()
	go BotFather.run()
}

func newBotFather() *botFather {
	return &botFather{
		botList:     make(map[*Bot]int),
		fromBot:     make(chan *Message, MAX_CHAN),
		fromManager: make(chan *Message, MAX_CHAN),
	}
}

func (bf *botFather) addBot(msg *Message) {

	// new bot adapter
	botAdapter := NewBotAdapter(msg.RoomId)
	if botAdapter == nil {
		return
	}

	// new bot
	bot := NewBot(msg.Data, botAdapter)
	bf.botList[bot] = msg.RoomId
}

func (bf *botFather) deleteBot(msg *Message) {

	// find bot
	var bot *Bot
	for b := range bf.botList {
		if b.botAdapter.nick == msg.Sender {
			bot = b
			break
		}
	}

	if bot == nil {
		return
	}

	// delete bot
	bot.botAdapter.socket.Close()
	delete(bf.botList, bot)
}

// bot 에게 받은 메시지를 Manager 에게 전달한다.
// Manager에게 받은 메시지를 bot 에게 전달한다.
func (bf *botFather) run() {
	for {
		select {
		case msg, ok := <-bf.fromBot:
			if !ok {
				Error.Println("run", "fromBot channel closed")
				return
			}

			switch msg.Action {
			case "leave-room":
				bf.deleteBot(msg)
			}

		case msg, ok := <-bf.fromManager:
			if !ok {
				Error.Println("run", "fromManager channel closed")
				return
			}

			switch msg.Action {
			case "add-bot":
				bf.addBot(msg)
			}
		}
	}
}
