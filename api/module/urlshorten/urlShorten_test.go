package urlshorten

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/lahuGunjal/url-shortner/api/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateURLRoute(t *testing.T) {
	e := echo.New()
	InitialiseMap()

	t.Run("Success", func(t *testing.T) {
		userJSON := `{
			"url":"https://www.youtube.com/watch?v=QwVWcvmcfuk",
			"domainName":"http://localhost:1323"
		}`
		req, _ := http.NewRequest(echo.POST, "/url/create", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		if assert.NoError(t, CreateURLRoute(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var responseURL string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseURL))
			assert.NotEmpty(t, responseURL)
		}
	})
	t.Run("Blank Url", func(t *testing.T) {
		userJSON := `{
			"url":"",
			"domainName":"http://localhost:1323"
		}`
		req, _ := http.NewRequest(echo.POST, "/url/create", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		CreateURLRoute(c)
		var responseURL string
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseURL))
		assert.Equal(t, "MISSING_URL", responseURL)
		assert.Equal(t, http.StatusExpectationFailed, rec.Code)
	})
	t.Run("Blank DomainName", func(t *testing.T) {
		userJSON := `{
			"url":"https://www.youtube.com/watch?v=QwVWcvmcfuk",
			"domainName":""
		}`
		req, _ := http.NewRequest(echo.POST, "/url/create", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		CreateURLRoute(c)
		var responseURL string
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseURL))
		// Assert the expected original URL
		assert.Equal(t, "MISSING_DOMAINNAME", responseURL)
		assert.Equal(t, http.StatusExpectationFailed, rec.Code)
	})
	t.Run("Parameter Bind Error", func(t *testing.T) {
		userJSON := `""{
			"org":"https://www.youtube.com/watch?v=QwVWcvmcfuk",
			"test":""
		}""`
		req, _ := http.NewRequest(echo.POST, "/url/create", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		CreateURLRoute(c)
		var responseURL string
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseURL))
		// Assert the expected original URL
		assert.Equal(t, "PARAMETER_BINDING_ERROR", responseURL)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetURLRoute(t *testing.T) {

	e := echo.New()
	InitialiseMap()
	req, _ := http.NewRequest(http.MethodGet, "/url/get/NTg3ODlhZDcxMTcx", nil)
	rec := httptest.NewRecorder()
	urlDetails := model.URLDetails{}
	urlDetails.DomainName = "http://localhost:1323/"
	urlDetails.HashValue = "NTg3ODlhZDcxMTcx"
	urlDetails.OriginalURL = "https://www.youtube.com/watch?v=QwVWcvmcfuk"
	urlDetails.ShortenedURL = "http://localhost:1323/NTg3ODlhZDcxMTcx"
	StoreURLToMap(&urlDetails)
	c := e.NewContext(req, rec)
	t.Run("Success", func(t *testing.T) {
		c.SetParamNames("url")
		c.SetParamValues("NTg3ODlhZDcxMTcx")
		if assert.NoError(t, GetURLRoute(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// Parse the response JSON
			var responseURL string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseURL))
			// Assert the expected original URL
			assert.Equal(t, urlDetails.OriginalURL, responseURL)
		}
	})
	t.Run("missing url", func(t *testing.T) {
		if assert.NoError(t, GetURLRoute(c)) {
			assert.Equal(t, http.StatusExpectationFailed, rec.Code)
			// Parse the response JSON
			var response string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			// Assert the expected original URL
			assert.Equal(t, "MISSING_URL", response)
		}
	})

	t.Run("URLNotFound", func(t *testing.T) {
		c.SetParamNames("url")
		c.SetParamValues("gkggjlvkjvjnknk")
		if assert.NoError(t, GetURLRoute(c)) {
			assert.Equal(t, http.StatusExpectationFailed, rec.Code)
			// Parse the response JSON
			var response string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			// Assert the expected original URL
			assert.Equal(t, "URL_NOT_FOUND", response)
		}
	})

}

func TestGetDomainStatsRoute(t *testing.T) {
	e := echo.New()
	InitialiseMap()
	req, _ := http.NewRequest(http.MethodGet, "/domainstats", nil)
	rec := httptest.NewRecorder()
	stats.Data["www.youtube.com"] = 3
	stats.Data["www.wikipedia.com"] = 4
	stats.Data["www.google.com"] = 6
	stats.Data["www.goplayground.com"] = 8
	topStats := []KeyValue{
		KeyValue{
			Key:   "www.goplayground.com",
			Value: 8,
		},
		KeyValue{
			Key:   "www.google.com",
			Value: 6,
		},
		KeyValue{
			Key:   "www.wikipedia.com",
			Value: 4,
		},
	}

	c := e.NewContext(req, rec)
	t.Run("Success", func(t *testing.T) {
		if assert.NoError(t, GetDomainStatsRoute(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// Parse the response JSON
			var response []KeyValue
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			// // Assert the expected original URL
			assert.Equal(t, topStats, response)
		}
	})
	t.Run("stats less than 3", func(t *testing.T) {
		topStats = topStats[:len(topStats)-1]
		delete(stats.Data, "www.goplayground.com")
		delete(stats.Data, "www.wikipedia.com")
		if assert.NoError(t, GetDomainStatsRoute(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// Parse the response JSON
			var response []KeyValue
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			// Assert the expected original URL
			assert.Equal(t, topStats, response)
		}
	})
	t.Run("no records", func(t *testing.T) {
		topStats = topStats[:0]
		delete(stats.Data, "www.goplayground.com")
		delete(stats.Data, "www.wikipedia.com")
		delete(stats.Data, "www.youtube.com")
		delete(stats.Data, "www.google.com")
		if assert.NoError(t, GetDomainStatsRoute(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// Parse the response JSON
			var response string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			// Assert the expected original URL
			assert.Equal(t, "No data found", response)
		}
	})
}
