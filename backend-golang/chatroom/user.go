package chatroom

import (
	model "app/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	model.User
	Socket *websocket.Conn
	Server *ChatServer
	Room   *string
}

var (
	pongWait = 60 * time.Second
)

func (client *Client) ReadMessage() {
	defer func() {
		client.Server.OfflineUsers[*client.User.Username] = *client
		client.Socket.Close()
	}()
	client.Socket.SetPongHandler(func(string) error { client.Socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msg := model.Message{}
		err := client.Socket.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		client.Server.NewMessage <- &msg
	}

}

func (client *Client) WriteMessage() {
	defer func() {
		client.Socket.Close()
	}()
	messages := make([]*model.Message, 0)
	for {
		select {
		case message := <-client.Server.NewMessage:
			messages = append(messages, message)
			for _, m := range client.Server.OnlineUsers {
				m.Send(messages)
			}
		default:
			return
		}
	}
}

func (client *Client) Send(messages []*model.Message) {
	client.Socket.WriteJSON(messages)
}

// Client has a new message to broadcast
func (client *Client) NewMessage(message model.Message) {
	message.MessageType = "user-message"
	client.Server.AddMessage(message)
}

// Exiting out
func (client *Client) Exit() {
	client.Server.Leave(*client.Username, *client.Room)
}
