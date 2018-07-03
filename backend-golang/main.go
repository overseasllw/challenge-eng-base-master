package main

import (
	"app/chatroom"
	m "app/message"
	rm "app/room"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // not checking origin
}

func main() {

	/*http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		var result string
		if err := common.DB.QueryRow(`SELECT col FROM test`).Scan(&result); err != nil {
			log.Panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"result":  result,
			"backend": "go",
		}); err != nil {
			log.Panic(err)
		}
	})*/
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", m.GetMessageListHandler)
	e.GET("/test", m.GetMessageListHandler)
	// message api
	//e.POST("/api/v1/messages/", chat.PostMessageHandler)

	server := chatroom.NewServer()
	go server.Init()
	e.GET("/ws", func(c echo.Context) error {
		chatroom.Listen(server, c)
		return nil
	})
	e.GET("/api/v1/messages/", func(c echo.Context) error {
		//	server.GetUserList()
		err := m.GetMessageListHandler(c)
		return err
	})

	e.POST("/api/v1/rooms/", rm.CreateNewRoomHandler)
	e.GET("/api/v1/rooms/", rm.GetRoomListHandler)
	e.Logger.Fatal(e.Start(":8000"))

	/*	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}*/
}
