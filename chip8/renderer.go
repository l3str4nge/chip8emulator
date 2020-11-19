package chip8

import "fmt"

func (chip8 *Chip8) drawSprite(x,y, nbytes int) {
	fmt.Println("XY", x, y)
	for col := 0; col < 8 ; col ++ {
		for row :=0; row < nbytes; row++ {
			value := chip8.Memory[chip8.IndexRegister + uint16(row)]
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
			
			pixelAt := chip8.Pixels[xx][yy]
			if pixValue == 0x1 {	

				if pixelAt == 0x1 {
					fmt.Println("COllision")
					chip8.Registers[0xF] = 1
				} else {
					chip8.Registers[0xF] = 0
				}

				chip8.Pixels[xx][yy] = pixelAt ^ 0xFF
			}			
		}
	}
}
