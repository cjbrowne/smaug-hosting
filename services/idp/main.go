package main

import (
	"bitbucket.org/smaug-hosting/services/database"
	"bitbucket.org/smaug-hosting/services/idp/bge_crypto"
	"bitbucket.org/smaug-hosting/services/idp/tokens"
	"bitbucket.org/smaug-hosting/services/idp/users"
	"bitbucket.org/smaug-hosting/services/idp/verify"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"bitbucket.org/smaug-hosting/services/libhttp/middleware"
	"bitbucket.org/smaug-hosting/services/logging"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	database.Setup()
	logging.Setup()

	mildRateLimit := middleware.RateLimit{Requests: 3, Per: time.Minute, BlockTime: time.Minute * 15}

	users.Dummy = users.User{
		Email:        "email@example.com",
		PasswordHash: bge_crypto.Encrypt("monkey1"),
	}

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     verify.HandleVerifyEmail,
		Pattern:     "/verify/{token}/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "GET",
		Description: "Verifies a user's email address using the given token",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     tokens.TokenHandler,
		Pattern:     "/token/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "GET",
		Description: "Get a token",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     tokens.RefreshTokenHandler,
		Pattern:     "/refresh/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "POST",
		Description: "Refresh your access token",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     tokens.TokenHandler,
		Pattern:     "/token/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "POST",
		Description: "Create a new token",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     users.UserHandler,
		Pattern:     "/user/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}, mildRateLimit},
		Method:      "POST",
		Description: "Create a new user",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     libhttp.SendHeaders,
		Pattern:     ".*",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "OPTIONS",
		Description: "Respond to pre-flight OPTIONS requests with CORS headers",
	})

	logrus.Info("Listening on :45000")
	err := http.ListenAndServe(":45000", libhttp.ServeMux())
	if err != nil {
		logrus.Errorf("Could not listen on port 45000: %s", err)
	} else {
		logrus.Infof("Server shut down gracefully")
	}
}
