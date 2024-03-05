package api

import (
	"github.com/labstack/echo"
	"github.com/lahuGunjal/url-shortner/api/middleware"
	"github.com/lahuGunjal/url-shortner/api/module/urlshorten"
)

func Init(e *echo.Echo) {
	r := e.Group("/r")
	o := e.Group("/o")
	middleware.Init(e, r, o)
	urlshorten.Init(o, r)

}
