# ebitenphysics

[日本語](./README.md)

`ebitenphysics` is a 2D physics engine wrapper library designed for Ebitengine.
It allows you to seamlessly switch between lightweight `Arcade` (AABB separation for platformers/RPGs) and full rigid-body `Box2D` engines using a unified, intuitive API.

---

## Features

- **Pure & Lightweight Architecture**: The core package contains no heavy physics engine code. You only import the specific engine adapters (`arcade`, `box2d`, etc.) you need.
- **Two Body Creation Styles**:
  - **Simple Style**: Create a single box or circle body cleanly using `Shape: physics.BoxShape{...}`.
  - **Advanced Style**: Define collision callbacks, compound shapes (`Shapes`), and collision layers using `BodyOptions`.
- **Flexible Drawing Options**: Fetch positions (`Position()`) to draw manually, or pair bodies with images using `AddRenderable` for automated batch drawing.

---

## Installation

Install the core package:

```bash
go get github.com/sofia-gros/ebiten/physics
```

Next, install your preferred adapter:

| Adapter | Description | Underlying Engine (Install if required) |
| :--- | :--- | :--- |
| `arcade` | Lightweight AABB (separation-only) engine for platformers & RPGs. | None |
| `box2d` | Full 2D rigid-body physics simulation (Box2D v2). | `go get github.com/ByteArena/box2d` |
| `box2dgo` | Pure Go port of Box2D v3. | `go get github.com/oliverbestmann/box2d-go` |
| `phygo` | Lightweight 2D physics engine. | `go get github.com/ab-dek/Phygo-2D` |

```bash
# Example 1: Arcade AABB engine (No external dependency required)
go get github.com/sofia-gros/ebiten/physics/adapters/arcade

# Example 2: Box2D (v3) engine (Install underlying engine as well)
go get github.com/oliverbestmann/box2d-go
go get github.com/sofia-gros/ebiten/physics/adapters/box2dgo
```

---

## Usage

### Quick Start

Create a single box body (`BoxShape`) with minimal code and draw it manually.


```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
)

type Game struct {
	phys *physics.Manager
	body physics.Body
}

func (g *Game) Init() {
	g.phys = physics.NewManager()
	g.phys.SetWorld(arcade.NewWorld()) // Inject Arcade adapter

	// Create a 32x32px dynamic box body cleanly with Shape
	g.body = g.phys.CreateBody(physics.BodyOptions{
		Type:  physics.BodyTypeDynamic,
		X:     100,
		Y:     100,
		Shape: physics.BoxShape{Width: 32, Height: 32},
	})
}

func (g *Game) Update() error {
	g.body.SetVelocity(100, 0)
	g.phys.Update(1.0 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fetch coordinates for custom drawing
	x, y := g.body.Position()
	_ = x
	_ = y

	// Debug shape rendering
	g.phys.DrawDebug(screen)
}
```

---

### Full Usage

Define collision groups, callbacks (`OnCollisionBegin`), and batch rendering with sprite images.


```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
)

const (
	GroupPlayer physics.Group = 1
	GroupEnemy  physics.Group = 2
)

type Game struct {
	phys       *physics.Manager
	playerBody physics.Body
	playerImg  *ebiten.Image
}

func (g *Game) Init() {
	g.phys = physics.NewManager()
	g.phys.SetWorld(arcade.NewWorld())

	// Advanced style: Full options with callbacks
	g.playerBody = g.phys.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    100,
		Y:    100,
		Shape: physics.BoxShape{Width: 32, Height: 32},
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupEnemy {
				fmt.Println("Collided with Enemy!")
			}
		},
	})

	g.playerBody.SetGroup(GroupPlayer)

	// Pair body with image for automated batch rendering
	g.phys.AddRenderable(g.playerBody, g.playerImg)
}

func (g *Game) Update() error {
	g.phys.Update(1.0 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Batch draw all paired bodies
	g.phys.Draw(screen)
	g.phys.DrawDebug(screen)
}
```

---

## License

MIT License
