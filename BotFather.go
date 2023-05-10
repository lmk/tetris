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
		fromBot:     make(chan *Message),
		fromManager: make(chan *Message),
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

	go bot.run()

}

// bot 에게 받은 메시지를 Manager 에게 전달한다.
// Manager에게 받은 메시지를 bot 에게 전달한다.
func (bf *botFather) run() {
	for {
		select {
		case msg := <-bf.fromBot:
			switch msg.Action {
			case "block-drop", "block-rotate", "block-left", "block-right", "block-down":
				// TODO: send to manager
			}

		case msg := <-bf.fromManager:
			switch msg.Action {
			case "add-bot":
				bf.addBot(msg)
			}
		}
	}
}
