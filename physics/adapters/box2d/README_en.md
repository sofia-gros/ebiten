# box2d (ebitenphysics adapter)

[日本語](./README.md)

An adapter for a full-scale rigid body physics engine powered by `github.com/ByteArena/box2d` under the hood.
It fully supports complex collision detection and physical simulations such as friction and restitution.

## Supported Features

Supports all features of `BodyOptions`.

### `BodyOptions` Support Status

| Property | Supported | Notes |
| :--- | :--- | :--- |
| `Type` | ◯ | `BodyTypeStatic` (StaticBody), `BodyTypeDynamic` (DynamicBody) |
| `X`, `Y` | ◯ | |
| `Shapes` | ◯ | Supports BoxShape, CircleShape. Passing multiple creates compound shapes (multiple Fixtures). |
| `Friction` | ◯ | Functions as coefficient of friction. |
| `Restitution`| ◯ | Functions as coefficient of restitution (bounciness). |
| `Density` | ◯ | Functions as density (used for mass calculation). |
| `IsSensor` | ◯ | Makes the Fixture a sensor. |
| `OnCollisionBegin` | ◯ | Implemented via ContactListener |
| `OnCollisionEnd` | ◯ | Implemented via ContactListener |
| `OnOverlap` | ◯ | Implemented via PreSolve or ContactListener |

## Installation Instructions

The Box2D core is also required.

```bash
go get github.com/ByteArena/box2d
go get github.com/sofia-gros/ebiten/physics/adapters/box2d
```
