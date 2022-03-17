package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func RestrictiveCors(origins []string) http.HandlerFunc {
	return func (w http.ResponseWriter, r * http.Request) {
		logrus.Tracef("Adding CORS header")
		w.Header().Add("Access-Control-Allow-Origin", strings.Join(origins, ","))
	}
}

type Cors struct{}

func (c Cors) Run(w http.ResponseWriter, r * http.Request) bool {
	logrus.Tracef("Adding CORS headers")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PATCH, OPTIONS")
	return false
}