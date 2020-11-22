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
    handleKey(chip8, ebiten.KeyQ, 0x1)
    handleKey(chip8, ebiten.KeyW, 0x2)
    handleKey(chip8, ebiten.KeyE, 0x3)
    handleKey(chip8, ebiten.KeyR, 0x4)
    handleKey(chip8, ebiten.KeyT, 0x5)
    handleKey(chip8, ebiten.KeyY, 0x6)
    handleKey(chip8, ebiten.KeyU, 0x7)
    handleKey(chip8, ebiten.KeyI, 0x8)
    handleKey(chip8, ebiten.KeyO, 0x9)
    handleKey(chip8, ebiten.KeyP, 0xA)
    handleKey(chip8, ebiten.KeyA, 0xB)
    handleKey(chip8, ebiten.KeyS, 0xC)
    handleKey(chip8, ebiten.KeyD, 0xD)
    handleKey(chip8, ebiten.KeyF, 0xE)
    handleKey(chip8, ebiten.KeyG, 0xF)
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
        ebiten.KeyQ,
        ebiten.KeyW,
        ebiten.KeyE,
        ebiten.KeyR,
        ebiten.KeyT,
        ebiten.KeyY,
        ebiten.KeyU,
        ebiten.KeyI,
        ebiten.KeyO,
        ebiten.KeyP,
        ebiten.KeyA,
        ebiten.KeyS,
        ebiten.KeyD,
        ebiten.KeyF,
        ebiten.KeyG,

    }

    for n, k := range keys {
        if ebiten.IsKeyPressed(k){
            return true, n
        }
    }

    return false, -1
}