package mos6502

const (
	amABS = iota
	amABX
	amABY
	amIND
	amIMM
	amIMP
	amIZX
	amIZY
	amREL
	amZP0
	amZPX
	amZPY
)

func ABS() int {
	opAddressMode = amABS

	lo := uint16(read(PC))
	PC++
	hi := uint16(read(PC))
	PC++

	absoluteAddress = hi<<8 | lo
	return 0
}

func ABX() int {
	opAddressMode = amABX

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
	opAddressMode = amABY

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

func IND() int {
	opAddressMode = amIND

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

func IMM() int {
	opAddressMode = amIMM
	absoluteAddress = PC
	PC++
	return 0
}

func IMP() int {
	opAddressMode = amIMP
	fetched = A
	return 0
}

func IZX() int {
	opAddressMode = amIZX

	tempAddress := uint16(read(PC))
	PC++

	lo := uint16(read(tempAddress & 0x00FF))
	hi := uint16(read((tempAddress + 1) & 0x00FF))

	absoluteAddress = hi<<8 | lo
	absoluteAddress += uint16(X)

	if (absoluteAddress & 0xFF00) != (hi << 8) {
		return 1
	}
	return 0
}

func IZY() int {
	opAddressMode = amIZY

	tempAddress := uint16(read(PC))
	PC++

	lo := uint16(read(tempAddress & 0x00FF))
	hi := uint16(read((tempAddress + 1) & 0x00FF))

	absoluteAddress = hi<<8 | lo
	absoluteAddress += uint16(Y)

	if (absoluteAddress & 0xFF00) != (hi << 8) {
		return 1
	}
	return 0
}

func REL() int {
	opAddressMode = amREL

	relativeAddress = uint16(read(PC))
	PC++
	if relativeAddress&0x80 != 0 {
		relativeAddress |= 0xFF00
	}
	return 0
}

func ZP0() int {
	opAddressMode = amZP0

	absoluteAddress = uint16(read(PC))
	PC++
	absoluteAddress &= 0x00FF
	return 0
}

func ZPX() int {
	opAddressMode = amZPX

	absoluteAddress = uint16(read(PC)) + uint16(X)
	PC++
	absoluteAddress &= 0x00FF
	return 0
}

func ZPY() int {
	opAddressMode = amZPY

	absoluteAddress = uint16(read(PC)) + uint16(Y)
	PC++
	absoluteAddress &= 0x00FF
	return 0
}
