package main

import (
	"fmt"

	"github.com/KaiWalter/go6502/mos6502"
)

func main() {
	fmt.Println("6502 emulation testing")
	fmt.Println(mos6502.Test())
}
