package verify

import (
	"bitbucket.org/smaug-hosting/services/idp/users"
	"bitbucket.org/smaug-hosting/services/libhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HandleVerifyEmail(response http.ResponseWriter, request *http.Request) {
	user, err := users.UserRepository{}.FindByVerificationToken(request.Context().Value("token").(string))
	if err != nil {
		logrus.Errorf("Could not verify user token: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not verify token", response)
		return
	}

	if user == nil {
		logrus.Errorf("Could not verify user: no user found for token")
		libhttp.SendError(http.StatusBadRequest, "User not found", response)
		return
	}

	user.VerificationToken = ""
	user.Verified = true
	err = users.UserRepository{}.Verify(*user)
	if err != nil {
		logrus.Errorf("Could not verify user: %s", err)
		libhttp.SendError(http.StatusInternalServerError, "Could not mark user as verified", response)
		return
	}

	// if all good, send 200
	libhttp.SendJson(struct{}{}, response)
}
