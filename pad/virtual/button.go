package virtual

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Button はバーチャル UI ボタンを表します。
type Button struct {
	x, y   float64
	radius float64

	pressed bool
	touchID ebiten.TouchID
	locked  bool
}

// SetPosition はボタンの中心位置を設定します。
func (b *Button) SetPosition(x, y float64) *Button {
	b.x, b.y = x, y
	return b
}

// SetRadius はボタンの半径を設定します。
func (b *Button) SetRadius(r float64) *Button {
	b.radius = r
	return b
}

// Pressed はボタンが現在押されているかどうかを返します。
func (b *Button) Pressed() bool {
	return b.pressed
}

// Update はボタンの状態を更新します。
func (b *Button) Update(touches []ebiten.TouchID) {
	// 前フレームの状態のリセット
	if !b.locked {
		b.pressed = false
	}

	// ロックされている場合、そのタッチ ID がまだ存在するか確認
	if b.locked {
		found := false
		for _, id := range touches {
			if id == b.touchID {
				found = true
				tx, ty := ebiten.TouchPosition(id)
				if b.isInside(float64(tx), float64(ty)) {
					b.pressed = true
				} else {
					// 範囲外に出た場合は解放（設計によるが、ここではロック解除）
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

	// 新しいタッチの検索
	for _, id := range touches {
		tx, ty := ebiten.TouchPosition(id)
		if b.isInside(float64(tx), float64(ty)) {
			b.pressed = true
			b.touchID = id
			b.locked = true
			break
		}
	}

	// マウス入力のサポート (簡易実装)
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

// Draw はデバッグ用にボタンを描画します。
func (b *Button) Draw(screen *ebiten.Image) {
	c := color.RGBA{200, 200, 200, 128}
	if b.pressed {
		c = color.RGBA{255, 255, 255, 200}
	}
	vector.DrawFilledCircle(screen, float32(b.x), float32(b.y), float32(b.radius), c, true)
}
