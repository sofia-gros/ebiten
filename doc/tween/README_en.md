# ebitentween

[日本語](./README.md)

`ebitentween` is a safe, feature-rich, and intuitive easing and animation interpolation library designed for Ebitengine.
It replaces direct pointer mutations with type-safe `OnUpdate(func(val))` and `OnRun(func(progress))` callbacks.

---

## Features

- **Safe Callback Values (`OnUpdate` / `OnRun`)**:
  - `OnUpdate(func(val float64))`: Receives interpolated values for any struct field or setter (`cam.SetZoom(val)`).
  - `OnRun(func(progress float64))`: Receives raw progress from `0.0` to `1.0`.
- **Separation of Definition and Playback**:
  - Define animation parameters with `tween.New(&Option{...})`.
  - Trigger playback by calling or chaining `.Play()`.
- **Dynamic Control (`Pause`, `Resume`, `Restart`, `Stop`)**:
  - Control individual active tweens in real time.
- **Category-Based Grouping (`Group`)**:
  - Group tweens by UI, enemy, or visual effects with `Group` and control them via `Group.PauseAll()`, `Group.Clear()`.
- **Comprehensive Easings**:
  - Linear, Quad, Cubic, Quart, Quint, Sine, Expo, Circ, Elastic, Back, Bounce (In / Out / InOut).

---

## Installation

```bash
go get github.com/sofia-gros/ebiten/tween
```

---

## Usage

### Quick Start

One-off animations defined and played instantly.


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/tween"
)

type Game struct {
	playerX float64
	tw      *tween.Tween
}

func (g *Game) Init() {
	g.tw = tween.New(&tween.Option{
		Start:    0.0,
		End:      500.0,
		Duration: 1.5,
		Ease:     tween.EaseOutBounce,
		Yoyo:     true,
		Loop:     -1,
	}).OnUpdate(func(val float64) {
		g.playerX = val
	}).Play()
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0
	tween.Update(dt)
	return nil
}
```

---

### Full Usage

Group animations by category (e.g. UI) and control them using `Group.PauseAll()` or `Group.Clear()`.


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/tween"
)

type Game struct {
	uiGroup *tween.Group
	uiAlpha float64
	panelX  float64
}

func (g *Game) Init() {
	g.uiGroup = tween.NewGroup()

	g.uiGroup.New(tween.Option{
		Start:    0.0,
		End:      1.0,
		Duration: 0.5,
		Ease:     tween.EaseInOutSine,
	}).OnUpdate(func(val float64) {
		g.uiAlpha = val
	}).Play()

	g.uiGroup.New(tween.Option{
		Start:    -200.0,
		End:      100.0,
		Duration: 0.8,
		Ease:     tween.EaseOutBack,
	}).OnUpdate(func(val float64) {
		g.panelX = val
	}).OnComplete(func() {
		fmt.Println("Slide-in Complete!")
	}).Play()
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0
	g.uiGroup.Update(dt)
	return nil
}

func (g *Game) OnOpenDialog() {
	g.uiGroup.PauseAll()
}

func (g *Game) OnCloseScene() {
	g.uiGroup.Clear()
}
```

---

## License

MIT License
