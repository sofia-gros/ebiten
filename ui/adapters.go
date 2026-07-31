package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/pad/input"
)

// EbitenInputMapping は Ebitengine の Key を直接バインドする構造体です。
type EbitenInputMapping struct {
	UpKey     ebiten.Key
	DownKey   ebiten.Key
	LeftKey   ebiten.Key
	RightKey  ebiten.Key
	SubmitKey ebiten.Key
}

// PadInputMapping は pad/input.Manager の Action (数値 uint) をバインドする構造体です。
type PadInputMapping struct {
	UpAction     input.Action
	DownAction   input.Action
	LeftAction   input.Action
	RightAction  input.Action
	SubmitAction input.Action
}

// BindCamera は UI コンテナへ Camera を接続し、カメラが振動・ズームしても UI を画面固定に追従固定します。
func (c *Container) BindCamera(cam *camera.Camera) {
	c.boundCamera = cam
}

// BindEbitenInput は Ebitengine 標準キーを UI のフォーカス＆決定操作に直接バインドします。
func (c *Container) BindEbitenInput(mapping EbitenInputMapping) {
	c.ebitenMapping = &mapping
}

// BindPadInput は pad/input.Manager の Action (数値) を UI のフォーカス＆決定操作にバインドします。
func (c *Container) BindPadInput(in *input.Input, mapping PadInputMapping) {
	c.padInput = in
	c.padMapping = &mapping
}
