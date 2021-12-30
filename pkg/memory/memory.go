package memory

type Memory struct {
	AddressOffset uint16
	AddressSpace  []byte
}

func (r *Memory) Size() int {
	return len(r.AddressSpace)
}

func (r *Memory) Read(addr uint16) byte {
	return r.AddressSpace[addr-r.AddressOffset]
}

func (r *Memory) Write(addr uint16, data byte) {
	r.AddressSpace[addr-r.AddressOffset] = data
}
