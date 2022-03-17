package µ

import (
	"github.com/sirupsen/logrus"
	"os"
)

func GetEnvDefault(env string, def string) string {
	rv := os.Getenv(env)
	if rv == "" {
		rv = def
	}
	return rv
}

func MustGetEnv(env string) string {
	rv := os.Getenv(env)
	if rv == "" {
		logrus.Fatalf("Please set the environment variable %s", env)
	}
	return rv
}
