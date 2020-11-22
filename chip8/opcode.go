package chip8


import "fmt"
import "math/rand"
import "time"

// opcodes mapping to function

var opcodes = map[uint16] func(*Chip8, uint16){
	// 1st digit unique 
	0x0: (*Chip8).OP_0000,
	0x1: (*Chip8).OP_1NNN,
	0x2: (*Chip8).OP_2NNN,
	0x3: (*Chip8).OP_3XNN,
	0x4: (*Chip8).OP_4XNN,
	0x5: (*Chip8).OP_5XY0,
	0x6: (*Chip8).OP_6XNN,
	0x7: (*Chip8).OP_7XNN,
	0x8: (*Chip8).OP_8XXX,
	0x9: (*Chip8).OP_9XY0,
	0xA: (*Chip8).OP_ANNN,
	0xB: (*Chip8).OP_BNNN,
	0xC: (*Chip8).OP_CXNN,
	0XD: (*Chip8).OP_DXYN,
	0xE: (*Chip8).OP_EXXX,
	0xF: (*Chip8).OP_FXXX,
}

var NULL_OPCODES = map[uint16] func(*Chip8, uint16){
	0x00E0: (*Chip8).OP_0NNN,
	0x00EE: (*Chip8).OP_00EE,
}

var F_OPCODES = map[uint16] func(*Chip8, uint16){
	0x07: (*Chip8).OP_FX07,
	0x0A: (*Chip8).OP_FX0A,
	0x15: (*Chip8).OP_FX15,
	0x18: (*Chip8).OP_FX18,
	0x1E: (*Chip8).OP_FX1E,
	0x29: (*Chip8).OP_FX29,
	0x33: (*Chip8).OP_FX33,
	0x55: (*Chip8).OP_FX55,
	0x65: (*Chip8).OP_FX65,	
}

var E_OPCODES = map[uint16] func(*Chip8, uint16){
	0x9E: (*Chip8).OP_EX9E,
	0xA1: (*Chip8).OP_EXA1,
}

var OPCODES_8 = map[uint16] func(*Chip8, uint16){
	0x0: (*Chip8).OP_8XY0,
	0x1: (*Chip8).OP_8XY1,
	0x2: (*Chip8).OP_8XY2,
	0x3: (*Chip8).OP_8XY3,
	0x4: (*Chip8).OP_8XY4,
	0x5: (*Chip8).OP_8XY5,
	0x6: (*Chip8).OP_8XY6,
	0x7: (*Chip8).OP_8XY7,
	0xE: (*Chip8).OP_8XYE,
}

// opcodes implementation
func (chip8 *Chip8) OP_0000(opcode uint16){
	fType := opcode & 0x00FF
	NULL_OPCODES[fType](chip8, opcode)
}

func (chip8 *Chip8) OP_0NNN(opcode uint16){
	chip8.ClearScreenFlag = true
}

func (chip8 *Chip8) OP_00EE(opcode uint16){
	chip8.StackPointer--
	chip8.ProgramCounter = chip8.Stack[chip8.StackPointer]
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "RETURN FROM SUBRUTINE")

}

func (chip8 *Chip8) OP_1NNN(opcode uint16){
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
	chip8.ProgramCounter = opcode & 0x0FFF
}

func (chip8 *Chip8) OP_2NNN(opcode uint16){
	// call subroutine
	subroutineAddr := opcode & 0x0FFF
	
	chip8.Stack[chip8.StackPointer] = chip8.ProgramCounter
	chip8.StackPointer++
	chip8.ProgramCounter = subroutineAddr
	fmt.Println("Called subroutine at: ", fmt.Sprintf("%X", subroutineAddr))
}

func (chip8 *Chip8) OP_3XNN(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF

	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "skipping instruction if vx ==nn")	
	if chip8.Registers[vx] == byte(value) {
		fmt.Println("skipped")
		chip8.ProgramCounter += 2
	} else {
		fmt.Println("cannot skip", chip8.Registers[vx], byte(value))
	}
}

