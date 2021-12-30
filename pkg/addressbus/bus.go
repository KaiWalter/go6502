package addressbus

type BusAddressing interface {
	Read(addr uint16) byte
	Write(addr uint16, data byte)
}

var (
	addressMap map[uint16]BusAddressing
	blockSize  int
	components = []BusAddressing{}
)

func InitBus(addressMapBlockSize int) {

	blockSize = addressMapBlockSize
	addressMap = make(map[uint16]BusAddressing)

}

func RegisterComponent(addressFrom int, addressTo int, component BusAddressing) {

	components = append(components, component)

	blockFrom := uint16(addressFrom / blockSize)
	blockTo := uint16(addressTo / blockSize)

	for b := blockFrom; b <= blockTo; b++ {
		addressMap[b] = component
	}

}

func Read(addr uint16) (byte, error) {

	block := addr / uint16(blockSize)

	component, exists := addressMap[block]
	if exists {
		return component.Read(addr), nil
	}

	return 0, &addressingError{Op: "Read", Address: addr}
}

func Write(addr uint16, data byte) error {

	block := addr / uint16(blockSize)

	component, exists := addressMap[block]
	if exists {
		component.Write(addr, data)
		return nil
	}

	return &addressingError{Op: "Write", Address: addr}
}
