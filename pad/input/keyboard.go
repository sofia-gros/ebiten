package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// keyBinding はアクションにバインドされた単一のキーを表します。
type keyBinding struct {
	controller Controller
	action     Action
	key        ebiten.Key
}

// keyAxisBinding はアクションにバインドされた4方向のキーを表します。
type keyAxisBinding struct {
	controller Controller
	action     Action
	left       ebiten.Key
	right      ebiten.Key
	up         ebiten.Key
	down       ebiten.Key
}

// keyboardManager はキーボード入力を管理します。
type keyboardManager struct {
	keys []keyBinding
	axes []keyAxisBinding
}

// newKeyboardManager は新しい keyboardManager を作成します。
func newKeyboardManager() *keyboardManager {
	return &keyboardManager{
		keys: []keyBinding{},
		axes: []keyAxisBinding{},
	}
}

// update はキーボード入力をポーリングし、各アクションの状態を更新します。
func (m *keyboardManager) update(actions map[Controller]map[Action]*ActionState, scanner KeyboardScanner) {
	for _, b := range m.keys {
		state := getOrInitState(actions, b.controller, b.action)
		if scanner.IsKeyPressed(b.key) {
			state.Pressed = true
			state.Strength = 1.0
		}
	}

	for _, b := range m.axes {
		state := getOrInitState(actions, b.controller, b.action)
		var dx, dy float64
		if scanner.IsKeyPressed(b.left) {
			dx -= 1.0
		}
		if scanner.IsKeyPressed(b.right) {
			dx += 1.0
		}
		if scanner.IsKeyPressed(b.up) {
			dy -= 1.0
		}
		if scanner.IsKeyPressed(b.down) {
			dy += 1.0
		}

		if dx != 0 || dy != 0 {
			state.Pressed = true
			// すでに入力がある場合は合成する（簡易的に大きい方を採用）
			if math.Abs(dx) > math.Abs(state.X) {
				state.X = dx
			}
			if math.Abs(dy) > math.Abs(state.Y) {
				state.Y = dy
			}
			// 実際のベクトル長を Strength として計算（斜め入力でも 1.0 が上限）
			strength := math.Min(math.Sqrt(dx*dx+dy*dy), 1.0)
			if strength > state.Strength {
				state.Strength = strength
			}
		}
	}
}

func getOrInitState(actions map[Controller]map[Action]*ActionState, controller Controller, action Action) *ActionState {
	controllerActions, ok := actions[controller]
	if !ok {
		controllerActions = make(map[Action]*ActionState)
		actions[controller] = controllerActions
	}
	state, ok := controllerActions[action]
	if !ok {
		state = &ActionState{}
		controllerActions[action] = state
	}
	return state
}
