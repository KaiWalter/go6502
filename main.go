package main

import (
	"fmt"
	"log"

	"github.com/KaiWalter/go6502/mos6502"
)

func main() {
	fmt.Println("6502 emulation testing")

	var err error

	ram, err = RetrieveROM("6502_functional_test.bin")
	if err != nil {
		log.Fatalf("could not retrieve ROM: %v", err)
	}

	fmt.Printf("%x\n", ram[0x400])

	fmt.Printf("%x\n", mos6502.TryCallback(readMem, writeMem))
}
