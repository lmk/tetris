package main

import (
	"strconv"
	"sync"
)

type Player struct {
	RoomId int
	Client *Client
}

type manager struct {
	players map[string]*Player // players list (nick, player)
	ch      chan *Message      // Game -> Manager
	mu      sync.RWMutex       // protects players map
}

var Manager *manager

func init() {
	Manager = newGameManager()
	go Manager.run()
}

func newGameManager() *manager {
	return &manager{
		players: make(map[string]*Player),
		ch:      make(chan *Message, MAX_CHAN),
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

	gm.mu.Lock()
	gm.players[client.Nick] = player
	gm.mu.Unlock()
}

func (gm *manager) Unregister(client *Client) {

	// remove player
	gm.mu.Lock()
	if gm.players[client.Nick] != nil {
		delete(gm.players, client.Nick)
	}
	gm.mu.Unlock()

	if client.Game != nil && !client.Game.IsGameOver() {
		client.Game.Stop()
	}
}

func (gm *manager) NewGame(roomId int, client *Client) *Game {

	// new game
	game := NewGame(gm.ch, client.Nick)
	client.Game = game

	return game
}

func (gm *manager) run() {
	for {
		msg, ok := <-gm.ch
		if !ok {
			Error.Println("run", "channel closed")
			return
		}

		gm.HandleMessage(msg)

		//time.Sleep(1 * time.Millisecond)
	}
}

func (gm *manager) getRankList(count int) (rankList []Rank) {

	rankList, err := ReadRankList(count)
	if err != nil {
		Error.Println("getRankList fail", err)
		return nil
	}

	return rankList
}

func (gm *manager) SaveTop(nick string, score int) int {
	rank, err := SaveTopRank(nick, score)
	if err != nil {
		Error.Println("SaveTop fail", err)
		return -1
	}

	return rank
}

func (gm *manager) startGame(msg *Message) {
	gm.mu.RLock()
	player, ok := gm.players[msg.Sender]
	gm.mu.RUnlock()

	if ok {
		msg.RoomId = player.RoomId
		player.Client.send <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) overGame(msg *Message) {

	Info.Printf("overGame: %s", msg.Sender)

	gm.mu.RLock()
	player, ok := gm.players[msg.Sender]
	gm.mu.RUnlock()

	if ok {
		msg.RoomId = player.RoomId
		player.Client.ws.broadcast <- msg

		var winner *Client = nil

		if len(player.Client.ws.rooms[player.RoomId].Clients) == 1 {
			// single play
			winner = player.Client
		} else {
			//multi play

			// 방안의 사용자 중에 한명만 state가 playing이면, 그 사용자도 게임을 중지 시킨다.
			for _, client := range player.Client.ws.rooms[player.RoomId].Clients {

				if client.Game != nil && client.Game.IsPlaying() {
					if winner == nil {
						winner = client
					} else {
						// 두명 이상
						winner = nil
						break
					}
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
	gm.mu.RLock()
	player, ok := gm.players[nick]
	gm.mu.RUnlock()

	if ok {

		if player.Client.Game == nil {
			Warning.Println(nick, "Game is nil")
			return
		}

		player.Client.Game.Stop()

		// 승자에게 플레이어수 x 100점 추가
		player.Client.Game.AddScore(WINNER_BONUS * (len(player.Client.ws.rooms[player.RoomId].Clients) - 1))

		// top 안에 있으면 추가
		rank := gm.SaveTop(nick, player.Client.Game.Score)

		msg := &Message{
			Action: "end-game",
			Sender: nick,
			RoomId: player.RoomId,
			Score:  player.Client.Game.Score,
		}
		msg.RoomList = appendRoomInfo(msg.RoomList, player.Client.ws.rooms[player.RoomId])

		if rank > 0 {
			msg.Data = strconv.Itoa(rank)
		}

		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("Unknown player:", nick)
	}
}

func (gm *manager) syncGame(msg *Message) {
	gm.mu.RLock()
	player, ok := gm.players[msg.Sender]
	gm.mu.RUnlock()

	if ok {
		msg.RoomId = player.RoomId
		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("[syncGame] Unknown player:", msg.Sender, "room:", msg.RoomId)
	}
}

func (gm *manager) eraseBlocks(msg *Message) {
	gm.mu.RLock()
	player, ok := gm.players[msg.Sender]
	gm.mu.RUnlock()

	if ok {
		msg.RoomId = player.RoomId
		player.Client.send <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) giftFullBlocks(msg *Message) {
	gm.mu.RLock()
	player, ok := gm.players[msg.Sender]
	gm.mu.RUnlock()

	if ok {
		msg.RoomId = player.RoomId

		player.Client.ws.broadcast <- msg

		// 방안의 사용자들에게 full block을 전달한다.
		for _, client := range player.Client.ws.rooms[player.RoomId].Clients {
			if client.Game != nil && client.Nick != msg.Sender && client.Game.IsPlaying() {
				client.Game.Ch <- msg
			}
		}

	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) getGame(nick string) *Game {
	gm.mu.RLock()
	player, ok := gm.players[nick]
	gm.mu.RUnlock()

	if ok {
		return player.Client.Game
	} else {
		return nil
	}
}

// getRoomId
// return -1 if not found
func (gm *manager) getRoomId(nick string) int {
	gm.mu.RLock()
	player, ok := gm.players[nick]
	gm.mu.RUnlock()

	if ok {
		return player.RoomId
	} else {
		return -1
	}
}

func (gm *manager) getClient(nick string) *Client {
	gm.mu.RLock()
	player, ok := gm.players[nick]
	gm.mu.RUnlock()

	if ok {
		return player.Client
	} else {
		return nil
	}
}

// Game 에서 발생한 이벤트 처리
func (gm *manager) HandleMessage(msg *Message) {
	switch msg.Action {
	case "start-game":
		gm.startGame(msg)

	case "over-game":
		gm.overGame(msg)

	case "sync-game":
		gm.syncGame(msg)

	case "gift-full-blocks":
		gm.giftFullBlocks(msg)

	case "erase-blocks":
		gm.eraseBlocks(msg)

	default:
		Error.Fatalln("Unknown action:", msg)
	}
}
