package addressbus

type BusAddressingExternal interface {
	Read(addr uint16) (byte, error)
	Write(addr uint16, data byte) error
}

type BusAddressingInternal interface {
	Read(addr uint16) byte
	Write(addr uint16, data byte)
}
