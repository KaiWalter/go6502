package mc6821

type Signal byte

const (
	Fall Signal = iota
	Rise
)

type InteruptSignal byte

const (
	NoSignal InteruptSignal = iota
	IRQ
	NMI
	BRK
)

type MC6821 struct {
	Name         string
	StartAddress uint16
	EndAddress   uint16
	initialized  bool

	ORA      byte   // Output register A
	IRA      byte   // Input register A
	DDRA     byte   // data direction register A             (Output=1, Input=0)
	DDRA_neg byte   // negative data direction register A    (Output=0, Input=1)
	CA1      Signal // control line A1
	CA2      Signal // control line A2

	CRA                        byte // control register A
	CRA_Bit0_EnableIRQA1       bool
	CRA_Bit1_CA1_PositiveTrans bool
	CRA_Bit2_WritePort         bool
	CRA_Bit3_EnableIRQA2       bool
	CRA_Bit3_PulseOutput       bool
	CRA_Bit3_CA2_set_high      bool
	CRA_Bit4_CA2_PositiveTrans bool
	CRA_Bit4_ManualOutput      bool
	CRA_Bit5_OutputMode        bool

	ORB      byte   // Output register B
	IRB      byte   // Input register B
	DDRB     byte   // data direction register B             (Output=1, Input=0)
	DDRB_neg byte   // negative data direction register B    (Output=0, Input=1)
	CB1      Signal // control line B1
	CB2      Signal // control line B2

	CRB                        byte // control register B
	CRB_Bit0_EnableIRQB1       bool
	CRB_Bit1_CB1_PositiveTrans bool
	CRB_Bit2_WritePort         bool
	CRB_Bit3_EnableIRQB2       bool
	CRB_Bit3_PulseOutput       bool
	CRB_Bit3_CB2_set_high      bool
	CRB_Bit4_CB2_PositiveTrans bool
	CRB_Bit4_ManualOutput      bool
	CRB_Bit5_OutputMode        bool

	sendInterrupt chan<- InteruptSignal
	sendOutputA   chan<- byte
	sendOutputB   chan<- byte

	receiveInputA <-chan byte
	receiveInputB <-chan byte
	receiveCA1    <-chan Signal
	receiveCA2    <-chan Signal
	receiveCB1    <-chan Signal
	receiveCB2    <-chan Signal
}

func (mc *MC6821) init() {

	mc.IRA = 0xFF
	mc.IRB = 0

	mc.CA1, mc.CA2 = Rise, Rise

	mc.CRA, mc.CRB, mc.ORA, mc.ORB = 0, 0, 0, 0

	mc.DDRA, mc.DDRB = 0, 0
	mc.DDRA_neg, mc.DDRB_neg = 0xFF, 0xFF

	mc.updateControlRegisters()

	mc.initialized = true
}

func (mc *MC6821) updateControlRegisters() {

	// section A -----------------------------------------
	mc.CRA_Bit0_EnableIRQA1 = (mc.CRA & 0x01) == 0x01
	mc.CRA_Bit1_CA1_PositiveTrans = (mc.CRA & 0x02) == 0x02
	mc.CRA_Bit2_WritePort = (mc.CRA & 0x04) == 0x04
	mc.CRA_Bit5_OutputMode = (mc.CRA & 0x20) == 0x20

	mc.CRA_Bit3_EnableIRQA2 = false
	mc.CRA_Bit3_PulseOutput = false
	mc.CRA_Bit3_CA2_set_high = false
	mc.CRA_Bit4_CA2_PositiveTrans = false
	mc.CRA_Bit4_ManualOutput = false

	if mc.CRA_Bit5_OutputMode {
		mc.CRA_Bit4_ManualOutput = (mc.CRA & 0x10) == 0x10
		if mc.CRA_Bit4_ManualOutput {
			mc.CRA_Bit3_CA2_set_high = (mc.CRA & 0x08) == 0x08
			if mc.CRA_Bit3_CA2_set_high {
				mc.CA2 = Rise
			} else {
				mc.CA2 = Fall
			}
		} else {
			mc.CRA_Bit3_PulseOutput = (mc.CRA & 0x08) == 0x08
		}
	} else {
		mc.CRA_Bit3_EnableIRQA2 = (mc.CRA & 0x08) == 0x08
		mc.CRA_Bit4_CA2_PositiveTrans = (mc.CRA & 0x10) == 0x10
	}

	// section B -----------------------------------------
	mc.CRB_Bit0_EnableIRQB1 = (mc.CRB & 0x01) == 0x01
	mc.CRB_Bit1_CB1_PositiveTrans = (mc.CRB & 0x02) == 0x02
	mc.CRB_Bit2_WritePort = (mc.CRB & 0x04) == 0x04
	mc.CRB_Bit5_OutputMode = (mc.CRB & 0x20) == 0x20

	mc.CRB_Bit3_EnableIRQB2 = false
	mc.CRB_Bit3_PulseOutput = false
	mc.CRB_Bit3_CB2_set_high = false
	mc.CRB_Bit4_CB2_PositiveTrans = false
	mc.CRB_Bit4_ManualOutput = false

	if mc.CRB_Bit5_OutputMode {
		mc.CRB_Bit4_ManualOutput = (mc.CRB & 0x10) == 0x10
		if mc.CRB_Bit4_ManualOutput {
			mc.CRB_Bit3_CB2_set_high = (mc.CRB & 0x08) == 0x08
			if mc.CRB_Bit3_CB2_set_high {
				mc.CB2 = Rise
			} else {
				mc.CB2 = Fall
			}
		} else {
			mc.CRB_Bit3_PulseOutput = (mc.CRB & 0x08) == 0x08
		}

	} else {
		mc.CRB_Bit3_EnableIRQB2 = (mc.CRB & 0x08) == 0x08
		mc.CRB_Bit4_CB2_PositiveTrans = (mc.CRB & 0x10) == 0x10
	}
}

