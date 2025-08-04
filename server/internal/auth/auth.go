package auth

import (
	"errors"
	"net/http"
)

func GetUrlQueryToken(r *http.Request) (string, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return "", errors.New("token not found in URL query")
	}
	return token, nil
}
