package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScrollBox は指定領域 (width, height) 内に子要素をスクロール表示する切り抜きコンテナです。
type ScrollBox struct {
	baseElement
	contentX    float64
	contentY    float64
	elements    []Element
	maxScrollY  float64
}

// NewScrollBox は表示枠の幅と高さを指定して ScrollBox を作成します。
func NewScrollBox(width, height float64) *ScrollBox {
	return &ScrollBox{
		baseElement: newBaseElement(0, 0, width, height),
		elements:    make([]Element, 0),
	}
}

func (s *ScrollBox) Add(elem Element) {
	if elem != nil {
		s.elements = append(s.elements, elem)
	}
}

func (s *ScrollBox) Update() {
	if !s.visible || !s.enabled {
		return
	}

	// マウスホイールによる上下スクロール
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		s.contentY += wheelY * 20.0
		if s.contentY > 0 {
			s.contentY = 0
		}
	}

	for _, elem := range s.elements {
		elem.Update()
	}
}

func (s *ScrollBox) Draw(screen *ebiten.Image) {
	if !s.visible || screen == nil {
		return
	}

	// 描画領域を ScrollBox の矩形枠にクリッピング
	vx, vy := int(s.x), int(s.y)
	vw, vh := int(s.width), int(s.height)

	subScreen := screen.SubImage(image.Rect(vx, vy, vx+vw, vy+vh)).(*ebiten.Image)

	for _, elem := range s.elements {
		// 子要素にオフセット位置を反映して描画
		origX, origY := elem.Pos()
		elem.SetPos(origX+s.contentX, origY+s.contentY)
		elem.Draw(subScreen)
		elem.SetPos(origX, origY)
	}
}
