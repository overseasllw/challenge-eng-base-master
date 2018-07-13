package chatroom

import (
	model "app/models"
	"log"
	"strings"
	//	"math"
	m "app/message"
	u "app/user"
	//	rm "app/room"

	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/twinj/uuid"
)

type ChatServer struct {
	OnlineUsers  map[string]Client
	NewMessage   chan *model.Message
	OfflineUsers map[string]Client
	NewUser      chan *Client
	RoomUserList map[string][]Client
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
		//	AllMessages:  []*model.Message{},
		NewMessage:   make(chan *model.Message, 5),
		OnlineUsers:  make(map[string]Client),
		OfflineUsers: make(map[string]Client),
		NewUser:      make(chan *Client, 5),
		RoomUserList: make(map[string][]Client),
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
			User:   model.User{Username: &temp, UserID: msg.UserID},
			Socket: conn,
			Server: server,
		}
		return client
	}

	if msg.Register != nil && *msg.Register == true {
		gue := server.OnlineUsers[*msg.Guestname]
		gue.Username = msg.Username
		gue.Socket = conn
		gue.Server = server
		server.OnlineUsers[*msg.Username] = gue
		//u.Username = msg.Username
		//delete(server.OnlineUsers, *msg.Guestname)
		//	server.RoomUserList[*msg.Room] = append(server.RoomUserList[*msg.Room], gue)
		server.AddMessage(model.Message{
			UUID:           uuid.NewV4().String(),
			MessageType:    "system-message",
			CreatedAt:      time.Now(),
			MessageContent: getPointer(*msg.Username + " has joined the chat."),
			User:           model.User{UserID: msg.UserID, Username: getPointer("system")},
			Room:           msg.Room,
			RoomId:         msg.RoomId,
		})
	}

	if _, exists := server.OnlineUsers[*msg.Username]; exists && len(*msg.Username) >= 3 {
		u := server.OnlineUsers[*msg.Username]
		server.updateOnlineUserList(&u)
		return &u
	} else if _, exists := server.OfflineUsers[*msg.Username]; exists && len(*msg.Username) >= 3 {
		u := server.OfflineUsers[*msg.Username]
		delete(server.OfflineUsers, *msg.Username)
		server.OnlineUsers[*msg.Username] = u
		server.updateOnlineUserList(&u)
		return &u
	}
	client := &Client{
		User:   model.User{Username: msg.Username, UserID: msg.UserID},
		Socket: conn,
		Server: server,
	}

	server.OnlineUsers[*msg.Username] = *client
	server.updateOnlineUserList(client)
	///server.RoomUserList[*msg.Room] = append(server.RoomUserList[*msg.Room], *client)
	server.AddMessage(model.Message{
		UUID:           uuid.NewV4().String(),
		MessageType:    "system-message",
		CreatedAt:      time.Now(),
		MessageContent: getPointer(*msg.Username + " has joined the chat."),
		User:           model.User{UserID: msg.UserID, Username: getPointer("system")},
		Room:           msg.Room,
		RoomId:         msg.RoomId,
	})

	client.Send([]*model.Message{
		&msg,
	})

	return client
}

// Leaving the chatroom
func (server *ChatServer) Leave(name string, room string) {
	server.OfflineUsers[name] = server.OnlineUsers[name]
	delete(server.OnlineUsers, name)

	server.AddMessage(
		model.Message{
			UUID:           uuid.NewV4().String(),
			MessageType:    "system-message",
			CreatedAt:      time.Now(),
			MessageContent: getPointer(name + " has left the chat."),
			User:           model.User{UserID: 0, Username: getPointer("system")},
			Room:           &room,
			//	RoomId:         msg.RoomId,
		})
}

func getPointer(s string) *string {
	return &s
}

// Adding message to queue
func (server *ChatServer) AddMessage(message model.Message) {
	server.NewMessage <- &message
	if message.Username != nil {
		if message.UserID != 0 {
			uid := u.CheckUserExist(*message.Username)
			if uid != nil {
				message.UserID = *uid

			}
		}
		m.CreateNewMessage(&message)
		//m.ReadMessage(&message)
	}
}

