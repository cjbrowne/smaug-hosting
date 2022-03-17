package services

import (
	"bitbucket.org/smaug-hosting/services/libhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func handleGet(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	capability := query.Get("capability")

	page, err := strconv.ParseUint(query.Get("page"), 10, 64)
	if err != nil {
		logrus.Debugf("Invalid page number: %s", query.Get("page"))
		page = 0
	}

	pageSize, err := strconv.ParseUint(query.Get("size"), 10, 64)
	if err != nil {
		logrus.Debugf("Invalid page size: %s", query.Get("size"))
		pageSize = 50
	}

	res := ServiceRepository{}.Find(ServiceSearchQuery{
		Page:       page,
		PageSize:   pageSize,
		Capability: capability,
	})
	if res.Error == nil {
		libhttp.SendJson(res, response)
	} else {
		logrus.Errorf("Error querying the database: %s", res.Error)
		libhttp.SendError(http.StatusInternalServerError, "Error while querying the database for services", response)
	}
}

func Handler(response http.ResponseWriter, request *http.Request) {
	switch strings.ToLower(request.Method) {
	case "get":
		handleGet(response, request)
	}
}
