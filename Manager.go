package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	RoomId int
	Client *Client
}

type manager struct {
	players map[string]*Player // players list (nick, player)
	ch      chan *Message      // Game -> Manager
}

var Manager *manager

func init() {
	Manager = newGameManager()
	go Manager.run()
}

func newGameManager() *manager {
	return &manager{
		players: make(map[string]*Player),
		ch:      make(chan *Message),
	}
}

func (p *Player) Send(msg *Message) {
	p.Client.send <- msg
}

func (gm *manager) Register(roomId int, client *Client) {

	// new player
	player := &Player{
		RoomId: roomId,
		Client: client,
	}

	gm.players[client.Nick] = player
}

func (gm *manager) Unregister(client *Client) {

	// remove player
	if gm.players[client.Nick] != nil {
		delete(gm.players, client.Nick)
	}

	if client.game != nil && !client.game.IsGameOver() {
		client.game.Stop()
	}
}

func (gm *manager) NewGame(roomId int, client *Client) *Game {

	// new game
	game := NewGame(gm.ch, client.Nick)
	client.game = game

	return game
}

func (gm *manager) run() {
	for {
		msg := <-gm.ch
		gm.handleMessage(msg)

		//time.Sleep(1 * time.Millisecond)
	}
}

// top 20
func (gm *manager) SaveTop20(nick string, score int) int {
	// 파일을 읽는다.
	file, err := os.OpenFile("top20.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		Error.Println(err)
		return -1
	}
	defer file.Close()

	buf := ""

	// 한줄씩 읽어서, nick과 score를 비교한다.
	// score가 더 크면, 그 줄을 지우고, 새로운 줄을 삽입한다.
	// 그렇지 않으면, 그냥 넘어간다.
	// 20개가 넘으면, 마지막 줄을 지운다.

	rank := -1

	reader := bufio.NewReader(file)
	for i := 0; i < 20; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		info := strings.Split(string(line), ",")
		if len(info) != 2 {
			Error.Printf("Invalid line %d: %s", i, string(line))
			break
		}

		s, err := strconv.Atoi(info[1])
		if err != nil {
			Error.Printf("Invalid score %d: %s", i, string(line))
			break
		}

		if rank == -1 && score > s {
			buf += fmt.Sprintf("%s,%d,%s\n", nick, score, time.Now().Format("2006-01-02T15:04:05"))
			i++
			rank = i
		}

		if i < 20 {
			buf += string(line) + "\n"
		}
	}

	if rank != -1 {
		// 파일을 다시 쓴다.
		file.Seek(0, 0)
		n, err := file.WriteString(buf)
		if n > 20 || err != nil {
			Error.Printf("Invalid write: %d, %s", n, err)
		}
	}

	file.Close()

	return rank
}

func (gm *manager) startGame(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.RoomId = player.RoomId
		player.Client.send <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) overGame(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.RoomId = player.RoomId
		player.Client.ws.broadcast <- msg

		winner := &Client{}

		// 방안의 사용자 중에 한명만 state가 playing이면, 그 사용자도 게임을 중지 시킨다.
		for _, client := range player.Client.ws.rooms[player.RoomId].Clients {

			if client.Nick == msg.Sender {
				continue
			}

			if client.game != nil && client.game.IsPlaying() {
				if winner == nil {
					winner = client
				} else {
					// 두명 이상
					winner = nil
					break
				}
			}
		}

		if winner != nil {
			gm.endGame(winner.Nick)
		}

	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) endGame(nick string) {
	if player, ok := gm.players[nick]; ok {

		player.Client.game.Stop()

		// 승자에게 플레이어수 x 100점 추가
		player.Client.game.AddScore(100 * len(player.Client.ws.rooms[player.RoomId].Clients))

		// top 20 안에 있으면 추가
		rank := gm.SaveTop20(nick, player.Client.game.score)

		msg := &Message{
			Action: "end-game",
			Sender: nick,
			RoomId: player.RoomId,
			Score:  player.Client.game.score,
			Data:   strconv.Itoa(rank),
		}
		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("Unknown player:", nick)
	}
}

func (gm *manager) syncGame(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.RoomId = player.RoomId
		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) giftFullBlocks(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.RoomId = player.RoomId
		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) getGame(nick string) *Game {
	if player, ok := gm.players[nick]; ok {
		return player.Client.game
	} else {
		return nil
	}
}

// getRoomId
// return -1 if not found
func (gm *manager) getRoomId(nick string) int {
	if player, ok := gm.players[nick]; ok {
		return player.RoomId
	} else {
		return -1
	}
}

func (gm *manager) getClient(nick string) *Client {
	if player, ok := gm.players[nick]; ok {
		return player.Client
	} else {
		return nil
	}
}

// Game 에서 발생한 이벤트 처리
func (gm *manager) handleMessage(msg *Message) {
	switch msg.Action {
	case "start-game":
		gm.startGame(msg)

	case "over-game":
		gm.overGame(msg)

	case "sync-game":
		gm.syncGame(msg)

	case "gift-full-blocks":
		gm.giftFullBlocks(msg)

	default:
		Error.Fatalln("Unknown action:", msg)
	}
}
