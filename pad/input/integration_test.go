package input

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type mockScanner struct {
	pressedKeys    map[ebiten.Key]bool
	pressedButtons map[ebiten.StandardGamepadButton]bool
	gamepadAxes    map[int]float64
}

func (s *mockScanner) IsKeyPressed(key ebiten.Key) bool {
	return s.pressedKeys[key]
}

func (s *mockScanner) AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID {
	return append(ids, 0)
}

func (s *mockScanner) IsStandardGamepadButtonPressed(id ebiten.GamepadID, button ebiten.StandardGamepadButton) bool {
	return s.pressedButtons[button]
}

func (s *mockScanner) StandardGamepadAxisValue(id ebiten.GamepadID, axis ebiten.StandardGamepadAxis) float64 {
	return s.gamepadAxes[int(axis)]
}

func TestMultipleInputSourcesIntegration(t *testing.T) {
	const jump Action = 1
	input := NewInput()

	mock := &mockScanner{
		pressedKeys:    make(map[ebiten.Key]bool),
		pressedButtons: make(map[ebiten.StandardGamepadButton]bool),
		gamepadAxes:    make(map[int]float64),
	}
	input.keyboardScanner = mock
	input.gamepadScanner = mock

	// バインド設定
	input.BindKey(jump, ebiten.KeySpace)
	input.BindGamepadButton(jump, ebiten.StandardGamepadButtonCenterLeft)

	// フレーム 1: キーボードのみ
	mock.pressedKeys[ebiten.KeySpace] = true
	input.Update()

	if !input.Pressed(jump) {
		t.Error("Frame 1: Action should be Pressed (Keyboard)")
	}
	if !input.JustPressed(jump) {
		t.Error("Frame 1: Action should be JustPressed")
	}

	// フレーム 2: キーボード継続 + ゲームパッド追加
	mock.pressedButtons[ebiten.StandardGamepadButtonCenterLeft] = true
	input.Update()

	if !input.Pressed(jump) {
		t.Error("Frame 2: Action should be Pressed (Both)")
	}
	if input.JustPressed(jump) {
		t.Error("Frame 2: Action should NOT be JustPressed (Already Pressed)")
	}

	// フレーム 3: キーボード離す + ゲームパッド継続
	mock.pressedKeys[ebiten.KeySpace] = false
	input.Update()

	if !input.Pressed(jump) {
		t.Error("Frame 3: Action should be Pressed (Gamepad)")
	}
	if input.JustReleased(jump) {
		t.Error("Frame 3: Action should NOT be JustReleased (Still Pressed by gamepad)")
	}

	// フレーム 4: 全て離す
	mock.pressedButtons[ebiten.StandardGamepadButtonCenterLeft] = false
	input.Update()

	if input.Pressed(jump) {
		t.Error("Frame 4: Action should be released")
	}
	if !input.JustReleased(jump) {
		t.Error("Frame 4: Action should be JustReleased")
	}
}

func TestGamepadAxisStrength(t *testing.T) {
	const move Action = 1
	in := NewInput()

	mock := &mockScanner{
		pressedKeys:    make(map[ebiten.Key]bool),
		pressedButtons: make(map[ebiten.StandardGamepadButton]bool),
		gamepadAxes:    make(map[int]float64),
	}
	in.keyboardScanner = mock
	in.gamepadScanner = mock

	in.BindGamepadAxis(move, 0, 1)

	// 単軸の最大入力: x=1.0, y=0.0 → Strength=1.0
	mock.gamepadAxes[0] = 1.0
	mock.gamepadAxes[1] = 0.0
	in.Update()
	state, ok := in.GetActionState(move)
	if !ok {
		t.Fatal("ActionState が取得できませんでした")
	}
	if math.Abs(state.Strength-1.0) > 1e-9 {
		t.Errorf("単軸最大入力の Strength は 1.0 であるべきです。実際: %f", state.Strength)
	}

	// 斜め 45 度入力: x=0.5, y=0.5 → Strength=√(0.5²+0.5²)≒0.707
	mock.gamepadAxes[0] = 0.5
	mock.gamepadAxes[1] = 0.5
	in.Update()
	state, _ = in.GetActionState(move)
	expected := math.Sqrt(0.5*0.5 + 0.5*0.5)
	if math.Abs(state.Strength-expected) > 1e-9 {
		t.Errorf("斜め入力の Strength は %.4f であるべきです。実際: %f", expected, state.Strength)
	}

	// 入力なし: Strength=0.0
	mock.gamepadAxes[0] = 0.0
	mock.gamepadAxes[1] = 0.0
	in.Update()
	state, _ = in.GetActionState(move)
	if state.Strength != 0.0 {
		t.Errorf("入力なしの Strength は 0.0 であるべきです。実際: %f", state.Strength)
	}
}

