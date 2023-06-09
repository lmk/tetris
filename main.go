// 웹 기반의 대전 테트리스 게임을 만들어보자.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var conf Config

const (
	MAX_CHAN = 10
)

func init() {
	rand.Seed(time.Now().UnixNano())

	gin.SetMode(gin.ReleaseMode)

	initFlag()
}

// 웹 서버를 실행하는 함수
func runServer() {

	// 웹소켓 핸들러 생성한다.
	wsServer := NewWebsocketServer()
	go wsServer.Run()
	go wsServer.Report()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	wsURI := fmt.Sprintf("ws://%s:%d/ws", conf.Domain, conf.Port)
	securityPolicy := ""
	if conf.Https {
		wsURI = fmt.Sprintf("ws://%s/ws", conf.Domain)
		securityPolicy = "upgrade-insecure-requests"
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"wsServer":       wsURI,
			"securityPolicy": securityPolicy,
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		serveWs(c, WAITITNG_ROOM, wsServer)
	})

	r.Use(static.Serve("/", static.LocalFile("public", true)))

	r.Run(fmt.Sprintf(":%d", conf.Port))
}

func main() {

	initConf()

	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	Info.Printf("config: %s", conf.makePretty())

	runServer()
}

func usage() {
	fmt.Printf("Usage: %s -config=<config file>\n", os.Args[0])
	flag.PrintDefaults()
}
