package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SliderOption は Slider の初期化オプション構造体です。
type SliderOption struct {
	Min       float64
	Max       float64
	Value     float64
	Width     float64
	Height    float64
	Grayscale bool
}

// Slider は数値の段階・アタッチ変更が可能なスライダー UI コンポーネントです。
type Slider struct {
	baseElement
	min        float64
	max        float64
	value      float64
	isDragging bool
	onChange   func(val float64)
}

// NewSlider は Option 構造体を指定して Slider を作成します。
func NewSlider(opts ...SliderOption) *Slider {
	opt := SliderOption{
		Min:    0.0,
		Max:    1.0,
		Value:  0.5,
		Width:  150,
		Height: 20,
	}
	if len(opts) > 0 {
		userOpt := opts[0]
		opt.Min = userOpt.Min
		opt.Max = userOpt.Max
		opt.Value = userOpt.Value
		if userOpt.Width > 0 {
			opt.Width = userOpt.Width
		}
		if userOpt.Height > 0 {
			opt.Height = userOpt.Height
		}
		opt.Grayscale = userOpt.Grayscale
	}

	s := &Slider{
		baseElement: newBaseElement(0, 0, opt.Width, opt.Height),
		min:         opt.Min,
		max:         opt.Max,
		value:       opt.Value,
	}
	s.SetGrayscale(opt.Grayscale)
	return s
}

func (s *Slider) SetValue(v float64) {
	if v < s.min {
		v = s.min
	}
	if v > s.max {
		v = s.max
	}
	s.value = v
}

func (s *Slider) Value() float64 { return s.value }

func (s *Slider) OnChange(fn func(val float64)) {
	s.onChange = fn
}

func (s *Slider) Update() {
	if !s.visible || !s.enabled {
		s.isDragging = false
		return
	}

	mx, my := ebiten.CursorPosition()
	fx, fy := float64(mx), float64(my)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if fx >= s.x && fx <= s.x+s.width && fy >= s.y && fy <= s.y+s.height {
				s.isDragging = true
			}
		}

		if s.isDragging {
			relX := fx - s.x
			if relX < 0 {
				relX = 0
			}
			if relX > s.width {
				relX = s.width
			}

			ratio := relX / s.width
			newVal := s.min + (s.max-s.min)*ratio
			if newVal != s.value {
				s.value = newVal
				if s.onChange != nil {
					s.onChange(s.value)
				}
			}
		}
	} else {
		s.isDragging = false
	}
}

func (s *Slider) Draw(screen *ebiten.Image) {
	if !s.visible || screen == nil {
		return
	}

	// 1. スライダー背景バーの描画
	barHeight := 6.0
	barY := s.y + (s.height-barHeight)/2

	bgBar := ebiten.NewImage(int(s.width), int(barHeight))
	bgBar.Fill(color.RGBA{60, 60, 80, 255})

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(s.x, barY)
	screen.DrawImage(bgBar, opts)

	// 2. 進行分アクティブバーの描画
	ratio := (s.value - s.min) / (s.max - s.min)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	activeW := s.width * ratio
	if activeW > 0 {
		activeBar := ebiten.NewImage(int(activeW), int(barHeight))
		activeBar.Fill(color.RGBA{100, 180, 255, 255})

		aOpts := &ebiten.DrawImageOptions{}
		aOpts.GeoM.Translate(s.x, barY)
		screen.DrawImage(activeBar, aOpts)
	}

	// 3. ノブ (つまみ) の描画
	knobSize := 14.0
	knobX := s.x + activeW - knobSize/2
	knobY := s.y + (s.height-knobSize)/2

	knobImg := ebiten.NewImage(int(knobSize), int(knobSize))
	knobImg.Fill(color.White)

	kOpts := &ebiten.DrawImageOptions{}
	kOpts.GeoM.Translate(knobX, knobY)
	screen.DrawImage(knobImg, kOpts)
}
