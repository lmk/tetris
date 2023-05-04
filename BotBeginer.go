package main

type BotBigenner struct {
	cycle int
}

func NewBotBigenner() *BotBigenner {
	return &BotBigenner{
		cycle: 1000,
	}
}

func (b BotBigenner) getNextAction() string {
	return "move-left"
}

func (b BotBigenner) getCycle() int {
	return b.cycle
}
