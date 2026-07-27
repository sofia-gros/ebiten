package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestBindGamepadButton(t *testing.T) {
	const jump Action = 1
	input := NewInput()

	input.BindGamepadButton(jump, ebiten.StandardGamepadButtonCenterLeft)

	if len(input.gamepad.buttons) != 1 {
		t.Errorf("BindGamepadButton 縺梧ｭ｣縺励￥螳溯｡後＆繧後∪縺帙ｓ縺ｧ縺励◆縲よ悄蠕・､: 1, 螳滄圀: %d", len(input.gamepad.buttons))
	}

	if input.gamepad.buttons[0].action != jump || input.gamepad.buttons[0].button != ebiten.StandardGamepadButtonCenterLeft {
		t.Error("繝舌う繝ｳ繝峨＆繧後◆繧｢繧ｯ繧ｷ繝ｧ繝ｳ縺ｾ縺溘・繝懊ち繝ｳ縺梧ｭ｣縺励￥縺ゅｊ縺ｾ縺帙ｓ")
	}
}

func TestBindGamepadAxis(t *testing.T) {
	const move Action = 1
	input := NewInput()

	input.BindGamepadAxis(move, 0, 1)

	if len(input.gamepad.axes) != 1 {
		t.Errorf("BindGamepadAxis 縺梧ｭ｣縺励￥螳溯｡後＆繧後∪縺帙ｓ縺ｧ縺励◆縲よ悄蠕・､: 1, 螳滄圀: %d", len(input.gamepad.axes))
	}

	axis := input.gamepad.axes[0]
	if axis.action != move || axis.axisX != 0 || axis.axisY != 1 {
		t.Error("繝舌う繝ｳ繝峨＆繧後◆繧｢繧ｯ繧ｷ繝ｧ繝ｳ縺ｾ縺溘・霆ｸ縺梧ｭ｣縺励￥縺ゅｊ縺ｾ縺帙ｓ")
	}
}
