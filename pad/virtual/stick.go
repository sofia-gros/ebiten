package virtual

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Stick はバーチャルアナログスティックを表します。
type Stick struct {
	x, y   float64
	radius float64

	inputX, inputY float64
	strength       float64

	touchID     ebiten.TouchID
	touchLocked bool
	mouseLocked bool
}

// SetPosition はスティックの中心位置を設定します。
func (s *Stick) SetPosition(x, y float64) *Stick {
	s.x, s.y = x, y
	return s
}

// SetRadius はスティックの稼働半径を設定します。
func (s *Stick) SetRadius(r float64) *Stick {
	s.radius = r
	return s
}

// Vector は現在の入力ベクトル (-1.0 ~ 1.0) を返します。
func (s *Stick) Vector() (x, y float64) {
	return s.inputX, s.inputY
}

// Strength は入力の強さ (0.0 ~ 1.0) を返します。
func (s *Stick) Strength() float64 {
	return s.strength
}

// Update はスティックの状態を更新します。
func (s *Stick) Update(touches []ebiten.TouchID) {
	// ロック中でなければ毅フレームの状態をリセット
	if !s.touchLocked && !s.mouseLocked {
		s.inputX = 0
		s.inputY = 0
		s.strength = 0
	}

	// タッチ入力処理
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
		// 新しいタッチ
		for _, id := range touches {
			tx, ty := ebiten.TouchPosition(id)
			if s.isInside(float64(tx), float64(ty)) {
				s.touchID = id
				s.touchLocked = true
				s.mouseLocked = false // タッチが優先
				s.updateInput(float64(tx), float64(ty))
				break
			}
		}
	}

	// タッチにロックされている場合はマウス入力を無視
	if s.touchLocked {
		return
	}

	// マウス入力（ドラッグでロック）
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
	return math.Sqrt(dx*dx+dy*dy) <= s.radius*1.5 // 遊びを持たせる
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

	// 正規化
	s.strength = math.Min(dist/s.radius, 1.0)
	s.inputX = (dx / dist) * s.strength
	s.inputY = (dy / dist) * s.strength
}

// Draw はデバッグ用にスティックを描画します。
func (s *Stick) Draw(screen *ebiten.Image) {
	// 背景の円
	vector.DrawFilledCircle(screen, float32(s.x), float32(s.y), float32(s.radius), color.RGBA{100, 100, 100, 128}, true)
	// 指の位置
	ix := s.x + s.inputX*s.radius
	iy := s.y + s.inputY*s.radius
	vector.DrawFilledCircle(screen, float32(ix), float32(iy), float32(s.radius*0.4), color.RGBA{255, 255, 255, 200}, true)
}
