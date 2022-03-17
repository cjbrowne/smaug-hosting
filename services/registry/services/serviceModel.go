package services

import (
	"net/url"
	"time"
)

type Protocol string

const (
	ProtoHttp Protocol = "http"
	ProtoAmqp          = "amqp"
)

type HealthStatus struct {
	Up           bool
	ResponseTime time.Duration
	// load is expressed as a value between 0.0 and 1.0 and, while it could in theory be used for LB, is more often
	// used by monitoring tools and auto-scalers.
	SelfReportedLoad float64
}

type Service struct {
	Id                int
	Name              string
	Protocols         []Protocol
	HttpUrl           url.URL `db:"http_url"`
	HealthCheck       url.URL `db:"health_check"`
	Health            HealthStatus `db:"-"`
	ResponsiblePerson string `db:"responsible_person"`
}
