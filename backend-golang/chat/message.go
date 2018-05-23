package chat

import (
	"app/common"
	model "app/models"
	"database/sql"
	"strings"

	"github.com/labstack/echo"
)

// Post new message
func PostMessageHandler(c echo.Context) (err error) {
	m := new(model.Message)
	if err = c.Bind(m); err != nil {
		return c.JSON(400, err.Error())
	}
	if m.MessageContent == nil || strings.TrimSpace(*m.MessageContent) == "" {
		return c.JSON(400, model.EmptyMessageErr.Error())
	}

	if m.UserID == 0 && m.Username != nil {
		if uId := checkUserExist(*m.Username); uId != nil {
			m.UserID = *uId
		}
	}
	if m.Username == nil {
		uname := "guest_" + RandomString(8)
		m.Username = &uname
	}
	if m.UserID == 0 {
		m.UserID = createNewUser(*m.Username)
	}
	ins, err := common.DB.Prepare(`insert into message(user_id,message_content,created_at)
		values(?,?,CURRENT_TIMESTAMP())`)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	_, err = ins.Exec(m.UserID, m.MessageContent)
	if err != nil {
		return c.JSON(400, err.Error())
	}
	return c.JSON(200, m)
}

func GetMessageListHandler(c echo.Context) (err error) {
	rows, err := common.DB.Query(`select m.message_id,m.user_id,m.message_content,m.created_at,
		u.username from message m 
		join user_ u on u.user_id =  m.user_id
		order by created_at asc limit 100`)
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(400, err.Error())
	}
	defer rows.Close()
	messages := []model.Message{}
	for rows.Next() {
		m := model.Message{}
		rows.Scan(&m.MessageID, &m.UserID, &m.MessageContent,
			&m.CreatedAt, &m.Username)
		messages = append(messages, m)
	}
	return c.JSON(200, messages)
}
