# ebitenphysics

[日本語](./README.md)

ebitenphysics is a physics engine wrapper library for Ebitengine.
It allows you to seamlessly switch between a lightweight Arcade mode (AABB collision without penetration) for RPGs and platformers, and a full-fledged Box2D mode (circles and rigid body physics) for puzzle games, using the exact same API.

## Features

- **Pure Design**: The ebitenphysics core does not include any specific physics engine. By importing only the engine (adapter) you want to use, no unnecessary code is injected into your game.
- **DI-based**: Designed to operate by passing an engine such as `arcade` or `box2d` into the Manager.
- **Flexible Drawing**: You can fetch `x, y` and draw elements yourself, or you can register them with images for batch drawing. It does not conflict with existing scene management or UI.

## Update History

Version 1.0.0 2026/07/27
- Initial release

## Installation

First, install the core library.

```bash
go get github.com/sofia-gros/ebiten/physics
```

Next, install the adapter for the physics engine you wish to use.

> **⚠️ Note: Installing External Engines**
> This library (ebitenphysics) itself does not include any external physics engine code.
> When using external engines like `box2d`, in addition to the adapter above, **you must also download (`go get`) the base physics engine yourself**.

| Adapter | Description | Dependency Engine (Please download yourself) |
| :--- | :--- | :--- |
| `arcade` | Custom lightweight AABB engine. | None |
| `box2d` | Full-scale 2D rigid body simulation (v2). | `go get github.com/ByteArena/box2d` |
| `box2dgo` | Go port of the latest Box2D (v3). | `go get github.com/oliverbestmann/box2d-go` |
| `phygo` | Lightweight 2D physics engine. | `go get github.com/ab-dek/Phygo-2D` |

```bash
# Example: Using the Box2D (v3) engine
go get github.com/oliverbestmann/box2d-go
go get github.com/sofia-gros/ebiten/physics/adapters/box2dgo
```

## Basic Usage

```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/physics"

	// Import the physics engine you want to use
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
)

// Collision group definitions (similar to layer numbers)
const (
	GroupPlayer physics.Group = 1
	GroupEnemy  physics.Group = 2
)

type GameScene struct {
	physManager *physics.Manager
	playerBody  physics.Body
	playerImage *ebiten.Image
}

func (s *GameScene) Init() {
	s.physManager = physics.NewManager()
	s.physManager.SetWorld(arcade.NewWorld())

	// Create a Body (entity for physical calculation)
	s.playerBody = s.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    100,
		Y:    100,
		// Basic rectangular collision
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		// Collision callback
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupEnemy {
				fmt.Println("Hit an enemy!")
			}
		},
	})

	// Assign group number
	s.playerBody.SetGroup(GroupPlayer)

	// Register for batch drawing
	s.physManager.AddRenderable(s.playerBody, s.playerImage)
}

func (s *GameScene) Update() error {
	s.physManager.Update(1.0 / 60.0)
	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	// When drawing manually
	// x, y := s.playerBody.Position()
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(x, y)
	// screen.DrawImage(s.playerImage, op)

	// When leaving it to the manager
	s.physManager.Draw(screen)

	// Draw debug wireframes for collision detection
	s.physManager.DrawDebug(screen)
}
```

## API & Function Details

Here is a list of the main features provided by this library's interface (common to all physics engines).
*Note: Which parameters, such as friction or restitution, are actually effective depends on the injected **physics engine adapter** (refer to its README).*

### 1. `physics.Manager` Methods

The central object overseeing the physical space and rendering.

- **`SetWorld(world World)`**: Injects the physics engine instance (adapter) to use. (Required)
- **`SetGravity(gx, gy float64)`**: Sets global gravity. (e.g., `SetGravity(0, 100)`)
- **`Gravity() (float64, float64)`**: Gets current gravity.
- **`CreateBody(options BodyOptions) Body`**: Creates an object in the space.
- **`RemoveBody(body Body)`**: Completely removes an object from space and the render list.
- **`Update(dt float64)`**: Advances physics simulation by `dt` seconds.
- **`AddRenderable(body Body, img *ebiten.Image)`**: Registers an image-body pair for batch drawing.
- **`Draw(screen *ebiten.Image)`**: Batch draws all pairs registered with `AddRenderable`.
- **`DrawDebug(screen *ebiten.Image)`**: Visualizes collision shapes (boxes or circles) on screen for debugging.

### 2. `physics.Body` Methods

The physics object returned by `CreateBody`. You can get or set its state at any time.

#### State Retrieval (Get)

- **`Position() (x, y float64)`**: Gets current X, Y coordinates.
- **`Rotation() float64`**: Gets current rotation angle (radians).
- **`Velocity() (vx, vy float64)`**: Gets current X, Y velocity.
- **`AngularVelocity() float64`**: Gets current angular velocity.
- **`Group() physics.Group`**: Gets the group number (numeric) used for collision layers.
- **`Data() interface{}`**: Gets arbitrary user data (like a custom struct pointer) attached to the body.

#### State Modification (Set)

- **`SetPosition(x, y float64)`**: Forcibly moves coordinates (warp).
- **`SetRotation(angle float64)`**: Forcibly changes angle.
- **`SetVelocity(vx, vy float64)`**: Directly sets velocity (e.g., pixels/sec). Useful for jumps and dashes.
- **`SetAngularVelocity(omega float64)`**: Sets angular velocity.
- **`ApplyForce(fx, fy float64)`**: Gradually applies force to current velocity (ideal for continuous forces like wind).
- **`SetGroup(group physics.Group)`**: Assigns object type numerically (Player=1, Enemy=2). Faster than string comparison.
- **`SetData(data interface{})`**: Allows the body to hold custom data (like status structs) for detailed processing during collisions.

### 3. `physics.BodyOptions` (Creation Settings)

- **`Type`**: Specify `physics.BodyTypeDynamic` (moving object) or `physics.BodyTypeStatic` (unmoving walls, etc.).
- **`Shapes []ShapeDef`**: The shapes composing the object. Multiple shapes allow for "compound shapes" with independent arms or legs forming a single body.
- **`Friction`**: Coefficient of friction (resistance to sliding).
- **`Restitution`**: Coefficient of restitution (bounciness).
- **`Density`**: Density, used for mass calculation.
- **`IsSensor`**: If `true`, collision callbacks are triggered, but the object is no longer physically repelled (perfect for item pickup areas or poison swamps).

#### Collision Callbacks

- **`OnCollisionBegin func(other Body)`**: Called exactly once at the "moment of contact".
- **`OnOverlap func(other Body)`**: Called every frame while "overlapping or continuing contact".
- **`OnCollisionEnd func(other Body)`**: Called exactly once at the "moment of separation".

### 4. Shape Definition

Using `ShapeDef`, you can offset a shape relative to the body's center coordinates.

```go
// Example of a square attached to the side of a central circle
Shapes: []physics.ShapeDef{
	{Shape: physics.CircleShape{Radius: 16}, OffsetX: 0, OffsetY: 0},
	{Shape: physics.BoxShape{Width: 32, Height: 8}, OffsetX: 32, OffsetY: 0},
}
```

## License

MIT License
