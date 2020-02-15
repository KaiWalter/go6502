package mos6502

// Status Flags
const (
	C uint8 = (1 << 0) // Carry Bit
	Z uint8 = (1 << 1) // Zero
	I uint8 = (1 << 2) // Disable Interrupts
	D uint8 = (1 << 3) // Decimal Mode
	B uint8 = (1 << 4) // Break
	U uint8 = (1 << 5) // Unused
	V uint8 = (1 << 6) // Overflow
	N uint8 = (1 << 7) // Negative
)

// SetFlag sets or unsets a flag on the Status register
func SetFlag(flag uint8, value bool) {
	if value {
		Status |= flag
	} else {
		reverse := 0xFF ^ flag
		Status &= reverse
	}
}

// GetFlag get a flag from the Status Register
func GetFlag(flag uint8) bool {
	return (Status & flag) != 0
}

// GetFlag get a flag from the Status Register as numeric value
func GetFlagN(flag uint8) uint16 {
	if (Status & flag) != 0 {
		return 1
	}
	return 0
}