func (mc *MC6821) updateIRQ() {
	if (mc.CRA_Bit0_EnableIRQA1 && (mc.CRA&0x80) == 0x80) ||
		(mc.CRA_Bit3_EnableIRQA2 && (mc.CRA&0x40) == 0x40) ||
		(mc.CRB_Bit0_EnableIRQB1 && (mc.CRB&0x80) == 0x80) ||
		(mc.CRB_Bit3_EnableIRQB2 && (mc.CRB&0x40) == 0x40) {

		// non blocking send
		select {
		case mc.sendInterrupt <- IRQ:
		default:
		}
	}
}

func (mc *MC6821) Read(addr uint16) byte {
	var reg = addr & 0x03
	var data byte = 0
	switch reg {

	case 0: // PA

		mc.CRA &= 0x3F // IRQ flags implicitly cleared by a read

		// mix input and output
		data |= mc.ORA & mc.DDRA
		data |= mc.IRA & mc.DDRA_neg

	case 1: // CRA
		data = mc.CRA

	case 2: // PB

		mc.CRB &= 0x3F // IRQ flags implicitly cleared by a read

		// mix input and output
		data |= mc.ORB & mc.DDRB
		data |= mc.IRB & mc.DDRB_neg

	case 3: // CRB
		data = mc.CRB
	}

	return data
}

func (mc *MC6821) Write(addr uint16, data byte) {
	var reg = addr & 0x03

	switch reg {
	case 0: // DDRA / PA
		if mc.CRA_Bit2_WritePort {
			mc.ORA = data // into output register A
			// mix input and output
			var bOut byte = 0
			bOut |= mc.ORA & mc.DDRA
			bOut |= mc.IRA & mc.DDRA_neg
			// non blocking send
			select {
			case mc.sendOutputA <- bOut:
			default:
			}
		} else {
			mc.DDRA = data // into data direction register A
			mc.DDRA_neg = ^data
		}

	case 1: // CRA
		mc.CRA = (mc.CRA & 0xC0) | (data & 0x3F) // do not change IRQ flags
		mc.updateControlRegisters()
		mc.updateIRQ()

	case 2: // DDRB / PB
		if mc.CRB_Bit2_WritePort {
			mc.ORB = data // into output register B
			// mix input and output
			var bOut byte = 0
			bOut |= mc.ORB & mc.DDRB
			bOut |= mc.IRB & mc.DDRB_neg
			// non blocking send
			select {
			case mc.sendOutputB <- bOut:
			default:
			}

			if mc.CRB_Bit5_OutputMode && !mc.CRB_Bit4_ManualOutput { // handshake on write mode
				mc.CB2 = Fall
				if mc.CRB_Bit3_PulseOutput {
					mc.CB2 = Rise
				}
			}
		} else {
			mc.DDRB = data // into data direction register B
			mc.DDRB_neg = ^data
		}

	case 3: // CRB
		mc.CRB = (mc.CRB & 0xC0) | (data & 0x3F) // do not change IRQ flags
		mc.updateControlRegisters()
		mc.updateIRQ()
	}
}

