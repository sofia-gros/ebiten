# ebiten

[日本語](./README.md)

A suite of game development libraries designed to make working with Ebitengine (ebiten) more powerful and intuitive.
It provides an environment where you can focus on implementing your game logic by eliminating common boilerplate code in game development.

To maintain a clean module architecture, all documentation files are centralized under the root `doc/` directory, keeping library package folders clean with only pure Go source code and unit tests.

## Components Provided

| Component | Role | Status |
| :--- | :--- | :--- |
| [**pad**](./doc/pad/README_en.md) | Action-based input abstraction & virtual pad UI | **Stable (v1.1.0)** |
| [**scene**](./doc/scene/README_en.md) | Scene management & stack-based overlay control | **Stable (v1.0.0)** |
| [**physics**](./doc/physics/README_en.md) | 2D physics engine wrapper with interchangeable Arcade/Box2D backends | **Stable (v1.0.1)** |
| [**emit**](./doc/emit/README_en.md) | Lightweight & type-safe event dispatcher / emitter | **Stable (v1.0.0)** |
| [**asset**](./doc/asset/README_en.md) | Unified manager for images, audio, data, and map assets | **Stable (v1.0.0)** |
| [**sound**](./doc/sound/README_en.md) | Sound manager for BGM/SE/Voice & 2D positional audio | **Stable (v1.0.0)** |
| [**save**](./doc/save/README_en.md) | Flat & encrypted save data manager with auto-hole fill slots | **Stable (v1.0.0)** |
| [**camera**](./doc/camera/README_en.md) | 2D camera, multi-camera, viewport & custom shader controller | **Stable (v1.0.0)** |
| [**tween**](./doc/tween/README_en.md) | Easing, dynamic animation control & group animation manager | **Stable (v1.0.0)** |
| [**ui**](./doc/ui/README_en.md) | 2D game UI with 9-slice, state buttons, scrollboxes & layout | **Stable (v1.0.0)** |
| [**animation**](./doc/animation/README_en.md) | 2D sprite animator & generic JSON (Aseprite, etc.) auto-loader | **Stable (v1.0.0)** |
| [**tilemap**](./doc/tilemap/README_en.md) | Tiled JSON auto-importer, viewport culling, 3 map types & area queries | **Stable (v1.0.0)** |

---

## Component Overview

### pad (Input Management)

Integrates and abstracts virtual pads (sticks and buttons) via keyboard, gamepad, and multi-touch into single logical "actions". It natively supports multiplayer (multiple controllers) setups, such as split-screen local multiplayer.

[Learn more](./doc/pad/README_en.md)

### scene (Scene Manager)

Manages each screen of the game (title, stage, menu, etc.) separately as independent components. It includes an "overlay stack" feature that allows pausing and layering screens, keeping Ebitengine's screen transition code clean.

[Learn more](./doc/scene/README_en.md)

### physics (Common Physics Engine)

A dedicated physics wrapping library for Ebitengine that allows you to seamlessly switch between a lightweight `Arcade` mode (AABB collision without penetration) for RPGs or platformers, and a full-fledged `Box2D` mode (circles and rigid body physics) for puzzle games and more, using the exact same API.

[Learn more](./doc/physics/README_en.md)

### emit (Event Dispatcher)

A lightweight and type-safe event dispatcher powered by Go generics. Operates independently with zero external dependencies, making event emission, listening, and unsubscription seamless.

[Learn more](./doc/emit/README_en.md)

### asset (Asset Manager)

An asset management library providing Phaser-like intuitive API. Manages single images, spritesheets, audio, tilemaps, animation definitions, and JSON/TOML/YAML/CSV data loading, type-safe conversion, and cache management.

[Learn more](./doc/asset/README_en.md)

### sound (Sound Manager)

Intuitive sound management library for Ebitengine. Features flexible sound type expansion (BGM, SE, Voice, Env), crossfading, real-time 2D positional sound tracking, and bulk fade/pause controls.

[Learn more](./doc/sound/README_en.md)

### save (Save Data Manager)

Intuitive and secure save data manager. Features plain text / binary sub-directory saving, automatic slot reuse (`SaveNewSlot`), AES encryption, tamper detection, and crash auto-recovery.

[Learn more](./doc/save/README_en.md)

### camera (2D Camera Manager)

Lightweight and feature-rich 2D camera library. Supports pure position control (`SetPos`), screen shake, bounds clamping, closure rendering (`cam.Render`), viewport clipping, custom shaders, and multi-camera Z-indexing.

[Learn more](./doc/camera/README_en.md)

### tween (Easing & Interpolation)

Safe and feature-rich easing & interpolation library. Provides explicit `Play()` control, `OnUpdate` callbacks, individual dynamic methods (`Pause`, `Resume`, `Stop`), group controls, and extensive easing functions.

[Learn more](./doc/tween/README_en.md)

### ui (Game UI Components)

Intuitive 2D game UI library. Features 9-slice rendering, state buttons, text inputs, scroll boxes, sliders, auto-layout containers (`VBox/HBox`), unified getters/setters, and optional `camera` / `pad` integration.

[Learn more](./doc/ui/README_en.md)

### animation (Sprite Animation)

Intuitive 2D sprite animation library. Supports manual frame definitions as well as fully automatic clip generation from generic JSON formats (Aseprite, TexturePacker). Features `Manager` for batch control / slow-mo and `Animator` for individual control.

[Learn more](./doc/animation/README_en.md)

### tilemap (2D Tilemap Manager)

Fully independent 2D tilemap library. Features Tiled JSON auto-import, code-first 2D array generation, 3 tilemap types (`StaticTilemap`, `AnimatedTilemap`, `InfiniteTilemap`), viewport culling (`DrawRegion`), area queries (`GetArea`), and AABB collision box generation (`CreateCollisionBoxes`).

[Learn more](./doc/tilemap/README_en.md)

---

## License

MIT License
