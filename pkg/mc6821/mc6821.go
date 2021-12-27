package mc6821

type Signal uint8

const (
	Fall Signal = iota
	Rise
)

type InteruptSignal uint8

const (
	NoSignal InteruptSignal = iota
	IRQ
	NMI
	BRK
)

var (
	nORA      uint8  // Output register A
	nIRA      uint8  // Input register A
	nDDRA     uint8  // data direction register A             (Output=1, Input=0)
	nDDRA_neg uint8  // negative data direction register A    (Output=0, Input=1)
	nCA1      Signal // control line A1
	nCA2      Signal // control line A2

	nCRA                        uint8 // control register A
	bCRA_Bit0_EnableIRQA1       bool
	bCRA_Bit1_CA1_PositiveTrans bool
	bCRA_Bit2_WritePort         bool
	bCRA_Bit3_EnableIRQA2       bool
	bCRA_Bit3_PulseOutput       bool
	bCRA_Bit3_CA2_set_high      bool
	bCRA_Bit4_CA2_PositiveTrans bool
	bCRA_Bit4_ManualOutput      bool
	bCRA_Bit5_OutputMode        bool

	nORB      uint8  // Output register B
	nIRB      uint8  // Input register B
	nDDRB     uint8  // data direction register B             (Output=1, Input=0)
	nDDRB_neg uint8  // negative data direction register B    (Output=0, Input=1)
	nCB1      Signal // control line B1
	nCB2      Signal // control line B2

	nCRB                        uint8 // control register B
	bCRB_Bit0_EnableIRQB1       bool
	bCRB_Bit1_CB1_PositiveTrans bool
	bCRB_Bit2_WritePort         bool
	bCRB_Bit3_EnableIRQB2       bool
	bCRB_Bit3_PulseOutput       bool
	bCRB_Bit3_CB2_set_high      bool
	bCRB_Bit4_CB2_PositiveTrans bool
	bCRB_Bit4_ManualOutput      bool
	bCRB_Bit5_OutputMode        bool

	sendInterrupt chan<- InteruptSignal
	sendOutputA   chan<- uint8
	sendOutputB   chan<- uint8
)

func init() {
	nIRA = 0xFF
	nIRB = 0

	nCA1, nCA2 = Rise, Rise

	nCRA, nCRB, nORA, nORB = 0, 0, 0, 0

	nDDRA, nDDRB = 0, 0
	nDDRA_neg, nDDRB_neg = 0xFF, 0xFF

	updateControlRegisters()
}

func updateControlRegisters() {

	// section A -----------------------------------------
	bCRA_Bit0_EnableIRQA1 = (nCRA & 0x01) == 0x01
	bCRA_Bit1_CA1_PositiveTrans = (nCRA & 0x02) == 0x02
	bCRA_Bit2_WritePort = (nCRA & 0x04) == 0x04
	bCRA_Bit5_OutputMode = (nCRA & 0x20) == 0x20

	bCRA_Bit3_EnableIRQA2 = false
	bCRA_Bit3_PulseOutput = false
	bCRA_Bit3_CA2_set_high = false
	bCRA_Bit4_CA2_PositiveTrans = false
	bCRA_Bit4_ManualOutput = false

	if bCRA_Bit5_OutputMode {
		bCRA_Bit4_ManualOutput = (nCRA & 0x10) == 0x10
		if bCRA_Bit4_ManualOutput {
			bCRA_Bit3_CA2_set_high = (nCRA & 0x08) == 0x08
			if bCRA_Bit3_CA2_set_high {
				nCA2 = Rise
			} else {
				nCA2 = Fall
			}
		} else {
			bCRA_Bit3_PulseOutput = (nCRA & 0x08) == 0x08
		}
	} else {
		bCRA_Bit3_EnableIRQA2 = (nCRA & 0x08) == 0x08
		bCRA_Bit4_CA2_PositiveTrans = (nCRA & 0x10) == 0x10
	}

	// section B -----------------------------------------
	bCRB_Bit0_EnableIRQB1 = (nCRB & 0x01) == 0x01
	bCRB_Bit1_CB1_PositiveTrans = (nCRB & 0x02) == 0x02
	bCRB_Bit2_WritePort = (nCRB & 0x04) == 0x04
	bCRB_Bit5_OutputMode = (nCRB & 0x20) == 0x20

	bCRB_Bit3_EnableIRQB2 = false
	bCRB_Bit3_PulseOutput = false
	bCRB_Bit3_CB2_set_high = false
	bCRB_Bit4_CB2_PositiveTrans = false
	bCRB_Bit4_ManualOutput = false

	if bCRB_Bit5_OutputMode {
		bCRB_Bit4_ManualOutput = (nCRB & 0x10) == 0x10
		if bCRB_Bit4_ManualOutput {
			bCRB_Bit3_CB2_set_high = (nCRB & 0x08) == 0x08
			if bCRB_Bit3_CB2_set_high {
				nCB2 = Rise
			} else {
				nCB2 = Fall
			}
		} else {
			bCRB_Bit3_PulseOutput = (nCRB & 0x08) == 0x08
		}

	} else {
		bCRB_Bit3_EnableIRQB2 = (nCRB & 0x08) == 0x08
		bCRB_Bit4_CB2_PositiveTrans = (nCRB & 0x10) == 0x10
	}
}

