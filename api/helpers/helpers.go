package helpers

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

func RemoveDomainNameError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newURL := strings.Replace(url, "http://", "", -1)
	newURL = strings.Replace(newURL, "https://", "", -1)
	newURL = strings.Replace(newURL, "www.", "", -1)
	newURL = strings.Split(newURL, "/")[0]

	return newURL != os.Getenv("DOMAIN")
}

func EnforceHTTP(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return url
}

func GenerateID() string {
	return uuid.New().String()[:6]
}
