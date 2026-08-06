package virtual

import (
	"testing"
)

func TestButtonState(t *testing.T) {
	btn := &Button{}
	btn.SetPosition(20, 20)
	btn.SetRadius(15)

	if btn.Pressed() {
		t.Errorf("expected button not pressed initially")
	}
}

func TestStickState(t *testing.T) {
	stick := &Stick{}
	stick.SetPosition(100, 100)
	stick.SetRadius(40)

	vx, vy := stick.Vector()
	if vx != 0 || vy != 0 {
		t.Errorf("expected zero vector initially, got (%.2f, %.2f)", vx, vy)
	}

	if stick.Strength() != 0 {
		t.Errorf("expected zero strength initially")
	}
}

func TestVirtualPadContainer(t *testing.T) {
	vp := NewVirtualPad()
	btn := vp.AddButton()
	stick := vp.AddStick()

	if btn == nil || stick == nil {
		t.Fatalf("failed to add button or stick to virtual pad")
	}

	// Update と Draw メソッドがパニックしないか検証
	vp.Update()
}

