package httpkit

import (
	"fmt"
	"net/http"
	"strconv"
)

// GetRequiredIntParam ...
func GetRequiredIntParam(key string, r *http.Request) (int64, error) {
	param, err := getQueryParam(key, r)
	if err != nil {
		return 0, fmt.Errorf("param %s is required", key)
	}

	return parseIntParam(key, param)
}

// GetIntParam ...
func GetIntParam(key string, r *http.Request) (int64, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return 0, nil
	}

	return parseIntParam(key, param)
}

// GetStrRequiredParam ...
func GetStrRequiredParam(key string, r *http.Request) (string, error) {
	return getQueryParam(key, r)
}

// GetStrParam ...
func GetStrParam(key string, r *http.Request) string {
	return r.URL.Query().Get(key)
}

func getQueryParam(key string, r *http.Request) (string, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return "", fmt.Errorf("param %s is missing", key)
	}
	return param, nil
}

func parseIntParam(key, param string) (int64, error) {
	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s param=%s", key, param)
	}
	return value, nil
}
