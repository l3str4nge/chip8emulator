package main


import "github.com/mateuszz0000/chip8emulator/chip8"
import "os"
import "fmt"
import "strconv"
import "io/ioutil"

func main(){

	args := os.Args[1:]

	if len(args) == 0{
		fmt.Println("Please specify game or game and scale")
		return
	}

	gameName := ""
	gameScale := 1
	if len(args) == 1{
		gameName = args[0]

	}else if len(args) >= 2{
		gameName = args[0]
		gameScale, _ = strconv.Atoi(args[1])
	}

	_, err := os.Stat("c8games/" + gameName)

	if os.IsNotExist(err) {
		fmt.Println("Game does not exists. Please choose from following list:")
		files, _ := ioutil.ReadDir("c8games")
		for i, f := range files {
            fmt.Println(i, f.Name())
    }
		return
    }
	

	memory := chip8.LoadRomToMemory("c8games/" + gameName)
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
			[16]byte{},
			false,
		    gameScale}

	chip8.Emulator.Run()
}