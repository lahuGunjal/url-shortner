package urlshorten

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"sync"

	"github.com/lahuGunjal/url-shortner/api/model"
)

// URLDetailsMap represents the state of the URL shortener.
type URLDetailsMap struct {
	Data map[string]*model.URLDetails
	Mux  *sync.Mutex
}

// WebUrl Stats For most used urls
type Stats struct {
	Data map[string]int
	Mux  *sync.Mutex
}

type KeyValue struct {
	Key   string
	Value int
}

var stats *Stats
var urlMap *URLDetailsMap

// InitialiseMap initializes the URLDetailsMap for the URL shortener.
func InitialiseMap() {
	urlMap = &URLDetailsMap{
		Data: make(map[string]*model.URLDetails),
		Mux:  &sync.Mutex{},
	}
	stats = &Stats{
		Data: make(map[string]int),
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

func StoreStats(webUrl string) {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	stats.Data[webUrl] = 1
}
func UpdateStats(webUrl string, count int) {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	stats.Data[webUrl] = count
}

func LoadStats(webUrl string) int {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	return stats.Data[webUrl]
}

// createURLService generates a shortened URL based on the provided URL details.
func createURLService(reqURLDetails model.RequestURLData) (string, error) {
	urlDetails := model.URLDetails{}
	//check if the url ia allready available in map
	shortURL := CheckIfURLAvailable(reqURLDetails.URL)
	if shortURL != "" {
		return shortURL, nil
	}
	//gnerate code uniq code
	genCode, err := GenerateCryptoID()
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
	StoreURLToMap(&urlDetails)
	return urlDetails.ShortenedURL, nil
}

// Encode a string to Base64
func EncodeToString(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func CheckIfURLAvailable(originalURL string) string {
	urlMap.Mux.Lock()
	defer urlMap.Mux.Unlock()
	for _, urlDetails := range urlMap.Data {
		if urlDetails.OriginalURL == originalURL {
			return urlDetails.ShortenedURL
		}
	}
	return ""
}

func GenerateCryptoID() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("ERROR_WHILE_GNERATING_UNIQ_ID")
		return "", err
	}
	return EncodeToString(hex.EncodeToString(bytes)), nil
}

func UpdateDomainStats(originalURL string) error {
	domainName := ExtractDomainName(originalURL)
	if domainName == "" {
		log.Println("domain name not found in url")
		return errors.New("DOMAIN_NAME_MISSING_IN_URL")
	}
	if count := LoadStats(domainName); count > 0 {
		UpdateStats(domainName, count+1)
	} else {
		StoreStats(domainName)
	}
	return nil
}

func GetStatsService() []KeyValue {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	if len(stats.Data) == 0 {
		return []KeyValue{}
	}
	// Convert the map to a slice of KeyValue
	keyValueSlice := []KeyValue{}
	for key, value := range stats.Data {
		keyValueSlice = append(keyValueSlice, KeyValue{Key: key, Value: value})
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
