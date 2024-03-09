package urlshorten

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/lahuGunjal/url-shortner/api/model"
	"github.com/lahuGunjal/url-shortner/api/module/encryption"
	"github.com/lahuGunjal/url-shortner/api/module/storage"
)

var stats *storage.DomainStats
var urlMap *storage.URLInfoMap

func InitialiseStorage() {
	urlMap, stats = storage.NewInMemoryStorage()
}

// createURLService generates a shortened URL based on the provided URL details.
func createURLService(reqURLDetails model.RequestURLData) (string, error) {
	urlDetails := model.URLDetails{}
	//check if the url is allready available in map
	shortURL := CheckIfURLAvailable(reqURLDetails.URL)
	if shortURL != "" {
		return shortURL, nil
	}
	//gnerate code uniq code
	genCode, err := encryption.GenerateCryptoID()
	if err != nil {
		return "", err
	}
	//Prepare model
	urlDetails.HashValue = genCode
	urlDetails.DomainName = reqURLDetails.DomainName
	urlDetails.OriginalURL = reqURLDetails.URL
	urlDetails.ShortenedURL = fmt.Sprintf("%s/%s", urlDetails.DomainName, urlDetails.HashValue)
	//Store details in memory map
	if err := UpdateDomainStats(reqURLDetails.URL); err != nil {
		return "", err
	}
	urlMap.StoreURLToMap(&urlDetails)
	return urlDetails.ShortenedURL, nil
}

func UpdateDomainStats(originalURL string) error {
	domainName := ExtractDomainName(originalURL)
	if domainName == "" {
		log.Println("domain name not found in url")
		return errors.New("DOMAIN_NAME_MISSING_IN_URL")
	}
	if count := stats.LoadStats(domainName); count > 0 {
		stats.UpdateStats(domainName, count+1)
	} else {
		stats.StoreStats(domainName)
	}
	return nil
}

func GetStatsService() []storage.KeyValue {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	if len(stats.Data) == 0 {
		return []storage.KeyValue{}
	}
	// Convert the map to a slice of KeyValue
	keyValueSlice := []storage.KeyValue{}
	for key, value := range stats.Data {
		keyValueSlice = append(keyValueSlice, storage.KeyValue{Key: key, Value: value})
	}

	// Sort the slice by values
	sort.Slice(keyValueSlice, func(i, j int) bool {
		return keyValueSlice[i].Value > keyValueSlice[j].Value
	})
	if len(keyValueSlice) < 3 {
		return keyValueSlice
	}
	return keyValueSlice[:3]
}

func GetUrlService(shortURL string) (string, error) {
	urlDetails := urlMap.GetURLFromMap(shortURL)
	if urlDetails.OriginalURL == "" {
		return "", errors.New("URL_NOT_FOUND")
	}
	return urlDetails.OriginalURL, nil
}

func URLRedirectService(shortURL string) (string, error) {
	urlDetails := urlMap.GetURLFromMap(shortURL)
	if urlDetails.OriginalURL != "" {
		return urlDetails.OriginalURL, nil
	} else {
		log.Println("Info: OUT Redirect route")
		return "", errors.New("URL_not_found")
	}
}
func ExtractDomainName(originalURL string) string {
	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	// Find the domain name using the regular expression
	match := re.FindStringSubmatch(originalURL)

	// Extract the domain name from the matched substring
	if len(match) > 1 {
		domainName := match[1]
		return domainName
	} else {
		return ""
	}
}

// CheckIfURLAvailable url allready availble in memory
func CheckIfURLAvailable(originalURL string) string {
	return urlMap.GetOriginalURL(originalURL)
}