func (chip8 *Chip8) OP_4XNN(opcode uint16){
	nn := opcode & 0x00FF
	vx := (opcode >> 8) & 0x0F

	if chip8.Registers[vx] != byte(nn) {
		chip8.ProgramCounter += 2
	}
}

func (chip8 *Chip8) OP_5XY0(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x0F
	if chip8.Registers[vx] == chip8.Registers[vy]{
		chip8.ProgramCounter += 2
	}
}

func (chip8 *Chip8) OP_6XNN(opcode uint16){
	/// Sets Vx to NN
	registerNr := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode),
	 "REG:", registerNr, "VALUE", value)
	chip8.Registers[registerNr] = byte(value)
}

func (chip8 *Chip8) OP_7XNN(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF
	chip8.Registers[vx] += byte(value)
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
}

func (chip8 *Chip8) OP_8XXX(opcode uint16){
	_type := opcode & 0x000F
	OPCODES_8[_type](chip8, opcode)
}

func (chip8 *Chip8) OP_8XY0(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F
	chip8.Registers[vx] = chip8.Registers[vy]
}

func (chip8 *Chip8) OP_8XY1(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F
	chip8.Registers[vx] |= chip8.Registers[vy]
}

func (chip8 *Chip8) OP_8XY2(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F
	chip8.Registers[vx] &= chip8.Registers[vy]

}

func (chip8 *Chip8) OP_8XY3(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F
	chip8.Registers[vx] ^= chip8.Registers[vy]

}

func (chip8 *Chip8) OP_8XY4(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F

	vyvalue := chip8.Registers[vy]

	sum := chip8.Registers[vx] + vyvalue

	if sum > 0xFF {
		chip8.Registers[0xF] = 1
	} else {
		chip8.Registers[0xF] = 0
	}

	chip8.Registers[vx] = sum 
}

func (chip8 *Chip8) OP_8XY5(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F

	if chip8.Registers[vx] > chip8.Registers[vy]{
		chip8.Registers[0xF] = 1
	}else {
		chip8.Registers[0xF] = 0
	}

	chip8.Registers[vx] -= chip8.Registers[vy]
}

func (chip8 *Chip8) OP_8XY6(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	val := chip8.Registers[vx]

	chip8.Registers[0xF] = byte(val & 0x1)
	chip8.Registers[vx] = val >> 1
}

func (chip8 *Chip8) OP_8XY7(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4) & 0x00F

	if chip8.Registers[vy] > chip8.Registers[vx] {
		chip8.Registers[0xF] = 1
	} else {
		chip8.Registers[0xF] = 0
	}

	chip8.Registers[vx] = chip8.Registers[vy] - chip8.Registers[vx]
}

func (chip8 *Chip8) OP_8XYE(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	val := chip8.Registers[vx]

	chip8.Registers[0xF] = byte((val >> 7) & 0x1)
	chip8.Registers[vx] = val << 1
}

func (chip8 *Chip8) OP_9XY0(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	vy := (opcode >> 4)& 0x00F

	if chip8.Registers[vx] != chip8.Registers[vy]{
		chip8.ProgramCounter += 2
	}
}

func (chip8 *Chip8) OP_ANNN(opcode uint16){
	// Sets NNN to I register
	chip8.IndexRegister = opcode & 0x0FFF
	fmt.Println("index register value", fmt.Sprintf("%X", chip8.IndexRegister))
}

func (chip8 *Chip8) OP_BNNN(opcode uint16){
	chip8.ProgramCounter = (opcode & 0x0FFF) + uint16(chip8.Registers[0])
}

func (chip8 *Chip8) OP_CXNN(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	value := opcode & 0x00FF
	rand.Seed(time.Now().Unix())

	elo := rand.Intn(0xFF)

	fmt.Println("EEEEEEEEEEEEEEEEEELo", elo)


	chip8.Registers[vx] = byte(elo) & byte(value)
}

