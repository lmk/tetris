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
	Owner           string        `json:"-"`     // owner of game
	Ch              chan *Message `json:"-"`     // request from manager
	ManagerCh       chan *Message `json:"-"`     // response to manager
	State           string        `json:"state"` // ready, playing, over
	Cell            [][]int       `json:"-"`     // 0: empty, other: block index
	Score           int           `json:"score,omitempty"`
	CurrentBlock    *Block        `json:"-"`
	NextBlockIndexs []int         `json:"-"` // next block indexs 10
	CycleMs         int           `json:"-"` // cycle time in ms
	DurationTime    time.Time     `json:"-"` // duration time
}

// NewGame create a new game
// ch: request channel
// id: game id
func NewGame(ch chan *Message, owner string) *Game {
	game := Game{
		Owner:        owner,
		Ch:           make(chan *Message),
		ManagerCh:    ch,
		State:        "ready",
		Score:        0,
		CurrentBlock: NewBlock((rand.Intn(len(SHAPES)) + 1)),
		CycleMs:      1000,
		Cell:         make([][]int, BOARD_ROW),
	}

	for i := 0; i < BOARD_ROW; i++ {
		game.Cell[i] = make([]int, BOARD_COLUMN)
	}

	game.CreateNextBlock()

	return &game
}

func (g *Game) reset() {
	g.Score = 0
	g.CurrentBlock = NewBlock((rand.Intn(len(SHAPES)) + 1))
	g.CycleMs = 1000
	g.Cell = make([][]int, BOARD_ROW)
	g.NextBlockIndexs = []int{}

	for i := 0; i < BOARD_ROW; i++ {
		g.Cell[i] = make([]int, BOARD_COLUMN)
	}

	g.CreateNextBlock()
}

func (g *Game) CreateNextBlock() {

	count := 10 - len(g.NextBlockIndexs)

	for i := 0; i < count; i++ {
		g.NextBlockIndexs = append(g.NextBlockIndexs, (rand.Intn(len(SHAPES)) + 1))
	}
}

func (g *Game) isSafeToRoteate() bool {
	block := Block{}
	block.Clone(g.CurrentBlock)
	block.Rotate()

	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if block.Shape[r][c] != EMPTY {
				y := r + block.Row
				x := c + block.Col

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.Cell[y][x] != EMPTY {
					return false
				}
			}
		}
	}

	return true
}

