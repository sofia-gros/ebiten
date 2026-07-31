package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/pad/input"
)

// Container は複数 Element の保持、自動配置レイアウト、GetAll / Get アクセス、一括非表示、カメラ・入力バインドを管理する最上位クラスです。
type Container struct {
	baseElement
	elements      []Element
	boundCamera   *camera.Camera
	ebitenMapping *EbitenInputMapping
	padInput      *input.Input
	padMapping    *PadInputMapping
}

// NewContainer は新しい Container を作成します。
func NewContainer() *Container {
	return &Container{
		baseElement: newBaseElement(0, 0, 640, 480),
		elements:    make([]Element, 0),
	}
}

// Add は Element をコンテナに追加します。
func (c *Container) Add(elem Element) {
	if elem != nil {
		c.elements = append(c.elements, elem)
	}
}

// GetAll はコンテナ内のすべての Element スライスを返します。
func (c *Container) GetAll() []Element {
	return c.elements
}

// Get はインデックス指定で子 Element を取得します。
func (c *Container) Get(index int) Element {
	if index < 0 || index >= len(c.elements) {
		return nil
	}
	return c.elements[index]
}

// Remove は指定した Element をコンテナから削除します。
func (c *Container) Remove(elem Element) {
	if elem == nil {
		return
	}
	newElems := make([]Element, 0, len(c.elements))
	for _, e := range c.elements {
		if e != elem {
			newElems = append(newElems, e)
		}
	}
	c.elements = newElems
}

// SetAllVisible はコンテナ内の全 Element を一括で表示 / 非表示切り替えします。
func (c *Container) SetAllVisible(visible bool) {
	for _, e := range c.elements {
		e.SetVisible(visible)
	}
}

// SetAllEnabled はコンテナ内の全 Element を一括で有効 / 無効切り替えします。
func (c *Container) SetAllEnabled(enabled bool) {
	for _, e := range c.elements {
		e.SetEnabled(enabled)
	}
}

// Clear はコンテナ内の全 Element を一括消去します。
func (c *Container) Clear() {
	c.elements = make([]Element, 0)
}

func (c *Container) Update() {
	if !c.visible || !c.enabled {
		return
	}
	for _, elem := range c.elements {
		elem.Update()
	}
}

func (c *Container) Draw(screen *ebiten.Image) {
	if !c.visible || screen == nil {
		return
	}
	for _, elem := range c.elements {
		elem.Draw(screen)
	}
}

// --- VBox (縦並びレイアウトコンテナ) ---

type VBox struct {
	Container
	spacing float64
	padding float64
}

func NewVBox() *VBox {
	return &VBox{
		Container: *NewContainer(),
		spacing:   10,
		padding:   10,
	}
}

func (v *VBox) SetSpacing(spacing float64) { v.spacing = spacing }
func (v *VBox) SetPadding(padding float64) { v.padding = padding }

func (v *VBox) Update() {
	if !v.visible || !v.enabled {
		return
	}

	// 縦並び整列計算
	currY := v.y + v.padding
	for _, elem := range v.elements {
		elem.SetPos(v.x+v.padding, currY)
		_, h := elem.Size()
		currY += h + v.spacing
		elem.Update()
	}
}
