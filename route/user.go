package route

import (
	"net/http"
	"strconv"

	"cndf.order.was/model"
	"cndf.order.was/storage"
	"github.com/labstack/echo"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("cndf.order.was")

// SetUserRouters Controller 역활
func SetUserRouters(e *echo.Echo) {
	e.GET("/user", userList)
	e.GET("/user/:id", userGet)
	e.PUT("/user", userPut)
	e.DELETE("/user/:id", userDelete)
}

// Service 함수 역활
func userList(c echo.Context) error {
	log.Debug("called userList")
	db, err := storage.ConnectDB()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "DB 연결 실패")
	}
	var id int64
	var name, loginID string
	var users []model.UserInfo
	rows, err := storage.Select(db, "SELECT id, name, login_id FROM users")
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "사용자 조회 실패")
	}
	for rows.Next() {
		err := rows.Scan(&id, &name, &loginID)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, "사용자 조회 실패")
		}
		log.Debug("Select Users ", id, "_", name, "_", loginID)
		// 배열에 저장
		users = append(users, model.UserInfo{ID: id, Name: name, LoginID: loginID})
	}

	return c.JSON(http.StatusOK, users)
}

// Service 함수 역활
func userGet(c echo.Context) error {
	log.Debug("called userGet")

	searchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 파라미터가 전달되었습니다")
	}
	db, err := storage.ConnectDB()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "DB 연결 실패")
	}
	var id int64
	var name, loginID string
	row := storage.SelectOne(db, "SELECT id, name, login_id FROM users WHERE id = ?", searchID)
	err = row.Scan(&id, &name, &loginID)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "해당 유저가 없습니다.")
	}
	log.Info("SelectOne ", id, "_", name, "_", loginID)

	return c.JSON(http.StatusOK, model.UserInfo{ID: id, Name: name, LoginID: loginID})
}

// Service 함수 역활
func userPut(c echo.Context) error {
	log.Debug("called userPut")
	u := &model.UserInfo{}
	if err := c.Bind(u); err != nil {
		return err
	}
	db, err := storage.ConnectDB()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "DB 연결 실패")
	}
	res, err := storage.Execute(db, "INSERT INTO users (name, login_id) VALUES(?,?)", u.Name, u.LoginID)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, "Duplicate Users")
	}
	// 적용된 res 개수를 가져와서 0건이면 에러
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "사용자 추가 실패")
	}
	if cnt == 0 {
		log.Error("Insert RowsAffected is 0")
	}
	lastestID, err := res.LastInsertId()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "사용자 추가 시 오류가 발생하였습니다.")
	}
	if lastestID == 0 {
		log.Error("Insert pk is cannot be 0")
		return echo.NewHTTPError(http.StatusBadRequest, "사용자 추가 시 오류가 발생하였습니다.")
	}
	u.ID = lastestID
	log.Info("insertId ", lastestID)

	return c.JSON(http.StatusOK, u)
}

// Service 함수 역활
func userDelete(c echo.Context) error {
	log.Debug("called userDelete")
	id, _ := strconv.Atoi(c.Param("id"))

	// DB 버전
	db, err := storage.ConnectDB()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "DB 연결 실패")
	}
	res, err := storage.Execute(db, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "사용자 삭제 실패")
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "삭제 건수 가져오기 실패")
	}
	if cnt == 0 {
		log.Error("Delete result count is 0")
		return echo.NewHTTPError(http.StatusBadRequest, "삭제할 사용자가 없습니다.")
	}
	return c.NoContent(http.StatusNoContent)
}
