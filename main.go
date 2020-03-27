// 패키지명이 반드시 main인 파일만 실행파일로 생성 및 run 이 가능하다
package main

import (
	"net/http"
	"strconv"

	"cndf.order.was/route"
	"github.com/labstack/echo"
)

// PORT 기본 포트는 81
const PORT = 81

// main 프로그램의 시작점
func main() {

	e := echo.New()

	e.Logger.Info("Attempting to start HTTP Server.")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello, World!")
	})
	e.GET("/user", route.UserList)
	e.GET("/user/:id", route.UserGet)
	e.PUT("/user", route.UserPut)
	e.DELETE("/user/:id", route.UserDelete)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
