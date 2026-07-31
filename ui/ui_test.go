package ui_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/ui"
)

func TestElementSettersAndGetters(t *testing.T) {
	btn := ui.NewButton(ui.ButtonOption{
		Text:      "Shop",
		Width:     150,
		Height:    40,
		Grayscale: true,
	})

	if btn.Text() != "Shop" {
		t.Errorf("expected text 'Shop', got '%s'", btn.Text())
	}

	btn.SetText("NewShop")
	if btn.Text() != "NewShop" {
		t.Errorf("expected text 'NewShop', got '%s'", btn.Text())
	}

	btn.SetPos(100, 200)
	x, y := btn.Pos()
	if x != 100 || y != 200 {
		t.Errorf("expected Pos (100, 200), got (%f, %f)", x, y)
	}

	btn.SetSize(200, 60)
	w, h := btn.Size()
	if w != 200 || h != 60 {
		t.Errorf("expected Size (200, 60), got (%f, %f)", w, h)
	}

	if !btn.IsGrayscale() {
		t.Errorf("expected IsGrayscale to be true")
	}

	btn.SetGrayscale(false)
	if btn.IsGrayscale() {
		t.Errorf("expected IsGrayscale to be false after SetGrayscale(false)")
	}
}

func TestButtonStateImages(t *testing.T) {
	normImg := ebiten.NewImage(100, 40)
	hoverImg := ebiten.NewImage(100, 40)
	pressedImg := ebiten.NewImage(100, 40)
	disabledImg := ebiten.NewImage(100, 40)

	btn := ui.NewButton()
	btn.SetNormalImage(normImg)
	btn.SetHoverImage(hoverImg)
	btn.SetPressedImage(pressedImg)
	btn.SetDisabledImage(disabledImg)

	if btn.NormalImage() != normImg {
		t.Errorf("NormalImage mismatch")
	}
	if btn.HoverImage() != hoverImg {
		t.Errorf("HoverImage mismatch")
	}
	if btn.PressedImage() != pressedImg {
		t.Errorf("PressedImage mismatch")
	}
	if btn.DisabledImage() != disabledImg {
		t.Errorf("DisabledImage mismatch")
	}
}

func TestContainerGetAllAndGet(t *testing.T) {
	container := ui.NewContainer()
	lbl1 := ui.NewLabel("Label 1")
	lbl2 := ui.NewLabel("Label 2")

	container.Add(lbl1)
	container.Add(lbl2)

	all := container.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 elements in GetAll(), got %d", len(all))
	}

	if container.Get(0) != lbl1 || container.Get(1) != lbl2 {
		t.Errorf("Get(index) mismatch")
	}

	container.SetAllVisible(false)
	if lbl1.Visible() || lbl2.Visible() {
		t.Errorf("expected SetAllVisible(false) to hide all elements")
	}

	container.Remove(lbl1)
	if len(container.GetAll()) != 1 {
		t.Errorf("expected 1 element after Remove, got %d", len(container.GetAll()))
	}
}

func TestVBoxAutoLayout(t *testing.T) {
	vbox := ui.NewVBox()
	vbox.SetPos(50, 50)
	vbox.SetSpacing(10)

	lbl1 := ui.NewLabel("Item 1")
	lbl1.SetSize(100, 20)

	lbl2 := ui.NewLabel("Item 2")
	lbl2.SetSize(100, 30)

	vbox.Add(lbl1)
	vbox.Add(lbl2)

	vbox.Update()

	// lbl1 ➔ (50+10, 50+10) = (60, 60)
	x1, y1 := lbl1.Pos()
	if x1 != 60 || y1 != 60 {
		t.Errorf("lbl1 position mismatch: got (%f, %f)", x1, y1)
	}

	// lbl2 ➔ Y = 60 + 20(高さ) + 10(Spacing) = 90
	x2, y2 := lbl2.Pos()
	if x2 != 60 || y2 != 90 {
		t.Errorf("lbl2 position mismatch: got (%f, %f)", x2, y2)
	}
}

func TestNineSliceDrawing(t *testing.T) {
	img := ebiten.NewImage(30, 30)
	img.Fill(color.RGBA{200, 100, 50, 255})

	ns := ui.NewNineSlice(img, 5, 5, 5, 5)
	target := ebiten.NewImage(100, 100)

	ns.Draw(target, 10, 10, 80, 80)
}
