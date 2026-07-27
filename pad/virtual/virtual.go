package virtual

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// VirtualComponent 縺ｯ繝舌・繝√Ε繝ｫ UI 繧ｳ繝ｳ繝昴・繝阪Φ繝医・蜈ｱ騾壹う繝ｳ繧ｿ繝ｼ繝輔ぉ繝ｼ繧ｹ縺ｧ縺吶・
type VirtualComponent interface {
	Update(touches []ebiten.TouchID)
	Draw(screen *ebiten.Image)
}

// VirtualPad 縺ｯ繝舌・繝√Ε繝ｫ繧ｹ繝・ぅ繝・け縺ｨ繝懊ち繝ｳ縺ｮ髮・粋繧堤ｮ｡逅・＠縺ｾ縺吶・
type VirtualPad struct {
	buttons []*Button
	sticks  []*Stick
}

// NewVirtualPad 縺ｯ譁ｰ縺励＞ VirtualPad 繧剃ｽ懈・縺励∪縺吶・
func NewVirtualPad() *VirtualPad {
	return &VirtualPad{
		buttons: []*Button{},
		sticks:  []*Stick{},
	}
}

// AddButton 縺ｯ譁ｰ縺励＞繝懊ち繝ｳ繧定ｿｽ蜉縺励※霑斐＠縺ｾ縺吶・
func (v *VirtualPad) AddButton() *Button {
	b := &Button{}
	v.buttons = append(v.buttons, b)
	return b
}

// AddStick 縺ｯ譁ｰ縺励＞繧ｹ繝・ぅ繝・け繧定ｿｽ蜉縺励※霑斐＠縺ｾ縺吶・
func (v *VirtualPad) AddStick() *Stick {
	s := &Stick{}
	v.sticks = append(v.sticks, s)
	return s
}

// Update 縺ｯ縺吶∋縺ｦ縺ｮ繧ｳ繝ｳ繝昴・繝阪Φ繝医・迥ｶ諷九ｒ譖ｴ譁ｰ縺励∪縺吶・
func (v *VirtualPad) Update() {
	touchIDs := ebiten.AppendTouchIDs(nil)
	
	for _, b := range v.buttons {
		b.Update(touchIDs)
	}
	for _, s := range v.sticks {
		s.Update(touchIDs)
	}
}

// Draw 縺ｯ縺吶∋縺ｦ縺ｮ繧ｳ繝ｳ繝昴・繝阪Φ繝医ｒ謠冗判縺励∪縺吶・
func (v *VirtualPad) Draw(screen *ebiten.Image) {
	for _, b := range v.buttons {
		b.Draw(screen)
	}
	for _, s := range v.sticks {
		s.Draw(screen)
	}
}
