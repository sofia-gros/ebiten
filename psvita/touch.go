package psvita

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TouchPoint はタッチ位置とタッチ識別情報構造体です。
type TouchPoint struct {
	ID    int
	X     int
	Y     int
	Force float64 // タッチ圧力 (0.0 〜 1.0)
}

// TouchFront は前面タッチパネル上のアクティブなタッチポイント一覧を返します。
func TouchFront() []TouchPoint {
	return getTouchFrontImpl()
}

// TouchBack は背面タッチパッド上のアクティブなタッチポイント一覧を返します。
func TouchBack() []TouchPoint {
	return getTouchBackImpl()
}

// フォールバック実装 (PC/Web標準用)
func getTouchFrontImpl() []TouchPoint {
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) == 0 {
		ids = ebiten.AppendTouchIDs(nil)
	}

	points := make([]TouchPoint, 0, len(ids))
	for _, id := range ids {
		x, y := ebiten.TouchPosition(id)
		points = append(points, TouchPoint{
			ID:    int(id),
			X:     x,
			Y:     y,
			Force: 1.0,
		})
	}

	// マウス操作の代替（タッチ入力が無い場合）
	if len(points) == 0 && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		points = append(points, TouchPoint{
			ID:    999,
			X:     mx,
			Y:     my,
			Force: 1.0,
		})
	}

	return points
}

func getTouchBackImpl() []TouchPoint {
	points := make([]TouchPoint, 0)
	// PC上での背面タッチ模擬: 右クリックまたは Shift キーを押しながらの左クリック
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) ||
		(ebiten.IsKeyPressed(ebiten.KeyShift) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
		mx, my := ebiten.CursorPosition()
		points = append(points, TouchPoint{
			ID:    888,
			X:     mx,
			Y:     my,
			Force: 0.8,
		})
	}
	return points
}

func isVitaTarget() bool {
	return false
}

func isButtonPressedImpl(button Button) bool {
	// PC環境でのボタン模擬
	switch button {
	case ButtonL1:
		return ebiten.IsKeyPressed(ebiten.KeyQ)
	case ButtonR1:
		return ebiten.IsKeyPressed(ebiten.KeyE)
	case ButtonSelect:
		return ebiten.IsKeyPressed(ebiten.KeyBackspace)
	case ButtonStart:
		return ebiten.IsKeyPressed(ebiten.KeyEnter)
	case ButtonUp:
		return ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)
	case ButtonRight:
		return ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	case ButtonDown:
		return ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)
	case ButtonLeft:
		return ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	case ButtonTriangle:
		return ebiten.IsKeyPressed(ebiten.KeyI)
	case ButtonCircle:
		return ebiten.IsKeyPressed(ebiten.KeyL)
	case ButtonCross:
		return ebiten.IsKeyPressed(ebiten.KeyK)
	case ButtonSquare:
		return ebiten.IsKeyPressed(ebiten.KeyJ)
	case ButtonRearTouchPad:
		return len(getTouchBackImpl()) > 0
	default:
		return false
	}
}
