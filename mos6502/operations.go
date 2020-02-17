package mos6502

func ADC() int {
	fetch()

	temp := uint16(A) + uint16(fetched) + GetFlagN(C)

	SetFlag(Z, (temp&0x00FF) == 0)

	if GetFlag(D) {
		if ((uint16(A) & 0xF) + (uint16(fetched) & 0xF) + GetFlagN(C)) > 9 {
			temp += 6
		}

		SetFlag(N, temp&0x80 != 0)

		SetFlag(V, (^(uint16(A)^uint16(fetched))&(uint16(A)^temp))&0x0080 != 0)

		if temp > 0x99 {
			temp += 96
		}
		SetFlag(C, temp > 0x99)
	} else {
		SetFlag(C, temp > 255)

		SetFlag(V, (^(uint16(A)^uint16(fetched))&(uint16(A)^temp))&0x0080 != 0)

		SetFlag(N, temp&0x80 != 0)
	}

	A = uint8(temp & 0x00FF)

	return 1
}

func AND() int {
	fetch()

	A = A & fetched

	SetFlag(Z, A == 0x00)
	SetFlag(N, A&0x80 != 0)

	return 0
}

func ASL() int {
	fetch()
	temp := uint16(fetched) << 1
	SetFlag(C, (temp&0xFF00) > 0)
	SetFlag(Z, (temp&0x00FF) == 0x00)
	SetFlag(N, temp&0x80 != 0)
	if opAddressMode == amIMP {
		A = uint8(temp & 0x00FF)
	} else {
		write(absoluteAddress, uint8(temp&0x00FF))
	}
	return 0
}

