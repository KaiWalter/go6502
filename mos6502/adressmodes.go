package mos6502

func ABS() int {
	lo := uint16(read(PC))
	PC++
	hi := uint16(read(PC))
	PC++
	absoluteAddress = hi<<8 | lo
	return 0
}

func ABX() int {
	lo := uint16(read(PC))
	PC++
	hi := uint16(read(PC))
	PC++
	absoluteAddress = hi<<8 | lo
	absoluteAddress += uint16(X)

	if (absoluteAddress & 0xFF00) != (hi << 8) {
		return 1
	}

	return 0
}

func ABY() int {
	lo := uint16(read(PC))
	PC++
	hi := uint16(read(PC))
	PC++
	absoluteAddress = hi<<8 | lo
	absoluteAddress += uint16(Y)

	if (absoluteAddress & 0xFF00) != (hi << 8) {
		return 1
	}

	return 0
}

func IMM() int {
	absoluteAddress = PC
	PC++
	return 0
}

func IND() int {
	pointerLo := uint16(read(PC))
	PC++
	pointerHi := uint16(read(PC))
	PC++
	pointer := pointerHi<<8 | pointerLo

	if pointerLo == 0x00FF {
		absoluteAddress = (uint16(read(pointer&0xFF00)) << 8) | uint16(read(pointer+0))
	} else {
		absoluteAddress = (uint16(read(pointer+1)) << 8) | uint16(read(pointer+0))
	}
	return 0
}

func IMP() int {
	fetched = A
	return 0
}

func IZX() int {
	tempAddress := uint16(read(PC))
	PC++

	lo := uint16(read((tempAddress + uint16(X)) & 0x00FF))
	hi := uint16(read((tempAddress + uint16(X) + 1) & 0x00FF))

	absoluteAddress = hi<<8 | lo

	return 0
}

func IZY() int {
	tempAddress := uint16(read(PC))
	PC++

	lo := uint16(read((tempAddress + uint16(Y)) & 0x00FF))
	hi := uint16(read((tempAddress + uint16(Y) + 1) & 0x00FF))

	absoluteAddress = hi<<8 | lo

	return 0
}

func REL() int {
	relativeAddress = uint16(read(PC))
	PC++
	if relativeAddress&0x80 != 0 {
		relativeAddress |= 0xFF00
	}
	return 0
}

func ZP0() int {
	absoluteAddress = uint16(read(PC))
	PC++
	absoluteAddress &= 0x00FF
	return 0
}

func ZPX() int {
	absoluteAddress = uint16(read(PC)) + uint16(X)
	PC++
	absoluteAddress &= 0x00FF
	return 0
}

func ZPY() int {
	absoluteAddress = uint16(read(PC)) + uint16(Y)
	PC++
	absoluteAddress &= 0x00FF
	return 0
}
