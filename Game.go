package main

import (
	"math/rand"
	"time"
)

const (
	BOARD_ROW    = 15
	BOARD_COLUMN = 10
	BOARD_CENTER = (BOARD_COLUMN - 3) / 2

	BLOCK_ROW    = 4
	BLOCK_COLUMN = 4

	EMPTY = 0
)

type Game struct {
	owner           string        // owner of game
	ch              chan *Message // request from manager
	managerCh       chan *Message // response to manager
	Id              int           // game id
	status          string        // ready-game, playing-game, over-game
	cell            [][]int       // 0: empty, other: block index
	score           int
	currentBlock    Block
	nextBlockIndexs []int // next block indexs 10
	cycleMs         int   // cycle time in ms
}

// NewGame create a new game
// ch: request channel
// id: game id
func NewGame(ch chan *Message, id int, owner string) *Game {
	game := Game{
		owner:        owner,
		ch:           make(chan *Message),
		managerCh:    ch,
		Id:           id,
		status:       "ready-game",
		score:        0,
		currentBlock: NewBlock((rand.Intn(len(SHAPES)) + 1)),
		cycleMs:      1000,
		cell:         make([][]int, BOARD_ROW),
	}

	for i := 0; i < BOARD_ROW; i++ {
		game.cell[i] = make([]int, BOARD_COLUMN)
	}

	game.CreateNextBlock()

	return &game
}

func (g *Game) reset() {
	g.score = 0
	g.currentBlock = NewBlock((rand.Intn(len(SHAPES)) + 1))
	g.cycleMs = 1000
	g.cell = make([][]int, BOARD_ROW)
	g.nextBlockIndexs = []int{}

	for i := 0; i < BOARD_ROW; i++ {
		g.cell[i] = make([]int, BOARD_COLUMN)
	}

	g.CreateNextBlock()
}

func (g *Game) CreateNextBlock() {

	count := 10 - len(g.nextBlockIndexs)

	for i := 0; i < count; i++ {
		g.nextBlockIndexs = append(g.nextBlockIndexs, (rand.Intn(len(SHAPES)) + 1))
	}
}