func (chip8 *Chip8) OP_DXYN(opcode uint16){
	x := (opcode >> 8) & 0x0F
	y := (opcode >> 4)& 0x00F
	nbytes := opcode & 0x000F
	chip8.drawSprite(int(chip8.Registers[x]), int(chip8.Registers[y]), int(nbytes))
}

func (chip8 *Chip8) OP_EXXX(opcode uint16){
	fType := opcode & 0x00FF
	E_OPCODES[fType](chip8, opcode)
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
}

func (chip8 *Chip8) OP_EXA1(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	key := chip8.Registers[vx]

	if chip8.Keyboard[key] == 0{
		chip8.ProgramCounter += 2
	}
}

func (chip8 *Chip8) OP_EX9E(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	key := chip8.Registers[vx]

	if chip8.Keyboard[key] == 1{
		chip8.ProgramCounter += 2
	}
}

func (chip8 *Chip8) OP_FXXX(opcode uint16) {
	fType := opcode & 0x00FF
	F_OPCODES[fType](chip8, opcode)
}

func (chip8 *Chip8) OP_FX07(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	chip8.Registers[vx] = chip8.Timer
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "set timer to Vx", vx)
}

func (chip8 *Chip8) OP_FX0A(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	// chip8.Registers[vx] = 
	// TODO: how to wait for key ??
	// all operatiin must be blocked until key not pressed
	pressed, keyNumber := IsKeyPressed()

	if pressed {
		chip8.Registers[vx] = byte(keyNumber)
	}else {
		chip8.ProgramCounter -= 2
	}

}

func (chip8 *Chip8) OP_FX15(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	chip8.Timer = chip8.Registers[vx]
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "timer set")

}

func (chip8 *Chip8) OP_FX18(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	chip8.SoundTimer = chip8.Registers[vx]
}

func (chip8 *Chip8) OP_FX1E(opcode uint16){
	vx := (opcode >> 8) & 0x0F
	chip8.IndexRegister += uint16(chip8.Registers[vx])
}

func (chip8 *Chip8) OP_FX29(opcode uint16){
	vx := (opcode >> 8) & 0x0F

	// TODO: make offset to get location of fonts loaded
	fontRawAddr := chip8.Registers[vx]
	fontAdr := (5 * fontRawAddr) + 0x50 // each font has 5 bytes long
	chip8.IndexRegister = uint16(fontAdr)

	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode))
	fmt.Println("I:", fmt.Sprintf("%X", chip8.IndexRegister))

}

func (chip8 *Chip8) OP_FX33(opcode uint16){
	registerNr := (opcode >> 8) & 0x0F
	value := chip8.Registers[registerNr]

	hundreds := byte((value / 100) % 10)
	tens := byte((value / 10) % 10)
	ones := byte(value % 10)
	
	chip8.Memory[chip8.IndexRegister] = hundreds
	chip8.Memory[chip8.IndexRegister + 1] = tens
	chip8.Memory[chip8.IndexRegister + 2] = ones

	fmt.Println("FX33, opcode", fmt.Sprintf("%X", opcode), "val", value, hundreds, tens, ones, registerNr)
}

func (chip8 *Chip8) OP_FX55(opcode uint16){
	registerThreshold := (opcode >> 8) & 0x0F
	for nr := uint16(0); nr <= registerThreshold; nr++ {
		chip8.Memory[chip8.IndexRegister + nr] = chip8.Registers[nr] 
	}
}

func (chip8 *Chip8) OP_FX65(opcode uint16){
	registerThreshold := (opcode >> 8) & 0x0F
	for nr := uint16(0); nr <= registerThreshold; nr++ {
		chip8.Registers[nr] = chip8.Memory[chip8.IndexRegister + nr]
	}
	fmt.Println("OPCODE", fmt.Sprintf("%X", opcode), "maxREG", registerThreshold)
}
