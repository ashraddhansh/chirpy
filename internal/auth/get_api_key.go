package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	
	fields := strings.Fields(header)
	if len(fields) != 2 || !strings.EqualFold(fields[0], "ApiKey") {
		return "", errors.New("invalid Authorization header")
	}
	return fields[1], nil

}
