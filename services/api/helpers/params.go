package helpers

import (
	"log"

	"github.com/labstack/echo/v4"
)

func StringParams(context echo.Context) (params map[string]string) {
	params = make(map[string]string)
	for k, v := range context.QueryParams() {
		params[k] = v[0]
	}
	values, _ := context.FormParams()
	for k, v := range values {
		params[k] = v[0]
	}
	log.Println("params: ", params)
	return
}
