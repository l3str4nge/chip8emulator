package chip8

import "github.com/hajimehoshi/ebiten"
import "image/color"
import "fmt"


var WIDTH = 64
var HEIGHT = 32


type Chip8 struct {
	Memory [4096]byte
	ProgramCounter uint16
	Registers [16]byte
	IndexRegister uint16
	Stack [16]uint16
	StackPointer byte 
	Screen *ebiten.Image
	Timer byte
	SoundTimer byte
	Pixels [64][32]byte
	Keyboard [16]byte
	ClearScreenFlag bool
	Scale int
}

var Emulator (*Chip8) // global pointer to struct

var square *ebiten.Image


func (chip8 *Chip8) Run() {
	square, _ = ebiten.NewImage(1 * chip8.Scale, 1 * chip8.Scale, ebiten.FilterNearest)
	square.Fill(color.White)
	fmt.Println("lol", square)	
    if err := ebiten.Run(update, WIDTH * chip8.Scale, HEIGHT * chip8.Scale, 1, "Hello world!"); err != nil {
        panic(err)
    }
}

func update(screen *ebiten.Image) error {
	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})
	for tick :=0; tick <10; tick++ {
		Emulator.runCycle()
		Emulator.checkClearScreen(screen)
		Emulator.render(screen)
		CheckKeyboard(Emulator)
}
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

	if chip8.SoundTimer > 0 {
		chip8.SoundTimer--
		// TODO: beep buzzer
	}
}

func (chip8 *Chip8) render(screen *ebiten.Image) {
	for x:=0 ; x < WIDTH; x++ {
		for y:=0 ; y< HEIGHT; y++{
			if chip8.Pixels[x][y] == 0xFF{ 
				
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x * chip8.Scale), float64(y*chip8.Scale))
				screen.DrawImage(square, opts)
			}
			
		}
	}
}

func (chip8 *Chip8) checkClearScreen(screen *ebiten.Image) {
	if chip8.ClearScreenFlag{
		for x:=0 ; x < WIDTH; x++ {
			for y:=0 ; y< HEIGHT; y++{
				chip8.Pixels[x][y] = 0x00
			}
		}

		screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})
		chip8.ClearScreenFlag = false
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

