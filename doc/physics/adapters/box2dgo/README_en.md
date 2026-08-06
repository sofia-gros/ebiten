# box2dgo (ebitenphysics adapter)

[日本語](./README.md)

An adapter for a rigid body physics engine powered by `github.com/oliverbestmann/box2d-go` (a Go port of Box2D v3).
Compared to the traditional Box2D (v2) adapter, it utilizes an engine with a faster and more modern internal design. It supports complex collision detection and physical simulations like friction and restitution.

## Supported Features

Supports all features of `BodyOptions`.

### `BodyOptions` Support Status

| Property | Supported | Notes |
| :--- | :--- | :--- |
| `Type` | ◯ | `BodyTypeStatic` (StaticBody), `BodyTypeDynamic` (DynamicBody) |
| `X`, `Y` | ◯ | |
| `Shapes` | ◯ | Supports BoxShape, CircleShape. Passing multiple creates compound shapes (multiple Shapes). |
| `Friction` | ◯ | Functions as coefficient of friction. |
| `Restitution`| ◯ | Functions as coefficient of restitution (bounciness). |
| `Density` | ◯ | Functions as density (used for mass calculation). |
| `IsSensor` | ◯ | Makes the shape a sensor. |
| `OnCollisionBegin` | ◯ | Implemented via `ContactEvents.BeginEvents` |
| `OnCollisionEnd` | ◯ | Implemented via `ContactEvents.EndEvents` |
| `OnOverlap` | △ | (Currently unimplemented: Planned for support in future updates via sensor events, etc.) |

## Installation Instructions

The Box2D-Go core is also required.

```bash
go get github.com/oliverbestmann/box2d-go
go get github.com/sofia-gros/ebiten/physics/adapters/box2dgo
```
