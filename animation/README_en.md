# ebitenanimation

[Japanese](./README.md)

ebitenanimation is an intuitive and feature-rich 2D sprite animation library for Ebitengine.
In addition to manual code-based frame definitions, it supports **fully automatic animation generation from generic JSON formats such as Aseprite and TexturePacker**.

It has zero external dependencies and relies solely on `ebiten/v2`.

---

## Features

- **Dual-Level Control (Manager & Individual)**:
  - `Manager`: Batch update, batch pause/resume, and global time-scale (slow-mo/fast-forward) management.
  - `Animator`: Individual play, pause, resume, stop, reset, and speed control for each object.
- **Generic JSON Auto-Loading (`CreateAnimatorFromJSON`)**:
  - Automatically builds clips from Aseprite/TexturePacker JSON (`frameTags` and frame `duration`).
- **Flexible Playback Options & Dynamic Toggles**:
  - `Loop`, `Reverse`, `PingPong`, and custom frame durations.
  - Supports initial configuration via option structs and dynamic methods (`SetLoop(bool)`, `SetReverse(bool)`).
- **Event Callbacks**:
  - `OnComplete` (when a clip finishes) and `OnFrame` (when reaching a specific frame index).
