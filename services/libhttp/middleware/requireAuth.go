package middleware

import (
	"bitbucket.org/smaug-hosting/services/idp/tokens"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func failAuth(response http.ResponseWriter) {
	libhttp.SendError(http.StatusUnauthorized, "You are not permitted to perform this request", response)
}

type RequireAuth struct{}

func (ra RequireAuth) Run(response http.ResponseWriter, request *http.Request) bool {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		failAuth(response)
		return true
	}

	token := strings.Split(authHeader, " ")[1]
	if token == "" {
		failAuth(response)
		return true
	}

	claims := tokens.TokenClaims{}

	err := tokens.ParseToken(token, &claims)
	if err != nil {
		logrus.Debugf("Error parsing token: %s", err)
		failAuth(response)
		return true
	}

	*request = *request.WithContext(context.WithValue(request.Context(), "token_claims", claims))

	return false
}
