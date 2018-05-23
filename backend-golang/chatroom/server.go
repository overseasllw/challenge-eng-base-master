package chatroom

import (
	model "app/models"
	"encoding/json"
	"log"
	//	"math"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type ChatServer struct {
	OnlineUsers  map[string]Client
	NewMessage   chan *model.Message
	OfflineUsers map[string]Client
	AllMessages  []*model.Message
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // not checking origin
	}
)

func NewServer() (server *ChatServer) {
	return &ChatServer{
		AllMessages:  []*model.Message{},
		NewMessage:   make(chan *model.Message, 5),
		OnlineUsers:  make(map[string]Client),
		OfflineUsers: make(map[string]Client),
	}
}

// Initializing the chatroom
func (server *ChatServer) Init() {
	go func() {
		for {
			server.BroadCast()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (server *ChatServer) Join(msg model.Message, conn *websocket.Conn) *Client {
	if msg.Username == nil {
		temp := "guest" + time.Now().String()
		client := &Client{
			User:   model.User{Username: &temp},
			Socket: conn,
			Server: server,
		}
		return client
	}
	if _, exists := server.OnlineUsers[*msg.Username]; exists && len(*msg.Username) >= 3 {
		u := server.OnlineUsers[*msg.Username]
		return &u
	} else if _, exists := server.OfflineUsers[*msg.Username]; exists && len(*msg.Username) >= 3 {
		u := server.OfflineUsers[*msg.Username]
		delete(server.OfflineUsers, *msg.Username)
		server.OnlineUsers[*msg.Username] = u
		return &u
	}
	client := &Client{
		User:   model.User{Username: msg.Username},
		Socket: conn,
		Server: server,
	}
	server.OnlineUsers[*msg.Username] = *client

	server.AddMessage(model.Message{
		MessageID:      3232,
		MessageType:    "system-message",
		CreatedAt:      time.Now(),
		MessageContent: getPointer(*msg.Username + " has joined the chat."),
		//	User:           model.User{Username: name},
	})
	client.Send([]*model.Message{
		&msg,
	})

	return client
}

// Leaving the chatroom
func (server *ChatServer) Leave(name string) {
	server.OfflineUsers[name] = server.OnlineUsers[name]
	delete(server.OnlineUsers, name)

	server.AddMessage(
		model.Message{
			MessageID:      332,
			MessageType:    "system-message",
			CreatedAt:      time.Now(),
			MessageContent: getPointer(name + " has left the chat."),
		})
}

func getPointer(s string) *string {
	return &s
}

// Adding message to queue
func (server *ChatServer) AddMessage(message model.Message) {
	//message.MessageType = "user-message"
	server.NewMessage <- &message
	if message.Username != nil {
		//StoreMessage(message)
		server.AllMessages = append(server.AllMessages, &message)
	}
}

func Listen(server *ChatServer, c echo.Context) error {
	//c.GET("/ws", server.Listen)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	if err != nil {
		//	log.Print(err)
		ws.Close()
		return err
	}
	defer ws.Close()
	msg := model.Message{}
	err = ws.ReadJSON(&msg)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		} else {
			log.Print(err)
		}
		ws.Close()
		return err
	}

	user := server.Join(msg, ws)

	if user == nil {
		log.Print(err)
		return err
	}
	for {
		msg := model.Message{}
		err = ws.ReadJSON(&msg)
		j, _ := json.Marshal(msg)
		log.Print(string(j))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			} else {
				log.Print(err)
			}
			user.Exit()
			return err
		}

		// Write
		user.NewMessage(msg)
	}
	//return err
}

// Broadcasting all the messages in the queue in one block
func (server *ChatServer) BroadCast() {

	messages := make([]*model.Message, 0)

InfiLoop:
	for {
		select {
		case message := <-server.NewMessage:
			messages = append(messages, message)
		default:
			break InfiLoop
		}
	}

	if len(messages) > 0 {
		for _, client := range server.OnlineUsers {
			client.Send(messages)
		}
	}
}