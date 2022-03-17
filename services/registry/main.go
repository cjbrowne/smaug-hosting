package main

import (
	"bitbucket.org/smaug-hosting/services/database"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"bitbucket.org/smaug-hosting/services/libhttp/middleware"
	"bitbucket.org/smaug-hosting/services/logging"
	"bitbucket.org/smaug-hosting/services/registry/health"
	"bitbucket.org/smaug-hosting/services/registry/services"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	logging.Setup()
	database.Setup()
	log := logrus.New()

	health.HealthCheckMain()

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     services.Handler,
		Pattern:     "/services/",
		Middleware:  nil,
		Method:      "GET",
		Description: "Get a list of all registered services",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     services.Handler,
		Pattern:     "/services/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "POST",
		Description: "Register a new service",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     services.Handler,
		Pattern:     "/services/{id}",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "PUT",
		Description: "Overwrite an existing service with new metadata",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     libhttp.SendDocs,
		Pattern:     "/",
		Middleware:  nil,
		Method:      "GET",
		Description: "Get a list of all endpoints served by this API",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     libhttp.SendHeaders,
		Pattern:     "*",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "OPTIONS",
		Description: "Respond to pre-flight OPTIONS requests with CORS headers",
	})

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":4000"
	}

	logrus.Infof("Now listening on %s", addr)
	err := http.ListenAndServe(addr, libhttp.ServeMux())
	if err != nil {
		log.Errorf("Could not listen: %s", err)
	} else {
		log.Infof("Shut down http server")
	}
}
