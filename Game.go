package main

import (
	"math/rand"
	"sync"
	"time"
)

const (
	BOARD_ROW    = 15
	BOARD_COLUMN = 10
	BOARD_CENTER = (BOARD_COLUMN - 3) / 2

	BLOCK_ROW    = 4
	BLOCK_COLUMN = 4

	EMPTY = 0

	MIN_CYCLE_MS     = 100  // minimum cycle time in ms
	INITIAL_CYCLE_MS = 1000 // initial cycle time in ms
	CYCLE_DECREASE   = 100  // cycle time decrease per minute
	SCORE_PER_LINE   = 10   // score per removed line
	WINNER_BONUS     = 100  // bonus score per opponent for winner
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
	mu              sync.RWMutex  `json:"-"` // protects Cell and game state
}

// NewGame create a new game
// ch: request channel
// id: game id
func NewGame(ch chan *Message, owner string) *Game {
	game := Game{
		Owner:        owner,
		Ch:           make(chan *Message, MAX_CHAN),
		ManagerCh:    ch,
		State:        "ready",
		Score:        0,
		CurrentBlock: NewBlock((rand.Intn(len(SHAPES)) + 1)),
		CycleMs:      INITIAL_CYCLE_MS,
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
	g.CycleMs = INITIAL_CYCLE_MS
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
// currentBlock을 board에 복사한 후 full line을 검사
// full line이면 row를 삭제하고, 위 row를 한칸씩 내림
func (g *Game) procFullLine() ([][]int, []int) {

	// 먼저 현재 블록을 보드에 복사
	g.currentBlockToBoard()

	removedRowIndexs := []int{}
	removedLines := [][]int{}

	// 현재 블록이 영향을 미친 행만 검사
	minRow := g.CurrentBlock.Row
	maxRow := BOARD_ROW
	if g.CurrentBlock.Row+BLOCK_ROW < BOARD_ROW {
		maxRow = g.CurrentBlock.Row + BLOCK_ROW
	}

	// 아래에서 위로 검사하면서 가득 찬 라인 삭제
	// 라인을 삭제하면 위의 라인들이 내려오므로, 같은 행을 다시 검사해야 함
	for r := minRow; r < maxRow; r++ {
		isFull := true

		// 해당 행이 가득 찼는지 확인
		for c := 0; c < BOARD_COLUMN; c++ {
			if g.Cell[r][c] == EMPTY {
				isFull = false
				break
			}
		}

		if isFull {
			removedLines = appendRow(removedLines, g.Cell[r])
			removedRowIndexs = append(removedRowIndexs, r)

			g.shiftDownCell(r)

			// 라인을 삭제하면 위의 라인들이 내려오므로
			// 같은 행(r)을 다시 검사하기 위해 인덱스를 하나 감소
			r--
			maxRow--

			// score는 삭제된 line 수 가중치를 줘서 계산
			g.AddScore(SCORE_PER_LINE * len(removedLines))
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
	// 새로운 블록을 받으면 아래에서 밀어올려지므로,
	// 위쪽의 행들이 밀려나게 됨

	// 먼저 밀려날 위쪽 행들에 블록이 있는지 확인 (게임오버 체크)
	for r := 0; r < len(blocks); r++ {
		for c := 0; c < BOARD_COLUMN; c++ {
			if g.Cell[r][c] != EMPTY {
				return false
			}
		}
	}

	// 위쪽 len(blocks)개 행을 제거하고 아래에 새 블록 추가
	g.Cell = append(g.Cell[len(blocks):], blocks...)

	return true
}

func (g *Game) gameOver() {
	g.State = "over"
}

func (g *Game) setObserver() {
	g.State = "observer"
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
		msg, ok := <-g.Ch
		if !ok {
			Error.Println("channel closed")
			return
		}

		Trace.Println("run game", g.Owner, msg)
		switch msg.Action {
		// case "start-game":
		// 	g.Start()

		case "stop-game":
			g.Stop()

		case "observer":
			g.setObserver()

		case "block-rotate":
			g.toRotate()
			g.SendSyncGame(msg.Action)

		case "block-left":
			g.toLeft()
			g.SendSyncGame(msg.Action)

		case "block-right":
			g.toRight()
			g.SendSyncGame(msg.Action)

		case "block-down":
			if !g.toDown() {
				if !g.nextTern() {
					return
				}
			}
			g.SendSyncGame(msg.Action)

		case "block-drop":
			if !g.toDrop() {
				// 블록을 놓을 수 없으면 게임 오버
				g.gameOver()
				g.currentBlockToBoard()
				g.SendGameOver()
				Trace.Println("block-drop: game over (cannot drop)")
				return
			}
			// 다음 턴으로 (procFullLine에서 currentBlockToBoard 호출됨)
			if !g.nextTern() {
				return
			}
			g.SendSyncGame(msg.Action)

		case "gift-full-blocks":
			if !g.receiveFullBlocks(msg.Cells) {
				g.gameOver()
				g.currentBlockToBoard()
				g.SendGameOver()
				Trace.Println("gift-full-blocks")
				return
			}
			g.SendSyncGame(msg.Action)

		case "auto-down":
			if !g.toDown() {
				if !g.nextTern() {
					// Game over
					g.gameOver()
					g.currentBlockToBoard()
					g.SendGameOver()
					return
				}
			}
			Debug.Println("auto-down s ", g.Owner, len(g.ManagerCh))
			g.SendSyncGame(msg.Action)
			Debug.Println("auto-down e ", g.Owner, len(g.ManagerCh))

		}
	}

	Trace.Println("run game over")
}

func (g *Game) autoDown() {
	for !g.IsGameOver() {

		time.Sleep(time.Millisecond * time.Duration(g.CycleMs))

		Debug.Println("autoDown s:", g.Owner, len(g.Ch))
		g.Ch <- &Message{
			Action: "auto-down",
		}
		Debug.Println("autoDown e:", g.Owner, len(g.Ch))

		// Faster every minute
		duration := time.Since(g.DurationTime)
		if duration > time.Minute*1 {
			if g.CycleMs > MIN_CYCLE_MS {
				g.CycleMs -= CYCLE_DECREASE
				// CycleMs가 최소값 이하로 떨어지지 않도록 보장
				if g.CycleMs < MIN_CYCLE_MS {
					g.CycleMs = MIN_CYCLE_MS
				}
			}
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

func (g *Game) SendSyncGame(data string) {
	g.ManagerCh <- &Message{
		Action:       "sync-game",
		Sender:       g.Owner,
		Data:         data,
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
