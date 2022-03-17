package libhttp

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
)

type Middleware interface {
	Run(response http.ResponseWriter, request* http.Request) bool
}

type Endpoint struct {
	Handler     http.HandlerFunc `json:"-"`
	Pattern     string
	Middleware  []Middleware `json:"-"`
	Method      string
	Description string
}

var endpoints []Endpoint

func RegisterEndpoint(endpoint Endpoint) {
	if endpoints == nil {
		endpoints = make([]Endpoint, 0)
	}

	endpoints = append(endpoints, endpoint)
}

func match(pattern string, path string) bool {
	re, err := regexp.Compile("/\\{.*\\}/")
	if err != nil {
		logrus.Errorf("Could not compile regexp: %s", err)
		return false
	}
	pattern = string(re.ReplaceAll([]byte(pattern), []byte(".*")))
	matched, err := regexp.Match(pattern, []byte(path))
	if err != nil {
		logrus.Errorf("Could not run regexp: %s", err)
		return false
	}
	return matched
}

func runAllMiddleware(mwArr []Middleware, w http.ResponseWriter, r *http.Request) bool {
	if len(mwArr) == 0 {
		logrus.Infof("No middlewares defined for this endpoint")
		return false
	}

	for _, mw := range mwArr {
		if mw.Run(w, r) {
			return true
		}
	}

	return false
}

func NoopHandler(http.ResponseWriter, *http.Request) {

}

func ServeMux() http.Handler {
	return http.HandlerFunc(func (response http.ResponseWriter, request *http.Request) {
		responded := false
		logrus.Infof("Received request for %s %s", request.Method, request.URL.Path)
		for _, ep := range endpoints {
			if ep.Method == request.Method || ep.Method == "ALL" || ep.Method == "*" || ep.Method == "" {
				if match(ep.Pattern, request.URL.Path) {
					updateRequestContextWithPathParams(ep.Pattern, request)
					logrus.Infof("Matching to handler with description \"%s\"", ep.Description)
					if !runAllMiddleware(ep.Middleware, response, request) {
						if ep.Handler == nil {
							logrus.Warnf("No handler defined for endpoint, serving HTTP status 501")
							SendError(http.StatusNotImplemented, "Endpoint not yet implemented", response)
						} else {
							ep.Handler.ServeHTTP(response, request)
						}
					}
					responded = true
					break // <-- important! otherwise we will run *every* handler that matches, instead of only the *first* match
				}
			}
		}
		if !responded {
			SendError(http.StatusNotFound, "not found", response)
		}
	})
}

func updateRequestContextWithPathParams(pattern string, request *http.Request) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(request.URL.Path, "/")
	if len(patternParts) != len(pathParts) {
		logrus.Warnf("Path & Pattern parts lengths don't match (%d vs %d)", len(pathParts), len(patternParts))
		return
	}

	for idx, part := range patternParts {
		if strings.HasPrefix(part, "{") {
			part = strings.Trim(part, "{}")
			logrus.Tracef("Adding context variable: %s=%s", part, pathParts[idx])
			*request = *request.WithContext(context.WithValue(request.Context(), part, pathParts[idx]))
		}
	}
}