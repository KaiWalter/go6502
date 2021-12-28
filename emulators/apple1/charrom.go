package apple1

import (
	"bufio"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	rom          [][]uint8
	rom_inverted [][]uint8
)

func loadRoms() {
	var err error

	// load character ROMs
	rom, rom_inverted, err = loadCharacterRom("./roms/Apple1_charmap.rom", false)
	if err != nil {
		panic(err)
	}

}

func loadCharacterRom(filename string, bInvert bool) ([][]uint8, [][]uint8, error) {

	romfile, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %v", err)
	}

	defer romfile.Close()

	buffer := make([]uint8, 256*8)

	bufferreader := bufio.NewReader(romfile)

	_, err = bufferreader.Read(buffer)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading file: %v", err)
	}

	// feed into character map
	// flip/reverse bits from right-to-left to left-to-right
	char_buffer := make([][]byte, 256)
	char_buffer_inv := make([][]byte, 256)

	for char_index := range char_buffer {
		char_buffer[char_index] = make([]byte, 8)
		char_buffer_inv[char_index] = make([]byte, 8)
		for line_index := range char_buffer[char_index] {
			rom_byte := buffer[(char_index*8)+line_index]
			var converted_byte uint8 = 0
			if rom_byte != 0 {
				var fromMask uint8 = 0x80
				var toMask uint8 = 0x01
				for i := 0; i < 8; i++ {
					if rom_byte&fromMask == fromMask {
						converted_byte |= toMask
					}
					fromMask >>= 1
					toMask <<= 1
				}
			}

			char_buffer[char_index][line_index] = converted_byte
			char_buffer_inv[char_index][line_index] = ^converted_byte
		}
	}

	return char_buffer, char_buffer_inv, err
}

func renderCharacter(x uint8, y uint8, charno uint8, bInvert bool) {

	var charmasks []uint8
	if bInvert {
		charmasks = rom_inverted[charno]
	} else {
		charmasks = rom[charno]
	}

	scanline := int(y) * nCharHeight
	linepos := int(x) * nCharWidth

	rects_on := []sdl.Rect{}
	rects_off := []sdl.Rect{}

	for r := 0; r < nCharHeight; r++ {
		mask := charmasks[r]
		for c := nCharWidth; c > 0; c-- {
			rect := sdl.Rect{X: int32(linepos+c) * nPixelSize, Y: int32(scanline+r) * nPixelSize, W: nPixelSize, H: nPixelSize}
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
