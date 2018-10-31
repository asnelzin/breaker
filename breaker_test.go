package breaker

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBreakerInitialStat(t *testing.T) {
	b := prep()
	assert.Equal(t, Closed, b.GetState())
}

func TestBreakerOpened(t *testing.T) {
	b := prep()
	_, err := b.Call(errorFn)
	assert.Equal(t, Closed, b.GetState())
	assert.EqualError(t, err, "error")
	_, err = b.Call(errorFn)
	assert.Equal(t, Open, b.GetState())
	assert.EqualError(t, err, "error")

	_, err = b.Call(errorFn)
	assert.Equal(t, ErrBreakerOpen, err)
}

func TestBreakerHalfOpened(t *testing.T) {
	b := prep()

	b.Call(errorFn)
	b.Call(errorFn)

	time.Sleep(3 * time.Second)
	assert.Equal(t, HalfOpen, b.GetState())

	_, err := b.Call(errorFn)
	assert.EqualError(t, err, "error")
	assert.Equal(t, Open, b.GetState())
}

func TestBreakerReset(t *testing.T) {
	b := prep()

	b.Call(errorFn)
	b.Call(errorFn)

	time.Sleep(3 * time.Second)
	assert.Equal(t, HalfOpen, b.GetState())

	b.Call(func() (interface{}, error) {
		return nil, nil
	})
	assert.Equal(t, Closed, b.GetState())
}

func TestBreakerTimeout(t *testing.T) {
	b := prep()
	start := time.Now()
	_, err := b.Call(func() (interface{}, error) {
		time.Sleep(3 * time.Second)
		return nil, nil
	})
	end := time.Now()
	assert.Equal(t, ErrTimeout, err)
	assert.Equal(t, 2*time.Second, end.Sub(start).Truncate(time.Second))
}

func prep() Breaker {
	return Breaker{
		Threshold:         2,
		InvocationTimeout: 2 * time.Second,
		ResetTimeout:      3 * time.Second,
	}
}

func errorFn() (interface{}, error) {
	return nil, fmt.Errorf("error")
}
