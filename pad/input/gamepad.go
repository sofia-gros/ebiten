package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// gamepadButtonBinding はアクションにバインドされたゲームパッドのボタンを表します。
type gamepadButtonBinding struct {
	controller Controller
	action     Action
	button     ebiten.StandardGamepadButton
}

// gamepadAxisBinding はアクションにバインドされたゲームパッドの軸を表します。
type gamepadAxisBinding struct {
	controller Controller
	action     Action
	axisX      int
	axisY      int
	deadzone   float64 // 0.0 のときデッドゾーンなし
}

// gamepadManager はゲームパッド入力を管理します。
type gamepadManager struct {
	buttons []gamepadButtonBinding
	axes    []gamepadAxisBinding
}

// newGamepadManager は新しい gamepadManager を作成します。
func newGamepadManager() *gamepadManager {
	return &gamepadManager{
		buttons: []gamepadButtonBinding{},
		axes:    []gamepadAxisBinding{},
	}
}

// update はゲームパッド入力をポーリングし、各アクションの状態を更新します。
// Controller の値をそのまま gamepad ID として使用します。
func (m *gamepadManager) update(actions map[Controller]map[Action]*ActionState, scanner GamepadScanner) {
	ids := scanner.AppendGamepadIDs(nil)
	if len(ids) == 0 {
		return
	}

	// 接続済み gamepad ID をセットで管理
	connected := make(map[ebiten.GamepadID]bool, len(ids))
	for _, id := range ids {
		connected[id] = true
	}

	for _, b := range m.buttons {
		gamepadID := ebiten.GamepadID(b.controller)
		if !connected[gamepadID] {
			continue
		}
		state := getOrInitState(actions, b.controller, b.action)
		if scanner.IsStandardGamepadButtonPressed(gamepadID, b.button) {
			state.Pressed = true
			state.Strength = 1.0
		}
	}

	for _, b := range m.axes {
		gamepadID := ebiten.GamepadID(b.controller)
		if !connected[gamepadID] {
			continue
		}
		state := getOrInitState(actions, b.controller, b.action)
		x := scanner.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxis(b.axisX))
		y := scanner.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxis(b.axisY))

		// デッドゾーン処理: 指定値以下の微小な入力を無視する
		strength := math.Sqrt(x*x + y*y)
		if strength <= b.deadzone {
			continue
		}

		if x != 0 || y != 0 {
			state.Pressed = true
			// すでに入力がある場合は合成する（大きい方を採用）
			if math.Abs(x) > math.Abs(state.X) {
				state.X = x
			}
			if math.Abs(y) > math.Abs(state.Y) {
				state.Y = y
			}
			// 実際のベクトル長を Strength として計算（1.0 が上限）
			strength = math.Min(strength, 1.0)
			if strength > state.Strength {
				state.Strength = strength
			}
		}
	}
}
