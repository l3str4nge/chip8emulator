package chip8

import "fmt"
import "os"
import "bufio"


func LoadRomToMemory(rom_path string) [4096]byte{
	var memory [4096]byte

	file, err := os.Open(rom_path)

	if err != nil {
		panic(err)
	}
	defer file.Close()


	info, _ := file.Stat()

	var size int64 = info.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	
	var programStartAddr int = 0x200
	fmt.Println("program starts at", programStartAddr)
	for i := 0; i <=  int(size) -1; i++ {
		memory[programStartAddr + i] = bytes[i]
	}

	return memory
}
