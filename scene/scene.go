package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene は ebitenscene で管理される各画面（シーン）が実装すべきインターフェースです。
type Scene interface {
	// Update は毎フレーム呼ばれます。
	// Context を通じて他のシーンへの遷移や、自身の一時停止などを行います。
	Update(ctx *Context) error

	// Draw は画面の描画時に呼ばれます。
	Draw(screen *ebiten.Image)
}
