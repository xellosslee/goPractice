package configs

import (
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	logging "github.com/op/go-logging"
)

type StaticLog struct {
	Log *logging.Logger // 로그 설정할때 사용하는 값
	//Format logging.Formatter // 로그 설정하는 for설정
	FpLog   *os.File // access
	InfoLog *os.File // 로그 출력
	Defung  *os.File // 디버그
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}
