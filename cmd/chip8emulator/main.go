package main

import "os"
import "fmt"
//import "io"
import "bufio"


import (
	"github.com/hajimehoshi/ebiten"
	"math/rand"
	//"encoding/binary"
	"image/color"
	// "github.com/hajimehoshi/ebiten/i"
    //"github.com/hajimehoshi/ebiten/ebitenutil" // This is required to draw debug texts.
)
// ebiten

var BLACK = color.RGBA{255, 255, 255, 0}
//0xDAB6


type Chip8 struct {
	memory [4096]byte
	programCounter uint16
	registers [16]byte
	indexRegister uint16
	stack [16]uint16
	stackPointer byte 
	screen *ebiten.Image
	timer byte

}

var instructions = map[uint16] func(*Chip8, uint16){
	// 1st digit unique 
	0x0: (*Chip8).OP_0000,
	0x3: (*Chip8).OP_3XNN,
	0x6: (*Chip8).OP_6XNN,
	0xA: (*Chip8).OP_ANNN,
	0XD: (*Chip8).OP_DXYN,
	0x2: (*Chip8).OP_2NNN,
	0xF: (*Chip8).OP_FXXX,
	0x7: (*Chip8).OP_7XNN,
	0xC: (*Chip8).OP_CXNN,
	0xE: (*Chip8).OP_EXXX,
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

var emulator (*Chip8)

func (chip8 *Chip8) OP_EXXX(opcode uint16){
	fType := opcode & 0x00FF
	E_OPCODES[fType](chip8, opcode)
}

func (chip8 *Chip8) OP_EXA1(opcode uint16){
	fmt.Println("skip if key stored in vx isn;t pressed")
    // TODO: finish it
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

	chip8.drawSprite(int(x), int(y), int(nbytes), chip8.screen)
}

func (chip8 *Chip8) drawSprite(x,y, nbytes int, screen *ebiten.Image) {
	
	for col := 0; col < 8 ; col ++ {
		for row :=0; row < nbytes; row++ {
			value := chip8.memory[chip8.indexRegister + uint16(row)]
		// for row, value := range rows {
			pixValue := (value >> (7 - col)) & 1
			
			if pixValue == 0x1 {	
				pixelAt := screen.At(col + x, row + y)
				r, g, b, _ := pixelAt.RGBA()
				

				if pixelAt == BLACK {
					fmt.Println("COllision")
				}

				// XOR on each value of RGB
				r ^= 255
				g ^= 255
				b ^= 255
				screen.Set(col + x, row + y, color.RGBA{uint8(r), uint8(g), uint8(b), 0})
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
	return nil
}

func game() {
	ebiten.SetMaxTPS(60)
    if err := ebiten.Run(update, 64, 32, 20, "Hello world!"); err != nil {
        panic(err)
    }
}



func main(){
	var memory [4096]byte

	file, err := os.Open("c8games/pong.rom")

	if err != nil {
		panic(err)
	}
	defer file.Close()


	info, _ := file.Stat()

	var size int64 = info.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	
	// fmt.Println(buffer)
	// fmt.Println("size", size)
	// fmt.Println("BYTEEEEEEEEEEEEES", bytes)
	
	var programStartAddr int = 0x200
	fmt.Println("program starts at", programStartAddr)
	for i := 0; i <=  int(size) -1; i++ {
		memory[programStartAddr + i] = bytes[i]
	}

	fmt.Println("PROGRAM LOADS TO MEMORY:", memory)

	// var pc uint16
	for b := programStartAddr; b <= programStartAddr + int(size); b += 2 {
		//fmt.Println(memory[b], memory[b+1])
		opcode := uint16(memory[b]) << 8 | uint16(memory[b+1])
		hex := fmt.Sprintf("%X", opcode)
		fmt.Println("Opcode: ", hex)
	}

	pc1 := memory[0x200]
	pc2 := memory[0x200 + 1]
	opcode := uint16(pc1) << 8 | uint16(pc2)
	fmt.Println(pc1, pc2, pc1 + pc2, opcode)

	var dr uint16 = 0xDAB6
	cord := dr & 0x0FFF
	fmt.Println(cord)
	x := cord >> 8
	fmt.Println("x", x)

	y := (cord >> 4) & 0x0F
	fmt.Println("y", y)

	bb := cord & 0x00F
	fmt.Println("bytes", bb)

	// u := memory[0x2EA]

	for n := 0x2EA ; n <= 0x2EA + int(bb); n++ {
		fmt.Println(memory[n])
	}

	// gamesdl2()
	//Opcode:  DAB6 -> X: A, Y: B, 6bytes, 
	//Opcode:  DCD6  -> X: D, Y: 6
	/*
		memory [4096]byte
	programCounter uint16
	registers [16]byte
	indexRegister uint16
	stack [16]uint16
	stackPointer byte 

	*/

	// load fonts
	var fontset = [80]byte {
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80}  // F

	for x:=0; x < 80; x++ {
		memory[0x50 + x] = fontset[x]
	}


	emulator = &Chip8{memory, 0x200, [16]byte{}, 0, [16]uint16{}, 0, nil, 0}
	game()
	

}