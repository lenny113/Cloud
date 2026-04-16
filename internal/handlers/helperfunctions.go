package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
)

// GetAndHashAPIKey gets the api key from the request header, hashes it if it exists and returns it
func GetAndHashAPIKey(r *http.Request) string {
	apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))
	//not key provided, return nothing
	if apiKey == "" {
		return ""
	}

	//key provided so we hash the key with the SHA256 algorithm
	apiKeyHash := sha256.Sum256([]byte(apiKey))
	//returning hashed key as a string
	return fmt.Sprintf("%x", apiKeyHash)
}
