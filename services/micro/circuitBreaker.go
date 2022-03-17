package µ

import "github.com/sirupsen/logrus"

type CircuitBreaker struct {
	MaxErrors int
	Name      string
	tripped   bool
	errors    int
}

func (cb *CircuitBreaker) Reset() {
	cb.tripped = false
}

func (cb *CircuitBreaker) IsTripped() bool {
	return cb.tripped
}

func (cb *CircuitBreaker) RegisterError() {
	cb.errors++
	if cb.errors > cb.MaxErrors {
		logrus.Errorf("Tripping circuit-breaker %s", cb.Name)
		cb.tripped = true
	}
}
