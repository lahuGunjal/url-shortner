package urlshorten

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/lahuGunjal/url-shortner/api/model"
)

func Init(e *echo.Echo) {
	e.POST("/url/create", CreateURLRoute)
	e.GET("/url/get/:url", GetURLRoute)
	e.GET("/:url", RedirectRoute)
	e.GET("/domainstats", GetDomainStatsRoute)
	InitialiseMap()
}

func CreateURLRoute(c echo.Context) error {
	log.Println("Info: IN CreateURLRoute route")
	reqURLDetails := model.RequestURLData{}
	bindErr := c.Bind(&reqURLDetails)
	if bindErr != nil {
		log.Println("PARAMETER_BINDING_ERROR", bindErr)
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusBadRequest, "PARAMETER_BINDING_ERROR")
	}
	if reqURLDetails.URL == "" {
		log.Println("URL_SHOULD_NOT_BE_BLANK")
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusExpectationFailed, "MISSING_URL")
	}
	if reqURLDetails.DomainName == "" {
		log.Println("DomainName_SHOULD_NOT_BE_BLANK")
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusExpectationFailed, "MISSING_DOMAINNAME")
	}
	url, err := createURLService(reqURLDetails)
	if err != nil {
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusOK, err)
	}
	log.Println("Info: OUT CreateURLRoute route")
	return c.JSON(http.StatusOK, url)
}

func GetURLRoute(c echo.Context) error {
	log.Println("Info: IN GetURLRoute route")
	shortURL := c.Param("url")
	urlDetails := GetURLFromMap(shortURL)
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
		return c.JSON(http.StatusNotFound, "URL not found")
	}
	log.Println("Info: OUT Redirect route")
	return c.JSON(http.StatusOK, nil)
}

func GetDomainStatsRoute(c echo.Context) error {
	log.Println("Info: IN GetURLRoute route")
	keyValueSlice := GetStatsService()
	if len(keyValueSlice) == 0 {
		return c.JSON(http.StatusOK, "No data found")
	}
	log.Println("Info: OUT GetURLRoute route")
	return c.JSON(http.StatusOK, keyValueSlice)
}
