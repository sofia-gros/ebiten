# ebitentilemap

[Japanese](./README.md)

ebitentilemap is a lightweight and feature-rich 2D tilemap library for Ebitengine.
It supports fully automatic import of Tiled Map Editor (`.json` / `.tmj`) files, code-first 2D array map creation, 3 tilemap variants, viewport culling (`DrawRegion`), intuitive area queries, and collision box generation for physics engines.

It has zero dependencies on other custom libraries and relies solely on `ebiten/v2`.

---

## Features

- **Fully Independent Module**: Relies only on `ebiten/v2`. Usable as a standalone package.
- **3 Tilemap Variants**:
  - `StaticTilemap`: Fast rendering for fixed ground and walls.
  - `AnimatedTilemap`: Animated tiles for water, swaying grass, etc.
  - `InfiniteTilemap`: Chunk-based auto-generation for open worlds.
- **Tiled Map Editor Auto-Import (`ImportTiledJSON`)**:
  - Automatically parses layers, animated tiles, and `solid: true` collision properties.
- **Code-First Map Creation (`NewStaticFromData`)**:
  - Instantly create tilemap layers from 2D slices (`[][]int`).
- **Viewport Culling (`DrawRegion`)**:
  - Culls non-visible tiles dynamically based on view bounds (`viewX, viewY, viewW, viewH`).
- **Intuitive Area Queries (`GetArea`)**:
  - Query regions via `GetArea(x, y, w, h)` to chain methods like `area.FindTiles(id)` and `area.ReplaceTile(oldID, newID)`.
- **Physics Integration (`CreateCollisionBoxes`)**:
  - Automatically merges adjacent `solid: true` tiles into optimized AABB collision boxes.
