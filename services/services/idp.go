package services

import (
	idp "bitbucket.org/smaug-hosting/services/idp/services"
	"bitbucket.org/smaug-hosting/services/idp/tokens"
	"bitbucket.org/smaug-hosting/services/idp/users"
	µ "bitbucket.org/smaug-hosting/services/micro"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

type authService struct {
	baseUrl url.URL
}

const idpServiceBaseUrl = "http://localhost:45000/"

func AuthService() authService {
	baseUrl, err := url.Parse(idpServiceBaseUrl)
	if err != nil {
		logrus.Errorf("Could not parse base IDP url: %s", err)
	}
	return authService{
		baseUrl: *baseUrl,
	}
}

func (as authService) GetServiceToken(clientId string, clientSecret string) string {
	tokenUrl, err := as.baseUrl.Parse(fmt.Sprintf("/tokens/"))

	if err != nil {
		logrus.Errorf("Could not parse token URL: %s", err)
		return ""
	}

	loginBytes, err := json.Marshal(idp.ServiceLogin{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	})

	if err != nil {
		logrus.Errorf("Could not marshal JSON: %s", err)
		return ""
	}

	response, err := http.Post(tokenUrl.String(), "application/json", bytes.NewReader(loginBytes))
	if err != nil {
		logrus.Errorf("Request failed: %s", err)
		return ""
	}
	tokenBytes, err := ioutil.ReadAll(response.Body)

	tok := tokens.TokenResponse{}

	err = json.Unmarshal(tokenBytes, &tok)
	if err != nil {
		logrus.Errorf("Could not unmarshal JSON: %s", err)
		return ""
	}

	return tok.Token
}

func (as authService) GetUser(id int) (*users.User, error) {
	userUrl, err := as.baseUrl.Parse(fmt.Sprintf("/users/%d/", id))

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", userUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	authHeader := as.GetServiceToken(µ.MustGetEnv("CLIENT_ID"), µ.MustGetEnv("CLIENT_SECRET"))

	req.Header.Add("Authorization", authHeader)

	response, err := http.Get(userUrl.String())
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var user *users.User

	user = new(users.User)

	err = json.Unmarshal(bytes, user)

	return user, err
}