// input channel A handling

func (mc *MC6821) SetInputChannelA(ch <-chan byte) {
	if !mc.initialized {
		mc.init()
	}

	mc.receiveInputA = ch
	go mc.receiveFromInputA()
}

func (mc *MC6821) receiveFromInputA() {
	for b := range mc.receiveInputA {
		mc.IRA = b
	}
}

// input channel B handling

func (mc *MC6821) SetInputChannelB(ch <-chan byte) {
	if !mc.initialized {
		mc.init()
	}

	mc.receiveInputB = ch
	go mc.receiveFromInputB()
}

func (mc *MC6821) receiveFromInputB() {
	for b := range mc.receiveInputB {
		mc.IRB = b
	}
}

// output channel handling

func (mc *MC6821) SetOutputChannelA(ch chan<- byte) {
	if !mc.initialized {
		mc.init()
	}

	mc.sendOutputA = ch
}

func (mc *MC6821) SetOutputChannelB(ch chan<- byte) {
	if !mc.initialized {
		mc.init()
	}

	mc.sendOutputB = ch
}

func (mc *MC6821) SetInterruptChannelB(ch chan<- InteruptSignal) {
	if !mc.initialized {
		mc.init()
	}

	mc.sendInterrupt = ch
}

// control A1 handling

func (mc *MC6821) setCA1(b Signal) {

	var condition Signal

	if mc.CRA_Bit1_CA1_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if mc.CA1 != b && condition == b {
		mc.CRA |= 0x80 // set bit 7 IRQA1
		mc.updateIRQ()
		if mc.CRA_Bit5_OutputMode && !mc.CRA_Bit4_ManualOutput && !mc.CRA_Bit3_PulseOutput { // handshake mode
			mc.CA2 = Rise
		}
	}

	mc.CA1 = b
}

func (mc *MC6821) SetCA1Channel(ch <-chan Signal) {
	if !mc.initialized {
		mc.init()
	}

	mc.receiveCA1 = ch
	go mc.receiveFromCA1()
}

func (mc *MC6821) receiveFromCA1() {
	for s := range mc.receiveCA1 {
		mc.setCA1(s)
	}
}

// control A2 handling

func (mc *MC6821) setCA2(b Signal) {

	var condition Signal

	if mc.CRA_Bit4_CA2_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if mc.CA2 != b && condition == b {
		mc.CRA |= 0x40 // set bit 6 IRQA2
		mc.updateIRQ()
	}

	mc.CA2 = b
}

func (mc *MC6821) SetCA2Channel(ch <-chan Signal) {
	if !mc.initialized {
		mc.init()
	}

	mc.receiveCA2 = ch
	go mc.receiveFromCA2()
}

func (mc *MC6821) receiveFromCA2() {
	for s := range mc.receiveCA2 {
		mc.setCA2(s)
	}
}

// control B1 handling

func (mc *MC6821) setCB1(b Signal) {

	var condition Signal

	if mc.CRB_Bit1_CB1_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if mc.CB1 != b && condition == b {
		mc.CRB |= 0x80 // set bit 7 IRQB1
		mc.updateIRQ()
		if mc.CRB_Bit5_OutputMode && !mc.CRB_Bit4_ManualOutput && !mc.CRB_Bit3_PulseOutput { // handshake mode
			mc.CB2 = Rise
		}
	}

	mc.CB1 = b
}

func (mc *MC6821) SetCB1Channel(ch <-chan Signal) {
	if !mc.initialized {
		mc.init()
	}

	mc.receiveCB1 = ch
	go mc.receiveFromCB1()
}

func (mc *MC6821) receiveFromCB1() {
	for s := range mc.receiveCB1 {
		mc.setCB1(s)
	}
}

// control B2 handling

func (mc *MC6821) setCB2(b Signal) {

	var condition Signal

	if mc.CRB_Bit4_CB2_PositiveTrans {
		condition = Rise
	} else {
		condition = Fall
	}

	if mc.CB2 != b && condition == b {
		mc.CRB |= 0x40 // set bit 6 IRQB2
		mc.updateIRQ()
	}

	mc.CB2 = b
}

func (mc *MC6821) SetCB2Channel(ch <-chan Signal) {
	if !mc.initialized {
		mc.init()
	}

	mc.receiveCB2 = ch
	go mc.receiveFromCB2()
}

func (mc *MC6821) receiveFromCB2() {
	for s := range mc.receiveCB2 {
		mc.setCB2(s)
	}
}
