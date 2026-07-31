package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Panel は NineSlice 9スライス背景枠を持った UI コンテナです。
type Panel struct {
	baseElement
	nineSlice *NineSlice
	elements  []Element
}

// NewPanel は指定サイズと NineSlice を指定して Panel を作成します。
func NewPanel(width, height float64, nineSlice *NineSlice) *Panel {
	return &Panel{
		baseElement: newBaseElement(0, 0, width, height),
		nineSlice:   nineSlice,
		elements:    make([]Element, 0),
	}
}

func (p *Panel) Add(elem Element) {
	if elem != nil {
		p.elements = append(p.elements, elem)
	}
}

func (p *Panel) Update() {
	if !p.visible || !p.enabled {
		return
	}
	for _, elem := range p.elements {
		elem.Update()
	}
}

func (p *Panel) Draw(screen *ebiten.Image) {
	if !p.visible || screen == nil {
		return
	}

	// 1. NineSlice 背景枠の描画
	if p.nineSlice != nil {
		p.nineSlice.Draw(screen, p.x, p.y, p.width, p.height)
	}

	// 2. 子要素の描画
	for _, elem := range p.elements {
		elem.Draw(screen)
	}
}
