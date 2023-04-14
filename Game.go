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
	client          *Client       // player
	request         chan *Message // request from player
	status          string        // ready-game, playing-game, over-game
	cell            [][]int       // 0: empty, other: block index
	score           int
	currentBlock    Block
	nextBlockIndexs []int // next block indexs 10
	cycleMs         int   // cycle time in ms
}

// NewGame create a new game
// client: response channel to player
func NewGame(client *Client) *Game {
	game := Game{
		status:       "ready-game",
		score:        0,
		currentBlock: NewBlock((rand.Intn(len(SHAPES)) + 1)),
		cycleMs:      1000,
		request:      make(chan *Message),
		client:       client,
		cell:         make([][]int, BOARD_ROW),
	}

	for i := 0; i < BOARD_ROW; i++ {
		game.cell[i] = make([]int, BOARD_COLUMN)
	}

	game.CreateNextBlock()

	return &game
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
			if newBlock.Shape[r][c] != EMPTY {
				y := r + newBlock.Row
				x := c + newBlock.Col

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
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row
				x := c + g.currentBlock.Col + 1

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
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row - 1
				x := c + g.currentBlock.Col

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
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row + 1
				x := c + g.currentBlock.Col

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
			if g.currentBlock.Shape[r][c] != EMPTY {
				y := r + g.currentBlock.Row
				x := c + g.currentBlock.Col

				g.cell[y][x] = g.currentBlock.ShapeIndex
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

func (g *Game) IsGameOver() bool {
	return g.status == "over-game"
}

// NextTern : 다음 턴으로 넘어감
// 리턴값 : 한번에 지운 라인의 수
// -1 : 게임오버
func (g *Game) NextTern() int {
	count := g.removeFullLine()

	g.CreateNextBlock()

	if !g.isSafeNewBlock() {
		return -1
	}
	g.firstNextToCurrnetBlock()

	return count
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

func (g *Game) GameOver() {
	g.currentBlockToBoard()
	g.status = "over-game"
}

func (g *Game) Stop() {
	g.request <- &Message{Action: "over-game"}
}

func (g *Game) Run() {

	g.status = "playing-game"
	g.SendSyncGame()

	for {
		select {
		case msg := <-g.request:
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
					if g.NextTern() == -1 {
						g.GameOver()
						g.SendGameOver()
						return
					}
				}
				g.SendSyncGame()

			case "block-drop":
				if !g.toDrop() {
					g.GameOver()
					g.SendGameOver()
					return
				}
				g.SendSyncGame()

			case "over-game":
				g.GameOver()
				g.SendGameOver()
				return

			case "gift-full-blocks":
				if !g.receiveFullBlocks(msg.Cells) {
					g.GameOver()
					g.SendGameOver()
					return
				}
				g.SendSyncGame()
			}

		case <-time.After(time.Millisecond * time.Duration(g.cycleMs)):
			if !g.toDown() {
				if g.NextTern() == -1 {
					g.GameOver()
					g.SendGameOver()
					return
				}
			}
			g.SendSyncGame()
		}
	}
}
