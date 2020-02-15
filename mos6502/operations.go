package mos6502

func ADC() int {
	return 0
}

func AND() int {
	return 0
}

func ASL() int {
	return 0
}

func BCC() int {
	return 0
}

func BCS() int {
	return 0
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
	return 0
}

func BMI() int {
	return 0
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
	return 0
}

func BRK() int {
	return 0
}

func BVC() int {
	return 0
}

func BVS() int {
	return 0
}

func CLC() int {
	return 0
}

func CLD() int {
	return 0
}

func CLI() int {
	return 0
}

func CLV() int {
	return 0
}

func CMP() int {
	return 0
}

func CPX() int {
	return 0
}

func CPY() int {
	return 0
}

func DEC() int {
	return 0
}

func DEX() int {
	return 0
}

func DEY() int {
	return 0
}

func EOR() int {
	return 0
}

func INC() int {
	return 0
}

func INX() int {
	return 0
}

func INY() int {
	return 0
}

func JMP() int {
	PC = absoluteAddress
	return 0
}

func JSR() int {
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
	return 0
}

func LDY() int {
	return 0
}

func LSR() int {
	return 0
}

func NOP() int {
	return 0
}

func ORA() int {
	return 0
}

func PHA() int {
	return 0
}

func PHP() int {
	return 0
}

func PLA() int {
	return 0
}

func PLP() int {
	return 0
}

func ROL() int {
	return 0
}

func ROR() int {
	return 0
}

func RTI() int {
	return 0
}

func RTS() int {
	return 0
}

func SBC() int {
	return 0
}

func SEC() int {
	return 0
}

func SED() int {
	return 0
}

func SEI() int {
	return 0
}

func STA() int {
	write(absoluteAddress, A)
	return 0
}

func STX() int {
	return 0
}

func STY() int {
	return 0
}

func TAX() int {
	return 0
}

func TAY() int {
	return 0
}

func TSX() int {
	return 0
}

func TXA() int {
	return 0
}

func TXS() int {
	return 0
}

func TYA() int {
	return 0
}

func XXX() int {
	return 0
}
