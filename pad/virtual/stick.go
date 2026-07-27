package virtual

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Stick 縺ｯ繝舌・繝√Ε繝ｫ繧｢繝翫Ο繧ｰ繧ｹ繝・ぅ繝・け繧定｡ｨ縺励∪縺吶・
type Stick struct {
	x, y   float64
	radius float64

	inputX, inputY float64
	strength       float64

	touchID     ebiten.TouchID
	touchLocked bool
	mouseLocked bool
}

// SetPosition 縺ｯ繧ｹ繝・ぅ繝・け縺ｮ荳ｭ蠢・ｽ咲ｽｮ繧定ｨｭ螳壹＠縺ｾ縺吶・
func (s *Stick) SetPosition(x, y float64) *Stick {
	s.x, s.y = x, y
	return s
}

// SetRadius 縺ｯ繧ｹ繝・ぅ繝・け縺ｮ遞ｼ蜒榊濠蠕・ｒ險ｭ螳壹＠縺ｾ縺吶・
func (s *Stick) SetRadius(r float64) *Stick {
	s.radius = r
	return s
}

// Vector 縺ｯ迴ｾ蝨ｨ縺ｮ蜈･蜉帙・繧ｯ繝医Ν (-1.0 ~ 1.0) 繧定ｿ斐＠縺ｾ縺吶・
func (s *Stick) Vector() (x, y float64) {
	return s.inputX, s.inputY
}

// Strength 縺ｯ蜈･蜉帙・蠑ｷ縺・(0.0 ~ 1.0) 繧定ｿ斐＠縺ｾ縺吶・
func (s *Stick) Strength() float64 {
	return s.strength
}

// Update 縺ｯ繧ｹ繝・ぅ繝・け縺ｮ迥ｶ諷九ｒ譖ｴ譁ｰ縺励∪縺吶・
func (s *Stick) Update(touches []ebiten.TouchID) {
	// 繝ｭ繝・け荳ｭ縺ｧ縺ｪ縺代ｌ縺ｰ豈・ヵ繝ｬ繝ｼ繝縺ｮ迥ｶ諷九ｒ繝ｪ繧ｻ繝・ヨ
	if !s.touchLocked && !s.mouseLocked {
		s.inputX = 0
		s.inputY = 0
		s.strength = 0
	}

	// 繧ｿ繝・メ蜈･蜉帛・逅・
	if s.touchLocked {
		found := false
		for _, id := range touches {
			if id == s.touchID {
				found = true
				tx, ty := ebiten.TouchPosition(id)
				s.updateInput(float64(tx), float64(ty))
				break
			}
		}
		if !found {
			s.touchLocked = false
			s.inputX = 0
			s.inputY = 0
			s.strength = 0
		}
	} else {
		// 譁ｰ縺励＞繧ｿ繝・メ
		for _, id := range touches {
			tx, ty := ebiten.TouchPosition(id)
			if s.isInside(float64(tx), float64(ty)) {
				s.touchID = id
				s.touchLocked = true
				s.mouseLocked = false // 繧ｿ繝・メ縺悟━蜈・
				s.updateInput(float64(tx), float64(ty))
				break
			}
		}
	}

	// 繧ｿ繝・メ縺ｫ繝ｭ繝・け縺輔ｌ縺ｦ縺・ｋ蝣ｴ蜷医・繝槭え繧ｹ蜈･蜉帙ｒ辟｡隕・
	if s.touchLocked {
		return
	}

	// 繝槭え繧ｹ蜈･蜉幢ｼ医ラ繝ｩ繝・げ縺ｧ繝ｭ繝・け・・
	if s.mouseLocked {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.mouseLocked = false
			s.inputX = 0
			s.inputY = 0
			s.strength = 0
		} else {
			mx, my := ebiten.CursorPosition()
			s.updateInput(float64(mx), float64(my))
		}
	} else {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			if s.isInside(float64(mx), float64(my)) {
				s.mouseLocked = true
				s.updateInput(float64(mx), float64(my))
			}
		}
	}
}

func (s *Stick) isInside(x, y float64) bool {
	dx := x - s.x
	dy := y - s.y
	return math.Sqrt(dx*dx+dy*dy) <= s.radius*1.5 // 驕翫・繧呈戟縺溘○繧・
}

func (s *Stick) updateInput(tx, ty float64) {
	dx := tx - s.x
	dy := ty - s.y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist == 0 {
		s.inputX, s.inputY = 0, 0
		s.strength = 0
		return
	}

	// 豁｣隕丞喧
	s.strength = math.Min(dist/s.radius, 1.0)
	s.inputX = (dx / dist) * s.strength
	s.inputY = (dy / dist) * s.strength
}

// Draw 縺ｯ繝・ヰ繝・げ逕ｨ縺ｫ繧ｹ繝・ぅ繝・け繧呈緒逕ｻ縺励∪縺吶・
func (s *Stick) Draw(screen *ebiten.Image) {
	// 閭梧勹縺ｮ蜀・
	vector.DrawFilledCircle(screen, float32(s.x), float32(s.y), float32(s.radius), color.RGBA{100, 100, 100, 128}, true)
	// 謖・・菴咲ｽｮ
	ix := s.x + s.inputX*s.radius
	iy := s.y + s.inputY*s.radius
	vector.DrawFilledCircle(screen, float32(ix), float32(iy), float32(s.radius*0.4), color.RGBA{255, 255, 255, 200}, true)
}
