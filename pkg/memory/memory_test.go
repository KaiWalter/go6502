package memory

import (
	"testing"
)

func TestSimpleWriteRead(t *testing.T) {

	const memSize = 0x400

	// arrange
	ram := Memory{AddressOffset: 0xF000, AddressSpace: make([]byte, memSize)}
	var expected byte = 0xAB

	// act
	ram.Write(0xF001, expected)
	actual := ram.Read(0xF001)

	// assert
	if actual != expected {
		t.Errorf("actual %x != expected %x", actual, expected)
	}

}
