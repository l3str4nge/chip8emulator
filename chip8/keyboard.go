package chip8

import "github.com/hajimehoshi/ebiten"
// import "fmt"



func KeyPressed(chip8 *Chip8, keyNumber byte){
	chip8.Keyboard[keyNumber] = 0x01
}

func KeyReleased(chip8 *Chip8, keyNumber byte){
	chip8.Keyboard[keyNumber] = 0x00
}


func CheckKeyboard(chip8 *Chip8){
    handleKey(chip8, ebiten.Key1, 0x1)
    handleKey(chip8, ebiten.Key2, 0x2)
    handleKey(chip8, ebiten.Key3, 0x3)
    handleKey(chip8, ebiten.Key4, 0xC)
    handleKey(chip8, ebiten.KeyQ, 0x4)
    handleKey(chip8, ebiten.KeyW, 0x5)
    handleKey(chip8, ebiten.KeyE, 0x6)
    handleKey(chip8, ebiten.KeyR, 0xD)
    handleKey(chip8, ebiten.KeyA, 0x7)
    handleKey(chip8, ebiten.KeyS, 0x8)
    handleKey(chip8, ebiten.KeyD, 0x9)
    handleKey(chip8, ebiten.KeyF, 0xE)
    handleKey(chip8, ebiten.KeyZ, 0xA)
    handleKey(chip8, ebiten.KeyX, 0x0)
    handleKey(chip8, ebiten.KeyC, 0xB)
    handleKey(chip8, ebiten.KeyV, 0xF)
}


func handleKey(chip8 *Chip8, key ebiten.Key, number byte){
    if ebiten.IsKeyPressed(key) {
        KeyPressed(chip8, number)
    }else{
        KeyReleased(chip8, number)
    }
}

func IsKeyPressed() (bool, int) {
    keys := [16]ebiten.Key{
        ebiten.KeyX, // 0x0
        ebiten.Key1, // 0x1
        ebiten.Key2, // 0x2
        ebiten.Key3, // 0x3
        ebiten.KeyC, // 0x4
        ebiten.KeyW, // 0x5
        ebiten.KeyE, // 0x6
        ebiten.KeyA, // 0x7
        ebiten.KeyS, // 0x8
        ebiten.KeyD, // 0x9
        ebiten.KeyZ, // 0xA
        ebiten.KeyC, // 0xB
        ebiten.Key4, // 0xC
        ebiten.Key9, // 0xD
        ebiten.KeyF, // 0xE
        ebiten.KeyV, // 0xF
    }

    for n, k := range keys {
        if ebiten.IsKeyPressed(k){
            return true, n
        }
    }

    return false, -1
}