func updateIRQ() {
	if (bCRA_Bit0_EnableIRQA1 && (nCRA&0x80) == 0x80) ||
		(bCRA_Bit3_EnableIRQA2 && (nCRA&0x40) == 0x40) ||
		(bCRB_Bit0_EnableIRQB1 && (nCRB&0x80) == 0x80) ||
		(bCRB_Bit3_EnableIRQB2 && (nCRB&0x40) == 0x40) {

		// non blocking send
		select {
		case sendInterrupt <- IRQ:
		default:
		}
	}
}

func CpuRead(addr uint16) uint8 {
	var reg = addr & 0x03
	var data uint8 = 0
	switch reg {

	case 0: // PA

		nCRA &= 0x3F // IRQ flags implicitly cleared by a read

		// mix input and output
		data |= nORA & nDDRA
		data |= nIRA & nDDRA_neg

	case 1: // CRA
		data = nCRA

	case 2: // PB

		nCRB &= 0x3F // IRQ flags implicitly cleared by a read

		// mix input and output
		data |= nORB & nDDRB
		data |= nIRB & nDDRB_neg

	case 3: // CRB
		data = nCRB
	}

	return data
}

func CpuWrite(addr uint16, data uint8) {
	var reg = addr & 0x03

	switch reg {
	case 0: // DDRA / PA
		if bCRA_Bit2_WritePort {
			nORA = data // into output register A
			// mix input and output
			var bOut uint8 = 0
			bOut |= nORA & nDDRA
			bOut |= nIRA & nDDRA_neg
			// non blocking send
			select {
			case sendOutputA <- bOut:
			default:
			}
		} else {
			nDDRA = data // into data direction register A
			nDDRA_neg = ^data
		}

	case 1: // CRA
		nCRA = (nCRA & 0xC0) | (data & 0x3F) // do not change IRQ flags
		updateControlRegisters()
		updateIRQ()

	case 2: // DDRB / PB
		if bCRB_Bit2_WritePort {
			nORB = data // into output register B
			// mix input and output
			var bOut uint8 = 0
			bOut |= nORB & nDDRB
			bOut |= nIRB & nDDRB_neg
			// non blocking send
			select {
			case sendOutputB <- bOut:
			default:
			}

			if bCRB_Bit5_OutputMode && !bCRB_Bit4_ManualOutput { // handshake on write mode
				nCB2 = Fall
				if bCRB_Bit3_PulseOutput {
					nCB2 = Rise
				}
			}
		} else {
			nDDRB = data // into data direction register B
			nDDRB_neg = ^data
		}

	case 3: // CRB
		nCRB = (nCRB & 0xC0) | (data & 0x3F) // do not change IRQ flags
		updateControlRegisters()
		updateIRQ()
	}
}

func SetInputA(b uint8) {
	nIRA = b
}

func SetInputB(b uint8) {
	nIRB = b
}

func SetOutputChannelA(ch chan<- uint8) {
	sendOutputA = ch
}

func SetOutputChannelB(ch chan<- uint8) {
	sendOutputB = ch
}

func SetInterruptChannelB(ch chan<- InteruptSignal) {
	sendInterrupt = ch
}

func SetCA1(b Signal) {

	var condition Signal

	if bCRA_Bit1_CA1_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if nCA1 != b && condition == b {
		nCRA |= 0x80 // set bit 7 IRQA1
		updateIRQ()
		if bCRA_Bit5_OutputMode && !bCRA_Bit4_ManualOutput && !bCRA_Bit3_PulseOutput { // handshake mode
			nCA2 = Rise
		}
	}

	nCA1 = b
}

func GetCA1() Signal {
	return nCA1
}

func SetCA2(b Signal) {

	var condition Signal

	if bCRA_Bit4_CA2_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if nCA2 != b && condition == b {
		nCRA |= 0x40 // set bit 6 IRQA2
		updateIRQ()
	}

	nCA2 = b
}

func GetCA2() Signal {
	return nCA2
}

func SetCB1(b Signal) {

	var condition Signal

	if bCRB_Bit1_CB1_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if nCB1 != b && condition == b {
		nCRB |= 0x80 // set bit 7 IRQB1
		updateIRQ()
		if bCRB_Bit5_OutputMode && !bCRB_Bit4_ManualOutput && !bCRB_Bit3_PulseOutput { // handshake mode
			nCB2 = Rise
		}
	}

	nCB1 = b
}

func GetCB1() Signal {
	return nCB1
}

func SetCB2(b Signal) {

	var condition Signal

	if bCRB_Bit4_CB2_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if nCB2 != b && condition == b {
		nCRB |= 0x40 // set bit 6 IRQB2
		updateIRQ()
	}

	nCB2 = b
}

func GetCB2() Signal {
	return nCB2
}
