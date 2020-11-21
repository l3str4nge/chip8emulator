package main


import "github.com/mateuszz0000/chip8emulator/chip8"


func main(){

	memory := chip8.LoadRomToMemory("c8games/PONG")
	memory = chip8.LoadFontsToMemory(memory)

	chip8.Emulator = &chip8.Chip8{
			memory, 
			0x200, 
			[16]byte{}, 
			0, 
			[16]uint16{}, 
			0, 
			nil, 
			0, 
			0,
			[64][32]byte{},
			[16]byte{}}

	chip8.Emulator.Run()
}