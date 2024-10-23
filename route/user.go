package route

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"gs.lee.was/model"
	"gs.lee.was/storage/mysql"

	"github.com/labstack/echo/v4"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("gs.lee.was")

// SetUserRouters Controller 역활
func SetUserRouters(e *echo.Echo) {
	// 유저 목록 조회
	e.GET("/user", userList)
	// 유저 페이징 조회
	e.POST("/userPage", userPage)
	// 유저 PK id 로 조회
	e.GET("/user/:id", userGetID)
	// 유저 loginID 로 조회
	e.GET("/user/id/:loginID", userGetLoginID)
	// 유저 추가
	e.PUT("/user", userPut)
	// 유저 Pk id 로 삭제
	e.DELETE("/user/:id", userDelete)
}

// @Summary 유저를 가져온다.
// @Description 유저 목록을 가져온다.
// @Accept json
// @Produce json
// @Success 200 json []model.User
// @Router /user [get]
func userList(c echo.Context) error {
	log.Debug("called userList")

	var id int64
	var name, loginID string
	var users []model.User
	rows, err := mysql.Select("SELECT id, name, login_id FROM users")
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgNotExists, "@", "유저", 1))
	}
	for rows.Next() {
		err := rows.Scan(&id, &name, &loginID)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgNotExists, "@", "유저", 1))
		}
		log.Debug("Select Users ", id, "_", name, "_", loginID)
		// 배열에 저장
		users = append(users, model.User{ID: id, Name: name, LoginID: loginID})
	}

	return c.JSON(http.StatusOK, users)
}

func userGetID(c echo.Context) error {
	log.Debug("called userGet")

	searchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrMsgWorngParam)
	}
	var id int64
	var name, loginID string
	row := mysql.SelectOne("SELECT id, name, login_id FROM users WHERE id = ?", searchID)
	err = row.Scan(&id, &name, &loginID)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, strings.Replace(model.ErrMsgNotExists, "@", "유저", 1))
	}
	log.Info("SelectOne ", id, "_", name, "_", loginID)

	return c.JSON(http.StatusOK, model.User{ID: id, Name: name, LoginID: loginID})
}

func userGetLoginID(c echo.Context) error {
	log.Debug("called userGet")

	searchLoginID, err := strconv.Atoi(c.Param("loginID"))
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrMsgWorngParam)
	}
	var id int64
	var name, loginID string
	row := mysql.SelectOne("SELECT id, name, login_id FROM users WHERE login_id = ?", searchLoginID)
	err = row.Scan(&id, &name, &loginID)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, strings.Replace(model.ErrMsgNotExists, "@", "유저", 1))
	}
	log.Info("SelectOne ", id, "_", name, "_", loginID)

	return c.JSON(http.StatusOK, model.User{ID: id, Name: name, LoginID: loginID})
}

func userPut(c echo.Context) error {
	log.Debug("called userPut")
	u := &model.User{}
	if err := c.Bind(u); err != nil {
		return err
	}
	res, err := mysql.Execute("INSERT INTO users (name, login_id) VALUES(?,?)", u.Name, u.LoginID)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotAcceptable, strings.Replace(model.ErrMsgDuplicate, "@", "유저", 1))
	}
	// 적용된 res 개수를 가져와서 0건이면 에러
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgProcFail, "@", "유저등록", 1))
	}
	if cnt == 0 {
		log.Error("Insert RowsAffected is 0")
	}
	lastestID, err := res.LastInsertId()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgProcFail, "@", "유저등록", 1))
	}
	if lastestID == 0 {
		log.Error("Insert pk is cannot be 0")
		return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgProcFail, "@", "유저등록", 1))
	}
	u.ID = lastestID
	log.Info("insertId ", lastestID)

	return c.JSON(http.StatusOK, u)
}

func userDelete(c echo.Context) error {
	log.Debug("called userDelete")
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := mysql.Execute("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgProcFail, "@", "삭제", 1))
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgProcFail, "@", "삭제", 1))
	}
	if cnt == 0 {
		log.Error("Delete result count is 0")
		return echo.NewHTTPError(http.StatusBadRequest, strings.Replace(model.ErrMsgProcFail, "@", "삭제", 1))
	}

	return c.NoContent(http.StatusNoContent)
}

func userPage(c echo.Context) error {
	var id int64
	var name, loginID string
	// 파라미터가 넘어오지 않았을 때의 기본값
	p := &model.Page{
		PageType: "num", // id, num
		Page:     1,
		Num:      20,
		ID:       0,
	}
	// 유저가 전달한 POST 값 바인드
	if err := c.Bind(p); err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrMsgWorngParam)
	}

	var rows *sql.Rows
	var err error
	var users []model.User
	if p.PageType == "id" {
		// id 가 0인 경우는 없으므로 최초 배열부터 가져옴
		// 반드시 정렬 순서가 정해져 있어야 함
		rows, err = mysql.Select("SELECT id, name, login_id FROM users WHERE id > ? ORDER BY id ASC LIMIT ?", p.ID, p.Num)
	} else {
		rows, err = mysql.Select("SELECT id, name, login_id FROM users LIMIT ?, ?", (p.Page-1)*p.Num, p.Num)
	}
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrMsgCannotFound)
	}
	var pageCnt = 0
	for rows.Next() {
		pageCnt++
		err := rows.Scan(&id, &name, &loginID)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, model.ErrMsgCannotFound)
		}
		// 배열에 저장
		users = append(users, model.User{ID: id, Name: name, LoginID: loginID})

		if p.PageType == "id" {
			log.Debug("SELECT users Pagination By PrimaryKey Value ", p.ID, " user : ", id, "_", name, "_", loginID)
		} else {
			log.Debug("SELECT users Pagination By p.Page ", p.Page, " user : ", id, "_", name, "_", loginID)
		}
	}
	defer rows.Close()

	// 해당 항목의 전체 개수 및 페이지 수 가져오기
	var pageInfo model.PageInfo
	row := mysql.SelectOne("SELECT COUNT(*) totalCounts, CEIL(COUNT(*) / ?) totalPages FROM users", p.Num)
	err = row.Scan(&pageInfo.TotalCounts, &pageInfo.TotalPages)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, model.ErrMsgCannotFound)
	}
	// 페이지 정보와 결과배열을 객체에 담아서 JSON으로 리턴
	var result model.PageResult
	result.PageInfo = pageInfo
	result.Result = users

	return c.JSON(http.StatusOK, result)
}
