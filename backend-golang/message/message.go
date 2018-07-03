package chat

import (
	"app/common"
	model "app/models"
	"database/sql"
	"log"

	"github.com/labstack/echo"
)

func GetMessageListHandler(c echo.Context) (err error) {
	rows, err := common.DB.Query(`select m.message_uuid, m.message_id,m.user_id,m.message_content,m.created_at,
		u.username,message_type,room_id from message m 
		join user_ u on u.user_id =  m.user_id and m.user_id !=0
		where m.message_content !="" 
		order by created_at asc limit 100`)
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(400, err.Error())
	}
	defer rows.Close()
	messages := []model.Message{}
	for rows.Next() {
		m := model.Message{}
		rows.Scan(&m.UUID, &m.MessageID, &m.UserID, &m.MessageContent,
			&m.CreatedAt, &m.Username, &m.MessageType, &m.RoomId)
		messages = append(messages, m)
	}
	return c.JSON(200, messages)
}

func GetAllMessageList() (messages []*model.Message, err error) {
	rows, err := common.DB.Query(`select m.uuid,m.message_id,m.user_id,m.message_content,m.created_at,
		u.username,message_type,room_id from message m 
		join user_ u on u.user_id =  m.user_id 
		order by created_at asc limit 100`)
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer rows.Close()
	messages = []*model.Message{}
	for rows.Next() {
		m := model.Message{}
		rows.Scan(&m.UUID, &m.MessageID, &m.UserID, &m.MessageContent,
			&m.CreatedAt, &m.Username, &m.MessageType, &m.RoomId)
		messages = append(messages, &m)
	}
	return messages, err
}

func CreateNewMessage(message *model.Message) (err error) {
	ins, err := common.DB.Prepare(`insert into message(message_uuid,user_id,message_type,room_id,message_content,created_at)
		values(?,?,?,?,?CURRENT_TIMESTAMP())`)
	if err != nil {
		return err
	}
	_, err = ins.Exec(message.UUID, message.UserID, message.MessageType, message.RoomId, message.MessageContent)
	if err != nil {
		log.Print(err)
		return err
	}
	return
}
