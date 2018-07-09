package room

import (
	"app/common"
	model "app/models"
	"database/sql"
	"net/http"

	"github.com/labstack/echo"
)

func GetRoomListHandler(c echo.Context) (err error) {
	rows, err := common.DB.Query(`select room_id,room_uuid,room_name from room`)
	if err != nil {
		return err
	}
	defer rows.Close()
	rs := []model.Room{}
	for rows.Next() {
		r := model.Room{}
		rows.Scan(&r.Id, &r.UUID, &r.Name)
		rs = append(rs, r)
	}
	//	j, _ := json.Marshal(rs)
	//	log.Print(string(j))
	return c.JSON(http.StatusCreated, rs)
}

func CreateNewRoomHandler(c echo.Context) (err error) {
	room := new(model.Room)
	if err = c.Bind(room); err != nil {
		return err
	}
	CreateNewRoom(room)
	return c.JSON(http.StatusCreated, room)
}

func CreateNewRoom(room *model.Room) (err error) {
	rid := CheckRoomExist(*room.Name)

	if rid == nil {
		ins, err := common.DB.Prepare(`insert into room(room_uuid,room_name)
			values(?,?)`)
		if err != nil {
			return err
		}
		res, err := ins.Exec(room.UUID, room.Name)
		if err != nil {
			return err
		}

		id, _ := res.LastInsertId()
		room.Id = id
	}
	return err
}

func CheckRoomExist(room string) (roomId *int64) {

	err := common.DB.QueryRow(`select room_name from room where room_name=?`,
		room).Scan(&roomId)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	return
}
