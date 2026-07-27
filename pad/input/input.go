package input

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/pad/virtual"
)

// Action は論理的なゲームアクションを表すカスタム整数型です。
type Action int

// Controller はプレイヤーや入力グループを識別する型です。
// 値はそのまま gamepad のインデックスとして使われます。
type Controller int

// DefaultController は For() を省略した場合のデフォルトコントローラーです。
const DefaultController Controller = 0

// ActionState は論理的なゲームアクションの状態を表します。
type ActionState struct {
	Pressed      bool
	JustPressed  bool
	JustReleased bool

	X float64
	Y float64

	Strength float64

	lastPressed bool
}

// KeyboardScanner はキーボードの状態をスキャンするインターフェースです。
type KeyboardScanner interface {
	IsKeyPressed(key ebiten.Key) bool
}

// GamepadScanner はゲームパッドの状態をスキャンするインターフェースです。
type GamepadScanner interface {
	AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID
	IsStandardGamepadButtonPressed(id ebiten.GamepadID, button ebiten.StandardGamepadButton) bool
	StandardGamepadAxisValue(id ebiten.GamepadID, axis ebiten.StandardGamepadAxis) float64
}

// DefaultKeyboardScanner は ebiten の標準 API を使用するデフォルトのスキャナーです。
type DefaultKeyboardScanner struct{}

func (s *DefaultKeyboardScanner) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// DefaultGamepadScanner は ebiten の標準 API を使用するデフォルトのスキャナーです。
type DefaultGamepadScanner struct{}

func (s *DefaultGamepadScanner) AppendGamepadIDs(ids []ebiten.GamepadID) []ebiten.GamepadID {
	return ebiten.AppendGamepadIDs(ids)
}

func (s *DefaultGamepadScanner) IsStandardGamepadButtonPressed(id ebiten.GamepadID, button ebiten.StandardGamepadButton) bool {
	return ebiten.IsStandardGamepadButtonPressed(id, button)
}

func (s *DefaultGamepadScanner) StandardGamepadAxisValue(id ebiten.GamepadID, axis ebiten.StandardGamepadAxis) float64 {
	return ebiten.StandardGamepadAxisValue(id, axis)
}

type virtualButtonBinding struct {
	controller Controller
	action     Action
	button     *virtual.Button
}

type virtualStickBinding struct {
	controller Controller
	action     Action
	stick      *virtual.Stick
}

// Input はアクションベースの入力を管理するメインマネージャーです。
type Input struct {
	actions  map[Controller]map[Action]*ActionState
	keyboard *keyboardManager
	gamepad  *gamepadManager
	virtual  *virtual.VirtualPad

	virtualButtons []virtualButtonBinding
	virtualSticks  []virtualStickBinding

	keyboardScanner KeyboardScanner
	gamepadScanner  GamepadScanner
}

// ControllerView は特定のコントローラーにスコープされたビューです。
// in.For(P1) のように取得し、バインドやクエリに使います。
type ControllerView struct {
	input      *Input
	controller Controller
}

// For は指定されたコントローラーにスコープされた ControllerView を返します。
func (i *Input) For(c Controller) *ControllerView {
	return &ControllerView{input: i, controller: c}
}

// NewInput は新しい Input インスタンスを作成し、初期化します。
func NewInput() *Input {
	return &Input{
		actions:         make(map[Controller]map[Action]*ActionState),
		keyboard:        newKeyboardManager(),
		gamepad:         newGamepadManager(),
		virtual:         virtual.NewVirtualPad(),
		virtualButtons:  []virtualButtonBinding{},
		virtualSticks:   []virtualStickBinding{},
		keyboardScanner: &DefaultKeyboardScanner{},
		gamepadScanner:  &DefaultGamepadScanner{},
	}
}

// Virtual はバーチャル UI マネージャーを返します。
func (i *Input) Virtual() *virtual.VirtualPad {
	return i.virtual
}

// ---- ControllerView: バインドメソッド ----

