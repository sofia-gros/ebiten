package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TextInputOption は TextInput の初期化オプション構造体です。
type TextInputOption struct {
	Placeholder string
	Text        string
	Width       float64
	Height      float64
	MaxLength   int
	Grayscale   bool
}

// TextInput は キーボード文字入力・バックスペース・フォーカス状態を備えた入力 UI です。
type TextInput struct {
	baseElement
	placeholder string
	text        string
	maxLength   int
	isFocused   bool
	cursorTimer int
	onSubmit    func(text string)
	label       *Label
}

// NewTextInput は Option 構造体を指定して TextInput を作成します。
func NewTextInput(opts ...TextInputOption) *TextInput {
	opt := TextInputOption{
		Placeholder: "入力...",
		Width:       200,
		Height:      35,
		MaxLength:   30,
	}
	if len(opts) > 0 {
		userOpt := opts[0]
		opt.Placeholder = userOpt.Placeholder
		opt.Text = userOpt.Text
		if userOpt.Width > 0 {
			opt.Width = userOpt.Width
		}
		if userOpt.Height > 0 {
			opt.Height = userOpt.Height
		}
		if userOpt.MaxLength > 0 {
			opt.MaxLength = userOpt.MaxLength
		}
		opt.Grayscale = userOpt.Grayscale
	}

	ti := &TextInput{
		baseElement: newBaseElement(0, 0, opt.Width, opt.Height),
		placeholder: opt.Placeholder,
		text:        opt.Text,
		maxLength:   opt.MaxLength,
		label:       NewLabel(opt.Text),
	}
	ti.SetGrayscale(opt.Grayscale)
	return ti
}

func (t *TextInput) SetText(txt string) {
	t.text = txt
	if t.label != nil {
		t.label.SetText(txt)
	}
}

func (t *TextInput) Text() string { return t.text }

func (t *TextInput) OnSubmit(fn func(text string)) {
	t.onSubmit = fn
}

func (t *TextInput) Update() {
	if !t.visible || !t.enabled {
		t.isFocused = false
		return
	}

	t.cursorTimer++

	// 1. フォーカス判定
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		fx, fy := float64(mx), float64(my)
		if fx >= t.x && fx <= t.x+t.width && fy >= t.y && fy <= t.y+t.height {
			t.isFocused = true
		} else {
			t.isFocused = false
		}
	}

	if !t.isFocused {
		return
	}

	// 2. キーボード入力の受信
	runes := ebiten.AppendInputChars(nil)
	for _, r := range runes {
		if len(t.text) < t.maxLength {
			t.text += string(r)
		}
	}

	// 3. バックスペース
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) || ebiten.IsKeyPressed(ebiten.KeyBackspace) && t.cursorTimer%5 == 0 {
		if len(t.text) > 0 {
			t.text = t.text[:len(t.text)-1]
		}
	}

	// 4. Enter 送信
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if t.onSubmit != nil {
			t.onSubmit(t.text)
		}
	}

	if t.label != nil {
		t.label.SetText(t.text)
	}
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	if !t.visible || screen == nil {
		return
	}

	// 背景枠
	bgColor := color.RGBA{40, 40, 50, 255}
	if t.isFocused {
		bgColor = color.RGBA{60, 60, 80, 255}
	}

	bgImg := ebiten.NewImage(int(t.width), int(t.height))
	bgImg.Fill(bgColor)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(t.x, t.y)
	screen.DrawImage(bgImg, opts)

	// テキストまたはプレースホルダー描画
	displayText := t.text
	if displayText == "" && !t.isFocused {
		t.label.SetText(t.placeholder)
		t.label.SetColor(color.RGBA{130, 130, 150, 255})
	} else {
		t.label.SetText(displayText)
		t.label.SetColor(color.White)
	}

	t.label.SetPos(t.x+8, t.y+(t.height-16)/2)
	t.label.Draw(screen)
}
