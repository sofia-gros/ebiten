package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/pad/input"
	"github.com/sofia-gros/ebiten/scene"
)

type TitleScene struct {
	in *input.Input
}

func NewTitleScene(in *input.Input) *TitleScene {
	return &TitleScene{in: in}
}

func (s *TitleScene) Update(ctx *scene.Context) error {
	// タップまたはクリック、エンターキーなどでゲーム開始
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		// ゲームシーンへ遷移
		ctx.Start(NewGameScene(s.in))
		// GUIシーンをオーバーレイ（上に乗せる）
		ctx.Overlay(NewGUIScene(s.in))
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255})
	ebitenutil.DebugPrintAt(screen, "RPG Example", 640/2-40, 480/2-20)
	ebitenutil.DebugPrintAt(screen, "Tap or Click to Start", 640/2-70, 480/2+10)
}