// BindKey は単一のキーをアクションにバインドします。
func (v *ControllerView) BindKey(action Action, key ebiten.Key) {
	v.input.keyboard.keys = append(v.input.keyboard.keys, keyBinding{
		controller: v.controller,
		action:     action,
		key:        key,
	})
}

// BindKeyAxis は4つのキーをベクトルアクションとしてバインドします。
func (v *ControllerView) BindKeyAxis(action Action, left, right, up, down ebiten.Key) {
	v.input.keyboard.axes = append(v.input.keyboard.axes, keyAxisBinding{
		controller: v.controller,
		action:     action,
		left:       left,
		right:      right,
		up:         up,
		down:       down,
	})
}

// BindGamepadButton はゲームパッドのボタンをアクションにバインドします。
// Controller の値がそのまま gamepad インデックスとして使われます。
func (v *ControllerView) BindGamepadButton(action Action, button ebiten.StandardGamepadButton) {
	v.input.gamepad.buttons = append(v.input.gamepad.buttons, gamepadButtonBinding{
		controller: v.controller,
		action:     action,
		button:     button,
	})
}

// BindGamepadAxis はゲームパッドのアナログスティック軸をアクションにバインドします。
func (v *ControllerView) BindGamepadAxis(action Action, axisX, axisY int) {
	v.BindGamepadAxisWithDeadzone(action, axisX, axisY, 0.0)
}

// BindGamepadAxisWithDeadzone はデッドゾーン付きでゲームパッドの軸をバインドします。
// deadzone は 0.0〜1.0 の範囲で指定します。この値以下の入力は無視されます。
func (v *ControllerView) BindGamepadAxisWithDeadzone(action Action, axisX, axisY int, deadzone float64) {
	v.input.gamepad.axes = append(v.input.gamepad.axes, gamepadAxisBinding{
		controller: v.controller,
		action:     action,
		axisX:      axisX,
		axisY:      axisY,
		deadzone:   deadzone,
	})
}

// BindButton はバーチャルボタンをアクションにバインドします。
func (v *ControllerView) BindButton(action Action, button *virtual.Button) {
	v.input.virtualButtons = append(v.input.virtualButtons, virtualButtonBinding{
		controller: v.controller,
		action:     action,
		button:     button,
	})
}

// BindStick はバーチャルスティックをアクションにバインドします。
func (v *ControllerView) BindStick(action Action, stick *virtual.Stick) {
	v.input.virtualSticks = append(v.input.virtualSticks, virtualStickBinding{
		controller: v.controller,
		action:     action,
		stick:      stick,
	})
}

// ---- ControllerView: クエリメソッド ----

// GetActionState は指定されたアクションの状態のコピーを返します。
func (v *ControllerView) GetActionState(action Action) (ActionState, bool) {
	if controllerActions, ok := v.input.actions[v.controller]; ok {
		if state, ok := controllerActions[action]; ok {
			return *state, true
		}
	}
	return ActionState{}, false
}

// Pressed はアクションが現在押されている場合に true を返します。
func (v *ControllerView) Pressed(action Action) bool {
	if controllerActions, ok := v.input.actions[v.controller]; ok {
		if state, ok := controllerActions[action]; ok {
			return state.Pressed
		}
	}
	return false
}

// JustPressed は現在のフレームでアクションが押されたばかりの場合に true を返します。
func (v *ControllerView) JustPressed(action Action) bool {
	if controllerActions, ok := v.input.actions[v.controller]; ok {
		if state, ok := controllerActions[action]; ok {
			return state.JustPressed
		}
	}
	return false
}

// JustReleased は現在のフレームでアクションが離されたばかりの場合に true を返します。
func (v *ControllerView) JustReleased(action Action) bool {
	if controllerActions, ok := v.input.actions[v.controller]; ok {
		if state, ok := controllerActions[action]; ok {
			return state.JustReleased
		}
	}
	return false
}

// ---- Input: 後方互換メソッド (DefaultController に委譲) ----

// BindKey は単一のキーをアクションにバインドします。
func (i *Input) BindKey(action Action, key ebiten.Key) {
	i.For(DefaultController).BindKey(action, key)
}

