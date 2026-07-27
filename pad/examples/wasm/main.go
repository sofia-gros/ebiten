package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/pad/input"
)

const (
	ActionJump input.Action = 1
	ActionMove input.Action = 2
)

type Game struct {
	in *input.Input
}

func (g *Game) Update() error {
	g.in.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	msg := `ebitenpad WASM Example
Controls:
- Keyboard: [Space] to Jump, [WASD] to Move
- Gamepad: [Standard Button 0] to Jump, [Stick] to Move
- Virtual: [Button] to Jump, [Stick] to Move

Jump:
  Pressed: %v
  JustPressed: %v
  JustReleased: %v

Move:
  Vector: (%.2f, %.2f)
  Strength: %.2f
`
	vx, vy := 0.0, 0.0
	strength := 0.0
	if state, ok := g.in.GetActionState(ActionMove); ok {
		vx, vy = state.X, state.Y
		strength = state.Strength
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(msg,
		g.in.Pressed(ActionJump),
		g.in.JustPressed(ActionJump),
		g.in.JustReleased(ActionJump),
		vx, vy, strength))

	// バーチャル UI の描画
	g.in.Virtual().Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	in := input.NewInput()

	// Bind Keyboard
	in.BindKey(ActionJump, ebiten.KeySpace)
	in.BindKeyAxis(ActionMove, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	// Bind Gamepad
	in.BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom)
	in.BindGamepadAxis(ActionMove, 0, 1)

	// Virtual Pad の設定
	vpad := in.Virtual()
	jumpBtn := vpad.AddButton().SetPosition(550, 400).SetRadius(40)
	moveStick := vpad.AddStick().SetPosition(100, 380).SetRadius(60)

	in.BindButton(ActionJump, jumpBtn)
	in.BindStick(ActionMove, moveStick)

	g := &Game{in: in}

	ebiten.SetWindowTitle("ebitenpad WASM Example")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
