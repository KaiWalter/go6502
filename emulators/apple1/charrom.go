package apple1

import (
	"bufio"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	rom         [][]byte
	romInverted [][]byte
)

func loadCharacterRom(filename string, bInvert bool) ([][]byte, [][]byte, error) {

	romFile, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %v", err)
	}

	defer romFile.Close()

	buffer := make([]byte, 256*8)

	bufferReader := bufio.NewReader(romFile)

	_, err = bufferReader.Read(buffer)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	// feed into character map
	// flip/reverse bits from right-to-left to left-to-right
	charBuffer := make([][]byte, 256)
	charBufferInv := make([][]byte, 256)

	for char_index := range charBuffer {
		charBuffer[char_index] = make([]byte, 8)
		charBufferInv[char_index] = make([]byte, 8)
		for line_index := range charBuffer[char_index] {
			rom_byte := buffer[(char_index*8)+line_index]
			var converted_byte byte = 0
			if rom_byte != 0 {
				var fromMask byte = 0x80
				var toMask byte = 0x01
				for i := 0; i < 8; i++ {
					if rom_byte&fromMask == fromMask {
						converted_byte |= toMask
					}
					fromMask >>= 1
					toMask <<= 1
				}
			}

			charBuffer[char_index][line_index] = converted_byte
			charBufferInv[char_index][line_index] = ^converted_byte
		}
	}

	return charBuffer, charBufferInv, err
}

func renderCharacter(x byte, y byte, charno byte, bInvert bool) {

	var charmasks []byte
	if bInvert {
		charmasks = romInverted[charno]
	} else {
		charmasks = rom[charno]
	}

	scanline := int(y) * charHeight
	linepos := int(x) * charWidth

	rects_on := []sdl.Rect{}
	rects_off := []sdl.Rect{}

	for r := 0; r < charHeight; r++ {
		mask := charmasks[r]
		for c := charWidth; c > 0; c-- {
			rect := sdl.Rect{X: int32(linepos+c) * pixelSize, Y: int32(scanline+r) * pixelSize, W: pixelSize, H: pixelSize}
			if mask&1 == 1 {
				rects_on = append(rects_on, rect)
			} else {
				rects_off = append(rects_off, rect)
			}
			mask >>= 1
		}
	}

	if len(rects_on) > 0 {
		renderer.SetDrawColor(98, 143, 0, 255)
		renderer.FillRects(rects_on)

	}
	if len(rects_off) > 0 {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.FillRects(rects_off)
	}

	renderer.Present()
}
