package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ButtonOption は Button の初期化オプション構造体です。
type ButtonOption struct {
	Text          string
	NormalImage   *ebiten.Image
	HoverImage    *ebiten.Image
	PressedImage  *ebiten.Image
	DisabledImage *ebiten.Image
	IconImage     *ebiten.Image
	Width         float64
	Height        float64
	Grayscale     bool
}

// Button は 通常/ホバー/押下/無効 の画像切り替えとクリックイベントを備えたボタンコンポーネントです。
type Button struct {
	baseElement
	text          string
	normalImage   *ebiten.Image
	hoverImage    *ebiten.Image
	pressedImage  *ebiten.Image
	disabledImage *ebiten.Image
	iconImage     *ebiten.Image
	label         *Label

	isHovered bool
	isPressed bool
	onClick   func()
}

// NewButton は Option 構造体を指定して Button を作成します。
func NewButton(opts ...ButtonOption) *Button {
	opt := ButtonOption{
		Width:  120,
		Height: 40,
	}
	if len(opts) > 0 {
		userOpt := opts[0]
		opt.Text = userOpt.Text
		opt.NormalImage = userOpt.NormalImage
		opt.HoverImage = userOpt.HoverImage
		opt.PressedImage = userOpt.PressedImage
		opt.DisabledImage = userOpt.DisabledImage
		opt.IconImage = userOpt.IconImage
		if userOpt.Width > 0 {
			opt.Width = userOpt.Width
		}
		if userOpt.Height > 0 {
			opt.Height = userOpt.Height
		}
		opt.Grayscale = userOpt.Grayscale
	}

	btn := &Button{
		baseElement:   newBaseElement(0, 0, opt.Width, opt.Height),
		text:          opt.Text,
		normalImage:   opt.NormalImage,
		hoverImage:    opt.HoverImage,
		pressedImage:  opt.PressedImage,
		disabledImage: opt.DisabledImage,
		iconImage:     opt.IconImage,
		label:         NewLabel(opt.Text),
	}
	btn.SetGrayscale(opt.Grayscale)
	return btn
}

// --- ゲッター・セッター ---

func (b *Button) SetText(txt string) {
	b.text = txt
	if b.label != nil {
		b.label.SetText(txt)
	}
}

func (b *Button) Text() string { return b.text }

func (b *Button) SetNormalImage(img *ebiten.Image)   { b.normalImage = img }
func (b *Button) NormalImage() *ebiten.Image         { return b.normalImage }

func (b *Button) SetHoverImage(img *ebiten.Image)    { b.hoverImage = img }
func (b *Button) HoverImage() *ebiten.Image          { return b.hoverImage }

func (b *Button) SetPressedImage(img *ebiten.Image)  { b.pressedImage = img }
func (b *Button) PressedImage() *ebiten.Image        { return b.pressedImage }

func (b *Button) SetDisabledImage(img *ebiten.Image) { b.disabledImage = img }
func (b *Button) DisabledImage() *ebiten.Image       { return b.disabledImage }

func (b *Button) SetIconImage(img *ebiten.Image)     { b.iconImage = img }
func (b *Button) IconImage() *ebiten.Image           { return b.iconImage }

// OnClick はボタンがクリック・決定された際の発火コールバックを登録します。
func (b *Button) OnClick(fn func()) {
	b.onClick = fn
}

func (b *Button) Update() {
	if !b.visible || !b.enabled {
		b.isHovered = false
		b.isPressed = false
		return
	}

	mx, my := ebiten.CursorPosition()
	fx, fy := float64(mx), float64(my)

	// マウス判定 (AABB)
	if fx >= b.x && fx <= b.x+b.width && fy >= b.y && fy <= b.y+b.height {
		b.isHovered = true

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.isPressed = true
		} else {
			if b.isPressed && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				if b.onClick != nil {
					b.onClick()
				}
			}
			b.isPressed = false
		}
	} else {
		b.isHovered = false
		b.isPressed = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	if !b.visible || screen == nil {
		return
	}

	// 1. 表示するテクスチャ画像の判定
	var currentImg *ebiten.Image

	if !b.enabled {
		if b.disabledImage != nil {
			currentImg = b.disabledImage
		} else {
			currentImg = b.normalImage
		}
	} else if b.isPressed && b.pressedImage != nil {
		currentImg = b.pressedImage
	} else if b.isHovered && b.hoverImage != nil {
		currentImg = b.hoverImage
	} else {
		currentImg = b.normalImage
	}

	// 2. ボタン背景の描画
	if currentImg != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(b.x, b.y)

		if b.grayscale {
			opts.ColorScale.ScaleWithColor(color.RGBA{150, 150, 150, 255})
		}
		screen.DrawImage(currentImg, opts)
	} else {
		// 画像が指定されていない場合は標準カラーレクト描画
		bgColor := color.RGBA{80, 80, 100, 255}
		if !b.enabled {
			bgColor = color.RGBA{50, 50, 50, 255}
		} else if b.isPressed {
			bgColor = color.RGBA{40, 40, 60, 255}
		} else if b.isHovered {
			bgColor = color.RGBA{110, 110, 140, 255}
		}

		rectImg := ebiten.NewImage(int(b.width), int(b.height))
		rectImg.Fill(bgColor)

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(b.x, b.y)
		screen.DrawImage(rectImg, opts)
	}

	// 3. アイコンとテキストの描画
	if b.iconImage != nil {
		iconOpts := &ebiten.DrawImageOptions{}
		iconOpts.GeoM.Translate(b.x+8, b.y+(b.height-float64(b.iconImage.Bounds().Dy()))/2)
		screen.DrawImage(b.iconImage, iconOpts)
	}

	if b.label != nil && b.text != "" {
		b.label.SetPos(b.x+b.width*0.2, b.y+(b.height-16)/2)
		b.label.Draw(screen)
	}
}
