package addressbus

import (
	"bufio"
	"fmt"
	"os"
)

func retrieveROM(filename string) ([]byte, error) {
	romFile, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	defer romFile.Close()

	stats, statsErr := romFile.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	buffer := make([]byte, stats.Size())

	bufferreader := bufio.NewReader(romFile)

	_, err = bufferreader.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return buffer, err
}
