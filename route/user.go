package route

import (
	"net/http"
	"sort"
	"strconv"

	"cndf.order.was/model"
	"github.com/labstack/echo"
)

func UserList(c echo.Context) error {
	// sort.Slice(model.Users, func(i, j int) bool { return model.Users[i].ID < model.Users[j].ID })
	keys := make([]int, 0)
	for k, _ := range model.Users { // 첫번째 값인 ID 가 k로 넘어와서 해당 값을 배열에 넣음
		keys = append(keys, k)
	}
	sort.Ints(keys) // ID값을 기준으로 데이터 정렬
	var result []*model.User
	for _, k := range keys {
		result = append(result, model.Users[k])
	}

	return c.JSON(http.StatusOK, result)
}

func UserGet(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
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
	id, _ := strconv.Atoi(c.Param("id"))
	delete(model.Users, id)
	return c.NoContent(http.StatusNoContent)
}
