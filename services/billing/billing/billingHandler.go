package billing

import (
	"bitbucket.org/smaug-hosting/services/billing/transactions"
	"bitbucket.org/smaug-hosting/services/idp/tokens"
	"bitbucket.org/smaug-hosting/services/libhttp"
	µ "bitbucket.org/smaug-hosting/services/micro"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"net/http"
	"os"
	"strconv"
)

type TopupResponse struct {
	StripeSessionId string `json:"stripe_session_id"`
}

func HandleGetTopup(response http.ResponseWriter, request *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	frontendBaseUrlStr := os.Getenv("FRONTEND_BASE_URL")
	amountStr := request.Context().Value("amount").(string)
	amount, err := strconv.ParseInt(amountStr, 10, 64)

	if err != nil {
		logrus.Warnf("Invalid amount specified for topup request: %s (err=%s)", amountStr, err)
		libhttp.SendError(http.StatusBadRequest, "Invalid amount specified", response)
		return
	}

	claims := request.Context().Value("token_claims").(tokens.TokenClaims)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Name:        stripe.String("Smaug Hosting Credit"),
				Description: stripe.String("Top-up of server credits"),
				Amount:      stripe.Int64(amount),
				Currency:    stripe.String(string(stripe.CurrencyGBP)),
				Quantity:    stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(µ.ResolveUrl(frontendBaseUrlStr, "/topup/success")),
		CancelURL:  stripe.String(µ.ResolveUrl(frontendBaseUrlStr, "/topup/cancelled")),
	}

	newSession, err := session.New(params)
	if err != nil {
		logrus.Errorf("Could not create checkout session")
		libhttp.SendError(http.StatusInternalServerError, "Could not establish stripe session", response)
		return
	}

	transactions.PendingTransactionRepository{}.Save(transactions.PendingTransaction{
		UserId:     claims.UserId,
		Amount:     amount * 10, // convert from milli-pounds (stripe) to micro-pounds (used internally)
		CheckoutId: newSession.ID,
	})

	libhttp.SendJson(TopupResponse{
		StripeSessionId: newSession.ID,
	}, response)
}
