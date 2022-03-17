package services

import (
	"bitbucket.org/smaug-hosting/services/billing/pricing"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

type billingService struct {
	baseUrl url.URL
}

const BillingServiceUrl = "http://localhost:25000"

func BillingService() billingService {
	baseUrl, err := url.Parse(BillingServiceUrl)

	if err != nil {
		logrus.Errorf("Could not create billing service URL: %s", err)
	}

	return billingService{
		baseUrl: *baseUrl,
	}
}

func (b billingService) GetPrice(software string, tier int) (pricing.Price, error) {

	var price pricing.Price

	pricingUrl, err := b.baseUrl.Parse("/pricing/")

	if err != nil {
		return price, err
	}

	response, err := http.DefaultClient.Get(pricingUrl.String())
	if err != nil {
		return price, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return price, err
	}

	err = json.Unmarshal(bodyBytes, price)

	return price, err
}
