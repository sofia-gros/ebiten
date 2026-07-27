# ebitenpad

ebitenpad 縺ｯ Ebitengine (ebiten) 蜷代￠縺ｮ縲√い繧ｯ繧ｷ繝ｧ繝ｳ繝吶・繧ｹ縺ｮ蜈･蜉帷ｮ｡逅・Λ繧､繝悶Λ繝ｪ縺ｧ縺吶・
繧ｭ繝ｼ繝懊・繝峨√ご繝ｼ繝繝代ャ繝峨√◎縺励※繝舌・繝√Ε繝ｫ繝代ャ繝会ｼ医ち繝・メ謫堺ｽ懶ｼ峨ｒ荳縺､縺ｮ隲也炊逧・↑縲後い繧ｯ繧ｷ繝ｧ繝ｳ縲阪↓邨ｱ蜷医＠縺ｦ謇ｱ縺・％縺ｨ縺後〒縺阪∪縺吶・

## 迚ｹ蠕ｴ

- **繧｢繧ｯ繧ｷ繝ｧ繝ｳ繝吶・繧ｹ縺ｮ謚ｽ雎｡蛹・*: 繧ｲ繝ｼ繝繝ｭ繧ｸ繝・け繧偵ョ繝舌う繧ｹ・医く繝ｼ繝懊・繝峨√ご繝ｼ繝繝代ャ繝峨√ち繝・メ・峨°繧牙・繧企屬縺励∪縺吶・
- **繝舌・繝√Ε繝ｫ繝代ャ繝牙・阡ｵ**: 繧ｹ繝・ぅ繝・け繧・・繧ｿ繝ｳ縺ｪ縺ｩ縺ｮUI繧堤ｰ｡蜊倥↓菴懈・縺励√・繝ｫ繝√ち繝・メ縺ｧ謫堺ｽ懷庄閭ｽ縺ｧ縺吶・
- **繧ｷ繝ｳ繝励Ν縺ｪAPI**: `Pressed`, `JustPressed`, `JustReleased` 縺ｪ縺ｩ縺ｮ逶ｴ諢溽噪縺ｪ繝｡繧ｽ繝・ラ繧呈署萓帙＠縺ｾ縺吶・
- **繝・ヰ繧､繧ｹ邨ｱ蜷・*: 1縺､縺ｮ繧｢繧ｯ繧ｷ繝ｧ繝ｳ縺ｫ蟇ｾ縺励※隍・焚縺ｮ繝・ヰ繧､繧ｹ蜈･蜉帙ｒ繝舌う繝ｳ繝峨〒縺阪∪縺吶・

## 繧｢繝・・繝・・繝域ュ蝣ｱ

v1.1.0 2026/07/06
- 1逕ｻ髱｢蛻・牡蟇ｾ謌ｦ縺ｮ繧医≧縺ｫ縲√・繝ｬ繧､繝､繝ｼ縺斐→縺ｫ蜈･蜉帙ｒ蛻・￠繧区ｩ溯・繧定ｿｽ蜉
- Strength縺悟ｮ滄圀縺ｮ蜈･蜉帛､繧剃ｽｿ逕ｨ縺吶ｋ繧医≧縺ｫ螟画峩

## 繧､繝ｳ繧ｹ繝医・繝ｫ

```bash
go get github.com/sofia-gros/ebiten/pad
```

## 蝓ｺ譛ｬ逧・↑菴ｿ縺・婿

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/pad/input"
)

const (
	ActionJump input.Action = 1
	ActionMove input.Action = 2
)

type Game struct {
	in *input.Input
}

