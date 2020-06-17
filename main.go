// 패키지명이 반드시 main인 파일만 실행파일로 생성 및 run 이 가능하다
package main

import (
	"net/http"

	"cndf.order.was/route"
	"github.com/labstack/echo"
)

// main 프로그램의 시작점
func main() {

	e := echo.New()

	e.Logger.Info("Attempting to start HTTP Server.")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello, World!")
	})

	route.SetUserRouters(e)

	e.Logger.Fatal(e.Start(":80"))
}
