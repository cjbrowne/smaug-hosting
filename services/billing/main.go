package main

import (
	"bitbucket.org/smaug-hosting/services/billing/Internal"
	"bitbucket.org/smaug-hosting/services/billing/billing"
	"bitbucket.org/smaug-hosting/services/billing/transactions"
	"bitbucket.org/smaug-hosting/services/database"
	"bitbucket.org/smaug-hosting/services/idp/tokens"
	"bitbucket.org/smaug-hosting/services/idp/users"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"bitbucket.org/smaug-hosting/services/libhttp/middleware"
	"bitbucket.org/smaug-hosting/services/libws"
	"bitbucket.org/smaug-hosting/services/logging"
	µ "bitbucket.org/smaug-hosting/services/micro"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/event"
	"net/http"
	"os"
	"time"
)

func main() {
	logging.Setup()
	database.Setup()

	ws := libws.SetupWebsocket("/ws")

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     billing.HandleGetTopup,
		Pattern:     "/topup/{amount}/",
		Middleware:  []libhttp.Middleware{middleware.Cors{}, middleware.RequireAuth{}},
		Method:      "GET",
		Description: "Redirects to stripe page for payment processing",
	})

	libhttp.RegisterEndpoint(libhttp.Endpoint{
		Handler:     libhttp.NoopHandler,
		Pattern:     ".*",
		Middleware:  []libhttp.Middleware{middleware.Cors{}},
		Method:      "OPTIONS",
		Description: "Handle CORS preflight request",
	})

	// bill users every 1 minute on the dot
	go func() {
		// the time of the next clock minute
		for {
			nextRun := time.Now().Truncate(time.Minute).Add(time.Minute)
			time.Sleep(time.Until(nextRun))
			go billingint.BillAllUsers()
		}
	}()

	pollTime := time.Second

	eventPollingCircuitBreaker := new(µ.CircuitBreaker)
	eventPollingCircuitBreaker.MaxErrors = 0
	eventPollingCircuitBreaker.Name = "STRIPE_EVENT_POLL"

	eventHandlingCircuitBreaker := new(µ.CircuitBreaker)
	eventHandlingCircuitBreaker.MaxErrors = 0
	eventHandlingCircuitBreaker.Name = "STRIPE_HANDLE_EVENTS"

	// check for successful top-ups every 1 second
	go func() {
		for {

			// the time of the next clock second
			nextRun := time.Now().Truncate(pollTime).Add(pollTime)
			time.Sleep(time.Until(nextRun))
			if eventPollingCircuitBreaker.IsTripped() || eventHandlingCircuitBreaker.IsTripped() {
				logrus.WithField("notify", "chat,email").Fatalf("Circuit breaker tripped!")
			}
			stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

			params := &stripe.EventListParams{
				Type: stripe.String("checkout.session.completed"),
				CreatedRange: &stripe.RangeQueryParams{
					// check for events in the last day
					GreaterThan: time.Now().Unix() - 24*3600,
				},
			}
			i := event.List(params)

			if i.Err() != nil {
				realErr := i.Err().(*stripe.Error)
				if r := recover(); r != nil {
					logrus.Errorf("Some other stripe error happened which we can't really recover from but I'm gonna try anyway. (%s)", i.Err())
					continue
				}
				if realErr.Code == stripe.ErrorCodeRateLimit {
					pollTime *= 2 // exponential back-off when being rate-limited
					if pollTime > 60 {
						eventPollingCircuitBreaker.RegisterError()
						continue
					}
				} else {
					logrus.Errorf("Error fetching events from stripe: %s", i.Err())
					continue
				}
			}

			for i.Next() {
				nextEvent := i.Event()
				var session stripe.CheckoutSession
				err := json.Unmarshal(nextEvent.Data.Raw, &session)
				if err != nil {
					logrus.Errorf("Could not process checkout event: %s", err)
					eventHandlingCircuitBreaker.RegisterError()
					continue
				}

				transaction, err := transactions.PendingTransactionRepository{}.FindByCheckoutId(session.ID)
				if err != nil {
					// see if we have a completed transaction matching the session id
					completed, err := transactions.PendingTransactionRepository{}.FindCompletedByCheckoutId(session.ID)
					if err != nil || completed == nil {
						logrus.Errorf("Error while fetching completed transaction: %s", err)
						eventHandlingCircuitBreaker.RegisterError()
						continue
					} else {
						// this is a (common) case of an already-processed tx being seen again in the response from stripe
						continue
					}
				}

				user, err := users.UserRepository{}.Find(transaction.UserId)
				if err != nil {
					logrus.Errorf("Error while fetching user: %s", err)
					eventHandlingCircuitBreaker.RegisterError()
					continue
				}

				user.Balance += transaction.Amount

				err = users.UserRepository{}.Save(*user)
				if err != nil {
					logrus.Errorf("Error while saving user back to database: %s", err)
					eventHandlingCircuitBreaker.RegisterError()
					continue
				}

				err = transactions.PendingTransactionRepository{}.MarkAsCompleted(transaction)
				if err != nil {
					logrus.Errorf("Error while saving user back to database: %s", err)
					eventHandlingCircuitBreaker.RegisterError()
					continue
				}
			}
		}
	}()

	// update all users every 1 second
	go func() {
		for {
			// the time of the next clock second
			nextRun := time.Now().Truncate(time.Second).Add(time.Second)
			time.Sleep(time.Until(nextRun))
			for _, client := range ws.AllClients {

				claims := tokens.TokenClaims{}
				logrus.Tracef("Sending balance to client: %+v", client)
				err := tokens.ParseToken(client.Session.Token, &claims)
				if err != nil {
					logrus.Warnf("Websocket client with invalid token found: %s", err)
					continue
				}

				user, err := users.UserRepository{}.Find(claims.UserId)
				if err != nil {
					logrus.Warnf("Websocket client with invalid user (id: %d) found: %s", claims.UserId, err)
					continue
				}

				if user == nil {
					logrus.Warnf("No user found for id: %d", claims.UserId)
					continue
				}

				ws.Send() <- libws.WSMessage{
					Subject: "balance",
					Body: map[string]interface{}{
						"balance": user.Balance,
					},
					Connection: client.Connection,
				}
			}
		}
	}()

	logrus.Fatal(http.ListenAndServe(µ.GetEnvDefault("LISTEN_ADDR", ":25000"), libhttp.ServeMux()))
}
