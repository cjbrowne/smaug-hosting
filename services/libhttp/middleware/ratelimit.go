package middleware

import (
	"net/http"
	"time"
)

type RateLimit struct {
	Requests  int
	Per       time.Duration
	BlockTime time.Duration
}

// todo: implement rate limiting
func (rl RateLimit) Run(response http.ResponseWriter, request *http.Request) bool {
	return false
}
