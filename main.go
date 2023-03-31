// 웹 기반의 대전 테트리스 게임을 만들어보자.

package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

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

	// websocker 통신을 위한 라우팅을 추가한다.
	r.GET("/ws/:room", func(c *gin.Context) {
		serveWs(c, c.Param("room"), wsServer)
	})

	r.Run(":8080")
}

func main() {

	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	runServer()

}
