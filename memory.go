package main

var ram []byte

func readMem(addr uint16) uint8 {
	return ram[addr]
}

func writeMem(addr uint16, data uint8) {

	ram[addr] = data
}
