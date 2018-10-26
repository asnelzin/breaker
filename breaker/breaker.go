package breaker

import (
	"fmt"
	"time"
)

// State is a Breaker state.
type State int

// State variants.
const (
	Closed State = iota
	Open
	HalfOpen
)

type RecorderError struct {
	error
}

func (e RecorderError) Error() string {
	return e.Error()
}

// Breaker implements CircuitBreaker pattern.
type Breaker struct {
	Threshold         int
	InvocationTimeout time.Duration
	ResetTimeout      time.Duration

	failCount    int
	lastFailedAt time.Time
	lastError    RecorderError
}

func (b Breaker) GetState() State {
	switch {
	case (b.failCount >= b.Threshold) && (time.Now().Sub(b.lastFailedAt) > b.ResetTimeout):
		return HalfOpen
	case b.failCount >= b.Threshold:
		return Open
	default:
		return Closed
	}
}

// Call function with middleware logic
func (b Breaker) Call(f func() (interface{}, error)) (resp interface{}, err error) {
	switch b.GetState() {
	case Closed, HalfOpen:
		resp, err = b.withTimeout(f)
		if err != nil {
			b.recordFail()
		}
	case Open:
		return nil, fmt.Errorf("CB open")
	}
	return
}

func (b Breaker) withTimeout(f func() (interface{}, error)) (resp interface{}, err error) {
	ch := make(chan bool, 1)
	defer close(ch)

	timer := time.NewTimer(b.InvocationTimeout)
	defer timer.Stop()

	go func() {
		resp, err = f()
		ch <- true
	}()

	select {
	case <-ch:
	case <-timer.C:
		return nil, fmt.Errorf("Operation timed out")
	}
	return
}

func (b *Breaker) recordFail() {
	b.failCount++
	b.lastFailedAt = time.Now()
}
