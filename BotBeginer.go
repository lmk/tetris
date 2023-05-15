package main

const (
	NOT_FOUND int = -999
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
		incX:          NOT_FOUND,
		incRotate:     0,
		cycleMoving:   200,
		cycleThinking: 1500,
	}
}

func (b *BotBigenner) init() {
	b.cycle = 1600
	b.incX = NOT_FOUND
	b.incRotate = 0
}

// cells의 모양에 맞는 x 좌표를 찾는다.
func (b *BotBigenner) findCombine(cells [][]int, board [][]int) (int, int) {

	for r := BOARD_ROW - len(cells); r >= 0; r-- {
		for c := 0; c <= BOARD_COLUMN-len(cells[0]); c++ {
			if CanCombine(cells, board, c, r) {
				return c, r
			}
		}
	}

	return NOT_FOUND, -1
}

func (b *BotBigenner) findXforDown() {

	// drop 했을때 막히는 cell을 채워준다.
	board := fillTailToUp(b.bot.Cell)
	//Trace.Print("board\n", cellsToString(board))

	var block Block
	block.Clone(b.bot.CurrentBlock)

	// block 모양에 맞는 가장 낮은 위치를 찾는다.
	incX := NOT_FOUND
	yBest := -1
	rotateBest := 0
	for rotate := 0; rotate < 4; rotate++ {

		trimed, margin := TrimShape(block.Shape)
		trimed = fillTailToDown(trimed)

		//Trace.Printf("block %d, %d\n%s", block.Col, block.Row, cellsToString(trimed))
		//Trace.Println("margin:", margin)

		x, y := b.findCombine(trimed, board)
		if x != NOT_FOUND && yBest < y+len(trimed) {
			incX = x - b.bot.CurrentBlock.Col - margin.Left
			yBest = y + len(trimed)
			rotateBest = rotate
		}

		//Trace.Println("x, y, rotate, currentX:", x, y, rotate, block.Col)

		block.Rotate()
	}

	// 찾지 못하면
	if incX == NOT_FOUND {
		// board의 가장 낮은 곳을 찾는다.
		bx := findLowest(board, 0)

		// block의 가장 낮은 곳을 찾는다.
		cx := findLowest(b.bot.CurrentBlock.Shape, 1)

		incX = bx - cx

		//Trace.Println("not found bx, cx, xBest :", bx, cx, incX)

		rotateBest = 0
	}

	b.incX = incX
	b.incRotate = rotateBest

	//Trace.Println("incX, incRotate :", b.incX, b.incRotate)
}

func (b *BotBigenner) getNextAction() string {

	action := ""

	if b.incX == NOT_FOUND {
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
		b.incX = NOT_FOUND
		b.cycle = b.cycleThinking

		action = "block-drop"
	}

	return action
}

func (b *BotBigenner) getCycle() int {
	return b.cycle
}
