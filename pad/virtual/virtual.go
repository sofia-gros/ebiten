package virtual

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// VirtualComponent はバーチャル UI コンポーネントの共通インターフェースです。
type VirtualComponent interface {
	Update(touches []ebiten.TouchID)
	Draw(screen *ebiten.Image)
}

// VirtualPad はバーチャルスティックとボタンの集合を管理します。
type VirtualPad struct {
	buttons []*Button
	sticks  []*Stick
}

// NewVirtualPad は新しい VirtualPad を作成します。
func NewVirtualPad() *VirtualPad {
	return &VirtualPad{
		buttons: []*Button{},
		sticks:  []*Stick{},
	}
}

// AddButton は新しいボタンを追加して返します。
func (v *VirtualPad) AddButton() *Button {
	b := &Button{}
	v.buttons = append(v.buttons, b)
	return b
}

// AddStick は新しいスティックを追加して返します。
func (v *VirtualPad) AddStick() *Stick {
	s := &Stick{}
	v.sticks = append(v.sticks, s)
	return s
}

// Update はすべてのコンポーネントの状態を更新します。
func (v *VirtualPad) Update() {
	touchIDs := ebiten.AppendTouchIDs(nil)
	
	for _, b := range v.buttons {
		b.Update(touchIDs)
	}
	for _, s := range v.sticks {
		s.Update(touchIDs)
	}
}

// Draw はすべてのコンポーネントを描画します。
func (v *VirtualPad) Draw(screen *ebiten.Image) {
	for _, b := range v.buttons {
		b.Draw(screen)
	}
	for _, s := range v.sticks {
		s.Draw(screen)
	}
}
