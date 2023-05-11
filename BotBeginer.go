package main

const (
	NEED_FIND int = -999
)

// BotBigenner
// 한 줄씩 제거
type BotBigenner struct {
	cycle         int
	bot           *Bot
	incX          int // -999 면 탐색, 0이면 down, 0보다 크면 right, 0보다 작으면 left
	incRotate     int // 0이면 rotate 안함, 기타는 rotate
	cycleMoving   int // 움직일때 시간
	cycleThinking int // 생각할때 시간
}

func NewBotBigenner(bot *Bot) *BotBigenner {
	return &BotBigenner{
		cycle:         1600,
		bot:           bot,
		incX:          NEED_FIND,
		incRotate:     0,
		cycleMoving:   100,
		cycleThinking: 900,
	}
}

func (b *BotBigenner) init() {
	b.cycle = 1600
	b.incX = NEED_FIND
	b.incRotate = 0
}

// cells의 모양에 맞는 x 좌표를 찾는다.
func (b *BotBigenner) findCombine(cells [][]int, m *Margin) (int, int) {

	// drop 했을때 막히는 cell을 채워준다.
	board := fillTailToUp(b.bot.Cell)
	cells = fillTailToDown(cells)

	for r := BOARD_ROW - len(cells); r > b.bot.CurrentBlock.Row+m.Bottom; r-- {
		for c := 0; c < BOARD_COLUMN-len(cells[0])+m.Right; c++ {
			if CanCombine(cells, board, c, r) {
				if b.bot.CurrentBlock.Col > c {
					return c, r
				}

				return c + m.Left, r
			}
		}
	}

	return -1, -1
}

func (b *BotBigenner) findXforDown() {

	var block Block
	block.Clone(b.bot.CurrentBlock)

	// block 모양에 맞는 가장 낮은 위치를 찾는다.
	xBest := -1
	yBest := -1
	rotateBest := 0
	for rotate := 0; rotate < 4; rotate++ {

		trimed, margin := TrimShape(block.Shape)

		x, y := b.findCombine(trimed, &margin)
		if x != -1 && yBest < y {
			xBest = x
			yBest = y
			rotateBest = rotate
		}

		block.Rotate()
	}

	// 찾지 못하면
	if xBest == -1 {
		// board의 가장 낮은 곳을 찾는다.
		bx := findLowest(b.bot.Cell, 0)

		// block의 가장 낮은 곳을 찾는다.
		cx := findLowest(b.bot.CurrentBlock.Shape, 1)

		xBest = bx - cx

		rotateBest = 0
	}

	b.incX = xBest - block.Col
	b.incRotate = rotateBest
}

func (b *BotBigenner) getNextAction() string {

	action := ""

	if b.incX == NEED_FIND {
		// down 할 위치를 찾는다.
		b.findXforDown()

		b.cycle = b.cycleMoving
	}

	if b.incRotate > 0 {
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
		b.cycle = b.cycleThinking

		action = "block-drop"
	}

	return action
}

func (b *BotBigenner) getCycle() int {
	return b.cycle
}
