package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/ui"
)


type Game struct {
	uiRoot    *ui.Container
	cam       *camera.Camera
	statusMsg string
}

func (g *Game) Init() {
	g.uiRoot = ui.NewContainer()
	g.cam = camera.New(640, 480)

	// カメラバインド
	g.uiRoot.BindCamera(g.cam)

	// 1. パネル背景
	panel := ui.NewPanel(450, 360, nil)
	panel.SetPos(95, 60)

	// 2. 縦並びレイアウト (VBox)
	vbox := ui.NewVBox()
	vbox.SetPos(115, 80)
	vbox.SetSpacing(12)

	// 3. ラベル
	titleLabel := ui.NewLabel("ebitenui フル機能デモ", ui.LabelOption{FontSize: 18})

	// 4. ボタン (Normal / Hover / Pressed / Disabled)
	btn := ui.NewButton(ui.ButtonOption{
		Text:   "クリックしてください",
		Width:  200,
		Height: 36,
	})
	btn.OnClick(func() {
		g.statusMsg = "Button Clicked!"
	})

	// 5. 文字入力 TextInput
	textInput := ui.NewTextInput(ui.TextInputOption{
		Placeholder: "プレイヤー名を入力...",
		Width:       220,
		Height:      32,
	})
	textInput.OnSubmit(func(text string) {
		g.statusMsg = "Submitted Name: " + text
	})

	// 6. スライダー Slider
	slider := ui.NewSlider(ui.SliderOption{
		Min: 0, Max: 100, Value: 75, Width: 200, Height: 20,
	})
	slider.OnChange(func(val float64) {
		g.statusMsg = fmt.Sprintf("Slider Value: %.1f", val)
	})

	// 7. チェックボックス CheckBox
	checkBox := ui.NewCheckBox("BGMを有効化", true)
	checkBox.OnChange(func(checked bool) {
		g.statusMsg = fmt.Sprintf("CheckBox Enabled: %v", checked)
	})

	// 8. スクロールボックス ScrollBox
	scrollBox := ui.NewScrollBox(250, 100)
	scrollContent := ui.NewVBox()
	scrollContent.Add(ui.NewLabel("スクロール行 1"))
	scrollContent.Add(ui.NewLabel("スクロール行 2"))
	scrollContent.Add(ui.NewLabel("スクロール行 3"))
	scrollContent.Add(ui.NewLabel("スクロール行 4"))
	scrollContent.Add(ui.NewLabel("スクロール行 5"))
	scrollBox.Add(scrollContent)

	// Layout へ組み込み
	vbox.Add(titleLabel)
	vbox.Add(btn)
	vbox.Add(textInput)
	vbox.Add(slider)
	vbox.Add(checkBox)

	panel.Add(vbox)
	g.uiRoot.Add(panel)
	g.uiRoot.Add(scrollBox)

	g.statusMsg = "UI Ready. Try interacting with components!"
}

func (g *Game) Update() error {
	// H キー: コンテナ全体の一括非表示トグル (SetAllVisible)
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		current := g.uiRoot.Visible()
		g.uiRoot.SetAllVisible(!current)
		g.uiRoot.SetVisible(!current)
	}

	g.uiRoot.Update()
	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	// 背景グリッド描画
	screen.Fill(color.RGBA{30, 30, 40, 255})

	g.uiRoot.Draw(screen)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[UI DEMO]\nPress H: Toggle SetAllVisible | Status: %s", g.statusMsg), 10, 10)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("ui Demo - Button, TextInput, Slider, CheckBox, ScrollBox, Panel, VBox")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
