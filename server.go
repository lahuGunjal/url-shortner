package main

import (
	"github.com/labstack/echo"
	"github.com/lahuGunjal/url-shortner/api"
)

func main() {
	e := echo.New()
	api.Init(e)
	e.Logger.Fatal(e.Start(":1323"))
}
