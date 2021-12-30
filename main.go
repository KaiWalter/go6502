package main

import (
	"log"
	"os"

	"github.com/KaiWalter/go6502/emulators/apple1"
)

func main() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	apple1.Run()
}
