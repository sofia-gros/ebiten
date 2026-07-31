# ebitenasset

[日本語](./README.md)

`ebitenasset` is an intuitive asset loading and caching library for Ebitengine.
It provides a Phaser-inspired API to load, access, and unload images, sprite sheets, audio, tilemaps, animation definitions, and JSON/YAML data.

---

## Features

- **Simple Asset Management**: Register assets on a `Loader` and load them all at once.
- **Two Image Modes (`Image` vs `Sprite`)**:
  - `Image`: Textures for backgrounds or single illustrations.
  - `Sprite`: Sprite sheets with frame dimensions (`FrameWidth`, `FrameHeight`).
- **Type-Safe Data Unmarshaling (`DataAs[T]`)**:
  - Automatically unmarshal loaded JSON/YAML data files directly into Go structs.
- **Memory Management (`Unload` / `Clear`)**:
  - Easily unload individual assets or clear all cached memory upon scene transitions.

---

## Installation

```bash
go get github.com/sofia-gros/ebiten/asset
```

---

## Usage

### Quick Start

Load background images and JSON data files cleanly.


```go
package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/asset"
)

//go:embed assets/*
var assetsFS embed.FS

type GameConfig struct {
	Title  string  `json:"title"`
	Volume float64 `json:"volume"`
}

type Game struct {
	m *asset.Manager
}

func (g *Game) Init() {
	g.m = asset.NewManager(assetsFS)

	g.m.Image("bg", "assets/background.png")
	g.m.Data("config", "assets/config.json")

	if err := g.m.Load(); err != nil {
		panic(err)
	}

	var cfg GameConfig
	if err := asset.DataAs(g.m, "config", &cfg); err == nil {
		println("Title:", cfg.Title)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if bgImg := g.m.GetImage("bg"); bgImg != nil {
		screen.DrawImage(bgImg, nil)
	}
}
```

---

### Full Usage

Load sprite sheets with frame dimensions, audio tracks, and tilemaps.


```go
package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/asset"
)

//go:embed assets/*
var assetsFS embed.FS

type Game struct {
	m *asset.Manager
}

func (g *Game) Init() {
	g.m = asset.NewManager(assetsFS)

	// Sprite sheet with 32x32px frames
	g.m.Sprite("player", "assets/player_sheet.png", asset.SpriteOptions{
		FrameWidth:  32,
		FrameHeight: 32,
	})

	g.m.Animation("player_anim", "assets/player_anim.json")
	g.m.Tilemap("level1", "assets/level1.json")
	g.m.Audio("bgm_main", "assets/bgm.mp3")

	if err := g.m.Load(); err != nil {
		panic(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if sheet := g.m.GetSprite("player"); sheet != nil {
		frameImg := sheet.GetFrame(2) // Get frame index 2
		screen.DrawImage(frameImg, nil)
	}
}

func (g *Game) OnSceneChange() {
	g.m.Unload("level1")
}
```

---

## License

MIT License
