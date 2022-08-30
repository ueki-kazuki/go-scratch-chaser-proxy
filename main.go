package main

//
// https://zenn.dev/empenguin/articles/bcf95c19451020
//
import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

		var client *Client
		host := "127.0.0.1"
		port := 2009
		name := "User1"

		for {
			// Client からのメッセージを読み込む
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
			}
			log.Println(msg)
			args := strings.Split(msg, " ")
			cmd, args := args[0], args[1:]
			log.Println(cmd, args)

			var response string
			switch cmd {
			case "connect":
				if len(args) > 2 {
					host = args[0]
					name = args[2]
					port, err = strconv.Atoi(args[1])
					if err != nil {
						c.String(http.StatusInternalServerError, err.Error())
					}
				}

				client, err = NewClient(name, host, port)
				if err != nil {
					log.Println(err)
				}
				defer client.Close()
			case "gr":
				response, err = client.GetReady()
			default:
				response, err = client.Order(cmd)
			}
			if err != nil {
				log.Println(err)
			}
			log.Println(cmd, response)

			if client == nil {
				continue
			}

			if client.GameSet {
				err := websocket.Message.Send(ws, fmt.Sprintf("%s", "GameSet!!"))
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
				}
				break
			}

			// CHaserClient の結果を整形してWebClientにメッセージを送信する
			err := websocket.Message.Send(ws, fmt.Sprintf("%s", response))
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
