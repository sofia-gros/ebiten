# ebitenui

[日本語](./README.md)

`ebitenui` is a 2D game UI library for Ebitengine.
It operates 100% independently while providing optional bindings for `camera` (screen-fixed UI) and `pad` (gamepad/keyboard navigation using `input.Action` numbers).

---

## Features

- **Rich Component Suite**:
  - `NineSlice`: 9-patch border stretching.
  - `Button`: State textures (Normal, Hover, Pressed, Disabled) + Icon and Text compositing.
  - `TextInput`: Text entry, cursor blinking, backspace, `OnSubmit`.
  - `ScrollBox`: Clipped scroll container.
  - `Slider`, `CheckBox`, `Label`, `Panel`, `VBox`, `HBox`, `Container`.
- **1-to-1 Option Structs & Getters/Setters**:
  - Every property (`SetPos`, `Pos`, `SetSize`, `Size`, `SetText`, `Text`, `SetGrayscale`) can be set via `Option` structs or individual setter methods.
- **Explicit Grayscale (`SetGrayscale`)**:
  - User textures are rendered faithfully without forced auto-grayscale. Optional grayscale shader mode via `SetGrayscale(true)`.
- **Clean Element Access (`GetAll`, `Get`, `Remove`)**:
  - Access elements via `uiRoot.GetAll()` or `uiRoot.Get(index)`.

---

## Installation

```bash
go get github.com/sofia-gros/ebiten/ui
```

---

## Usage

### Quick Start

Build UI components using standard Ebitengine input without external dependencies.


```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/ui"
)

type Game struct {
	uiRoot *ui.Container
}

func (g *Game) Init() {
	g.uiRoot = ui.NewContainer()
	g.uiRoot.SetPos(50, 50)

	vbox := ui.NewVBox()
	vbox.SetSpacing(15)

	titleLabel := ui.NewLabel("Main Menu", ui.LabelOption{FontSize: 20})
	btn := ui.NewButton(ui.ButtonOption{
		Text:   "Start Game",
		Width:  200,
		Height: 45,
	})
	btn.OnClick(func() {
		fmt.Println("Start Game Clicked!")
	})

	vbox.Add(titleLabel)
	vbox.Add(btn)
	g.uiRoot.Add(vbox)
}

func (g *Game) Update() error {
	g.uiRoot.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.uiRoot.Draw(screen)
}
```

---

### Full Usage

Use text input boxes, scroll containers, camera anchoring, and pad action binding.


```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/pad/input"
	"github.com/sofia-gros/ebiten/ui"
)

const (
	ActionUp     input.Action = 1
	ActionDown   input.Action = 2
	ActionSubmit input.Action = 3
)

type Game struct {
	uiRoot *ui.Container
	cam    *camera.Camera
}

func (g *Game) Init() {
	g.uiRoot = ui.NewContainer()
	g.cam = camera.New(640, 480)

	nameInput := ui.NewTextInput(ui.TextInputOption{
		Placeholder: "Enter name...",
		Width:       220,
		Height:      35,
	})
	nameInput.OnSubmit(func(text string) {
		fmt.Println("Name:", text)
	})

	scrollBox := ui.NewScrollBox(200, 150)
	scrollBox.Add(nameInput)
	g.uiRoot.Add(scrollBox)

	// Anchoring UI to camera viewport
	g.uiRoot.BindCamera(g.cam)

	// Binding pad actions
	in := input.NewInput()
	g.uiRoot.BindPadInput(in, ui.PadInputMapping{
		UpAction:     ActionUp,
		DownAction:   ActionDown,
		SubmitAction: ActionSubmit,
	})
}

func (g *Game) Update() error {
	g.uiRoot.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.uiRoot.Draw(screen)
}
```

---

## License

MIT License
