package on_off_fan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldTurnFanOn(t *testing.T) {
	// Temperature is below threshold, fan is off, and last state change was more than 1 second ago, fan should stay off
	lastStateChange := time.Now().Add(-2 * time.Second)
	assert.False(t, shouldTurnFanOn(25, 30, false, time.Duration(1*time.Second), lastStateChange))

	// Temperature is above threshold, fan is off, and last state change was more than 1 second ago, fan should turn on
	lastStateChange = time.Now().Add(-2 * time.Second)
	assert.True(t, shouldTurnFanOn(30, 30, false, time.Duration(1*time.Second), lastStateChange))

	// Temperature is above threshold, fan is on, and last state change was more than 1 second ago, fan should stay on
	lastStateChange = time.Now().Add(-2 * time.Second)
	assert.False(t, shouldTurnFanOn(30, 30, true, time.Duration(1*time.Second), lastStateChange))

	// Temperature is above threshold, fan is off, and last state change was less than 1 second ago, fan should stay off
	lastStateChange = time.Now().Add(-500 * time.Millisecond)
	assert.False(t, shouldTurnFanOn(30, 30, false, time.Duration(1*time.Second), lastStateChange))
}

func TestShouldTurnFanOff(t *testing.T) {
	// Temperature is above threshold, fan is on, and last state change was more than 1 second ago, fan should stay on
	lastStateChange := time.Now().Add(-2 * time.Second)
	assert.False(t, shouldTurnFanOff(30, 30, true, time.Duration(1*time.Second), lastStateChange))

	// Temperature is below threshold, fan is on, and last state change was more than 1 second ago, fan should turn off
	lastStateChange = time.Now().Add(-2 * time.Second)
	assert.True(t, shouldTurnFanOff(25, 30, true, time.Duration(1*time.Second), lastStateChange))

	// Temperature is below threshold, fan is off, and last state change was more than 1 second ago, fan should stay off
	lastStateChange = time.Now().Add(-2 * time.Second)
	assert.False(t, shouldTurnFanOff(25, 30, false, time.Duration(1*time.Second), lastStateChange))

	// Temperature is below threshold, fan is on, and last state change was less than 1 second ago, fan should stay on
	lastStateChange = time.Now().Add(-500 * time.Millisecond)
	assert.False(t, shouldTurnFanOff(25, 30, true, time.Duration(1*time.Second), lastStateChange))
}
