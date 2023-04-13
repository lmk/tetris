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
	owner           string  // player nick name
	status          string  // ready, playing, gameover
	cell            [][]int // 0: empty, other: block index
	score           int
	currentBlock    Block
	nextBlockIndexs []int        // next block indexs 10
	cycleMs         int          // cycle time in ms
	recvCh          chan *Action // revc channel
	sendCh          chan *Action // send channel
}

func NewGame(ch chan *Action, user string) Game {
	game := Game{
		owner:        user,
		status:       "ready",
		score:        0,
		currentBlock: NewBlock(rand.Intn(len(SHAPES))),
		cycleMs:      1000,
		sendCh:       ch,
		recvCh:       make(chan *Action),
		cell:         make([][]int, BOARD_ROW),
	}

	for i := 0; i < BOARD_ROW; i++ {
		game.cell[i] = make([]int, BOARD_COLUMN)
	}

	game.CreateNextBlock(10)

	return game
}

func (g *Game) CreateNextBlock(count int) {
	for i := 0; i < count; i++ {
		g.nextBlockIndexs = append(g.nextBlockIndexs, rand.Intn(len(SHAPES)))
	}
}

func (g *Game) isSafeToRoteate() bool {
	block := Block{}
	block.Clone(g.currentBlock)
	block.Rotate()

	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if block.shape[r][c] != EMPTY {
				x := r + block.row
				y := c + block.col

				if x < 0 || x >= BOARD_ROW || y < 0 || y >= BOARD_COLUMN {
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
			if newBlock.shape[r][c] != EMPTY {
				x := r + newBlock.row
				y := c + newBlock.col

				if x < 0 || x >= BOARD_ROW || y < 0 || y >= BOARD_COLUMN {
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
			if g.currentBlock.shape[r][c] != EMPTY {
				x := r + g.currentBlock.row
				y := c + g.currentBlock.col + 1

				if x < 0 || x >= BOARD_ROW || y < 0 || y >= BOARD_COLUMN {
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
			if g.currentBlock.shape[r][c] != EMPTY {
				x := r + g.currentBlock.row - 1
				y := c + g.currentBlock.col

				if x < 0 || x >= BOARD_ROW || y < 0 || y >= BOARD_COLUMN {
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
			if g.currentBlock.shape[r][c] != EMPTY {
				x := r + g.currentBlock.row + 1
				y := c + g.currentBlock.col

				if x < 0 || x >= BOARD_ROW || y < 0 || y >= BOARD_COLUMN {
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
			if g.currentBlock.shape[r][c] != EMPTY {
				x := r + g.currentBlock.row
				y := c + g.currentBlock.col

				g.cell[y][x] = g.currentBlock.shapeIndex
			}
		}
	}
}

func (g *Game) isFullLine(row int) bool {
	for c := 0; c < BOARD_COLUMN; c++ {
		if g.cell[row][c] == EMPTY {
			return false
		}
	}
	return true
}

func (g *Game) removeLine(row int) {
	for r := row; r > 0; r-- {
		for c := 0; c < BOARD_COLUMN; c++ {
			g.cell[r][c] = g.cell[r-1][c]
		}
	}

	for c := 0; c < BOARD_COLUMN; c++ {
		g.cell[0][c] = EMPTY
	}
}

func (g *Game) removeFullLine() int {
	count := 0
	for r := 0; r < BOARD_ROW; r++ {
		if g.isFullLine(r) {
			g.removeLine(r)
			count++
			g.score += (10 * count)
		}
	}

	return count
}

func (g *Game) firstNextToCurrnetBlock() {
	g.currentBlock = NewBlock(g.nextBlockIndexs[0])
	g.nextBlockIndexs = g.nextBlockIndexs[1:]
}

func (g *Game) NextTern() bool {
	count := g.removeFullLine()

	g.CreateNextBlock(count + 1)

	if !g.isSafeNewBlock() {
		return false
	}
	g.firstNextToCurrnetBlock()

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

	g.currentBlock.col--

	return true
}

func (g *Game) toRight() bool {
	if !g.isSafeToRight() {
		return false
	}

	g.currentBlock.col++

	return true
}

func (g *Game) toDown() bool {
	if !g.isSafeToDown() {
		return false
	}

	g.currentBlock.row++

	return true
}

func (g *Game) toDrop() bool {
	for g.isSafeToDown() {
		g.currentBlock.row++
	}

	return true
}

func (g *Game) receiveGifts(gifts [][]int) bool {

	g.cell = append(g.cell, gifts...)

	// check gamevoer 밀린 위쪽 cell에 블럭이 있으면 게임오버
	for r := 0; r < len(gifts); r++ {
		for c := 0; c < len(g.cell[r]); c++ {
			if g.cell[r][c] != EMPTY {
				return false
			}
		}
	}

	g.cell = g.cell[len(gifts):]

	return true
}

func (g *Game) GameOver() {
	g.currentBlockToBoard()
	g.status = "gameover"
}

func (g *Game) Run() {

	g.status = "playing"

	for {
		select {
		case action := <-g.recvCh:
			switch action.action {
			case "rotate":
				g.toRotate()

			case "left":
				g.toLeft()

			case "right":
				g.toRight()

			case "down":
				if !g.toDown() {
					g.GameOver()
					return
				}

			case "drop":
				if !g.toDrop() {
					g.GameOver()
					return
				}

			case "stop":
				g.GameOver()
				return

			case "gift-blocks":
				if !g.receiveGifts(action.blocks) {
					g.GameOver()
					return
				}
			}

		case <-time.After(time.Millisecond * time.Duration(g.cycleMs)):
			if !g.toDown() {
				g.GameOver()
				return
			}
		}
	}
}

type Action struct {
	action string // rotate, left, right, down, drop, stop, gift-blocks
	nick   string
	blocks [][]int
}
