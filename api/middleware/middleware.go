package middleware

import (
	"github.com/lahuGunjal/url-shortner/api/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init(e *echo.Echo, r *echo.Group, o *echo.Group) {
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &model.JwtCustomClaims{},
		SigningKey: []byte(model.JwtKey),
	}))
}
