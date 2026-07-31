package ui

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var defaultFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err == nil {
		defaultFaceSource = s
	}
}


// LabelOption は Label の初期化オプション構造体です。
type LabelOption struct {
	Text      string
	Color     color.Color
	FontSize  float64
	Grayscale bool
}

// Label はテキストを表示する基本 UI コンポーネントです。
type Label struct {
	baseElement
	text     string
	textColor color.Color
	fontSize float64
}

// NewLabel はテキストとオプションを指定して Label を作成します。
func NewLabel(txt string, opts ...LabelOption) *Label {
	opt := LabelOption{
		Text:     txt,
		Color:    color.White,
		FontSize: 16,
	}
	if len(opts) > 0 {
		userOpt := opts[0]
		if userOpt.Text != "" {
			opt.Text = userOpt.Text
		}
		if userOpt.Color != nil {
			opt.Color = userOpt.Color
		}
		if userOpt.FontSize > 0 {
			opt.FontSize = userOpt.FontSize
		}
		opt.Grayscale = userOpt.Grayscale
	}

	lbl := &Label{
		baseElement: newBaseElement(0, 0, float64(len(opt.Text))*opt.FontSize*0.6, opt.FontSize*1.2),
		text:        opt.Text,
		textColor:   opt.Color,
		fontSize:    opt.FontSize,
	}
	lbl.SetGrayscale(opt.Grayscale)
	return lbl
}

// SetText はラベルのテキストを設定します。
func (l *Label) SetText(txt string) {
	l.text = txt
	l.width = float64(len(txt)) * l.fontSize * 0.6
}

// Text はラベルの現在のテキストを返します。
func (l *Label) Text() string {
	return l.text
}

// SetColor はテキストの色を設定します。
func (l *Label) SetColor(c color.Color) {
	l.textColor = c
}

// Color はテキストの色を返します。
func (l *Label) Color() color.Color {
	return l.textColor
}

func (l *Label) Update() {}

func (l *Label) Draw(screen *ebiten.Image) {
	if !l.visible || screen == nil || l.text == "" {
		return
	}

	// 簡易テキスト描画 (ebitenutil / text パッケージ)
	op := &text.DrawOptions{}
	op.GeoM.Translate(l.x, l.y)
	if l.grayscale {
		op.ColorScale.ScaleWithColor(color.RGBA{150, 150, 150, 255})
	} else {
		op.ColorScale.ScaleWithColor(l.textColor)
	}

	text.Draw(screen, l.text, &text.GoTextFace{
		Source: defaultFaceSource,
		Size:   l.fontSize,
	}, op)
}