func Listen(server *ChatServer, c echo.Context) error {
	//c.GET("/ws", server.Listen)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	if err != nil {
		ws.Close()
		return err
	}
	defer ws.Close()
	msg := model.Message{}
	err = ws.ReadJSON(&msg)
	msg.UUID = uuid.NewV4().String()
	//j, _ := json.Marshal(msg)
	//log.Print(string(j))

	if err != nil {

		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		} else {
			log.Print(err)
		}
		ws.Close()
		return err
	}
	uid := u.CreateNewUser(*msg.Username)
	msg.UserID = uid
	user := server.Join(msg, ws)
	user.Room = msg.Room
	if user == nil {
		user.Exit()
		return err
	}
	if msg.Room != nil {
		server.RoomUserList[*msg.Room] = append(server.RoomUserList[*msg.Room], *user)
	}
	user.UserID = uid
	//log.Print(msg.UserID)
	for {
		msg := model.Message{}
		err = ws.ReadJSON(&msg)
		msg.UUID = uuid.NewV4().String()
		msg.UserID = user.UserID

		if msg.Room != nil {
			if _, ok := server.RoomUserList[*msg.Room]; !ok {
				server.RoomUserList[*msg.Room] = append(server.RoomUserList[*msg.Room], *user)
			}
		}
		if err != nil {

			server.updateOnlineUserList(user)

			user.Exit()
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			} else {
				log.Print(err)
			}
			return err
		}
		if user.Username != nil && msg.Username != nil {
			if strings.TrimSpace(*user.Username) != strings.TrimSpace(*msg.Username) {
				u.UpdateUser(*user.Username, *msg.Username)
				delete(server.OnlineUsers, *msg.Guestname)
				user.Username = msg.Username
				user.Register = msg.Register
				user.Guestname = msg.Guestname
				server.OnlineUsers[*user.Username] = *user
				server.updateOnlineUserList(user)
			}

		}
		if msg.MessageType == "typing_indicator" {
			//	user.NewMessage(msg)
			server.NewMessage <- &model.Message{
				UUID:           "typing_indicator_" + *user.User.Username,
				MessageType:    "typing_indicator",
				CreatedAt:      time.Now(),
				MessageContent: getPointer(*user.User.Username + " is typing."),
				Room:           msg.Room,
				User:           model.User{UserID: 0, Username: user.User.Username},
			}
			//	return err
		} else if msg.MessageType == "stop_typing_indicator" {
			server.NewMessage <- &model.Message{
				UUID:           "typing_indicator_" + *user.User.Username,
				MessageType:    "stop_typing_indicator",
				CreatedAt:      time.Now(),
				MessageContent: getPointer(*user.User.Username + " is typing."),
				Room:           msg.Room,
				User:           model.User{UserID: 0, Username: user.User.Username},
			}
		} else {
			// Write
			user.NewMessage(msg)
		}
	}
	//return err
}

func (server *ChatServer) updateOnlineUserList(client *Client) {
	server.NewUser <- client
}

// Broadcasting all the messages in the queue in one block
func (server *ChatServer) BroadCast() {

	messages := make([]*model.Message, 0)
	userList := make(map[string]interface{})
InfiLoop:
	for {
		select {
		case message := <-server.NewMessage:
			messages = append(messages, message)
		case <-server.NewUser:
			//	user["username"] = *newUser.Username
			//	user["created_at"] = time.Now().String()
			userList["message_type"] = "user_list"
			us := []model.User{}
			for _, c := range server.OnlineUsers {
				us = append(us, c.User)
			}
			userList["list"] = us
		default:
			break InfiLoop
		}
	}
	if len(userList) > 0 {
		for _, client := range server.OnlineUsers {
			client.Socket.WriteJSON([]map[string]interface{}{
				userList,
			})
		}
	}

	if len(messages) > 0 {
		for _, rname := range server.RoomUserList[*messages[0].Room] {
			//	for _, client := range server.OnlineUsers {
			n := rname.User.Username
			if _, ok := server.OnlineUsers[*n]; ok {

				rname.Send(messages)
			}
			//	}
		}
	}
}
