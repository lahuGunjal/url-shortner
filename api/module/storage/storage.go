package storage

import (
	"sync"

	"github.com/lahuGunjal/url-shortner/api/model"
)

// URLDetailsMap represents the state of the URL shortener.
type URLInfoMap struct {
	Data map[string]*model.URLDetails
	Mux  *sync.Mutex
}

// WebUrl Stats For most used urls
type DomainStats struct {
	Data map[string]int
	Mux  *sync.Mutex
}

type KeyValue struct {
	Key   string
	Value int
}

// InitialiseMap initializes the URLDetailsMap for the URL shortener.
func NewInMemoryStorage() (*URLInfoMap, *DomainStats) {
	return &URLInfoMap{
			Data: make(map[string]*model.URLDetails),
			Mux:  &sync.Mutex{},
		}, &DomainStats{
			Data: make(map[string]int),
			Mux:  &sync.Mutex{},
		}
}

// StoreURLToMap stores the URL details in the URLDetailsMap.
func (urlMap *URLInfoMap) StoreURLToMap(urlDetails *model.URLDetails) {
	urlMap.Mux.Lock()
	defer urlMap.Mux.Unlock()
	urlMap.Data[urlDetails.HashValue] = urlDetails
}

// GetURLFromMap retrieves URL details from the URLDetailsMap based on the provided URL.
func (urlMap *URLInfoMap) GetURLFromMap(url string) *model.URLDetails {
	urlMap.Mux.Lock()
	defer urlMap.Mux.Unlock()
	if val, ok := urlMap.Data[url]; ok {
		return val
	}
	return &model.URLDetails{}
}

// GetOriginalURL to check if url is already added in the memory
func (urlMap *URLInfoMap) GetOriginalURL(originalURL string) string {
	urlMap.Mux.Lock()
	defer urlMap.Mux.Unlock()
	for _, urlDetails := range urlMap.Data {
		if urlDetails.OriginalURL == originalURL {
			return urlDetails.ShortenedURL
		}
	}
	return ""
}

// StoreStats for webDomain
func (stats *DomainStats) StoreStats(domain string) {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	stats.Data[domain] = 1
}

// UpdateStats of webDomain
func (stats *DomainStats) UpdateStats(domain string, count int) {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	stats.Data[domain] = count
}

// LoadStats get stats of webDomain
func (stats *DomainStats) LoadStats(domain string) int {
	stats.Mux.Lock()
	defer stats.Mux.Unlock()
	return stats.Data[domain]
}
