package urlshorten

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/lahuGunjal/url-shortner/api/model"
)

func Init(e *echo.Echo) {
	e.POST("/url/create", CreateURLRoute)
	e.GET("/url/get/:url", GetURLRoute)
	e.GET("/:url", URLRedirectRoute)
	e.GET("/domainstats", GetDomainStatsRoute)
	InitialiseStorage()
}

func CreateURLRoute(c echo.Context) error {
	log.Println("Info: IN CreateURLRoute route")
	reqURLDetails := model.RequestURLData{}
	bindErr := c.Bind(&reqURLDetails)
	if bindErr != nil {
		log.Println("PARAMETER_BINDING_ERROR", bindErr)
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusBadRequest, errors.New("PARAMETER_BINDING_ERROR").Error())
	}
	if reqURLDetails.URL == "" {
		log.Println("URL_SHOULD_NOT_BE_BLANK")
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusBadRequest, errors.New("MISSING_URL").Error())
	}
	if reqURLDetails.DomainName == "" {
		log.Println("DomainName_SHOULD_NOT_BE_BLANK")
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusBadRequest, errors.New("MISSING_DOMAINNAME").Error())
	}
	url, err := createURLService(reqURLDetails)
	if err != nil {
		log.Println("Info: OUT CreateURLRoute route")
		return c.JSON(http.StatusExpectationFailed, err)
	}
	log.Println("Info: OUT CreateURLRoute route")
	return c.JSON(http.StatusOK, url)
}

func GetURLRoute(c echo.Context) error {
	log.Println("Info: IN GetURLRoute route")
	shortURL := c.Param("url")
	if shortURL == "" {
		log.Println("Info: OUT GetURLRoute route")
		return c.JSON(http.StatusBadRequest, errors.New("MISSING_URL").Error())
	}
	originalURL, err := GetUrlService(shortURL)
	if err != nil {
		log.Println("Info: OUT GetURLRoute route")
		return c.JSON(http.StatusNotFound, err.Error())
	}
	log.Println("Info: OUT GetURLRoute route")
	return c.JSON(http.StatusOK, originalURL)
}

// RedirectRoute redirect to actual url
func URLRedirectRoute(c echo.Context) error {
	log.Println("Info: IN Redirect route")
	shortURL := c.Param("url")
	originalURL, err := URLRedirectService(shortURL)
	if err != nil {
		log.Println("Info: OUT Redirect route")
		c.JSON(http.StatusNotFound, err.Error())
	}
	http.Redirect(c.Response().Writer, c.Request(), originalURL, http.StatusMovedPermanently)
	log.Println("Info: OUT Redirect route")
	return c.JSON(http.StatusOK, nil)
}

// GetDomainStatsRoute to get the top 3 webDomain and there count occurances
func GetDomainStatsRoute(c echo.Context) error {
	log.Println("Info: IN GetURLRoute route")
	keyValueSlice := GetStatsService()
	if len(keyValueSlice) == 0 {
		log.Println("Info: OUT GetURLRoute route")
		return c.JSON(http.StatusNotFound, errors.New("NO_DATA_FOUND").Error())
	}
	log.Println("Info: OUT GetURLRoute route")
	return c.JSON(http.StatusOK, keyValueSlice)
}
