package api

import (
	"github.com/labstack/echo"
	"github.com/lahuGunjal/url-shortner/api/module/urlshorten"
)

func Init(e *echo.Echo) {
	urlshorten.Init(e)

}
