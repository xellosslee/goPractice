package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gs.lee.was/configs"
	"gs.lee.was/util"
)

// SetUserRouters Controller 역활
func SetWaitRoute(e *echo.Echo) {
	// 유저 목록 조회
	e.GET("/wait", getWaitNo)
}

func getWaitNo(c echo.Context) error {
	util.RdbSet("a", "bb")
	val := util.RdbGet("a")
	configs.Log.Info("key", val)
	return c.JSON(http.StatusOK, val)
}
