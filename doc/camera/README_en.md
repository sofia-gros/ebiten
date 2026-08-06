# ebitencamera

[日本語](./README.md)

`ebitencamera` is a 2D camera library designed for Ebitengine.
It provides pure coordinate control (`SetPos`), zoom, rotation, screen shaking, world bounds clipping, viewport cropping, automatic Custom Shader switching, and ZIndex sorting for multi-camera systems.

---

## Features

- **Pure Coordinate Control (`SetPos` / `Move` / `MoveTo`)**: Directly set or interpolate camera position via `cam.SetPos(x, y)` or `cam.MoveTo(x, y, speed)`.
- **Zero-Buffer Direct Rendering (`cam.Render`)**: Execute rendering inside the `cam.Render(screen, drawFunc)` closure. Compatible with custom render methods like `pad.Draw` or physics debug drawers.
- **Automatic Shader Mode Switching (`SetShader` / `ClearShader`)**: Call `cam.SetShader(shader, opts)` to switch rendering modes automatically without changing your draw calls.
- **Multi-Camera & ZIndex Automatic Sorting (`Group`)**: Manage multiple cameras (main screen, mini-map, split screen) and sort rendering order automatically via `ZIndex`.
- **Screen-to-World Coordinate Conversion & Culling**: Bidirectional coordinate conversion (`ScreenToWorld` / `WorldToScreen`) and visible world AABB bounds retrieval (`VisibleBounds()`).

---

## Installation

```bash
go get github.com/sofia-gros/ebiten/camera
```

---

## Usage

### Quick Start

Track a player and render within the camera closure.


```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
)

type Game struct {
	cam     *camera.Camera
	playerX float64
	playerY float64
}

func (g *Game) Init() {
	g.cam = camera.New(640, 480)
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0
	g.cam.MoveTo(g.playerX, g.playerY, 0.1)
	g.cam.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cam.Render(screen, func(target *ebiten.Image) {
		g.drawMap(target)
		g.drawPlayer(target)
	})
}
```

---

### Full Usage

Manage main and mini-map cameras using a `Group` with automatic shader switching.


```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
)

type Game struct {
	group      *camera.Group
	mainCam    *camera.Camera
	miniCam    *camera.Camera
	darkShader *ebiten.Shader
}

func (g *Game) Init() {
	g.mainCam = camera.New("main", 640, 480)
	g.mainCam.SetZIndex(0)

	g.miniCam = camera.New("mini", 640, 480)
	g.miniCam.SetViewport(460, 20, 160, 120)
	g.miniCam.SetZoom(0.2)
	g.miniCam.SetZIndex(10)

	g.group = camera.NewGroup(g.mainCam, g.miniCam)
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0
	g.mainCam.SetPos(playerX, playerY)
	g.miniCam.SetPos(playerX, playerY)

	if inDarkArea {
		g.mainCam.SetShader(g.darkShader)
	} else {
		g.mainCam.ClearShader()
	}

	g.group.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.group.Render(screen, func(cam *camera.Camera, target *ebiten.Image) {
		g.drawWorld(target)
	})
}
```

---

## API Reference

### `camera.Camera`
- `New(width, height)`: Creates a new camera with specified dimensions.
- `SetPos(x, y)` / `Pos() (x, y)`: Set / get camera center position.
- `Move(dx, dy)` / `MoveTo(targetX, targetY, speed)`: Move camera / smoothly interpolate towards target.
- `SetZoom(zoom)` / `Zoom()`: Set / get zoom level.
- `SetRotation(angle)` / `Rotation()`: Set / get rotation angle in radians.
- `Forward() (dirX, dirY)` / `Right() (rtX, rtY)`: Get camera forward / right unit vectors for WASD movement relative to camera orientation.
- `Direction() (dirX, dirY)`: Get camera facing unit vector.
- `SetBounds(minX, minY, maxX, maxY)`: Set world bounds for camera clamping.
- `SetViewport(x, y, w, h)`: Set output viewport rectangle (useful for mini-maps).
- `Shake(strength, durationSec)`: Trigger screen shake.
- `SetShader(shader, opts...)` / `ClearShader()`: Apply or clear custom shader.
- `ScreenToWorld(screenX, screenY)` / `WorldToScreen(worldX, worldY)`: Transform coordinates.
- `VisibleBounds()`: Get currently visible world bounds `(minX, minY, maxX, maxY)` for culling.
- `Render(screen, drawFunc)`: Execute draw closure with camera matrix and viewport applied.

---

## License


MIT License
