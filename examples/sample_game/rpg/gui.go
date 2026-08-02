package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/pad/input"
	"github.com/sofia-gros/ebiten/scene"
)

type GUIScene struct {
	in *input.Input
	showPad bool
}

func NewGUIScene(in *input.Input) *GUIScene {
	s := &GUIScene{
		in: in,
		showPad: true,
	}

	// バーチャルパッドのレイアウト設定
	vpad := in.Virtual()
	// 移動スティックを左下に
	moveStick := vpad.AddStick().SetPosition(100, 380).SetRadius(60)
	// アクションボタンを右下に
	hitBtn := vpad.AddButton().SetPosition(540, 380).SetRadius(40)
	
	// UI切り替えボタンを右上に（表示切替用）
	toggleBtn := vpad.AddButton().SetPosition(590, 50).SetRadius(25)

	in.BindStick(ActionMove, moveStick)
	in.BindButton(ActionHit, hitBtn)
	in.BindButton(ActionUI, toggleBtn)

	return s
}

func (s *GUIScene) Update(ctx *scene.Context) error {
	// 全ての入力更新はGUIシーンが担当する
	s.in.Update()

	// UIトグルボタンが押されたらパッドの表示非表示を切り替える
	if s.in.JustPressed(ActionUI) {
		s.showPad = !s.showPad
	}
	return nil
}

func (s *GUIScene) Draw(screen *ebiten.Image) {
	if s.showPad {
		s.in.Virtual().Draw(screen)
	}
}
