package mos6502

// ReadFunc defines a function where the CPU can read from RAM or Bus
type ReadFunc func(uint16) uint8

// WriteFunc defines a function where the CPU can write to RAM or Bus
type WriteFunc func(uint16, uint8)

var (
	// Read points to a function where the CPU can read from RAM or Bus
	Read ReadFunc

	// Write points to a function where the CPU can write to RAM or Bus
	Write WriteFunc

	// PC = ProgramCounter
	PC uint16

	// SP = Stack Pointer
	SP uint8

	// A Accumulator
	A uint8

	// X Index
	X uint8

	// Y Index
	Y uint8
)

// Test just returns a sample
func TryCallback(rf ReadFunc, wf WriteFunc) uint8 {

	wf(0, 0x42)

	return rf(0)
}
