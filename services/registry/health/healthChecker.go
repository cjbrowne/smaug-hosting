package health

import (
	"bitbucket.org/smaug-hosting/services/registry/services"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func runHealthCheck(service services.Service, wg *sync.WaitGroup) {
	defer wg.Done()

	log := logrus.WithField("service", service.Name)

	c := http.Client{}
	resp, err := c.Get(service.HealthCheck.String())
	if err != nil {
		log.Errorf("Could not GET health check: %s", err)
		return
	}

	healthCheckResponse := services.HealthStatus{}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Could not read response body from health check: %s", err)
		return
	}

	err = json.Unmarshal(b, &healthCheckResponse)
	if err != nil {
		log.Errorf("Could not unmarshal JSON from health check response body: %s", err)
		return
	}

	service.Health = healthCheckResponse
	_, err = services.ServiceRepository{}.Save(service)
	if err != nil {
		log.Errorf("Could not save service: %s", err)
	}
}

func HealthCheckMain() {
	interval := os.Getenv("HEALTH_CHECK_INTERVAL")
	if interval == "" {
		interval = "5s"
	}

	i, err := time.ParseDuration(interval)
	if err != nil {
		i = time.Second * 5
		logrus.Errorf("Could not parse duration for health check: %s", err)
	}

	// I am in awe at this feature.  Not only does it only take 2 LOC to create a clock and block the goroutine
	// until the clock ticks, but this is a *standard library feature*, and to top it all off the default behaviour
	// is smart - it won't "tick" while you're processing, so you don't build up a backlog of unprocessed ticks.
	// this is so elegant!!
	ticker := time.NewTicker(i)

	go func() {
		for range ticker.C {
			wg := sync.WaitGroup{}
			allServices := services.ServiceRepository{}.All()
			for _, sv := range allServices {
				wg.Add(1)
				go runHealthCheck(sv, &wg)
			}
			wg.Wait()
		}
	}()
}
