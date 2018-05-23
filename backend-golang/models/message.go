package models

import (
	"app/common"
	"database/sql"
	"time"
)

type (
	Message struct {
		MessageID      int64     `json:"message_id"`
		MessageType    string    `json:"message_type"`
		CreatedAt      time.Time `json:"created_at,omitempty"`
		MessageContent *string   `json:"message"`
		User
	}
)

func GetAllMessageList() (messages []*Message, err error) {
	rows, err := common.DB.Query(`select m.message_id,m.user_id,m.message_content,m.created_at,
		u.username from message m 
		join user_ u on u.user_id =  m.user_id
		order by created_at asc limit 100`)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer rows.Close()
	messages = []*Message{}
	for rows.Next() {
		m := Message{}
		rows.Scan(&m.MessageID, &m.UserID, &m.MessageContent,
			&m.CreatedAt, &m.Username)
		messages = append(messages, &m)
	}
	return messages, err
}

func CreateNewMessage(message *Message) (err error) {
	ins, err := common.DB.Prepare(`insert into message(user_id,message_content,created_at)
		values(?,?,CURRENT_TIMESTAMP())`)
	if err != nil {
		return err
	}
	_, err = ins.Exec(message.UserID, message.MessageContent)
	if err != nil {
		return err
	}
	return
}
