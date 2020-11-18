package main

// import "os"
import "fmt"
//import "io"
// import "bufio"


import (
	"github.com/hajimehoshi/ebiten"
	"math/rand"
	//"encoding/binary"
	"image/color"
	"github.com/mateuszz0000/chip8emulator/chip8"
	// "github.com/hajimehoshi/ebiten/i"
    //"github.com/hajimehoshi/ebiten/ebitenutil" // This is required to draw debug texts.
)



type Chip8 struct {
	memory [4096]byte
	programCounter uint16
	registers [16]byte
	indexRegister uint16
	stack [16]uint16
	stackPointer byte 
	screen *ebiten.Image
	timer byte
	pixels [64][32]byte

}

var instructions = map[uint16] func(*Chip8, uint16){
	// 1st digit unique 
	0x0: (*Chip8).OP_0000,
	0x1: (*Chip8).OP_1NNN,
	0x3: (*Chip8).OP_3XNN,
	0x4: (*Chip8).OP_4XNN,
	0x6: (*Chip8).OP_6XNN,
	0xA: (*Chip8).OP_ANNN,
	0XD: (*Chip8).OP_DXYN,
	0x2: (*Chip8).OP_2NNN,
	0xF: (*Chip8).OP_FXXX,
	0x7: (*Chip8).OP_7XNN,
	0xC: (*Chip8).OP_CXNN,
	0xE: (*Chip8).OP_EXXX,
	0x8: (*Chip8).OP_8XXX,
}

var NULL_OPCODES = map[uint16] func(*Chip8, uint16){
	0x00E0: (*Chip8).OP_0NNN,
	0x00EE: (*Chip8).OP_00EE,
}

var F_OPCODES = map[uint16] func(*Chip8, uint16){
	0x33: (*Chip8).OP_FX33,
	0x65: (*Chip8).OP_FX65,
	0x29: (*Chip8).OP_FX29,
	0x15: (*Chip8).OP_FX15,
	0x07: (*Chip8).OP_FX07,
}

var E_OPCODES = map[uint16] func(*Chip8, uint16){
	0xA1: (*Chip8).OP_EXA1,
}

var OPCODES_8 = map[uint16] func(*Chip8, uint16){
	0x2: (*Chip8).OP_8XY2,
	0x4: (*Chip8).OP_8XY4,
}

func (chip8 *Chip8) OP_8XXX(opcode uint16){
	_type := opcode & 0x000F
	OPCODES_8[_type](chip8, opcode)
}

func (chip8 *Chip8) OP_8XY2(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F
	chip8.registers[vx] &= chip8.registers[vy]

}

func (chip8 *Chip8) OP_8XY4(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F

	vyvalue := chip8.registers[vy]

	sum := chip8.registers[vx] + vyvalue

	if sum > 255 {
		chip8.registers[0xF] = 1
	} else {
		chip8.registers[0xF] = 0
	}

	chip8.registers[vx] = sum & 0xFF
}

var emulator *chip8.Chip8



func (chip8 *Chip8) OP_EXXX(opcode uint16){
	fType := opcode & 0x00FF
	E_OPCODES[fType](chip8, opcode)
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
}

func (chip8 *Chip8) OP_EXA1(opcode uint16){
	fmt.Println("skip if key stored in vx isn;t pressed")
	// TODO: finish it
	//chip8.programCounter += 2
}

func (chip8 *Chip8) OP_3XNN(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF

	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "skipping instruction if vx ==nn")	
	if chip8.registers[vx] == byte(value) {
		fmt.Println("skipped")
		chip8.programCounter += 2
	}
}

func (chip8 *Chip8) OP_4XNN(opcode uint16){
	nn := opcode & 0x00FF
	vx := (opcode >> 8) & 0x0F

	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
	if chip8.registers[vx] != byte(nn) {
		chip8.programCounter += 2
		fmt.Println("SKIPPED")
	}
}

func (chip8 *Chip8) OP_CXNN(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF

	chip8.registers[vx] = byte(rand.Intn(0xFF)) & byte(value)
}

func (chip8 *Chip8) OP_FX07(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	chip8.registers[vx] = chip8.timer
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "set timer to Vx")
}

func (chip8 *Chip8) runCycle(){
	opcode := chip8.getOpcode()

	fmt.Println("fetched op", fmt.Sprintf("%X", opcode))
	chip8.decodeAndRunInstruction(opcode)
	chip8.programCounter += 2

	if chip8.timer > 0 {
		chip8.timer--
	}

}

func (chip8 *Chip8) getOpcode() uint16{
	addr := chip8.programCounter
	return uint16(chip8.memory[addr]) << 8 | uint16(chip8.memory[addr + 1])
}

func (chip8 *Chip8) decodeAndRunInstruction(opcode uint16){
	q := opcode >> 12
	instructions[q](chip8, opcode)
}

func (chip8 *Chip8) OP_0000(opcode uint16){
	fType := opcode & 0x00FF
	NULL_OPCODES[fType](chip8, opcode)
}

func (chip8 *Chip8) OP_1NNN(opcode uint16){
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
	chip8.programCounter = opcode & 0x0FFF
}

func (chip8 *Chip8) OP_00EE(opcode uint16){
	chip8.stackPointer--
	chip8.programCounter = chip8.stack[chip8.stackPointer]
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "RETURN FROM SUBRUTINE")

}

