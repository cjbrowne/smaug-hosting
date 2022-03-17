package libhttp

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type StatusMessage struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func SendError(status int, message string, response http.ResponseWriter) {
	response.WriteHeader(status)

	responseStr, err := json.Marshal(StatusMessage{
		Message: message,
		Success: false,
	})

	if err != nil {
		logrus.Errorf("Could not marshal JSON for HTTP response: %s", err)
		responseStr = []byte("")
		return
	}

	_, err = response.Write(responseStr)
	if err != nil {
		logrus.Errorf("Could not write HTTP Response: %s", err)
		return
	}
}

func SendJsonWithStatus(status int, msg interface{}, response http.ResponseWriter) {
	responseStr, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("Could not marshal JSON for HTTP response: %s", err)
		SendError(http.StatusInternalServerError, "Error while trying to send JSON", response)
		return
	}
	response.WriteHeader(status)
	_, err = response.Write(responseStr)
	if err != nil {
		logrus.Errorf("Could not write HTTP response: %s", err)
		return
	}
}

func SendJson(msg interface{}, response http.ResponseWriter) {
	SendJsonWithStatus(http.StatusOK, msg, response)
}

func SendHeaders(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	_, _ = response.Write([]byte("OK"))
}

func SendDocs(response http.ResponseWriter, request *http.Request) {
	SendJson(endpoints, response)
}