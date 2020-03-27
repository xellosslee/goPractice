package route

import (
	"net/http"
	"strconv"

	"cndf.order.was/model"
	"github.com/labstack/echo"
)

func UserList(c echo.Context) error {
	return c.JSON(http.StatusOK, model.Users)
}

func UserGet(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	return c.JSON(http.StatusOK, model.Users[id])
}

func UserPut(c echo.Context) error {
	u := &model.User{
		ID: model.UserSeq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}

	// id값이 Client 로부터 전송된 경우 수정이므로 update처리
	if u.ID != model.UserSeq {
		model.Users[u.ID] = u
		model.UserSeq++
	} else {
		// id값이 넘어오지 않은 경우 초기화 된 UserSeq와 같은 값이므로 유저 추가
		model.Users[u.ID] = u
		model.UserSeq++
	}

	return c.JSON(http.StatusOK, u)
}

func UserDelete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	delete(model.Users, id)
	return c.NoContent(http.StatusNoContent)
}
