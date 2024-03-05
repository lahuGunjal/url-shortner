package urlshorten

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/lahuGunjal/url-shortner/api/model"
)

// URLDetailsMap represents the state of the URL shortener.
type URLDetailsMap struct {
	Data map[string]*model.URLDetails
	Mux  *sync.Mutex
}

var urlMap *URLDetailsMap

// InitialiseMap initializes the URLDetailsMap for the URL shortener.
func InitialiseMap() {
	urlMap = &URLDetailsMap{
		Data: make(map[string]*model.URLDetails),
		Mux:  &sync.Mutex{},
	}
}

// StoreURLToMap stores the URL details in the URLDetailsMap.
func StoreURLToMap(urlDetails *model.URLDetails) {
	urlMap.Mux.Lock()
	defer urlMap.Mux.Unlock()
	urlMap.Data[urlDetails.HashValue] = urlDetails
}

// GetURLFromMap retrieves URL details from the URLDetailsMap based on the provided URL.
func GetURLFromMap(url string) *model.URLDetails {
	urlMap.Mux.Lock()
	defer urlMap.Mux.Unlock()
	if val, ok := urlMap.Data[url]; ok {
		return val
	}
	return &model.URLDetails{}
}

// createURLService generates a shortened URL based on the provided URL details.
func createURLService(reqURLDetails model.RequestURLData) string {
	urlDetails := model.URLDetails{}
	urlDetails.HashValue = getHashValue(reqURLDetails.URL)
	urlData := GetURLFromMap(urlDetails.HashValue)
	if urlData.HashValue != "" {
		return fmt.Sprintf("%s/%s", urlData.DomainName, urlData.HashValue)
	}

	urlDetails.DomainName = reqURLDetails.DomainName
	urlDetails.OriginalURL = reqURLDetails.URL
	urlDetails.ShortenedURL = fmt.Sprintf("%s/%s", urlDetails.OriginalURL, urlDetails.HashValue)

	StoreURLToMap(&urlDetails)
	return fmt.Sprintf("%s/%s", urlDetails.DomainName, urlDetails.HashValue)
}

// getHashValue generates an MD5 hash value for the given URL.
func getHashValue(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hashValue := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashValue)
	return hashString
}

// validateURL checks the validity of the input URL.
func validateURL(shortURL string) error {
	splitPath := strings.Split(shortURL, "/")
	if len(splitPath) != 2 {
		log.Println("INVALID_URL")
		return errors.New("INVALID_URL")
	}
	return nil
}
