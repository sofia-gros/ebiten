package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/pad/input"
	"github.com/sofia-gros/ebiten/scene"
)

const (
	ActionMove input.Action = 1
	ActionHit  input.Action = 2
	ActionUI   input.Action = 3
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten RPG Example (Scene + Pad + Physics)")

	// シーンマネージャーの生成
	manager := scene.NewManager(640, 480)

	// 入力マネージャーの生成と初期設定
	in := input.NewInput()
	in.BindKeyAxis(ActionMove, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)
	in.BindKey(ActionHit, ebiten.KeySpace)
	in.BindGamepadAxis(ActionMove, 0, 1)
	in.BindGamepadButton(ActionHit, ebiten.StandardGamepadButtonRightBottom)

	// タイトル画面からスタート
	manager.Context().Start(NewTitleScene(in))

	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