func (chip8 *Chip8) OP_FXXX(opcode uint16) {
	fType := opcode & 0x00FF
	F_OPCODES[fType](chip8, opcode)
}

func (chip8 *Chip8) OP_FX33(opcode uint16){
	registerNr := (opcode >> 8) & 0x0F
	value := chip8.registers[registerNr]

	hundreds := byte((value / 100) % 10)
	tens := byte((value / 10) % 10)
	ones := byte(value % 10)

	chip8.memory[chip8.indexRegister] = hundreds
	chip8.memory[chip8.indexRegister + 1] = tens
	chip8.memory[chip8.indexRegister + 2] = ones

	fmt.Println("FX33, opcode", fmt.Sprintf("%X", opcode), "val", value, hundreds, tens, ones, registerNr)
}

func (chip8 *Chip8) OP_FX15(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	chip8.timer = chip8.registers[vx]
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "timer set")

}

func (chip8 *Chip8) OP_FX65(opcode uint16){
	registerThreshold := (opcode >> 8) & 0x0F
	for nr := uint16(0); nr <= registerThreshold; nr++ {
		chip8.registers[nr] = chip8.memory[chip8.indexRegister + nr]
	}
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "maxREG", registerThreshold)
}

func (chip8 *Chip8) OP_FX29(opcode uint16){
	vx := (opcode >> 8) & 0x0F

	// TODO: make offset to get location of fonts loaded
	fontRawAddr := chip8.registers[vx]
	fontAdr := (5 * fontRawAddr) + 0x50 // each font has 5 bytes long
	chip8.indexRegister = uint16(fontAdr)

	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
	fmt.Println("I:", fmt.Sprintf("%X", chip8.indexRegister))

}

func (chip8 *Chip8) OP_7XNN(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF
	chip8.registers[vx] += byte(value)
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
}

func (chip8 *Chip8) OP_0NNN(opcode uint16){
	/// clear the screen

}

func (chip8 *Chip8) OP_6XNN(opcode uint16){
	/// Sets Vx to NN
	registerNr := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode),
	 "REG:", registerNr, "VALUE", value)
	chip8.registers[registerNr] = byte(value)
}

func (chip8 *Chip8) OP_ANNN(opcode uint16){
	// Sets NNN to I register
	chip8.indexRegister = opcode & 0x0FFF
	fmt.Println("index register value", fmt.Sprintf("%X", chip8.indexRegister))
}

func (chip8 *Chip8) OP_DXYN(opcode uint16){
	x := (opcode >> 8) & 0x0F
	y := (opcode >> 4)& 0x00F
	nbytes := opcode & 0x000F

	fmt.Println("opcode", 
		fmt.Sprintf("%X", opcode), 
		"X:", fmt.Sprintf("%X", x), 
		"Y:", fmt.Sprintf("%X", y),
		"nbytes", nbytes)

	chip8.drawSprite(int(chip8.registers[x]), int(chip8.registers[y]), int(nbytes), chip8.screen)
}

func (chip8 *Chip8) drawSprite(x,y, nbytes int, screen *ebiten.Image) {
	fmt.Println("XY", x, y)
	for col := 0; col < 8 ; col ++ {
		for row :=0; row < nbytes; row++ {
			value := chip8.memory[chip8.indexRegister + uint16(row)]
			pixValue := (value >> (7 - col)) & 1

			fmt.Println(col, row)
			xx := col +x
			yy := row + y
			if xx > 63 {
				xx = xx - 63
			}

			if yy > 31 {
				yy = yy - 31
			}
			
			pixelAt := chip8.pixels[xx][yy]
			if pixValue == 0x1 {	

				if pixelAt == 0x1 {
					fmt.Println("COllision")
					chip8.registers[0xF] = 1
				} else {
					chip8.registers[0xF] = 0
				}

				chip8.pixels[xx][yy] = pixelAt ^ 0xFF
			}			
		}
	}
}


func (chip8 *Chip8) OP_2NNN(opcode uint16){
	// call subroutine
	subroutineAddr := opcode & 0x0FFF
	
	chip8.stack[chip8.stackPointer] = chip8.programCounter
	chip8.stackPointer++
	chip8.programCounter = subroutineAddr
	fmt.Println("Called subroutine at: ", fmt.Sprintf("%X", subroutineAddr))
}

func update(screen *ebiten.Image) error {
	emulator.screen = screen
	emulator.runCycle()

	emulator.updateScreen()

	return nil
}

func (chip8 *Chip8) updateScreen(){
	for x:=0 ; x < 64; x++ {
		for y:=0 ; y< 32; y++{
			if chip8.pixels[x][y] == 0 {
				chip8.screen.Set(x, y, color.Black)
			}else{
				chip8.screen.Set(x, y, color.White)
			}
			
		}
	}
	
}

func game() {
	ebiten.SetMaxTPS(60)
    if err := ebiten.Run(update, 64, 32, 20, "Hello world!"); err != nil {
        panic(err)
    }
}



func main(){

	memory := chip8.LoadRomToMemory("c8games/pong.rom")
	chip8.LoadFontsToMemory(memory)


	emulator = &chip8.Chip8{memory, 0x200, [16]byte{}, 0, [16]uint16{}, 0, nil, 0, [64][32]byte{}}
	game()
	

}