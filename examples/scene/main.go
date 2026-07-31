package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/scene"
)

type TitleScene struct{}

func (s *TitleScene) Update(ctx *scene.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// ゲームメイン画面へ切替
		ctx.Start(&GameScene{score: 100})
	}
	return nil
}
func (s *TitleScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "[SCENE DEMO: Title Scene]\n\nPress SPACE to Start Game Scene")
}

type GameScene struct {
	score int
	pause *PauseOverlay
}

func (s *GameScene) Update(ctx *scene.Context) error {
	s.score++
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if s.pause == nil {
			s.pause = &PauseOverlay{}
			ctx.Overlay(s.pause)
		} else {
			ctx.Show(s.pause)
		}
		ctx.Stop(s) // ゲーム画面の更新を一時停止
	}
	return nil
}
func (s *GameScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("[SCENE DEMO: Game Scene]\nScore: %d\n\nPress P to Overlay Pause Screen", s.score))
}

type PauseOverlay struct{}

func (s *PauseOverlay) Update(ctx *scene.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ctx.Hide(s)        // ポーズ画面を非表示化
		ctx.Run(&GameScene{}) // ゲーム画面の更新を再開
	}
	return nil
}
func (s *PauseOverlay) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(640, 480)
	overlay.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(overlay, nil)
	ebitenutil.DebugPrint(screen, "[SCENE DEMO: Pause Overlay]\n\nPress ESC to Hide Overlay & Resume Game")
}

func main() {
	sm := scene.NewManager(640, 480)
	sm.Context().Start(&TitleScene{})

	ebiten.SetWindowTitle("scene Demo - Start, Overlay, Stop, Run, Hide, Show")
	if err := ebiten.RunGame(sm); err != nil {
		panic(err)
	}
}