func (g *Game) Update() error {
	g.in.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vx, vy := 0.0, 0.0
	strength := 0.0
	if state, ok := g.in.GetActionState(ActionMove); ok {
		vx, vy = state.X, state.Y
		strength = state.Strength
	}

	// 繝舌・繝√Ε繝ｫ UI 縺ｮ謠冗判
	g.in.Virtual().Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	in := input.NewInput()

	// Bind Keyboard
	in.BindKey(ActionJump, ebiten.KeySpace)
	in.BindKeyAxis(ActionMove, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	// Bind Gamepad
	in.BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom)
	in.BindGamepadAxis(ActionMove, 0, 1)

	// Virtual Pad縺ｮ險ｭ螳・
	vpad := in.Virtual()
	jumpBtn := vpad.AddButton().SetPosition(550, 400).SetRadius(40)
	moveStick := vpad.AddStick().SetPosition(100, 380).SetRadius(60)

	in.BindButton(ActionJump, jumpBtn)
	in.BindStick(ActionMove, moveStick)

	g := &Game{in: in}

	ebiten.SetWindowTitle("ebitenpad WASM Example")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
```

## 繝ｩ繧､繝悶Λ繝ｪ縺ｮ讒区・

- `input`: 繝｡繧､繝ｳ縺ｮ蜈･蜉帙・繝阪・繧ｸ繝｣繝ｼ縲ゅい繧ｯ繧ｷ繝ｧ繝ｳ縺ｮ繝舌う繝ｳ繝峨→迥ｶ諷狗ｮ｡逅・ｒ陦後＞縺ｾ縺吶・
- `virtual`: 繝舌・繝√Ε繝ｫ繧ｹ繝・ぅ繝・け繧・・繧ｿ繝ｳ縺ｪ縺ｩ縺ｮUI繧ｳ繝ｳ繝昴・繝阪Φ繝医・
- `examples/wasm`: WebAssembly 縺ｧ蜍穂ｽ懊☆繧九ョ繝｢繧ｳ繝ｼ繝峨・

## 隍・焚繝励Ξ繧､繝､繝ｼ蟇ｾ蠢・(`Controller` / `For`)

1逕ｻ髱｢蛻・牡蟇ｾ謌ｦ縺ｮ繧医≧縺ｫ縲√・繝ｬ繧､繝､繝ｼ縺斐→縺ｫ蜈･蜉帙ｒ蛻・￠縺溘＞蝣ｴ蜷医・ `Controller` 蝙九→ `For()` 繧剃ｽｿ縺・∪縺吶・

```go
const (
    P1 input.Controller = 0 // gamepad 0 縺ｫ繧ょｯｾ蠢・
    P2 input.Controller = 1 // gamepad 1 縺ｫ繧ょｯｾ蠢・
)

const (
    ActionJump input.Action = 1
    ActionMove input.Action = 2
)

// 繧ｻ繝・ヨ繧｢繝・・
in.For(P1).BindKey(ActionJump, ebiten.KeySpace)
in.For(P1).BindKeyAxis(ActionMove, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)
in.For(P1).BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom) // gamepad 0

in.For(P2).BindKey(ActionJump, ebiten.KeyEnter)
in.For(P2).BindKeyAxis(ActionMove, ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyUp, ebiten.KeyDown)
in.For(P2).BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom) // gamepad 1

// Update 縺ｯ1蝗槭□縺・
in.Update()

// 迥ｶ諷句叙蠕・(P1/P2 縺ｯ螳悟・縺ｫ迢ｬ遶九ょ酔譎ょ・蜉帙＠縺ｦ繧ゆｸ頑嶌縺阪↑縺・
if state, ok := g.in.For(P1).GetActionState(ActionMove); ok {
    p1vx, p1vy = state.X, state.Y
}
if state, ok := g.in.For(P2).GetActionState(ActionMove); ok {
    p2vx, p2vy = state.X, state.Y
}
```

1莠ｺ繝励Ξ繧､縺ｮ蝣ｴ蜷医・ `For()` 繧堤怐逡･縺吶ｋ縺縺代〒蠕捺擂騾壹ｊ縺ｫ蜍輔″縺ｾ縺吶ＡController` 縺ｮ蛟､縺ｯ縺昴・縺ｾ縺ｾ gamepad 縺ｮ繧､繝ｳ繝・ャ繧ｯ繧ｹ縺ｨ縺励※謇ｱ繧上ｌ縺ｾ縺吶・

```go
// 1莠ｺ繝励Ξ繧､: 蠕捺擂騾壹ｊ
in.BindKey(ActionJump, ebiten.KeySpace)
in.GetActionState(ActionJump)
```

