# A 6502 emulator in Go

This project is just for me to learn Go by solving a familar problem - 6502 emulation. I'll will start with most simpliest form and then try to move my way up in the more sophisticated language constructs and toolings.

> developed on Windows 11 + WSL2g/Ubuntu + SDL2
> <br/>DISCLAIMER: `.devcontainer` is only suited for the pure 6502 emulation tests, not the Apple 1 emulation

## key features

- **Apple 1 emulation** converted from <https://github.com/KaiWalter/olcApple1> - still WORK IN PROGRESS
- [Functional tests](https://github.com/Klaus2m5/6502_65C02_functional_tests) implemented with the `testing` framework, run with e.g. `go test ./pkg/mos6502/ -test.v`

## open issues

- [ ] work on some basic cycle accuracy for the Apple 1 emulation
- [ ] rework variable naming conventions <https://talks.golang.org/2014/names.slide#5>

## installation

- follow instructions for SDL2 on <https://github.com/veandco/go-sdl2#requirements>