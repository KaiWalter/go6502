package mos6502

// Status Flags
const (
	C byte = (1 << 0) // Carry Bit
	Z byte = (1 << 1) // Zero
	I byte = (1 << 2) // Disable Interrupts
	D byte = (1 << 3) // Decimal Mode
	B byte = (1 << 4) // Break
	U byte = (1 << 5) // Unused
	V byte = (1 << 6) // Overflow
	N byte = (1 << 7) // Negative
)

// SetFlag sets or unsets a flag on the Status register
func SetFlag(flag byte, value bool) {
	if value {
		Status |= flag
	} else {
		reverse := 0xFF ^ flag
		Status &= reverse
	}
}

// GetFlag get a flag from the Status Register
func GetFlag(flag byte) bool {
	return (Status & flag) != 0
}

// GetFlag get a flag from the Status Register as numeric value
func GetFlagN(flag byte) uint16 {
	if (Status & flag) != 0 {
		return 1
	}
	return 0
}
