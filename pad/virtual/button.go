package virtual

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Button 縺ｯ繝舌・繝√Ε繝ｫ UI 繝懊ち繝ｳ繧定｡ｨ縺励∪縺吶・
type Button struct {
	x, y   float64
	radius float64

	pressed bool
	touchID ebiten.TouchID
	locked  bool
}

// SetPosition 縺ｯ繝懊ち繝ｳ縺ｮ荳ｭ蠢・ｽ咲ｽｮ繧定ｨｭ螳壹＠縺ｾ縺吶・
func (b *Button) SetPosition(x, y float64) *Button {
	b.x, b.y = x, y
	return b
}

// SetRadius 縺ｯ繝懊ち繝ｳ縺ｮ蜊雁ｾ・ｒ險ｭ螳壹＠縺ｾ縺吶・
func (b *Button) SetRadius(r float64) *Button {
	b.radius = r
	return b
}

// Pressed 縺ｯ繝懊ち繝ｳ縺檎樟蝨ｨ謚ｼ縺輔ｌ縺ｦ縺・ｋ縺九←縺・°繧定ｿ斐＠縺ｾ縺吶・
func (b *Button) Pressed() bool {
	return b.pressed
}

// Update 縺ｯ繝懊ち繝ｳ縺ｮ迥ｶ諷九ｒ譖ｴ譁ｰ縺励∪縺吶・
func (b *Button) Update(touches []ebiten.TouchID) {
	// 蜑阪ヵ繝ｬ繝ｼ繝縺ｮ迥ｶ諷九・繝ｪ繧ｻ繝・ヨ
	if !b.locked {
		b.pressed = false
	}

	// 繝ｭ繝・け縺輔ｌ縺ｦ縺・ｋ蝣ｴ蜷医√◎縺ｮ繧ｿ繝・メ ID 縺後∪縺蟄伜惠縺吶ｋ縺狗｢ｺ隱・
	if b.locked {
		found := false
		for _, id := range touches {
			if id == b.touchID {
				found = true
				tx, ty := ebiten.TouchPosition(id)
				if b.isInside(float64(tx), float64(ty)) {
					b.pressed = true
				} else {
					// 遽・峇螟悶↓蜃ｺ縺溷ｴ蜷医・隗｣謾ｾ・郁ｨｭ險医↓繧医ｋ縺後√％縺薙〒縺ｯ繝ｭ繝・け隗｣髯､・・
					b.pressed = false
					b.locked = false
				}
				break
			}
		}
		if !found {
			b.pressed = false
			b.locked = false
		}
		return
	}

	// 譁ｰ縺励＞繧ｿ繝・メ縺ｮ讀懃ｴ｢
	for _, id := range touches {
		tx, ty := ebiten.TouchPosition(id)
		if b.isInside(float64(tx), float64(ty)) {
			b.pressed = true
			b.touchID = id
			b.locked = true
			break
		}
	}

	// 繝槭え繧ｹ蜈･蜉帙・繧ｵ繝昴・繝・(邁｡譏灘ｮ溯｣・
	if !b.locked && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if b.isInside(float64(mx), float64(my)) {
			b.pressed = true
		}
	}
}

func (b *Button) isInside(x, y float64) bool {
	dx := x - b.x
	dy := y - b.y
	return math.Sqrt(dx*dx+dy*dy) <= b.radius
}

// Draw 縺ｯ繝・ヰ繝・げ逕ｨ縺ｫ繝懊ち繝ｳ繧呈緒逕ｻ縺励∪縺吶・
func (b *Button) Draw(screen *ebiten.Image) {
	c := color.RGBA{200, 200, 200, 128}
	if b.pressed {
		c = color.RGBA{255, 255, 255, 200}
	}
	vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), float32(b.radius), c, true)
}
