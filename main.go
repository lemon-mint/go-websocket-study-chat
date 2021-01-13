package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/ws/chat", wsChatEndPoint)
	e.GET("/chat/send", sendMsg)
	e.Static("/static", "static")
	e.File("/", "views/index.html")
	e.Logger.Fatal(e.Start(":54791"))
}

var chats []string = make([]string, 0)
var chatlock sync.Mutex

var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsChatEndPoint(c echo.Context) error {
	type serverMSG struct {
		Ts   int64  `json:"ts"`
		Hts  string `json:"hts"`
		Body string `json:"body"`
	}

	ws, err := wsUpgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil
	}
	defer ws.Close()
	clientRequest := new(struct {
		Pkt  string `json:"pkt"`
		Body string `json:"body"`
	})
	err = ws.ReadJSON(clientRequest)
	if err != nil {
		return nil
	}
	var curs int = 0
	if clientRequest.Pkt == "getchat" {
		cTime := time.Now().UTC()
		err := ws.WriteJSON(serverMSG{
			Ts:   cTime.UnixNano(),
			Hts:  cTime.String(),
			Body: "Welcome to the public chat room",
		})
		if err != nil {
			return nil
		}
		for {
			time.Sleep(time.Millisecond * 100)
			chatlock.Lock()
			for i := curs; i < len(chats); i++ {
				cTime := time.Now().UTC()
				err := ws.WriteJSON(serverMSG{
					Ts:   cTime.UnixNano(),
					Hts:  cTime.String(),
					Body: chats[i],
				})
				if err != nil {
					chatlock.Unlock()
					return nil
				}
			}
			chatlock.Unlock()
			curs = len(chats)
		}
	}
	return nil
}

func sendMsg(c echo.Context) error {
	chatlock.Lock()
	chats = append(chats, c.Request().URL.Query().Get("data"))
	chatlock.Unlock()
	return c.String(http.StatusOK, "OK")
}
