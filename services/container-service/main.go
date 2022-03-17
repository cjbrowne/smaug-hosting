package main

import (
	"bitbucket.org/smaug-hosting/services/cache"
	"bitbucket.org/smaug-hosting/services/container-service/containers"
	"bitbucket.org/smaug-hosting/services/database"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"bitbucket.org/smaug-hosting/services/libhttp/middleware"
	"bitbucket.org/smaug-hosting/services/logging"
	"bitbucket.org/smaug-hosting/services/micro"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	logging.Setup()
	database.Setup()
	cache.Setup()

	addr := µ.GetEnvDefault("LISTEN_ADDR", ":35000")

	mildRateLimit := middleware.RateLimit{Requests: 5, Per: time.Second, BlockTime: 30 * time.Second}

	// NB: always order endpoint registrations from most-specific toward least-specific.  This is to avoid
	// the more general routes matching the more specific routes first, and taking precedence.
	// In the long run the microframework should either handle this better (by ordering routes by "specificity")
	// and/or provide a "precedence" option to give the microframework an order "hint".

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     containers.HandleStartContainer,
		Pattern:     "/containers/{containerId}/start/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}, mildRateLimit, middleware.RequireAuth{}},
		Method:      "POST",
		Description: "Starts a stopped container (has no effect if container already started)",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     containers.HandleStopContainer,
		Pattern:     "/containers/{containerId}/stop/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}, mildRateLimit, middleware.RequireAuth{}},
		Method:      "POST",
		Description: "Stops a container",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     containers.HandleDeleteContainer,
		Pattern:     "/containers/{containerId}/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}, mildRateLimit, middleware.RequireAuth{}},
		Method:      "DELETE",
		Description: "Deletes a container",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     containers.HandleGetContainers,
		Pattern:     "/containers/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}, mildRateLimit, middleware.RequireAuth{}},
		Method:      "GET",
		Description: "Get a list of your own containers",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     containers.HandlePostContainer,
		Pattern:     "/containers/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}, mildRateLimit, middleware.RequireAuth{}},
		Method:      "POST",
		Description: "Create a new container",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     libhttp.NoopHandler,
		Pattern:     ".*",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "OPTIONS",
		Description: "Respond to OPTIONS request with CORS headers",
	})

	logrus.Fatal(http.ListenAndServe(addr, libhttp.ServeMux()))
}
