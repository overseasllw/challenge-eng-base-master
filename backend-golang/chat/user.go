package chat

import (
	"app/common"
	"database/sql"
	"net/http"

	model "app/models"

	"github.com/labstack/echo"
)

func PostUserHandler(c echo.Context) (err error) {
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if u.Username == nil {
		var uname = "guest_" + RandomString(8)
		u.Username = &uname
	}
	uid := checkUserExist(*u.Username)
	if uid == nil {
		ins, err := common.DB.Prepare(`insert into user_(username,created_at,last_login)
	values(?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP())`)
		if err != nil {
			return c.JSON(400, err.Error())
		}
		res, err := ins.Exec(u.Username)
		if err != nil {
			return c.JSON(400, err.Error())
		}
		u.UserID, _ = res.LastInsertId()
	} else {
		u.UserID = *uid
	}
	return c.JSON(http.StatusOK, u)
}

func checkUserExist(username string) (userId *int64) {

	err := common.DB.QueryRow(`select user_id from user_ where username=?`,
		username).Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	return
}

func createNewUser(username string) (userId int64) {
	ins, err := common.DB.Prepare(`insert into user_(username,created_at,last_login)
	values(?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP())`)
	if err != nil {
		return
	}
	res, err := ins.Exec(username)
	if err != nil {
		return
	}
	userId, _ = res.LastInsertId()
	return userId
}
