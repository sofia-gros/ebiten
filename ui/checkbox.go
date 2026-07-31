package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// CheckBox は ON/OFF のブール値を切り替える UI コンポーネントです。
type CheckBox struct {
	baseElement
	text     string
	checked  bool
	onChange func(checked bool)
	label    *Label
}

// NewCheckBox はテキストと初期状態を指定して CheckBox を作成します。
func NewCheckBox(txt string, checked bool) *CheckBox {
	cb := &CheckBox{
		baseElement: newBaseElement(0, 0, 150, 24),
		text:        txt,
		checked:     checked,
		label:       NewLabel(txt),
	}
	return cb
}

func (c *CheckBox) SetText(txt string) {
	c.text = txt
	if c.label != nil {
		c.label.SetText(txt)
	}
}

func (c *CheckBox) Text() string { return c.text }

func (c *CheckBox) SetChecked(checked bool) {
	c.checked = checked
}

func (c *CheckBox) Checked() bool { return c.checked }

func (c *CheckBox) OnChange(fn func(checked bool)) {
	c.onChange = fn
}

func (c *CheckBox) Update() {
	if !c.visible || !c.enabled {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		fx, fy := float64(mx), float64(my)
		if fx >= c.x && fx <= c.x+c.width && fy >= c.y && fy <= c.y+c.height {
			c.checked = !c.checked
			if c.onChange != nil {
				c.onChange(c.checked)
			}
		}
	}
}

func (c *CheckBox) Draw(screen *ebiten.Image) {
	if !c.visible || screen == nil {
		return
	}

	// チェックボックス枠の描画
	boxSize := 18.0
	boxImg := ebiten.NewImage(int(boxSize), int(boxSize))
	boxImg.Fill(color.RGBA{60, 60, 80, 255})

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(c.x, c.y+(c.height-boxSize)/2)
	screen.DrawImage(boxImg, opts)

	// ON 時のマーク
	if c.checked {
		markSize := 10.0
		markImg := ebiten.NewImage(int(markSize), int(markSize))
		markImg.Fill(color.RGBA{100, 220, 120, 255})

		mOpts := &ebiten.DrawImageOptions{}
		mOpts.GeoM.Translate(c.x+(boxSize-markSize)/2, c.y+(c.height-markSize)/2)
		screen.DrawImage(markImg, mOpts)
	}

	// テキストラベル
	if c.label != nil {
		c.label.SetPos(c.x+boxSize+8, c.y+(c.height-16)/2)
		c.label.Draw(screen)
	}
}
