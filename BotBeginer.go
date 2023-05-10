package main

const (
	NEED_FIND int = -999
)

// BotBigenner
// 한 줄씩 제거
type BotBigenner struct {
	cycle     int
	bot       *Bot
	incX      int // -999 면 탐색, 0이면 down, 0보다 크면 right, 0보다 작으면 left
	incRotate int // 0이면 rotate 안함, 기타는 rotate
}

func NewBotBigenner(bot *Bot) *BotBigenner {
	return &BotBigenner{
		cycle:     900,
		bot:       bot,
		incX:      NEED_FIND,
		incRotate: 0,
	}
}

func (b *BotBigenner) findCombine(cells [][]int, m *Margin) int {

	for r := BOARD_ROW - 1 - len(cells); r > b.bot.CurrentBlock.Row; r-- {
		for c := 0; c < BOARD_COLUMN-len(cells[0]); c++ {
			if CanCombine(cells, b.bot.Cell, BOARD_COLUMN-c, BOARD_ROW-r) {
				if b.bot.CurrentBlock.Col > c {
					return b.bot.CurrentBlock.Col - c
				}

				return b.bot.CurrentBlock.Col + m.Left
			}
		}
	}

	return -1
}

func (b *BotBigenner) findXforDown() {

	// block 모양에 맞는 가장 낮은 위치를 찾는다.
	x := -1
	rotate := 0
	for rotate = 0; rotate < 4; rotate++ {

		trimed, margin := TrimShape(b.bot.CurrentBlock.Shape)

		x = b.findCombine(trimed, &margin)
		if x != -1 {
			break
		}

		rotate++
		b.bot.CurrentBlock.Rotate()
	}

	if x == -1 {
		// 찾지 못하면 board의 가장 낮은 곳을 찾는다.
		x = findLowest(b.bot.Cell)
	}

	b.incX = x
	b.incRotate = rotate
}

func (b *BotBigenner) getNextAction() string {

	action := ""

	if b.incX == NEED_FIND {
		// down 할 위치를 찾는다.
		b.findXforDown()

		// cycle을 2/1로 줄인다.
		b.cycle /= 2
	} else if b.incRotate > 0 {
		action = "block-rotate"
		b.incRotate--
	} else if b.incX > 0 {
		action = "block-right"
		b.incX--
	} else if b.incX < 0 {
		action = "block-left"
		b.incX++
	} else if b.incX == 0 {
		// col에 도착하면 down 한다.
		// cycle을 *2로 늘린다.
		b.incX = NEED_FIND
		b.cycle *= 2
	}

	return action
}

func (b *BotBigenner) getCycle() int {
	return b.cycle
}
