package mc6821

import (
	"testing"
	"time"
)

const (
	kbd   uint16 = 0xd010 // read key
	kbdcr uint16 = 0xd011 // control port
	dsp   uint16 = 0xd012 // write ascii
	dspcr uint16 = 0xd013 // control port
)

func TestInputOutput(t *testing.T) {

	// arrange
	pia := MC6821{
		Name: "Testing",
	}

	screenOutputChannel := make(chan byte, 10)
	pia.SetOutputChannelB(screenOutputChannel)

	// act
	var expected byte = 0x5A
	var actual byte

	go func() {
		for b := range screenOutputChannel {
			actual = b
		}
	}()

	pia.Write(dsp, 0x7F)   // 01111111 -> DDRB : configure all bits except highest bit for output
	pia.Write(dspcr, 0x04) // 00000100 -> CRB  : write to output port B
	pia.Write(dsp, expected)

	close(screenOutputChannel)

	time.Sleep(50 * time.Millisecond)

	// assert
	if actual != expected {
		t.Errorf("actual %x != expected %x", actual, expected)
	}

}