func (g *Game) isSafeNewBlock() bool {

	newBlock := NewBlock(g.NextBlockIndexs[0])

	for r := 0; r < BLOCK_ROW; r++ {
		for c := 0; c < BLOCK_COLUMN; c++ {
			if newBlock.Shape[r][c] != EMPTY {
				y := r + newBlock.Row
				x := c + newBlock.Col

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.Cell[y][x] != EMPTY {
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
			if g.CurrentBlock.Shape[r][c] != EMPTY {
				y := r + g.CurrentBlock.Row + 1
				x := c + g.CurrentBlock.Col

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.Cell[y][x] != EMPTY {
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
			if g.CurrentBlock.Shape[r][c] != EMPTY {
				y := r + g.CurrentBlock.Row
				x := c + g.CurrentBlock.Col - 1

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.Cell[y][x] != EMPTY {
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
			if g.CurrentBlock.Shape[r][c] != EMPTY {
				y := r + g.CurrentBlock.Row
				x := c + g.CurrentBlock.Col + 1

				if x < 0 || x >= BOARD_COLUMN || y < 0 || y >= BOARD_ROW {
					return false
				}

				if g.Cell[y][x] != EMPTY {
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
			if g.CurrentBlock.Shape[r][c] != EMPTY {
				y := r + g.CurrentBlock.Row
				x := c + g.CurrentBlock.Col

				g.Cell[y][x] = g.CurrentBlock.ShapeIndex
			}
		}
	}
}

func (g *Game) shiftDownCell(row int) {
	for r := row; r > 0; r-- {
		for c := 0; c < BOARD_COLUMN; c++ {
			g.Cell[r][c] = g.Cell[r-1][c]
		}
	}

	for c := 0; c < BOARD_COLUMN; c++ {
		g.Cell[0][c] = EMPTY
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
func (g *Game) procFullLine() ([][]int, []int) {

	removedRowIndexs := []int{}
	removedLines := [][]int{}

	maxRow := BOARD_ROW
	if g.CurrentBlock.Row+BLOCK_ROW < BOARD_ROW {
		maxRow = g.CurrentBlock.Row + BLOCK_ROW
	}

	for r := g.CurrentBlock.Row; r < maxRow; r++ {

		y := r - g.CurrentBlock.Row
		x := 0

		isFull := true

		for c := 0; c < BOARD_COLUMN; c++ {

			x = c - g.CurrentBlock.Col

			// if cell is empty, check current block
			if g.Cell[r][c] == EMPTY {

				if g.CurrentBlock.inBlock(r, c) {

					if g.CurrentBlock.Shape[y][x] == EMPTY {
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

			removedLines = appendRow(removedLines, g.Cell[r])
			removedRowIndexs = append(removedRowIndexs, r)

			g.shiftDownCell(r)

			// score는 삭제된 line 수 가중치를 줘서 계산
			g.AddScore(10 * len(removedLines))

		} else {

			// currnet block to board
			for i := 0; i < BLOCK_COLUMN; i++ {
				if g.CurrentBlock.Shape[y][i] != EMPTY {
					g.Cell[r][g.CurrentBlock.Col+i] = g.CurrentBlock.ShapeIndex
				}
			}
		}
	}

	return removedLines, removedRowIndexs
}

func (g *Game) AddScore(score int) {
	g.Score += score
}

func (g *Game) firstNextToCurrnetBlock() {
	g.CurrentBlock = NewBlock(g.NextBlockIndexs[0])
	g.NextBlockIndexs = g.NextBlockIndexs[1:]
}

func (g *Game) IsGameOver() bool {
	return g.State == "over"
}

func (g *Game) IsPlaying() bool {
	return g.State == "playing"
}

func (g *Game) IsObserver() bool {
	return g.State == "observer"
}

// nextTern : 다음 턴으로 넘어감
// 삭제된 라인 처리
// 삭제된 라인이 2줄 이상이면, 경쟁자들에게 선물을 보냄
// 새 블럭을 생성
// 새 블럭을 생성하지 못하면, 게임오버
func (g *Game) nextTern() bool {
	removedLines, removedRowIndexs := g.procFullLine()
	if len(removedRowIndexs) > 0 {
		g.SendEraseBlocks(removedRowIndexs)
	}

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

	g.CurrentBlock.Rotate()

	return true
}

func (g *Game) toLeft() bool {
	if !g.isSafeToLeft() {
		return false
	}

	g.CurrentBlock.Col--

	return true
}

func (g *Game) toRight() bool {
	if !g.isSafeToRight() {
		return false
	}

	g.CurrentBlock.Col++

	return true
}

func (g *Game) toDown() bool {
	if !g.isSafeToDown() {
		return false
	}

	g.CurrentBlock.Row++

	return true
}

func (g *Game) toDrop() bool {

	if !g.isSafeToDown() {
		return false
	}

	for g.isSafeToDown() {
		g.CurrentBlock.Row++
	}

	return true
}

func (g *Game) receiveFullBlocks(blocks [][]int) bool {

	g.Cell = append(g.Cell, blocks...)

	// check gamevoer 밀린 위쪽 cell에 블럭이 있으면 게임오버
	for r := 0; r < len(blocks); r++ {
		for c := 0; c < len(g.Cell[r]); c++ {
			if g.Cell[r][c] != EMPTY {
				return false
			}
		}
	}

	g.Cell = g.Cell[len(blocks):]

	return true
}

func (g *Game) gameOver() {
	g.State = "over"
}

func (g *Game) observer() {
	g.State = "observer"
}

// Action : 게임에 메시지를 전달한다.
func (g *Game) Action(msg *Message) {
	g.Ch <- msg
}

// Stop : 게임을 종료시킨다.
func (g *Game) Stop() {
	g.gameOver()
	g.currentBlockToBoard()
	//	g.SendGameOver()
}

// Start : 게임을 시작한다.
func (g *Game) Start() {

	if g.IsPlaying() {
		return
	}

	g.reset()
	g.State = "playing"
	g.SendStartGame()
	g.DurationTime = time.Now()

	go g.run()
	go g.autoDown()
}

func (g *Game) run() {

	for !g.IsGameOver() {
		msg := <-g.Ch
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
				Trace.Print("block-drop")
			}
			if !g.nextTern() {
				return
			}
			g.SendSyncGame()

		case "gift-full-blocks":
			if !g.receiveFullBlocks(msg.Cells) {
				g.gameOver()
				g.currentBlockToBoard()
				g.SendGameOver()
				Trace.Print("gift-full-blocks")
				return
			}
			g.SendSyncGame()
		}

	}

	Trace.Println("run game over")
}

func (g *Game) autoDown() {
	for !g.IsGameOver() {
		Info.Println("autoDown", g.CycleMs, time.Now())
		time.Sleep(time.Millisecond * time.Duration(g.CycleMs))
		if !g.toDown() && !g.nextTern() {
			return
		}
		g.SendSyncGame()

		// Faster every minute
		duration := time.Since(g.DurationTime)
		if duration > time.Minute*1 {
			g.CycleMs -= 100
			g.DurationTime = time.Now()
		}

	}
	Trace.Println("autoDown game over")
}

func (g *Game) SendGameOver() {

	player, ok := Manager.players[g.Owner]
	if !ok {
		Error.Println(g.Owner, " player not found")
		return
	}

	roomList := make([]RoomInfo, 0)
	roomList = appendRoomInfo(roomList, player.Client.ws.rooms[player.RoomId])

	g.ManagerCh <- &Message{
		Action:       "over-game",
		Sender:       g.Owner,
		Cells:        g.Cell,
		CurrentBlock: g.CurrentBlock,
		Score:        g.Score,
		RoomList:     roomList,
	}
}

func (g *Game) SendStartGame() {

	g.ManagerCh <- &Message{
		Action:       "start-game",
		Sender:       g.Owner,
		BlockIndexs:  g.NextBlockIndexs,
		Cells:        g.Cell,
		CurrentBlock: g.CurrentBlock,
	}
}

func (g *Game) SendSyncGame() {
	g.ManagerCh <- &Message{
		Action:       "sync-game",
		Sender:       g.Owner,
		BlockIndexs:  g.NextBlockIndexs,
		Cells:        g.Cell,
		CurrentBlock: g.CurrentBlock,
		Score:        g.Score,
	}
}

func (g *Game) SendGiftFullBlocks(gift [][]int) {

	g.ManagerCh <- &Message{
		Action: "gift-full-blocks",
		Sender: g.Owner,
		Cells:  gift,
	}
}

func (g *Game) SendEraseBlocks(removedRowIndexs []int) {

	g.ManagerCh <- &Message{
		Action:      "erase-blocks",
		Sender:      g.Owner,
		BlockIndexs: removedRowIndexs,
	}
}
