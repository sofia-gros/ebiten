package ui

import "github.com/hajimehoshi/ebiten/v2"

// Element はすべての UI コンポーネントが実装する統一基底インターフェースです。
type Element interface {
	// 位置・サイズ操作
	SetPos(x, y float64)
	Pos() (float64, float64)
	SetSize(w, h float64)
	Size() (float64, float64)

	// 表示・有効・モノクロ制御
	SetVisible(visible bool)
	Visible() bool
	SetEnabled(enabled bool)
	Enabled() bool
	SetGrayscale(grayscale bool)
	IsGrayscale() bool

	// ライフサイクル & 描画
	Update()
	Draw(screen *ebiten.Image)
}

// baseElement は Element インターフェースの共通実装を提供する構造体です。
type baseElement struct {
	x         float64
	y         float64
	width     float64
	height    float64
	visible   bool
	enabled   bool
	grayscale bool
}

func newBaseElement(x, y, w, h float64) baseElement {
	return baseElement{
		x:         x,
		y:         y,
		width:     w,
		height:    h,
		visible:   true,
		enabled:   true,
		grayscale: false,
	}
}

func (b *baseElement) SetPos(x, y float64) {
	b.x = x
	b.y = y
}

func (b *baseElement) Pos() (float64, float64) {
	return b.x, b.y
}

func (b *baseElement) SetSize(w, h float64) {
	b.width = w
	b.height = h
}

func (b *baseElement) Size() (float64, float64) {
	return b.width, b.height
}

func (b *baseElement) SetVisible(visible bool) {
	b.visible = visible
}

func (b *baseElement) Visible() bool {
	return b.visible
}

func (b *baseElement) SetEnabled(enabled bool) {
	b.enabled = enabled
}

func (b *baseElement) Enabled() bool {
	return b.enabled
}

func (b *baseElement) SetGrayscale(grayscale bool) {
	b.grayscale = grayscale
}

func (b *baseElement) IsGrayscale() bool {
	return b.grayscale
}