// BindKeyAxis は4つのキーをベクトルアクションとしてバインドします。
func (i *Input) BindKeyAxis(action Action, left, right, up, down ebiten.Key) {
	i.For(DefaultController).BindKeyAxis(action, left, right, up, down)
}

// BindGamepadButton はゲームパッドのボタンをアクションにバインドします。
func (i *Input) BindGamepadButton(action Action, button ebiten.StandardGamepadButton) {
	i.For(DefaultController).BindGamepadButton(action, button)
}

// BindGamepadAxis はゲームパッドのアナログスティック軸をアクションにバインドします。
func (i *Input) BindGamepadAxis(action Action, axisX, axisY int) {
	i.For(DefaultController).BindGamepadAxis(action, axisX, axisY)
}

// BindGamepadAxisWithDeadzone はデッドゾーン付きでゲームパッドの軸をバインドします。
func (i *Input) BindGamepadAxisWithDeadzone(action Action, axisX, axisY int, deadzone float64) {
	i.For(DefaultController).BindGamepadAxisWithDeadzone(action, axisX, axisY, deadzone)
}

// BindButton はバーチャルボタンをアクションにバインドします。
func (i *Input) BindButton(action Action, button *virtual.Button) {
	i.For(DefaultController).BindButton(action, button)
}

// BindStick はバーチャルスティックをアクションにバインドします。
func (i *Input) BindStick(action Action, stick *virtual.Stick) {
	i.For(DefaultController).BindStick(action, stick)
}

// GetActionState は指定されたアクションの状態のコピーを返します。
func (i *Input) GetActionState(action Action) (ActionState, bool) {
	return i.For(DefaultController).GetActionState(action)
}

// Pressed はアクションが現在押されている場合に true を返します。
func (i *Input) Pressed(action Action) bool {
	return i.For(DefaultController).Pressed(action)
}

// JustPressed は現在のフレームでアクションが押されたばかりの場合に true を返します。
func (i *Input) JustPressed(action Action) bool {
	return i.For(DefaultController).JustPressed(action)
}

// JustReleased は現在のフレームでアクションが離されたばかりの場合に true を返します。
func (i *Input) JustReleased(action Action) bool {
	return i.For(DefaultController).JustReleased(action)
}

// ---- Update ----

// Update はすべてのアクションの状態を更新します。
func (i *Input) Update() {
	// 状態のリセット
	for _, controllerActions := range i.actions {
		for _, state := range controllerActions {
			lastPressed := state.Pressed
			state.Pressed = false
			state.JustPressed = false
			state.JustReleased = false
			state.X = 0
			state.Y = 0
			state.Strength = 0
			state.lastPressed = lastPressed
		}
	}

	// 各デバイスのポーリング
	i.keyboard.update(i.actions, i.keyboardScanner)
	i.gamepad.update(i.actions, i.gamepadScanner)
	i.updateVirtual()

	// JustPressed / JustReleased の確定
	for _, controllerActions := range i.actions {
		for _, state := range controllerActions {
			if state.Pressed && !state.lastPressed {
				state.JustPressed = true
			}
			if !state.Pressed && state.lastPressed {
				state.JustReleased = true
			}
		}
	}
}

func (i *Input) updateVirtual() {
	i.virtual.Update()

	for _, b := range i.virtualButtons {
		state := getOrInitState(i.actions, b.controller, b.action)
		if b.button.Pressed() {
			state.Pressed = true
			state.Strength = 1.0
		}
	}

	for _, s := range i.virtualSticks {
		state := getOrInitState(i.actions, s.controller, s.action)
		vx, vy := s.stick.Vector()
		if vx != 0 || vy != 0 {
			state.Pressed = true
			// 軸合成
			if math.Abs(vx) > math.Abs(state.X) {
				state.X = vx
			}
			if math.Abs(vy) > math.Abs(state.Y) {
				state.Y = vy
			}
			if s.stick.Strength() > state.Strength {
				state.Strength = s.stick.Strength()
			}
		}
	}
}
