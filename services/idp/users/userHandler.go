package users

import (
	"bitbucket.org/smaug-hosting/services/idp/bge_crypto"
	"bitbucket.org/smaug-hosting/services/idp/email"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"encoding/json"
	"github.com/docker/docker/pkg/stringutils"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		newUser := userRequest{}
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			logrus.Errorf("Could not read body: %s", err)
			libhttp.SendError(http.StatusInternalServerError, "Could not read body from HTTP request", response)
			return
		}

		err = json.Unmarshal(body, &newUser)
		if err != nil {
			logrus.Errorf("Could not parse JSON body into struct: %s", err)
			libhttp.SendError(http.StatusBadRequest, "Malformed JSON", response)
			return
		}

		user, err := UserRepository{}.FindByEmail(newUser.Email)
		if user != nil {
			logrus.Debugf("Duplicate user: %s", newUser.Email)
			libhttp.SendError(http.StatusConflict, "That email address is already in use", response)
			return
		}

		verificationToken := stringutils.GenerateRandomAlphaOnlyString(64)
		err = UserRepository{}.Save(User{
			Email:             newUser.Email,
			PasswordHash:      bge_crypto.Encrypt(newUser.Password),
			Roles:             []Role{},
			VerificationToken: verificationToken,
		})

		if err != nil {
			logrus.Errorf("Could not save user: %s", err)
			libhttp.SendError(http.StatusInternalServerError, "Could not save user", response)
			return
		}
		// send verification email
		err = email.SendVerificationEmail(newUser.Email, verificationToken)
		if err != nil {
			// don't actually send an error for this case, the user can still log in, they just need to verify their email
			logrus.Errorf("Could not send verification email: %s", err)
		}

		libhttp.SendJson(libhttp.StatusMessage{
			Success: true,
			Message: "Created user",
		}, response)
	}
}
