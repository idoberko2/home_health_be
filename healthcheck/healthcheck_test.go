package healthcheck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const somePassphrase = "somePassphrase"
const someHistoryLength = 10

func testEngine() HealthCheck {
	return New(HealthCheckConfig{
		HistoryLength: someHistoryLength,
		Passphrase:    somePassphrase,
		GracePeriod:   time.Second,
	})
}

func TestPingIncorrectPassphrase(t *testing.T) {
	engine := testEngine()
	err := engine.Ping("incorrectPassphrase")
	assert.NotNil(t, err)
	assert.Equal(t, ErrIncorrectPassphrase, err)
}

func TestIsHealthyNeverPinged(t *testing.T) {
	engine := testEngine()
	_, err := engine.IsHealthy()
	assert.NotNil(t, err)
	assert.Equal(t, ErrNeverPinged, err)
}

func TestIsHealthyTrue(t *testing.T) {
	engine := testEngine()
	err := engine.Ping(somePassphrase)
	assert.Nil(t, err)
	<-time.After(time.Second)
	err = engine.Ping(somePassphrase)
	assert.Nil(t, err)
	result, err := engine.IsHealthy()
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestIsHealthyFalse(t *testing.T) {
	engine := testEngine()
	err := engine.Ping(somePassphrase)
	assert.Nil(t, err)
	<-time.After(1100 * time.Millisecond)
	result, err := engine.IsHealthy()
	assert.Nil(t, err)
	assert.False(t, result)
}

func TestHistoryLength(t *testing.T) {
	engine := testEngine()

	for i := 0; i < someHistoryLength; i++ {
		err := engine.Ping(somePassphrase)
		assert.Nil(t, err)
	}

	<-time.After(time.Second)
	err := engine.Ping(somePassphrase)
	assert.Nil(t, err)

	result, err := engine.IsHealthy()
	assert.Nil(t, err)
	assert.True(t, result)

	history, err := engine.GetHistory()
	assert.Nil(t, err)
	assert.Len(t, history, someHistoryLength)
}
