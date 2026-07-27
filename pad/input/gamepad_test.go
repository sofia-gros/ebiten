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
		t.Errorf("BindGamepadButton が正しく実行されませんでした。期待値: 1, 実際: %d", len(input.gamepad.buttons))
	}

	if input.gamepad.buttons[0].action != jump || input.gamepad.buttons[0].button != ebiten.StandardGamepadButtonCenterLeft {
		t.Error("バインドされたアクションまたはボタンが正しくありません")
	}
}

func TestBindGamepadAxis(t *testing.T) {
	const move Action = 1
	input := NewInput()

	input.BindGamepadAxis(move, 0, 1)

	if len(input.gamepad.axes) != 1 {
		t.Errorf("BindGamepadAxis が正しく実行されませんでした。期待値: 1, 実際: %d", len(input.gamepad.axes))
	}

	axis := input.gamepad.axes[0]
	if axis.action != move || axis.axisX != 0 || axis.axisY != 1 {
		t.Error("バインドされたアクションまたは軸が正しくありません")
	}
}
