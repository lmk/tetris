// 웹 기반의 대전 테트리스 게임을 만들어보자.

package main

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// 웹 서버를 실행하는 함수
func runServer() {

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("public", true)))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 웹 소켓 서버를 생성한다.
	wsServer := NewWebsocketServer()

	// websocker 통신을 위한 라우팅을 추가한다.
	r.GET("/ws", func(c *gin.Context) {
		serveWs(c, "list", wsServer)
	})

	r.Run(":8080")
}

func main() {

	runServer()

}