func (g *Game) isSafeToRoteate() bool {
	block := Block{}
	block.Clone(g.currentBlock)
	block.Rotate()

	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if block.Shape[r][c] != EMPTY {
				y := r + block.Row
				x := c + block.Col

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.cell[y][x] != EMPTY {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) isSafeNewBlock() bool {

	newBlock := NewBlock(g.nextBlockIndexs[0])

	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if newBlock.Shape[r][c] != EMPTY {
				y := r + newBlock.Row
				x := c + newBlock.Col

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.cell[y][x] != EMPTY {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) isSafeToDown() bool {
	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row + 1
				x := c + g.currentBlock.Col

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.cell[y][x] != EMPTY {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) isSafeToLeft() bool {
	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row
				x := c + g.currentBlock.Col - 1

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.cell[y][x] != EMPTY {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) isSafeToRight() bool {
	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row
				x := c + g.currentBlock.Col + 1

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.cell[y][x] != EMPTY {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) currentBlockToBoard() {
	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row
				x := c + g.currentBlock.Col

				g.cell[y][x] = g.currentBlock.ShapeIndex
			}
		}
	}
}

func (g *Game) shiftDownCell(row int) {
	for r := row; r > 0; r-- {
		for c := 0; c < BOARD_COLUMN; c++ {
			g.cell[r][c] = g.cell[r-1][c]
		}
	}

	for c := 0; c < BOARD_COLUMN; c++ {
		g.cell[0][c] = EMPTY
	}
}

func appendRow(cell [][]int, row []int) [][]int {

	newLine := make([]int, len(row))
	copy(newLine, row)

	cell = append(cell, newLine)

	return cell
}

// procFullLine
// currentBlock 이 있는 row만 검사
// full line이면 row를 삭제하고, 위 row를 한칸씩 내림
// full line이 아니면, currentBlock를 board에 복사
func (g *Game) procFullLine() [][]int {

	removedLines := [][]int{}

	maxRow := BOARD_ROW
	if g.currentBlock.Row+BLOCK_ROW < BOARD_ROW {
		maxRow = g.currentBlock.Row + BLOCK_ROW
	}

	for r := g.currentBlock.Row; r < maxRow; r++ {

		y := r - g.currentBlock.Row
		x := 0

		isFull := true

		for c := 0; c < BOARD_COLUMN; c++ {

			x = c - g.currentBlock.Col

			// if cell is empty, check current block
			if g.cell[r][c] == EMPTY {

				if g.currentBlock.inBlock(r, c) {

					if g.currentBlock.Shape[y][x] == EMPTY {
						isFull = false
						break
					}
				} else {
					isFull = false
					break
				}
			}
		}

		if isFull {

			removedLines = appendRow(removedLines, g.cell[r])

			g.shiftDownCell(r)

			// score는 삭제된 line 수 가중치를 줘서 계산
			g.score += (10 * len(removedLines))

		} else {

			// currnet block to board
			for i := 0; i < BLOCK_COLUMN; i++ {
				if g.currentBlock.Shape[y][i] != EMPTY {
					g.cell[r][g.currentBlock.Col+i] = g.currentBlock.ShapeIndex
				}
			}
		}
	}

	return removedLines
}

func (g *Game) firstNextToCurrnetBlock() {
	g.currentBlock = NewBlock(g.nextBlockIndexs[0])
	g.nextBlockIndexs = g.nextBlockIndexs[1:]
}

func (g *Game) IsGameOver() bool {
	return g.status == "over-game"
}

// nextTern : 다음 턴으로 넘어감
// 삭제된 라인 처리
// 삭제된 라인이 2줄 이상이면, 경쟁자들에게 선물을 보냄
// 새 블럭을 생성
// 새 블럭을 생성하지 못하면, 게임오버
func (g *Game) nextTern() bool {
	removedLines := g.procFullLine()
	if len(removedLines) > 1 {
		g.SendGiftFullBlocks(removedLines)
	}

	if !g.isSafeNewBlock() {
		g.gameOver()
		g.SendGameOver()
		return false
	}
	g.firstNextToCurrnetBlock()

	g.CreateNextBlock()
	return true
}

func (g *Game) toRotate() bool {

	if !g.isSafeToRoteate() {
		return false
	}

	g.currentBlock.Rotate()

	return true
}

func (g *Game) toLeft() bool {
	if !g.isSafeToLeft() {
		return false
	}

	g.currentBlock.Col--

	return true
}

func (g *Game) toRight() bool {
	if !g.isSafeToRight() {
		return false
	}

	g.currentBlock.Col++

	return true
}

func (g *Game) toDown() bool {
	if !g.isSafeToDown() {
		return false
	}

	g.currentBlock.Row++

	return true
}

func (g *Game) toDrop() bool {

	if !g.isSafeToDown() {
		return false
	}

	for g.isSafeToDown() {
		g.currentBlock.Row++
	}

	return true
}

func (g *Game) receiveFullBlocks(blocks [][]int) bool {

	g.cell = append(g.cell, blocks...)

	// check gamevoer 밀린 위쪽 cell에 블럭이 있으면 게임오버
	for r := 0; r < len(blocks); r++ {
		for c := 0; c < len(g.cell[r]); c++ {
			if g.cell[r][c] != EMPTY {
				return false
			}
		}
	}

	g.cell = g.cell[len(blocks):]

	return true
}

func (g *Game) IsPlaying() bool {
	return g.status == "playing-game"
}

func (g *Game) gameOver() {
	g.status = "over-game"
}

// Action : 게임에 메시지를 전달한다.
func (g *Game) Action(msg *Message) {
	g.ch <- msg
}

// Stop : 게임을 종료시킨다.
func (g *Game) Stop() {
	g.ch <- &Message{Action: "over-game"}
}

// Start : 게임을 시작한다.
func (g *Game) Start() {

	if g.IsPlaying() {
		return
	}

	g.reset()
	g.status = "playing-game"
	g.SendStartGame()

	go g.run()
}

func (g *Game) run() {

	for !g.IsGameOver() {
		select {
		case msg := <-g.ch:
			switch msg.Action {
			case "block-rotate":
				g.toRotate()
				g.SendSyncGame()

			case "block-left":
				g.toLeft()
				g.SendSyncGame()

			case "block-right":
				g.toRight()
				g.SendSyncGame()

			case "block-down":
				if !g.toDown() {
					if !g.nextTern() {
						return
					}
				}
				g.SendSyncGame()

			case "block-drop":
				if !g.toDrop() {
					g.gameOver()
					g.currentBlockToBoard()
					g.SendGameOver()
				}
				if !g.nextTern() {
					return
				}
				g.SendSyncGame()

			case "over-game":
				g.gameOver()
				g.currentBlockToBoard()
				g.SendGameOver()
				return

			case "gift-full-blocks":
				if !g.receiveFullBlocks(msg.Cells) {
					g.gameOver()
					g.currentBlockToBoard()
					g.SendGameOver()
					return
				}
				g.SendSyncGame()
			}

		case <-time.After(time.Millisecond * time.Duration(g.cycleMs)):
			if !g.toDown() && !g.nextTern() {
				return
			}
			g.SendSyncGame()
		}
	}
}
