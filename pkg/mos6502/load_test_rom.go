package mos6502

import (
	"bufio"
	"fmt"
	"os"
)

// RetrieveROM retrieves contents of a file into memory
// https://github.com/Klaus2m5/6502_65C02_functional_tests/blob/master/bin_files/6502_functional_test.lst
func RetrieveROM(filename string) ([]byte, error) {
	romfile, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer romfile.Close()

	buffer := make([]byte, 0x10000)

	bufferreader := bufio.NewReader(romfile)

	_, err = bufferreader.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return buffer, err
}
