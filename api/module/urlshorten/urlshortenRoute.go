package urlshorten

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/lahuGunjal/url-shortner/api/model"
)

func Init(o *echo.Group, r *echo.Group) {
	o.POST("/url/create", CreateURLRoute)
	o.GET("/url/get/:url", GetURLRoute)
	o.GET("/:url", RedirectRoute)
	InitialiseMap()
}

func CreateURLRoute(c echo.Context) error {
	log.Println("Info: IN CreateURLRoute route")
	reqURLDetails := model.RequestURLData{}
	bindErr := c.Bind(&reqURLDetails)
	if bindErr != nil {
		log.Println("PARAMETER_BINDING_ERROR", bindErr)
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusExpectationFailed, errors.New("PARAMETER_BINDING_ERROR"))
	}
	if reqURLDetails.URL == "" {
		log.Println("URL_SHOULD_NOT_BE_BLANK")
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusExpectationFailed, errors.New("MISSING_URL"))
	}
	if reqURLDetails.DomainName == "" {
		log.Println("DomainName_SHOULD_NOT_BE_BLANK")
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusExpectationFailed, errors.New("MISSING_DOMAINNAME"))
	}
	url := createURLService(reqURLDetails)
	log.Println("Info: OUT CreateURLRoute route")
	return c.JSON(http.StatusOK, url)
}

func GetURLRoute(c echo.Context) error {
	log.Println("Info: IN GetURLRoute route")
	shortURL := c.Param("url")
	err := validateURL(shortURL)
	if err != nil {
		log.Println("Info: OUT GetURLRoute route")
		return err
	}
	hashValue := strings.Split(shortURL, "/")[1]
	urlDetails := GetURLFromMap(hashValue)
	log.Println("Info: OUT GetURLRoute route")
	return c.JSON(http.StatusOK, urlDetails.OriginalURL)
}

func RedirectRoute(c echo.Context) error {
	log.Println("Info: IN Redirect route")
	shortURL := c.Param("url")
	urlDetails := GetURLFromMap(shortURL)
	if urlDetails.OriginalURL != "" {
		http.Redirect(c.Response().Writer, c.Request(), urlDetails.OriginalURL, http.StatusMovedPermanently)
	} else {
		log.Println("Info: OUT Redirect route")
		return c.JSON(http.StatusNotFound, errors.New("URL not found"))
	}
	log.Println("Info: OUT Redirect route")
	return c.JSON(http.StatusOK, nil)
}
