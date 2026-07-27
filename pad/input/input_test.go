package input

import (
	"testing"
)

func TestNewInput(t *testing.T) {
	input := NewInput()
	if input == nil {
		t.Fatal("NewInput() が nil を返しました")
	}
	if input.actions == nil {
		t.Error("NewInput() が actions マップを初期化しませんでした")
	}
}

func TestActionStateInitialValues(t *testing.T) {
	state := ActionState{}
	if state.Pressed != false {
		t.Error("ActionState.Pressed のデフォルト値は false である必要があります")
	}
	if state.JustPressed != false {
		t.Error("ActionState.JustPressed のデフォルト値は false である必要があります")
	}
	if state.JustReleased != false {
		t.Error("ActionState.JustReleased のデフォルト値は false である必要があります")
	}
	if state.X != 0 {
		t.Errorf("ActionState.X のデフォルト値は 0 である必要があります。現在の値: %f", state.X)
	}
	if state.Y != 0 {
		t.Errorf("ActionState.Y のデフォルト値は 0 である必要があります。現在の値: %f", state.Y)
	}
	if state.Strength != 0 {
		t.Errorf("ActionState.Strength のデフォルト値は 0 である必要があります。現在の値: %f", state.Strength)
	}
}

func TestInputQueries(t *testing.T) {
	const jump Action = 1
	input := NewInput()

	// 初期状態ではすべてのクエリが false を返す必要があります
	if input.Pressed(jump) {
		t.Error("初期状態の Pressed() は false である必要があります")
	}
	if input.JustPressed(jump) {
		t.Error("初期状態の JustPressed() は false である必要があります")
	}
	if input.JustReleased(jump) {
		t.Error("初期状態の JustReleased() は false である必要があります")
	}

	// ジャンプアクションの状態をモックします
	if input.actions[DefaultController] == nil {
		input.actions[DefaultController] = make(map[Action]*ActionState)
	}
	input.actions[DefaultController][jump] = &ActionState{
		Pressed:      true,
		JustPressed:  true,
		JustReleased: false,
	}

	// モックされた状態でのクエリ結果を確認します
	if !input.Pressed(jump) {
		t.Error("state.Pressed が true の場合、Pressed() は true である必要があります")
	}
	if !input.JustPressed(jump) {
		t.Error("state.JustPressed が true の場合、JustPressed() は true である必要があります")
	}
	if input.JustReleased(jump) {
		t.Error("state.JustReleased が false の場合、JustReleased() は false である必要があります")
	}

	// 別の状態をモックします
	input.actions[DefaultController][jump].Pressed = false
	input.actions[DefaultController][jump].JustPressed = false
	input.actions[DefaultController][jump].JustReleased = true

	if input.Pressed(jump) {
		t.Error("state.Pressed が false の場合、Pressed() は false である必要があります")
	}
	if input.JustPressed(jump) {
		t.Error("state.JustPressed が false の場合、JustPressed() は false である必要があります")
	}
	if !input.JustReleased(jump) {
		t.Error("state.JustReleased が true の場合、JustReleased() は true である必要があります")
	}
}
