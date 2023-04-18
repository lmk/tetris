package main

import (
	"encoding/json"
)

type Player struct {
	RoomId int
	Client *Client
}

func (p *Player) Send(msg *Message) {
	p.Client.send <- msg
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

func (gm *manager) NewGame(roomId int, client *Client) *Game {

	// new game
	game := NewGame(gm.ch, len(gm.players)+1, client.nick)
	client.game = game

	// new player
	player := &Player{
		RoomId: roomId,
		Client: client,
	}

	gm.players[client.nick] = player

	// send next block indexs to client
	jsonNexts, _ := json.Marshal(game.nextBlockIndexs)
	client.send <- &Message{Action: "next-block", Data: string(jsonNexts)}

	return game
}

func (gm *manager) run() {
	for {
		msg := <-gm.ch
		gm.handleMessage(msg)

		//time.Sleep(1 * time.Millisecond)
	}
}

func (gm *manager) gameOver(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.roomId = player.RoomId
		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) syncGame(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.roomId = player.RoomId
		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) giftFullBlocks(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.roomId = player.RoomId
		player.Client.ws.broadcast <- msg
	} else {
		Warning.Println("Unknown player:", msg)
	}
}

func (gm *manager) nextBlock(msg *Message) {
	if player, ok := gm.players[msg.Sender]; ok {
		msg.roomId = player.RoomId
		player.Client.send <- msg
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

// Game 에서 발생한 이벤트 처리
func (gm *manager) handleMessage(msg *Message) {
	switch msg.Action {
	case "over-game":
		gm.gameOver(msg)

	case "sync-game":
		gm.syncGame(msg)

	case "gift-full-blocks":
		gm.giftFullBlocks(msg)

	case "next-block":
		gm.nextBlock(msg)

	default:
		Warning.Println("Unknown action:", msg)
	}
}
