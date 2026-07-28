# ebiten

[日本語](./README.md)

A suite of game development libraries designed to make working with Ebitengine (ebiten) more powerful and intuitive.
It provides an environment where you can focus on implementing your game logic by eliminating common boilerplate code in game development.

## Components Provided

| Component | Role | Status |
| :--- | :--- | :--- |
| [**pad**](./pad/README_en.md) | Action-based input abstraction & virtual pad UI | **Stable (v1.1.0)** |
| [**scene**](./scene/README_en.md) | Scene management & stack-based overlay control | **Stable (v1.0.0)** |
| [**physics**](./physics/README_en.md) | 2D physics engine wrapper with interchangeable Arcade/Box2D backends | **Stable (v1.0.0)** |

---

## Component Overview

### pad (Input Management)

Integrates and abstracts virtual pads (sticks and buttons) via keyboard, gamepad, and multi-touch into single logical "actions". It natively supports multiplayer (multiple controllers) setups, such as split-screen local multiplayer.

[Learn more](./pad/README_en.md)

### scene (Scene Manager)

Manages each screen of the game (title, stage, menu, etc.) separately as independent components. It includes an "overlay stack" feature that allows pausing and layering screens, keeping Ebitengine's screen transition code clean.

[Learn more](./scene/README_en.md)

### physics (Common Physics Engine)

A dedicated physics wrapping library for Ebitengine that allows you to seamlessly switch between a lightweight `Arcade` mode (AABB collision without penetration) for RPGs or platformers, and a full-fledged `Box2D` mode (circles and rigid body physics) for puzzle games and more, using the exact same API.

[Learn more](./physics/README_en.md)

---

## License

MIT License