func TestGamepadAxisDeadzone(t *testing.T) {
	const move Action = 1
	in := NewInput()

	mock := &mockScanner{
		pressedKeys:    make(map[ebiten.Key]bool),
		pressedButtons: make(map[ebiten.StandardGamepadButton]bool),
		gamepadAxes:    make(map[int]float64),
	}
	in.keyboardScanner = mock
	in.gamepadScanner = mock

	// デッドゾーン 0.2 でバインド
	in.BindGamepadAxisWithDeadzone(move, 0, 1, 0.2)

	// デッドゾーン以内の入力は無視される
	mock.gamepadAxes[0] = 0.1
	mock.gamepadAxes[1] = 0.0
	in.Update()
	if in.Pressed(move) {
		t.Error("デッドゾーン以内の入力は Pressed にならないべきです")
	}

	// デッドゾーンを超えた入力は有効
	mock.gamepadAxes[0] = 0.5
	mock.gamepadAxes[1] = 0.0
	in.Update()
	if !in.Pressed(move) {
		t.Error("デッドゾーンを超えた入力は Pressed になるべきです")
	}
	state, _ := in.GetActionState(move)
	if state.Strength <= 0 {
		t.Errorf("デッドゾーンを超えた入力の Strength は正の値であるべきです。実際: %f", state.Strength)
	}
}

// mockMultiGamepadScanner は複数のゲームパッドIDを返すモックです。
type mockMultiGamepadScanner struct {
	ids         []ebiten.GamepadID
	axesByID    map[ebiten.GamepadID]map[int]float64
	buttonsByID map[ebiten.GamepadID]map[ebiten.StandardGamepadButton]bool
}

func (s *mockMultiGamepadScanner) IsKeyPressed(_ ebiten.Key) bool { return false }
func (s *mockMultiGamepadScanner) AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID {
	return append(ids, s.ids...)
}
func (s *mockMultiGamepadScanner) IsStandardGamepadButtonPressed(id ebiten.GamepadID, button ebiten.StandardGamepadButton) bool {
	return s.buttonsByID[id][button]
}
func (s *mockMultiGamepadScanner) StandardGamepadAxisValue(id ebiten.GamepadID, axis ebiten.StandardGamepadAxis) float64 {
	if axes, ok := s.axesByID[id]; ok {
		return axes[int(axis)]
	}
	return 0
}

