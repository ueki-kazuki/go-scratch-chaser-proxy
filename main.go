package main

//
// https://zenn.dev/empenguin/articles/bcf95c19451020
//
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func handleWebSocket(c *gin.Context) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// 初回のメッセージを送信
		err := websocket.Message.Send(ws, "Server: Hello, Client!")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		for {
			// Client からのメッセージを読み込む
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			}

			// Client からのメッセージを元に返すメッセージを作成し送信する
			err := websocket.Message.Send(ws, fmt.Sprintf("Server: \"%s\" received!", msg))
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			}
		}
	}).ServeHTTP(c.Writer, c.Request)
}

func main() {
	r := gin.Default()
	r.GET("/ws", handleWebSocket)
	// r.Static("/", "./public/")
	// r.StaticFS("/", http.Dir("public"))
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/main.js", "./public/main.js")
	r.Run(":8080")
}
