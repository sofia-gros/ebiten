# ebiten

[Japanese](./README.md)

A collection of game development libraries designed to make Ebitengine (ebiten) more powerful and intuitive.
It eliminates common boilerplate in game development, providing an environment where you can focus on implementing your logic.

## Provided Components

| Component                                 | Role                                                                    | Status              |
| :---------------------------------------- | :---------------------------------------------------------------------- | :------------------ |
| [**pad**](./pad/README_en.md)             | Action-based input abstraction & virtual pad UI                         | **Stable (v1.1.0)** |
| [**scene**](./scene/README_en.md)         | Scene management & stack-based overlay control                          | **Stable (v1.0.0)** |
| [**physics**](./physics/README_en.md)     | 2D physics engine with swappable Arcade/Box2D backends                  | **Stable (v1.0.1)** |
| [**emit**](./emit/README_en.md)           | Lightweight and type-safe event dispatcher/emitter                      | **Stable (v1.0.0)** |
| [**asset**](./asset/README_en.md)         | Unified manager for images, audio, data, maps, etc.                     | **Stable (v1.0.0)** |
| [**sound**](./sound/README_en.md)         | Manager for sound, BGM/SE/voice, and 2D audio control                   | **Stable (v1.0.0)** |
| [**save**](./save/README_en.md)           | Flat storage, encryption, and auto-filling slot manager                 | **Stable (v1.0.0)** |
| [**camera**](./camera/README_en.md)       | 2D camera, multi-camera, viewport, and shader control                   | **Stable (v1.0.0)** |
| [**tween**](./tween/README_en.md)         | Easing, dynamic control, and group animations                           | **Stable (v1.0.0)** |
| [**ui**](./ui/README_en.md)               | Game UI, 9-slice, state-based image buttons, and layout                 | **Stable (v1.0.0)** |
| [**animation**](./animation/README_en.md) | Sprite animation player & auto-loader for generic JSON (Aseprite, etc.) | **Stable (v1.0.0)** |
| [**tilemap**](./tilemap/README_en.md)     | Tiled JSON auto-import, culling render, 3 map types, and area queries   | **Stable (v1.0.0)** |

---

## Component Overviews

### pad (Input Management)

Integrates and abstracts keyboards, gamepads, and multi-touch virtual pads (sticks and buttons) into a single logical "action". It also comes with standard support for multiplayer (multiple controllers) such as split-screen battles.

[Learn more](./pad/README_en.md)

### scene (Scene Manager)

Manages each game screen (titles, stages, menus, etc.) as independent, separated components. It features a built-in "overlay stack" that allows you to layer pause screens over the active scene, keeping your Ebitengine screen transition code clean.

[Learn more](./scene/README_en.md)

### physics (Common Physics Engine)

A dedicated Ebitengine physics wrapping library that allows you to seamlessly switch between a lightweight `Arcade` mode (AABB penetration prevention) for RPGs and platformers, and a full-scale `Box2D` mode (circles and rigid body physics) for puzzle games using the exact same API.

[Learn more](./physics/README_en.md)

### emit (Event Dispatcher)

A lightweight, type-safe event dispatcher leveraging Go generics. Operating entirely independently without any external library dependencies, it allows you to simply issue, receive, and unsubscribe from temporary animation events between components or within Tweens.

[Learn more](./emit/README_en.md)

### asset (Asset Manager)

An asset management library offering a simple, Phaser-like user experience. It centrally manages bulk loading, type-safe conversion, and cache clearing for images (single and grid sprites), audio, Tilemap map data, animation definitions, and JSON/TOML/YAML/CSV data.

[Learn more](./asset/README_en.md)

### sound (Sound Manager)

An intuitive sound management library tailored for Ebitengine. It centrally manages flexible sound type extensions using `iota` (BGM, SE, Voice, Env, Custom), crossfading, real-time positional sound (2D audio tracking), batch stopping for fade controls, and pausing.

[Learn more](./sound/README_en.md)

### save (Save Data Manager)

An intuitive and secure save data management library. It centrally manages secure subdirectory storage in plain text (JSON/YAML) or binary, automatic reuse hole-filling storage for deleted gaps via `SaveNewSlot`, AES encryption, tampering detection, and automatic backup recovery upon crashes.

[Learn more](./camera/README_en.md)

### camera (2D Camera Manager)

A lightweight, feature-rich 2D camera library for Ebitengine. It centrally manages pure coordinate control (`SetPos`), screen shaking, boundary limits, closure rendering for external libraries (pad, physics) via `cam.Render`, viewport clipping, automatic custom shader switching, and multi-camera batch rendering prioritized by `ZIndex`.

[Learn more](./camera/README_en.md)

### tween (Easing & Animation Interpolation)

A safe and feature-rich easing and animation interpolation library. It centrally manages safe definitions and explicit playback start via `tween.New(&Option{...}).Play()`, `OnUpdate` callbacks, individual dynamic controls (`Pause`, `Resume`, `Restart`, `Stop`), category-based batch controls via `Group`, and a rich set of easing functions.

[Learn more](./tween/README_en.md)

### ui (Game UI Components)

An intuitive and highly extensible 2D game UI library. It centrally manages 9-slices (`NineSlice`), state-based image switching buttons (`Button`), text inputs (`TextInput`), scroll frames (`ScrollBox`), sliders (`Slider`), automatic layout (`VBox/HBox`), unified getters/setters across all elements (`SetPos/Pos`, `SetSize/Size`, `SetText/Text`, `SetGrayscale`), and optional connections to `camera` and `pad`.

[Learn more](./ui/README_en.md)

### animation (Sprite Animation)

An intuitive 2D sprite animation control library. In addition to manual frame definitions, it supports fully automated animation construction from generic JSON formats like Aseprite and TexturePacker. It balances batch control and slow-motion effects using a `Manager` with individual controls for each standalone `Animator`.

[Learn more](./animation/README_en.md)

### tilemap (2D Tilemap Manager)

A fully independent library providing automatic Tiled Map Editor (`.json` / `.tmj`) importing, direct code generation from 2D arrays, three types of tilemap configurations (`StaticTilemap`, `AnimatedTilemap`, `InfiniteTilemap`), view frustum culling rendering (`DrawRegion`), range-specified queries (`GetArea`), and physics AABB combined rectangle calculations (`CreateCollisionBoxes`).

[Learn more](./tilemap/README_en.md)

---

## License

MIT License
