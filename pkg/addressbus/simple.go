// simple address bus implementation e.g. for unit tests

package addressbus

type SimpleBus struct {
	memory BusAddressingInternal
}

func (b *SimpleBus) InitBus(mem BusAddressingInternal) {

	b.memory = mem

}

func (b *SimpleBus) Read(addr uint16) (byte, error) {

	return b.memory.Read(addr), nil

}

func (b *SimpleBus) Write(addr uint16, data byte) error {

	b.memory.Write(addr, data)

	return nil

}
