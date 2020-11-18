package chip8

import "github.com/hajimehoshi/ebiten"

type Chip8 struct {
	Memory [4096]byte
	programCounter uint16
	Registers [16]byte
	indexRegister uint16
	stack [16]uint16
	stackPointer byte 
	screen *ebiten.Image
	timer byte
	pixels [64][32]byte

}

func (chip8 *Chip8) RunGame() {
   //TODO: finish it
}

func (chip8 *Chip8) Render() {

}



var chip8 (*Chip8)