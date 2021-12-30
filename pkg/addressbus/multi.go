// address bus implementation with multiple components (RAM,ROM,PIA) attached

package addressbus

type MultiBus struct {
	addressMap map[uint16]BusAddressingInternal
	blockSize  int
	components []BusAddressingInternal
}

func (b *MultiBus) InitBus(addressMapBlockSize int) {

	b.components = make([]BusAddressingInternal, 0)

	b.blockSize = addressMapBlockSize
	b.addressMap = make(map[uint16]BusAddressingInternal)

}

func (b *MultiBus) RegisterComponent(addressFrom int, addressTo int, component BusAddressingInternal) {

	b.components = append(b.components, component)

	blockFrom := uint16(addressFrom / b.blockSize)
	blockTo := uint16(addressTo / b.blockSize)

	for block := blockFrom; block <= blockTo; block++ {
		b.addressMap[block] = component
	}

}

func (b *MultiBus) Read(addr uint16) (byte, error) {

	block := addr / uint16(b.blockSize)

	component, exists := b.addressMap[block]
	if exists {
		return component.Read(addr), nil
	}

	return 0, &addressingError{Op: "Read", Address: addr}
}

func (b *MultiBus) Write(addr uint16, data byte) error {

	block := addr / uint16(b.blockSize)

	component, exists := b.addressMap[block]
	if exists {
		component.Write(addr, data)
		return nil
	}

	return &addressingError{Op: "Write", Address: addr}
}