## `state.Strength` 縺ｫ縺､縺・※

`ActionState.Strength` 縺ｯ繧｢繧ｯ繧ｷ繝ｧ繝ｳ縺ｮ蜈･蜉帙・蠑ｷ縺輔ｒ `0.0 縲・1.0` 縺ｧ陦ｨ縺励∪縺吶ゅョ繝舌う繧ｹ縺ｫ繧医▲縺ｦ險育ｮ玲婿豕輔′逡ｰ縺ｪ繧翫∪縺吶・

| 繝・ヰ繧､繧ｹ                   | 險育ｮ玲婿豕・                                                               |
| -------------------------- | ----------------------------------------------------------------------- |
| 繧ｭ繝ｼ繝懊・繝会ｼ亥腰繧ｭ繝ｼ・・      | 謚ｼ縺励※縺・ｌ縺ｰ蟶ｸ縺ｫ `1.0`                                                  |
| 繧ｭ繝ｼ繝懊・繝会ｼ郁ｻｸ・・          | `竏・dxﾂｲ + dyﾂｲ)` 繧・`1.0` 縺ｧ clamp縲よ万繧∝・蜉帙〒繧・`1.0` 縺御ｸ企剞             |
| 繧ｲ繝ｼ繝繝代ャ繝会ｼ医・繧ｿ繝ｳ・・    | 謚ｼ縺励※縺・ｌ縺ｰ蟶ｸ縺ｫ `1.0`                                                  |
| 繧ｲ繝ｼ繝繝代ャ繝会ｼ医い繝翫Ο繧ｰ霆ｸ・・| `竏・xﾂｲ + yﾂｲ)` 繧・`1.0` 縺ｧ clamp縲ゅい繝翫Ο繧ｰ蜈･蜉帙・蠑ｷ縺輔′螳滄圀縺ｫ蜿肴丐縺輔ｌ繧・   |
| 繝舌・繝√Ε繝ｫ繧ｹ繝・ぅ繝・け       | 荳ｭ蠢・°繧峨・霍晞屬繧貞濠蠕・〒蜑ｲ縺｣縺溷､縲よ欠繧貞ｰ代＠蜍輔°縺吶□縺代↑繧・`0.3` 縺ｨ縺九↓縺ｪ繧・|

## 繧ｲ繝ｼ繝繝代ャ繝峨・繝・ャ繝峨だ繝ｼ繝ｳ

繧ｹ繝・ぅ繝・け縺後ル繝･繝ｼ繝医Λ繝ｫ縺ｧ繧ょｾｮ螯吶↑蛟､縺悟・繧九・縺ｯ繧医￥縺ゅｋ隧ｱ縺ｧ縲～BindGamepadAxisWithDeadzone` 繧剃ｽｿ縺・→縺昴・霎ｺ繧貞・逅・〒縺阪∪縺吶・

```go
// Strength 縺・0.2 莉･荳九・蜈･蜉帙・辟｡隕悶☆繧・
in.BindGamepadAxisWithDeadzone(ActionMove, 0, 1, 0.2)
```

`BindGamepadAxis` 縺ｯ繝・ャ繝峨だ繝ｼ繝ｳ縺ｪ縺暦ｼ・0.0`・峨→蜷後§縺ｧ縺吶・

## 隍・焚繧ｲ繝ｼ繝繝代ャ繝・

謗･邯壹＆繧後※縺・ｋ繧ｲ繝ｼ繝繝代ャ繝峨ｒ蜈ｨ驛ｨ蟾｡蝗槭＠縺ｦ縲∽ｸ逡ｪ螟ｧ縺阪＞蜈･蜉帙ｒ謗｡逕ｨ縺励∪縺吶ゅさ繝ｳ繝医Ο繝ｼ繝ｩ繝ｼ繧・譛ｬ縺､縺ｪ縺・〒縺・※繧ゅ←縺｡繧峨°縺悟虚縺代・蜍輔￥縲√￥繧峨＞縺ｮ諢溘§縺ｧ縺吶・

## 繝ｩ繧､繧ｻ繝ｳ繧ｹ

MIT License