// TestMultipleGamepads は Controller ごとに別の gamepad ID が使われることを検証します。
func TestMultipleGamepads(t *testing.T) {
	const move Action = 1
	const P1 Controller = 0
	const P2 Controller = 1
	in := NewInput()

	mock := &mockMultiGamepadScanner{
		ids: []ebiten.GamepadID{0, 1},
		axesByID: map[ebiten.GamepadID]map[int]float64{
			0: {0: 0.8, 1: 0.0},   // P1 の gamepad
			1: {0: -0.5, 1: 0.6},  // P2 の gamepad
		},
		buttonsByID: map[ebiten.GamepadID]map[ebiten.StandardGamepadButton]bool{
			0: {},
			1: {},
		},
	}
	in.keyboardScanner = mock
	in.gamepadScanner = mock

	in.For(P1).BindGamepadAxis(move, 0, 1)
	in.For(P2).BindGamepadAxis(move, 0, 1)
	in.Update()

	// P1 は gamepad 0 の値を受け取る
	p1State, ok := in.For(P1).GetActionState(move)
	if !ok {
		t.Fatal("P1 の ActionState が取得できませんでした")
	}
	if math.Abs(p1State.X-0.8) > 1e-9 {
		t.Errorf("P1 の X は 0.8 であるべきです。実際: %f", p1State.X)
	}

	// P2 は gamepad 1 の値を受け取る
	p2State, ok := in.For(P2).GetActionState(move)
	if !ok {
		t.Fatal("P2 の ActionState が取得できませんでした")
	}
	if math.Abs(p2State.X-(-0.5)) > 1e-9 {
		t.Errorf("P2 の X は -0.5 であるべきです。実際: %f", p2State.X)
	}

	// P1 と P2 の状態は独立している
	if math.Abs(p1State.X-p2State.X) < 1e-9 {
		t.Error("P1 と P2 の X 値が同じになっています。独立しているべきです")
	}
}

// TestTwoPlayerInputSeparation は P1/P2 のキーボード入力が干渉しないことを検証します。
func TestTwoPlayerInputSeparation(t *testing.T) {
	const jump Action = 1
	const P1 Controller = 0
	const P2 Controller = 1
	in := NewInput()

	mock := &mockScanner{
		pressedKeys:    make(map[ebiten.Key]bool),
		pressedButtons: make(map[ebiten.StandardGamepadButton]bool),
		gamepadAxes:    make(map[int]float64),
	}
	in.keyboardScanner = mock
	in.gamepadScanner = mock

	// P1 は Space、P2 は Enter
	in.For(P1).BindKey(jump, ebiten.KeySpace)
	in.For(P2).BindKey(jump, ebiten.KeyEnter)

	// フレーム 1: P1 が Space を押す
	mock.pressedKeys[ebiten.KeySpace] = true
	in.Update()

	if !in.For(P1).Pressed(jump) {
		t.Error("Frame 1: P1 should be Pressed")
	}
	if in.For(P2).Pressed(jump) {
		t.Error("Frame 1: P2 should NOT be Pressed when only Space is held")
	}
	if !in.For(P1).JustPressed(jump) {
		t.Error("Frame 1: P1 should be JustPressed")
	}

	// フレーム 2: P2 も Enter を押す
	mock.pressedKeys[ebiten.KeyEnter] = true
	in.Update()

	if !in.For(P1).Pressed(jump) {
		t.Error("Frame 2: P1 should still be Pressed")
	}
	if !in.For(P2).Pressed(jump) {
		t.Error("Frame 2: P2 should be Pressed")
	}
	if !in.For(P2).JustPressed(jump) {
		t.Error("Frame 2: P2 should be JustPressed")
	}
	if in.For(P1).JustPressed(jump) {
		t.Error("Frame 2: P1 should NOT be JustPressed (already was pressed)")
	}

	// フレーム 3: P1 が Space を離す（P2 は継続）
	mock.pressedKeys[ebiten.KeySpace] = false
	in.Update()

	if in.For(P1).Pressed(jump) {
		t.Error("Frame 3: P1 should NOT be Pressed after Space release")
	}
	if !in.For(P1).JustReleased(jump) {
		t.Error("Frame 3: P1 should be JustReleased")
	}
	if !in.For(P2).Pressed(jump) {
		t.Error("Frame 3: P2 should still be Pressed")
	}
	if in.For(P2).JustReleased(jump) {
		t.Error("Frame 3: P2 should NOT be JustReleased")
	}
}
