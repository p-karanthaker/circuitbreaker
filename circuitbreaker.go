package main

import (
	"log"
	"time"
)

const (
	open     = "open"
	closed   = "closed"
	halfOpen = "half-open"
)

type CircuitBreaker struct {
	failureCount     int
	successCount     int
	thresholdToOpen  int           // failure threshold to open the circuit
	thresholdToClose int           // success threshold to close the circuit
	resetTimeout     time.Duration // timeout for open state to move to half-open
	lastErrorTime    time.Time
}

func NewCircuitBreaker(thresholdToOpen int, thresholdToClose int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		thresholdToOpen:  thresholdToOpen,
		thresholdToClose: thresholdToClose,
		resetTimeout:     resetTimeout,
	}
}

func (cb *CircuitBreaker) Execute(f func() error) {
	state := cb.state()
	if state == closed || state == halfOpen {
		if err := f(); err != nil {
			cb.failureCount++
			cb.lastErrorTime = time.Now()
			log.Printf("Call to service failed. %v", cb.failureCount)
			return
		}
		cb.reset()
	} else if state == open {
		log.Printf("Circuit is open")
		return
	}
}

func (cb *CircuitBreaker) reset() {
	cb.failureCount = 0
	cb.successCount = 0
	cb.lastErrorTime = time.Time{}
}

func (cb *CircuitBreaker) state() string {
	switch {
	case cb.failureCount >= cb.thresholdToOpen && time.Since(cb.lastErrorTime) > cb.resetTimeout:
		return halfOpen
	case cb.failureCount >= cb.thresholdToOpen:
		return open
	default:
		return closed
	}
}
