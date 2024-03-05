package model

type URLDetails struct {
	OriginalURL  string `json:"originalURL"`
	DomainName   string `json:"domainName"`
	ShortenedURL string `json:"shortenedURL"`
	HashValue    string `json:"hashValue"`
}

type RequestURLData struct {
	URL        string `json:"url"`
	DomainName string `json:"domainName"`
}
