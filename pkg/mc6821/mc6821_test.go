package mc6821

import (
	"testing"
	"time"
)

func TestInputOutput(t *testing.T) {

	// arrange
	pia := MC6821{
		Name: "Testing",
	}

	keyboardInputChannelA := make(chan byte, 10)
	pia.SetInputChannelA(keyboardInputChannelA)

	screenOutputChannel := make(chan byte, 10)
	pia.SetOutputChannelB(screenOutputChannel)

	piaCA1Channel := make(chan Signal, 10)
	pia.SetCA1Channel(piaCA1Channel)

	// act
	var expected byte = 0x40
	var actual byte

	go func() {
		for b := range screenOutputChannel {
			actual = b
		}
	}()

	go func() {
		piaCA1Channel <- Fall                      // bring keyboard strobe to low to force active transition
		keyboardInputChannelA <- (expected | 0x80) // bit 7 is constantly set (+5V)
		piaCA1Channel <- Rise                      // send only pulse
		piaCA1Channel <- Fall                      // 20 micro secs are not worth emulating
	}()

	time.Sleep(50 * time.Millisecond)

	// assert
	if actual != expected {
		t.Errorf("actual %x != expected %x", actual, expected)
	}

}
