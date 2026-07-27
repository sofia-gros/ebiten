package input

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestBindKey(t *testing.T) {
	const jump Action = 1
	input := NewInput()

	input.BindKey(jump, ebiten.KeySpace)

	if len(input.keyboard.keys) != 1 {
		t.Errorf("BindKey が正しく実行されませんでした。期待値: 1, 実際: %d", len(input.keyboard.keys))
	}

	if input.keyboard.keys[0].action != jump || input.keyboard.keys[0].key != ebiten.KeySpace {
		t.Error("バインドされたアクションまたはキーが正しくありません")
	}
}

func TestBindKeyAxis(t *testing.T) {
	const move Action = 1
	input := NewInput()

	input.BindKeyAxis(move, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	if len(input.keyboard.axes) != 1 {
		t.Errorf("BindKeyAxis が正しく実行されませんでした。期待値: 1, 実際: %d", len(input.keyboard.axes))
	}

	axis := input.keyboard.axes[0]
	if axis.action != move || axis.left != ebiten.KeyA || axis.right != ebiten.KeyD || axis.up != ebiten.KeyW || axis.down != ebiten.KeyS {
		t.Error("バインドされたアクションまたは軸キーが正しくありません")
	}
}

// mockKeyboardScanner はキーボードのモックスキャナーです。
type mockKeyboardScanner struct {
	pressedKeys map[ebiten.Key]bool
}

func (m *mockKeyboardScanner) IsKeyPressed(key ebiten.Key) bool {
	return m.pressedKeys[key]
}

func TestKeyAxisUpdateStrength(t *testing.T) {
	const move Action = 1
	in := NewInput()
	mock := &mockKeyboardScanner{pressedKeys: make(map[ebiten.Key]bool)}
	in.keyboardScanner = mock
	in.gamepadScanner = &mockNoGamepadScanner{}

	in.BindKeyAxis(move, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	// 単方向（右）: Strength は 1.0 であるべき
	mock.pressedKeys[ebiten.KeyD] = true
	in.Update()
	state, _ := in.GetActionState(move)
	if state.Strength != 1.0 {
		t.Errorf("単方向の Strength は 1.0 であるべきです。実際: %f", state.Strength)
	}
	if state.X != 1.0 {
		t.Errorf("右キーの X は 1.0 であるべきです。実際: %f", state.X)
	}

	// 斜め（右+下）: √(1²+1²) = √2 → clamp → 1.0
	mock.pressedKeys[ebiten.KeyS] = true
	in.Update()
	state, _ = in.GetActionState(move)
	if state.Strength != 1.0 {
		t.Errorf("斜め入力の Strength は 1.0 (clamp後) であるべきです。実際: %f", state.Strength)
	}

	// キーを離す: Strength は 0.0 に戻るべき
	mock.pressedKeys[ebiten.KeyD] = false
	mock.pressedKeys[ebiten.KeyS] = false
	in.Update()
	state, ok := in.GetActionState(move)
	if ok && state.Strength != 0.0 {
		t.Errorf("入力なしの Strength は 0.0 であるべきです。実際: %f", state.Strength)
	}
}

func TestKeyAxisStrengthSingleAxis(t *testing.T) {
	const move Action = 1
	in := NewInput()
	mock := &mockKeyboardScanner{pressedKeys: make(map[ebiten.Key]bool)}
	in.keyboardScanner = mock
	in.gamepadScanner = &mockNoGamepadScanner{}

	in.BindKeyAxis(move, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	// 上方向のみ: dx=0, dy=-1 → Strength=1.0
	mock.pressedKeys[ebiten.KeyW] = true
	in.Update()
	state, _ := in.GetActionState(move)
	if math.Abs(state.Strength-1.0) > 1e-9 {
		t.Errorf("上方向の Strength は 1.0 であるべきです。実際: %f", state.Strength)
	}
}

// mockNoGamepadScanner はゲームパッドが接続されていないスキャナーのモックです。
type mockNoGamepadScanner struct{}

func (m *mockNoGamepadScanner) AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID {
	return ids
}
func (m *mockNoGamepadScanner) IsStandardGamepadButtonPressed(_ ebiten.GamepadID, _ ebiten.StandardGamepadButton) bool {
	return false
}
func (m *mockNoGamepadScanner) StandardGamepadAxisValue(_ ebiten.GamepadID, _ ebiten.StandardGamepadAxis) float64 {
	return 0
}
