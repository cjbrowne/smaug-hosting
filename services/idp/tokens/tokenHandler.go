package tokens

import (
	"bitbucket.org/smaug-hosting/services/idp/bge_crypto"
	idp "bitbucket.org/smaug-hosting/services/idp/services"
	"bitbucket.org/smaug-hosting/services/idp/users"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func TokenHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		rawBody, err := ioutil.ReadAll(request.Body)
		if err != nil {
			libhttp.SendError(http.StatusBadRequest, "Could not read body from http POST request", response)
			return
		}

		userLogin := users.UserLogin{}
		serviceLogin := idp.ServiceLogin{}

		err = json.Unmarshal(rawBody, &userLogin)
		if err != nil {
			err = json.Unmarshal(rawBody, &serviceLogin)
			if err != nil {
				libhttp.SendError(http.StatusBadRequest, "Could not parse JSON from http POST request", response)
				return
			}

			srvc, err := idp.ServiceRepository{}.FindByClientId(serviceLogin.ClientId)
			if err != nil {
				libhttp.SendError(http.StatusBadRequest, "Could not find service with that client id", response)
				return
			}

			bge_crypto.Verify(serviceLogin.ClientSecret, srvc.ClientSecretHash)
		}

		u, err := users.UserRepository{}.FindByEmail(userLogin.Email)

		if err != nil && err != users.ErrUserNotFound {
			logrus.Errorf("Failed to fetch user from database: %s", err)
			libhttp.SendError(http.StatusInternalServerError, "Could not create token: could not fetch user from database", response)
			return
		}

		if u == nil {
			logrus.Debugf("Using dummy user (err = %s)", err)
			u = &users.Dummy
		}

		if !bge_crypto.Verify(userLogin.Password, u.PasswordHash) {
			libhttp.SendError(http.StatusUnauthorized, "Could not create token: email or password invalid", response)
			return
		}

		tokStr, err := GenerateToken(*u)

		if err != nil {
			logrus.Errorf("Could not generate token: %s", err)
			libhttp.SendError(http.StatusInternalServerError, "Could not generate token", response)
			return
		}

		logrus.Tracef("Generated token %s", tokStr)

		tok := Token{
			Token: tokStr,
		}

		err = TokenRepository{}.Save(tok)
		if err != nil {
			logrus.Errorf("Could not save token", err)
			libhttp.SendError(http.StatusInternalServerError, "Could not save token", response)
			return
		}

		libhttp.SendJson(TokenResponse{
			Token: tok.Token,
		}, response)
	case "GET":
		libhttp.SendError(http.StatusMethodNotAllowed, "You are not permitted to list tokens", response)
	}
}

type RefreshRequest struct {
	RefreshToken []byte
}

func RefreshTokenHandler(response http.ResponseWriter, request *http.Request) {
	refreshRequest := RefreshRequest{}
	err := libhttp.UnmarshalBody(request, response, &refreshRequest)
	if err != nil {
		return
	}

	token, err := TokenRepository{}.FindByRefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		logrus.Errorf("Could not find refresh token in database: %s", err)
		libhttp.SendError(http.StatusNotFound, "Could not find refresh token in database", response)
		return
	}
	claims := TokenClaims{}

	err = ParseToken(token.Token, &claims)
	if err != nil {
		logrus.Errorf("Could not parse token from database token string: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not parse token from database", response)
		return
	}

	user, err := users.UserRepository{}.Find(claims.UserId)
	if err != nil {
		logrus.Errorf("Invalid token, user id was invalid: %s", err)
		// todo: might want to have a further think about the most appropriate status code here
		libhttp.SendError(http.StatusInternalServerError, "User id in token was invalid", response)
		return
	}

	token.Token, err = GenerateToken(*user)

	err = TokenRepository{}.Save(token)
	if err != nil {
		logrus.Errorf("Could not save token in database: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not save token in database", response)
		return
	}

	libhttp.SendJson(TokenResponse{
		Token: token.Token,
	}, response)
}
