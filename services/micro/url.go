package µ

import (
	"github.com/sirupsen/logrus"
	"net/url"
)

func ResolveUrl(baseUrl string, parts ...string) string {
	base, err := url.Parse(baseUrl)
	if err != nil {
		logrus.Errorf("Could not parse url: %s", err)
		return ""
	}
	for _, part := range parts {
		base, err = base.Parse(part)
		if err != nil {
			logrus.Errorf("Could not parse url part: %s", err)
			break
		}
	}

	logrus.Debugf("Resolved URL %s from base %s and parts %+v", base.String(), baseUrl, parts)

	return base.String()
}
