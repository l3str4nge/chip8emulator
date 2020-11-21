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
    handleKey(chip8, ebiten.KeyQ, 0x01)
    handleKey(chip8, ebiten.KeyW, 0x02)
    handleKey(chip8, ebiten.KeyE, 0x03)
    handleKey(chip8, ebiten.KeyR, 0x04)
    handleKey(chip8, ebiten.KeyT, 0x05)
    handleKey(chip8, ebiten.KeyY, 0x06)
    handleKey(chip8, ebiten.KeyU, 0x07)
    handleKey(chip8, ebiten.KeyI, 0x08)
    handleKey(chip8, ebiten.KeyO, 0x09)
    handleKey(chip8, ebiten.KeyP, 0x0A)
    handleKey(chip8, ebiten.KeyA, 0x0B)
    handleKey(chip8, ebiten.KeyS, 0x0C)
    handleKey(chip8, ebiten.KeyD, 0x0D)
    handleKey(chip8, ebiten.KeyF, 0x0E)
    handleKey(chip8, ebiten.KeyG, 0x0F)
    // handleKey(chip8, ebiten.KeyH, 15)
}


func handleKey(chip8 *Chip8, key ebiten.Key, number byte){
    if ebiten.IsKeyPressed(key) {
        KeyPressed(chip8, number)
    }else{
        KeyReleased(chip8, number)
    }
}

func IsKeyPressed() bool {
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

    for _, k := range keys {
        if ebiten.IsKeyPressed(k){
            return true
        }
    }

    return false
}