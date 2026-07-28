# ebitenscene

[日本語](./README.md)

ebitenscene is a practical and robust scene manager library for Ebitengine (ebiten).
It features a `Hide/Show` mechanism to prevent memory overhead and bugs caused by unnecessary instance creation/destruction, and provides a safe way to pass data between scenes. It is perfectly suited for full-scale game development, such as RPGs and shoot 'em ups.

## Features

- **Type-Safe Scene Data Transmission**: By passing instances directly when launching the next scene, data like scores can be safely transferred with compile-time verification.
- **Memory-Friendly `Hide / Show`**: You can stop and resume just the update and draw cycles without destroying the scene, preserving its state (score, menu selection, etc.). This makes it highly resilient against rapid pause screen toggles.
- **Intuitive Method Set**: You can completely control all screen states using combinations of just `Start`, `Overlay`, `Hide`, `Show`, `Stop`, `Run`, and `Destroy`.
- **Safe Transitions via Lazy Evaluation**: Even if transition functions are called mid-update, the stack won't break immediately; the transition safely occurs on the next frame.

## Update History

Version 1.0.0 2026/07/27
- Initial release
- Introduced Manager-based type-safe scene management architecture

## Installation

```bash
go get github.com/sofia-gros/ebiten/scene
```

## Basic Usage

```go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/scene"
)

// --- 1. Title Scene ---
type TitleScene struct{}

func (s *TitleScene) Update(ctx *scene.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// Transition to main game (TitleScene gets destroyed)
		ctx.Start(&GameScene{})
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "TITLE SCENE\nPress ENTER to start")
}

// --- 2. Main Game Scene ---
type GameScene struct {
	Score int
	pause *PauseScene // Hold pointer to reuse
}

func (s *GameScene) Update(ctx *scene.Context) error {
	s.Score += 1 // Increment score

	// Toggle pause screen with P key
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if s.pause == nil {
			// Create on first time and Overlay on top
			s.pause = &PauseScene{}
			ctx.Overlay(s.pause)
		} else {
			// Second time onwards, reveal from hidden state and move to front
			ctx.Show(s.pause)
		}
		ctx.Stop(s) // Pause the game scene's Update
	}

	// Game over on X key, move to Result screen
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		// Type-safe because data is passed directly at launch!
		result := &ResultScene{FinalScore: s.Score}
		ctx.Start(result)
	}
	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "GAME SCENE\nPress P to Pause, X to Game Over")
}

// --- 3. Pause Scene ---
type PauseScene struct{}

func (s *PauseScene) Update(ctx *scene.Context) error {
	// Resume main game with Space key
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		ctx.Hide(s) // Hide the pause screen (instance preserved)
		
		// Find GameScene by type and resume its Update
		ctx.Run(&GameScene{})
	}
	return nil
}

func (s *PauseScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "\n\nPAUSE SCENE\nPress SPACE to resume")
}

// --- 4. Result Scene ---
type ResultScene struct {
	FinalScore int // Data passed from GameScene
}

func (s *ResultScene) Update(ctx *scene.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ctx.Start(&TitleScene{}) // Return to title
	}
	return nil
}

func (s *ResultScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "RESULT SCENE\nPress ESC to return Title")
}

// --- 5. Main Loop (Manager Setup) ---
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ebitenscene Example")

	// Create manager
	manager := scene.NewManager(640, 480)
	
	// Set the first scene at app launch
	manager.Context().Start(&TitleScene{})

	// Pass manager directly to ebiten.RunGame
	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
```

## Library Structure

- `Manager`: The central manager that integrates into Ebitengine's main loop (`ebiten.Game`). Handles scene stacks, draw order, and deferred commands.
- `Context`: The object passed to each scene's `Update`. Used by scenes to instruct screen transitions to the manager.
- `Scene`: The interface that each screen must implement (`Update`, `Draw`).

### Key Methods of Context

Screen operations called via `Context` are not reflected immediately; they are safely applied **at the end of that frame (specifically, before the next Update begins)**.

| Method | Description | Update State | Draw State | Primary Use Case |
| :--- | :--- | :---: | :---: | :--- |
| `Start(s)` | Destroys all existing scenes and launches `s` as the new root. | `true` | `true` | Complete screen transition, like title to main game. |
| `Overlay(s)` | Launches `s` layered on top while preserving the current scene. | `true` | `true` | Temporarily display pause screens or dialogs. |
| `Hide(s)` | Puts scene `s` in a "hidden" state (preserves data). | `false` | `false` | Hide an overlaid pause screen to return to game. |
| `Show(s)` | Resumes a hidden scene `s` and moves it to the front. | `true` | `true` | Call back a pause screen hidden by Hide. |
| `Stop(s)` | Pauses only the "movement (Update)" of scene `s`. | `false` | Unchanged | Pause the background game while showing a pause menu. |
| `Run(s)` | Resumes the "movement (Update)" of scene `s` stopped by Stop. | `true` | Unchanged | Resume background game when unpausing. |
| `Destroy(s)` | Completely removes and destroys scene `s` from the stack. | - | - | Remove a specific layer (scene) that is no longer needed. |

## License

MIT License
