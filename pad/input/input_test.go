package input

import (
	"testing"
)

func TestNewInput(t *testing.T) {
	input := NewInput()
	if input == nil {
		t.Fatal("NewInput() 縺・nil 繧定ｿ斐＠縺ｾ縺励◆")
	}
	if input.actions == nil {
		t.Error("NewInput() 縺・actions 繝槭ャ繝励ｒ蛻晄悄蛹悶＠縺ｾ縺帙ｓ縺ｧ縺励◆")
	}
}

func TestActionStateInitialValues(t *testing.T) {
	state := ActionState{}
	if state.Pressed != false {
		t.Error("ActionState.Pressed 縺ｮ繝・ヵ繧ｩ繝ｫ繝亥､縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if state.JustPressed != false {
		t.Error("ActionState.JustPressed 縺ｮ繝・ヵ繧ｩ繝ｫ繝亥､縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if state.JustReleased != false {
		t.Error("ActionState.JustReleased 縺ｮ繝・ヵ繧ｩ繝ｫ繝亥､縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if state.X != 0 {
		t.Errorf("ActionState.X 縺ｮ繝・ヵ繧ｩ繝ｫ繝亥､縺ｯ 0 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺吶ら樟蝨ｨ縺ｮ蛟､: %f", state.X)
	}
	if state.Y != 0 {
		t.Errorf("ActionState.Y 縺ｮ繝・ヵ繧ｩ繝ｫ繝亥､縺ｯ 0 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺吶ら樟蝨ｨ縺ｮ蛟､: %f", state.Y)
	}
	if state.Strength != 0 {
		t.Errorf("ActionState.Strength 縺ｮ繝・ヵ繧ｩ繝ｫ繝亥､縺ｯ 0 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺吶ら樟蝨ｨ縺ｮ蛟､: %f", state.Strength)
	}
}

func TestInputQueries(t *testing.T) {
	const jump Action = 1
	input := NewInput()

	// 蛻晄悄迥ｶ諷九〒縺ｯ縺吶∋縺ｦ縺ｮ繧ｯ繧ｨ繝ｪ縺・false 繧定ｿ斐☆蠢・ｦ√′縺ゅｊ縺ｾ縺・
	if input.Pressed(jump) {
		t.Error("蛻晄悄迥ｶ諷九・ Pressed() 縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if input.JustPressed(jump) {
		t.Error("蛻晄悄迥ｶ諷九・ JustPressed() 縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if input.JustReleased(jump) {
		t.Error("蛻晄悄迥ｶ諷九・ JustReleased() 縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}

	// 繧ｸ繝｣繝ｳ繝励い繧ｯ繧ｷ繝ｧ繝ｳ縺ｮ迥ｶ諷九ｒ繝｢繝・け縺励∪縺・
	if input.actions[DefaultController] == nil {
		input.actions[DefaultController] = make(map[Action]*ActionState)
	}
	input.actions[DefaultController][jump] = &ActionState{
		Pressed:      true,
		JustPressed:  true,
		JustReleased: false,
	}

	// 繝｢繝・け縺輔ｌ縺溽憾諷九〒縺ｮ繧ｯ繧ｨ繝ｪ邨先棡繧堤｢ｺ隱阪＠縺ｾ縺・
	if !input.Pressed(jump) {
		t.Error("state.Pressed 縺・true 縺ｮ蝣ｴ蜷医￣ressed() 縺ｯ true 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if !input.JustPressed(jump) {
		t.Error("state.JustPressed 縺・true 縺ｮ蝣ｴ蜷医゛ustPressed() 縺ｯ true 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if input.JustReleased(jump) {
		t.Error("state.JustReleased 縺・false 縺ｮ蝣ｴ蜷医゛ustReleased() 縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}

	// 蛻･縺ｮ迥ｶ諷九ｒ繝｢繝・け縺励∪縺・
	input.actions[DefaultController][jump].Pressed = false
	input.actions[DefaultController][jump].JustPressed = false
	input.actions[DefaultController][jump].JustReleased = true

	if input.Pressed(jump) {
		t.Error("state.Pressed 縺・false 縺ｮ蝣ｴ蜷医￣ressed() 縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if input.JustPressed(jump) {
		t.Error("state.JustPressed 縺・false 縺ｮ蝣ｴ蜷医゛ustPressed() 縺ｯ false 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
	if !input.JustReleased(jump) {
		t.Error("state.JustReleased 縺・true 縺ｮ蝣ｴ蜷医゛ustReleased() 縺ｯ true 縺ｧ縺ゅｋ蠢・ｦ√′縺ゅｊ縺ｾ縺・)
	}
}
