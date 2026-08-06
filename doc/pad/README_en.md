# ebitenpad

[日本語](./README.md)

ebitenpad is an action-based input management library for Ebitengine (ebiten).
It allows you to integrate keyboard, gamepad, and virtual pad (touch) controls into a single logical "action".

## Features

- **Action-based Abstraction**: Separates game logic from input devices (keyboard, gamepad, touch).
- **Built-in Virtual Pad**: Easily create UI elements like sticks and buttons, controllable via multi-touch.
- **Simple API**: Provides intuitive methods such as `Pressed`, `JustPressed`, and `JustReleased`.
- **Device Integration**: Bind multiple device inputs to a single action.

## Update History

v1.1.0 2026/07/06
- Added a feature to separate inputs per player, like in split-screen local multiplayer.
- Changed `Strength` to use actual input values.

## Installation

```bash
go get github.com/sofia-gros/ebiten/pad
```

## Basic Usage

```go
package main

import (
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
	vx, vy := 0.0, 0.0
	strength := 0.0
	if state, ok := g.in.GetActionState(ActionMove); ok {
		vx, vy = state.X, state.Y
		strength = state.Strength
	}

	// Draw Virtual UI
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

	// Virtual Pad setup
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
```

## Library Structure

- `input`: The main input manager. Handles binding actions and state management.
- `virtual`: UI components like virtual sticks and buttons.
- `examples/wasm`: Demo code running on WebAssembly.

## Multi-player Support (`Controller` / `For`)

To separate input per player, such as in a split-screen local multiplayer game, use the `Controller` type and the `For()` method.

```go
const (
    P1 input.Controller = 0 // also corresponds to gamepad 0
    P2 input.Controller = 1 // also corresponds to gamepad 1
)

const (
    ActionJump input.Action = 1
    ActionMove input.Action = 2
)

// Setup
in.For(P1).BindKey(ActionJump, ebiten.KeySpace)
in.For(P1).BindKeyAxis(ActionMove, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)
in.For(P1).BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom) // gamepad 0

in.For(P2).BindKey(ActionJump, ebiten.KeyEnter)
in.For(P2).BindKeyAxis(ActionMove, ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyUp, ebiten.KeyDown)
in.For(P2).BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom) // gamepad 1

// Call Update only once
in.Update()

// Getting state (P1 and P2 are completely independent. Concurrent inputs do not overwrite each other)
if state, ok := g.in.For(P1).GetActionState(ActionMove); ok {
    p1vx, p1vy = state.X, state.Y
}
if state, ok := g.in.For(P2).GetActionState(ActionMove); ok {
    p2vx, p2vy = state.X, state.Y
}
```

For single-player, just omit `For()` and it will work as usual. The `Controller` value acts as the gamepad index.

```go
// Single player: as usual
in.BindKey(ActionJump, ebiten.KeySpace)
in.GetActionState(ActionJump)
```

## About `state.Strength`

`ActionState.Strength` represents the input strength from `0.0` to `1.0`. The calculation method depends on the device.

| Device | Calculation Method |
| --- | --- |
| Keyboard (Single Key) | Always `1.0` while pressed. |
| Keyboard (Axis) | `√(dx² + dy²)` clamped at `1.0`. Diagonal input is capped at `1.0`. |
| Gamepad (Button) | Always `1.0` while pressed. |
| Gamepad (Analog Axis)| `√(x² + y²)` clamped at `1.0`. Reflects actual analog input strength. |
| Virtual Stick | Distance from the center divided by the radius. A slight touch might yield `0.3`. |

## Gamepad Deadzone

It's common for a stick to register small values even in a neutral position. You can handle this using `BindGamepadAxisWithDeadzone`.

```go
// Ignore input where Strength is 0.2 or less
in.BindGamepadAxisWithDeadzone(ActionMove, 0, 1, 0.2)
```

`BindGamepadAxis` is equivalent to no deadzone (`0.0`).

## Multiple Gamepads

It cycles through all connected gamepads and adopts the one with the strongest input. This means if you have two controllers plugged in, using either one will work.

## License

MIT License
