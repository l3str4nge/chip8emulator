package chip8

import "github.com/hajimehoshi/ebiten"
import "image/color"
import "fmt"


type Chip8 struct {
	Memory [4096]byte
	ProgramCounter uint16
	Registers [16]byte
	IndexRegister uint16
	Stack [16]uint16
	StackPointer byte 
	Screen *ebiten.Image
	Timer byte
	Pixels [64][32]byte

}

var Emulator (*Chip8) // global pointer to struct


func (chip8 *Chip8) Run() {
	ebiten.SetMaxTPS(600)
    if err := ebiten.Run(update, 64, 32, 20, "Hello world!"); err != nil {
        panic(err)
    }
}

func update(screen *ebiten.Image) error {
	Emulator.Screen = screen
	Emulator.runCycle()
	Emulator.render()
	return nil
}

func (chip8 *Chip8) runCycle(){
	opcode := chip8.getOpcode()

	fmt.Println("fetched op", fmt.Sprintf("%X", opcode))
	chip8.ProgramCounter += 2
	chip8.decodeAndRunInstruction(opcode)

	if chip8.Timer > 0 {
		chip8.Timer--
	}

}

func (chip8 *Chip8) render() {
	for x:=0 ; x < 64; x++ {
		for y:=0 ; y< 32; y++{
			if chip8.Pixels[x][y] == 0 {
				chip8.Screen.Set(x, y, color.Black)
			}else{
				chip8.Screen.Set(x, y, color.White)
			}
			
		}
	}
}




func (chip8 *Chip8) getOpcode() uint16{
	addr := chip8.ProgramCounter
	return uint16(chip8.Memory[addr]) << 8 | uint16(chip8.Memory[addr + 1])
}


func (chip8 *Chip8) decodeAndRunInstruction(opcode uint16){
	q := opcode >> 12
	opcodes[q](chip8, opcode)
}

