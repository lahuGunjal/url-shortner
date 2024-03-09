package storage

import (
	"sync"
	"testing"

	"github.com/lahuGunjal/url-shortner/api/model"
	"github.com/stretchr/testify/assert"
)

func TestStoreURLToMap(t *testing.T) {
	urlMap := &URLInfoMap{
		Data: make(map[string]*model.URLDetails),
		Mux:  &sync.Mutex{},
	}
	urlInfo := &model.URLDetails{
		OriginalURL:  "https://www.youtube.com/watch?v=QwVWcvmcfuk",
		DomainName:   "http:localhost:1323",
		HashValue:    "testJVHSvvhs",
		ShortenedURL: "http:localhost:1323/testJVHSvvhs",
	}
	t.Run("If valid values are passed it should store it in the map", func(t *testing.T) {
		urlMap.StoreURLToMap(urlInfo)
		urlDetails := urlMap.GetURLFromMap("testJVHSvvhs")
		assert.Equal(t, urlDetails.OriginalURL, urlInfo.OriginalURL)
		assert.Equal(t, urlDetails.DomainName, urlInfo.DomainName)
		assert.Equal(t, urlDetails.HashValue, urlInfo.HashValue)
		assert.Equal(t, urlDetails.ShortenedURL, urlInfo.ShortenedURL)
	})

}

func TestGetURLFromMap(t *testing.T) {
	urlMap := &URLInfoMap{
		Data: make(map[string]*model.URLDetails),
		Mux:  &sync.Mutex{},
	}
	urlInfo := &model.URLDetails{
		OriginalURL:  "https://www.youtube.com/watch?v=QwVWcvmcfuk",
		DomainName:   "http:localhost:1323",
		HashValue:    "testJVHSvvhs",
		ShortenedURL: "http:localhost:1323/testJVHSvvhs",
	}
	urlMap.Data["testJVHSvvhs"] = urlInfo
	t.Run("If valid key is passed it should load values from map", func(t *testing.T) {
		urlDetails := urlMap.GetURLFromMap("testJVHSvvhs")
		assert.Equal(t, urlInfo.OriginalURL, urlDetails.OriginalURL)
		assert.Equal(t, urlInfo.DomainName, urlDetails.DomainName)
		assert.Equal(t, urlInfo.HashValue, urlDetails.HashValue)
		assert.Equal(t, urlInfo.ShortenedURL, urlDetails.ShortenedURL)
	})
	t.Run("If InValid key is passed it should return struct with blank values", func(t *testing.T) {
		urlDetails := urlMap.GetURLFromMap("jfjfj")
		assert.NotEqual(t, urlInfo.OriginalURL, urlDetails.OriginalURL)
		assert.NotEqual(t, urlInfo.DomainName, urlDetails.DomainName)
		assert.NotEqual(t, urlInfo.HashValue, urlDetails.HashValue)
		assert.NotEqual(t, urlInfo.ShortenedURL, urlDetails.ShortenedURL)
	})

}

func TestGetOriginalURL(t *testing.T) {
	urlMap := &URLInfoMap{
		Data: make(map[string]*model.URLDetails),
		Mux:  &sync.Mutex{},
	}
	urlInfo := &model.URLDetails{
		OriginalURL:  "https://www.youtube.com/watch?v=QwVWcvmcfuk",
		DomainName:   "http:localhost:1323",
		HashValue:    "testJVHSvvhs",
		ShortenedURL: "http:localhost:1323/testJVHSvvhs",
	}
	urlMap.Data["testJVHSvvhs"] = urlInfo
	t.Run("If valid url is passed it should return originalurl", func(t *testing.T) {
		shortenedURL := urlMap.GetOriginalURL("https://www.youtube.com/watch?v=QwVWcvmcfuk")
		assert.Equal(t, urlInfo.ShortenedURL, shortenedURL)
	})
	t.Run("If invalid url is passed it should return blank string", func(t *testing.T) {
		originalURL := urlMap.GetOriginalURL("https://www.yojhskbsjkutube.com/watch?v=QwVWcvmcfuk")
		assert.Equal(t, "", originalURL)
	})
}
