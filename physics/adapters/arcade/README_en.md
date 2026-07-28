# arcade (ebitenphysics adapter)

[日本語](./README.md)

A custom, lightweight AABB (Axis-Aligned Bounding Box / penetration prevention) physics engine adapter.
It has no external dependencies and is specialized for simple rectangle and circle pushing/overlap detection, ideal for RPGs and platformers.

## Supported Features

Because this engine prioritizes being lightweight, it does not support advanced physics simulation features (like friction and restitution).

### `BodyOptions` Support Status

| Property | Supported | Notes |
| :--- | :--- | :--- |
| `Type` | ◯ | Only `BodyTypeStatic` and `BodyTypeDynamic` are supported. |
| `X`, `Y` | ◯ | |
| `Shapes` | ◯ | Supports `BoxShape` and `CircleShape`. Also supports compound shapes with offsets. |
| `Friction` | ✕ | Ignored. |
| `Restitution`| ✕ | Ignored. |
| `Density` | ✕ | Ignored. (All Dynamic bodies have equal push-back force) |
| `IsSensor` | ◯ | Triggers callbacks only, without physical push-back. |
| `OnCollisionBegin` | ◯ | |
| `OnCollisionEnd` | ◯ | |
| `OnOverlap` | ◯ | |

## Installation

```bash
go get github.com/sofia-gros/ebiten/physics/adapters/arcade
```