func BCC() int {
	cycles := 0

	if !GetFlag(C) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func BCS() int {
	cycles := 0

	if GetFlag(C) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func BEQ() int {
	cycles := 0

	if GetFlag(Z) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func BIT() int {
	fetch()

	temp := uint16(A & fetched)

	SetFlag(Z, (temp&0x00FF) == 0x00)
	SetFlag(N, fetched&(1<<7) != 0)
	SetFlag(V, fetched&(1<<6) != 0)

	return 0
}

func BMI() int {
	cycles := 0

	if GetFlag(N) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func BNE() int {
	cycles := 0

	if !GetFlag(Z) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func BPL() int {
	cycles := 0

	if !GetFlag(N) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func BRK() int {
	write(absoluteSP(), uint8((PC>>8)&0x00FF))
	SP--
	write(absoluteSP(), uint8(PC&0x00FF))
	SP--

	SetFlag(B, true)
	write(absoluteSP(), Status)
	SP--
	SetFlag(I, true)
	SetFlag(B, false)

	PC = uint16(read(0xFFFE)) | (uint16(read(0xFFFF)) << 8)

	return 0
}

func BVC() int {
	cycles := 0

	if !GetFlag(V) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func BVS() int {
	cycles := 0

	if GetFlag(V) {
		cycles++
		absoluteAddress = PC + relativeAddress

		if (absoluteAddress & 0xFF00) != (PC & 0xFF00) {
			cycles++
		}

		PC = absoluteAddress
	}

	return cycles
}

func CLC() int {
	SetFlag(C, false)
	return 0
}

func CLD() int {
	SetFlag(D, false)
	return 0
}

func CLI() int {
	SetFlag(I, false)
	return 0
}

func CLV() int {
	SetFlag(V, false)
	return 0
}

func CMP() int {
	fetch()
	temp := uint16(A) - uint16(fetched)
	SetFlag(C, A >= fetched)
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	return 1
}

func CPX() int {
	fetch()
	temp := uint16(X) - uint16(fetched)
	SetFlag(C, X >= fetched)
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	return 0
}

func CPY() int {
	fetch()
	temp := uint16(Y) - uint16(fetched)
	SetFlag(C, Y >= fetched)
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	return 0
}

func DEC() int {
	fetch()
	temp := uint16(fetched) - 1
	write(absoluteAddress, uint8(temp&0x00FF))
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	return 0
}

func DEX() int {
	X--
	SetFlag(Z, X == 0x00)
	SetFlag(N, X&0x80 != 0)
	return 0
}

func DEY() int {
	Y--
	SetFlag(Z, Y == 0x00)
	SetFlag(N, Y&0x80 != 0)
	return 0
}

func EOR() int {
	fetch()

	A = A ^ fetched

	SetFlag(Z, A == 0x00)
	SetFlag(N, A&0x80 != 0)

	return 1
}

func INC() int {
	fetch()
	temp := uint16(fetched) + 1
	write(absoluteAddress, uint8(temp&0x00FF))
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	return 0
}

func INX() int {
	X++
	SetFlag(Z, X == 0x00)
	SetFlag(N, X&0x80 != 0)
	return 0
}

func INY() int {
	Y++
	SetFlag(Z, Y == 0x00)
	SetFlag(N, Y&0x80 != 0)
	return 0
}

func JMP() int {
	PC = absoluteAddress
	return 0
}

func JSR() int {
	PC--

	write(absoluteSP(), uint8((PC>>8)&0x00FF))
	SP--
	write(absoluteSP(), uint8(PC&0x00FF))
	SP--

	PC = absoluteAddress
	return 0
}

func LDA() int {
	fetch()
	A = fetched
	SetFlag(Z, A == 0x00)
	SetFlag(N, A&0x80 != 0)
	return 0
}

func LDX() int {
	fetch()
	X = fetched
	SetFlag(Z, X == 0x00)
	SetFlag(N, X&0x80 != 0)
	return 0
}

func LDY() int {
	fetch()
	Y = fetched
	SetFlag(Z, Y == 0x00)
	SetFlag(N, Y&0x80 != 0)
	return 0
}

func LSR() int {
	fetch()
	SetFlag(C, fetched&0x0001 != 0)
	temp := uint16(fetched) >> 1
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	if opAddressMode == amIMP {
		A = uint8(temp & 0x00FF)
	} else {
		write(absoluteAddress, uint8(temp&0x00FF))
	}
	return 0
}

func NOP() int {
	switch opCode {
	case 0x1C:
	case 0x3C:
	case 0x5C:
	case 0x7C:
	case 0xDC:
	case 0xFC:
		return 1
		break
	}

	return 0
}

func ORA() int {
	fetch()
	A = A | fetched
	SetFlag(Z, A == 0x00)
	SetFlag(N, A&0x80 != 0)
	return 0
}

func PHA() int {
	write(absoluteSP(), A)
	SP--
	return 0
}

func PHP() int {
	write(absoluteSP(), Status|B|U)
	SetFlag(B, false)
	SP--
	return 0
}

func PLA() int {
	SP++
	A = read(absoluteSP())
	SetFlag(Z, A == 0x00)
	SetFlag(N, A&0x80 != 0)
	return 0
}

func PLP() int {
	SP++
	Status = read(absoluteSP())
	SetFlag(U, true)
	return 0
}

func ROL() int {
	fetch()
	temp := (uint16(fetched) << 1) | GetFlagN(C)
	SetFlag(C, temp&0xFF00 != 0)
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	if opAddressMode == amIMP {
		A = uint8(temp & 0x00FF)

	} else {
		write(absoluteAddress, uint8(temp&0x00FF))
	}
	return 0
}

func ROR() int {
	fetch()
	temp := (GetFlagN(C) << 7) | (uint16(fetched) >> 1)
	SetFlag(C, temp&0x01 != 0)
	SetFlag(Z, (temp&0x00FF) == 0x0000)
	SetFlag(N, temp&0x0080 != 0)
	if opAddressMode == amIMP {
		A = uint8(temp & 0x00FF)

	} else {
		write(absoluteAddress, uint8(temp&0x00FF))
	}
	return 0
}

func RTI() int {
	SP++
	Status = read(absoluteSP())
	Status &= ^B
	Status &= ^U

	SP++
	PC = uint16(read(absoluteSP()))
	SP++
	PC |= uint16(read(absoluteSP())) << 8
	return 0
}

func RTS() int {
	SP++
	PC = uint16(read(absoluteSP()))
	SP++
	PC |= uint16(read(absoluteSP())) << 8

	PC++
	return 0
}

func SBC() int {
	fetch()

	temp := uint16(A) - uint16(fetched) - (1 - GetFlagN(C))

	SetFlag(Z, (temp&0x00FF) == 0)
	SetFlag(V, (uint16(A)^temp)&0x0080 != 0 && (uint16(A)^uint16(fetched))&0x0080 != 0)

	if GetFlag(D) {
		if ((uint16(A) & 0xF) - (1 - GetFlagN(C))) < (uint16(fetched) & 0xF) {
			temp -= 6
		}

		if temp > 0x99 {
			temp -= 0x60
		}
	}
	SetFlag(C, temp < 0x100)

	A = uint8(temp & 0x00FF)

	return 1
}

func SEC() int {
	SetFlag(C, true)
	return 0
}

func SED() int {
	SetFlag(D, true)
	return 0
}

func SEI() int {
	SetFlag(I, true)
	return 0
}

func STA() int {
	write(absoluteAddress, A)
	return 0
}

func STX() int {
	write(absoluteAddress, X)
	return 0
}

func STY() int {
	write(absoluteAddress, Y)
	return 0
}

func TAX() int {
	X = A
	SetFlag(Z, X == 0x00)
	SetFlag(N, X&0x80 != 0)
	return 0
}

func TAY() int {
	Y = A
	SetFlag(Z, Y == 0x00)
	SetFlag(N, Y&0x80 != 0)
	return 0
}

func TSX() int {
	X = SP
	SetFlag(Z, X == 0x00)
	SetFlag(N, X&0x80 != 0)
	return 0
}

func TXA() int {
	A = X
	SetFlag(Z, A == 0x00)
	SetFlag(N, A&0x80 != 0)
	return 0
}

func TXS() int {
	SP = X
	return 0
}

func TYA() int {
	A = Y
	SetFlag(Z, A == 0x00)
	SetFlag(N, A&0x80 != 0)
	return 0
}

func XXX() int {
	return 0
}
