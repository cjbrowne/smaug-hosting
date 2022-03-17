package libhttp

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func UnmarshalBody(request *http.Request, response http.ResponseWriter, i interface{}) error {
	rawBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logrus.Errorf("Could not read request body")
		SendError(http.StatusBadRequest, "Could not read request body", response)
		return err
	}
	err = json.Unmarshal(rawBody, i)
	if err != nil {
		logrus.Errorf("Could not parse request body")
		SendError(http.StatusBadRequest, "Could not parse request body", response)
		return err
	}
	return nil
}