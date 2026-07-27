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
		t.Errorf("BindKey 縺梧ｭ｣縺励￥螳溯｡後＆繧後∪縺帙ｓ縺ｧ縺励◆縲よ悄蠕・､: 1, 螳滄圀: %d", len(input.keyboard.keys))
	}

	if input.keyboard.keys[0].action != jump || input.keyboard.keys[0].key != ebiten.KeySpace {
		t.Error("繝舌う繝ｳ繝峨＆繧後◆繧｢繧ｯ繧ｷ繝ｧ繝ｳ縺ｾ縺溘・繧ｭ繝ｼ縺梧ｭ｣縺励￥縺ゅｊ縺ｾ縺帙ｓ")
	}
}

func TestBindKeyAxis(t *testing.T) {
	const move Action = 1
	input := NewInput()

	input.BindKeyAxis(move, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	if len(input.keyboard.axes) != 1 {
		t.Errorf("BindKeyAxis 縺梧ｭ｣縺励￥螳溯｡後＆繧後∪縺帙ｓ縺ｧ縺励◆縲よ悄蠕・､: 1, 螳滄圀: %d", len(input.keyboard.axes))
	}

	axis := input.keyboard.axes[0]
	if axis.action != move || axis.left != ebiten.KeyA || axis.right != ebiten.KeyD || axis.up != ebiten.KeyW || axis.down != ebiten.KeyS {
		t.Error("繝舌う繝ｳ繝峨＆繧後◆繧｢繧ｯ繧ｷ繝ｧ繝ｳ縺ｾ縺溘・霆ｸ繧ｭ繝ｼ縺梧ｭ｣縺励￥縺ゅｊ縺ｾ縺帙ｓ")
	}
}

// mockKeyboardScanner 縺ｯ繧ｭ繝ｼ繝懊・繝峨・繝｢繝・け繧ｹ繧ｭ繝｣繝翫・縺ｧ縺吶・
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

	// 蜊俶婿蜷托ｼ亥承・・ Strength 縺ｯ 1.0 縺ｧ縺ゅｋ縺ｹ縺・
	mock.pressedKeys[ebiten.KeyD] = true
	in.Update()
	state, _ := in.GetActionState(move)
	if state.Strength != 1.0 {
		t.Errorf("蜊俶婿蜷代・ Strength 縺ｯ 1.0 縺ｧ縺ゅｋ縺ｹ縺阪〒縺吶ょｮ滄圀: %f", state.Strength)
	}
	if state.X != 1.0 {
		t.Errorf("蜿ｳ繧ｭ繝ｼ縺ｮ X 縺ｯ 1.0 縺ｧ縺ゅｋ縺ｹ縺阪〒縺吶ょｮ滄圀: %f", state.X)
	}

	// 譁懊ａ・亥承+荳具ｼ・ 竏・1ﾂｲ+1ﾂｲ) = 竏・ 竊・clamp 竊・1.0
	mock.pressedKeys[ebiten.KeyS] = true
	in.Update()
	state, _ = in.GetActionState(move)
	if state.Strength != 1.0 {
		t.Errorf("譁懊ａ蜈･蜉帙・ Strength 縺ｯ 1.0 (clamp蠕・ 縺ｧ縺ゅｋ縺ｹ縺阪〒縺吶ょｮ滄圀: %f", state.Strength)
	}

	// 繧ｭ繝ｼ繧帝屬縺・ Strength 縺ｯ 0.0 縺ｫ謌ｻ繧九∋縺・
	mock.pressedKeys[ebiten.KeyD] = false
	mock.pressedKeys[ebiten.KeyS] = false
	in.Update()
	state, ok := in.GetActionState(move)
	if ok && state.Strength != 0.0 {
		t.Errorf("蜈･蜉帙↑縺励・ Strength 縺ｯ 0.0 縺ｧ縺ゅｋ縺ｹ縺阪〒縺吶ょｮ滄圀: %f", state.Strength)
	}
}

func TestKeyAxisStrengthSingleAxis(t *testing.T) {
	const move Action = 1
	in := NewInput()
	mock := &mockKeyboardScanner{pressedKeys: make(map[ebiten.Key]bool)}
	in.keyboardScanner = mock
	in.gamepadScanner = &mockNoGamepadScanner{}

	in.BindKeyAxis(move, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	// 荳頑婿蜷代・縺ｿ: dx=0, dy=-1 竊・Strength=1.0
	mock.pressedKeys[ebiten.KeyW] = true
	in.Update()
	state, _ := in.GetActionState(move)
	if math.Abs(state.Strength-1.0) > 1e-9 {
		t.Errorf("荳頑婿蜷代・ Strength 縺ｯ 1.0 縺ｧ縺ゅｋ縺ｹ縺阪〒縺吶ょｮ滄圀: %f", state.Strength)
	}
}

// mockNoGamepadScanner 縺ｯ繧ｲ繝ｼ繝繝代ャ繝峨′謗･邯壹＆繧後※縺・↑縺・せ繧ｭ繝｣繝翫・縺ｮ繝｢繝・け縺ｧ縺吶・
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
