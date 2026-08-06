# ebitensound

[日本語](./README.md)

`ebitensound` is an intuitive audio management library for Ebitengine.
It handles BGM and SE playback, volume control, crossfading, real-time 2D positional audio tracking, and fade-outs.

---

## Features

- **Extensible Sound Types (`Type`)**:
  - Extend built-in types (`sound.TypeSE`, `sound.TypeBGM`, `sound.TypeVoice`, `sound.TypeEnv`) using `sound.TypeCustom + iota`.
- **Per-Type Volume Control (`SetVolume`)**:
  - Adjust BGM or SE master volumes independently for settings screens.
- **Fade Stop & Crossfading (`Stop`, `CrossFade`)**:
  - Smoothly fade out tracks over a specified duration or crossfade seamlessly between BGMs.
- **2D Positional Audio (`PlayAt` / `SetPosition`)**:
  - Dynamically update panning and volume attenuation based on distances between sound sources and the listener.

---

## Installation

```bash
go get github.com/sofia-gros/ebiten/sound
```

---

## Usage

### Quick Start

Play sound effects instantly or fade out BGMs over a specified duration.


```go
package main

import (
	"github.com/sofia-gros/ebiten/sound"
)

func main() {
	m := sound.NewManager(44100)

	// Instant SE playback
	seTrack := m.Play(seBytes)
	_ = seTrack

	// BGM with looping and volume option
	bgmTrack := m.Play(bgmBytes, sound.Option{
		Type:   sound.TypeBGM,
		Volume: 0.8,
		Loop:   true,
	})

	// Fade out BGM over 1.5 seconds
	m.Stop(bgmTrack, 1.5)
}
```

---

### Full Usage

Adjust master volumes per sound type, track moving sound sources in 2D space, and crossfade BGMs across scenes.


```go
package main

import (
	"github.com/sofia-gros/ebiten/sound"
)

const (
	TypeBattle sound.Type = sound.TypeCustom + iota
)

type Game struct {
	m          *sound.Manager
	enemyTrack *sound.Track
}

func (g *Game) Init() {
	g.m = sound.NewManager(44100)

	// Master volume per sound type
	g.m.SetVolume(sound.TypeSE, 0.7)
	g.m.SetVolume(sound.TypeBGM, 0.5)

	// Play 2D positional audio
	g.enemyTrack = g.m.PlayAt(engineSEBytes, 200, 100, 100, 100, sound.Option{
		Type: sound.TypeSE,
		Loop: true,
	})
}

func (g *Game) Update() error {
	// Update 2D positions of the enemy and the player dynamically
	enemyX, enemyY := 350.0, 100.0
	playerX, playerY := 100.0, 100.0
	g.enemyTrack.SetPosition(enemyX, enemyY, playerX, playerY)
	return nil
}

func (g *Game) ChangeToBossScene(nextBgmBytes []byte) {
	// Crossfade to new BGM over 2.0 seconds
	g.m.CrossFade(nextBgmBytes, 2.0, sound.Option{
		Type: sound.TypeBGM,
		Loop: true,
	})
}
```

---

## License

MIT License
