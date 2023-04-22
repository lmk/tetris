// 웹 기반의 대전 테트리스 게임을 만들어보자.

package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 웹 서버를 실행하는 함수
func runServer() {

	// 웹소켓 핸들러 생성한다.
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("public", true)))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// // websocker 통신을 위한 라우팅을 추가한다.
	// r.GET("/ws/:room", func(c *gin.Context) {
	// 	//string to int
	// 	roomId, err := strconv.Atoi(c.Param("room"))
	// 	if err != nil {
	// 		Error.Println("Invaild URI", err)
	// 	} else {
	// 		serveWs(c, roomId, wsServer)
	// 	}
	// })

	r.GET("/ws", func(c *gin.Context) {
		serveWs(c, WAITITNG_ROOM, wsServer)
	})

	r.Run(":8090")
}

func main() {

	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	runServer()

}
