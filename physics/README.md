# ebitenphysics

[English](./README_en.md)

ebitenphysics は Ebitengine 向けの 2D 物理演算ラッパーライブラリです。
RPGやプラットフォーマー向けの軽量な `Arcade`（AABBめり込み防止）モードと、本格的な `Box2D`（剛体物理・回転・跳ね返り）モードを、共通の直感的な API でシームレスに切り替えて利用できます。

---

## 特徴

- **ピュア＆軽量な設計**: `physics` 本体には特定物理エンジンの実装を含まず、使いたいアダプター（`arcade`, `box2d` 等）だけを注入して利用。
- **シンプルな2通りのボディ作成**:
  - **簡易スタイル**: 1行で四角や円の物理ボディを作成。
  - **本格スタイル**: `BodyOptions` で衝突コールバック・複合当たり判定・グループ番号等を詳細設定。
- **柔軟な描画**: 位置座標 (`x, y`) を取得して手動描画するアプローチと、画像と紐付けて `physManager.Draw(screen)` で一括描画するアプローチの両方に対応。

---

## インストール

本体をインストールします：

```bash
go get github.com/sofia-gros/ebiten/physics
```

次に、利用したい物理エンジンのアダプターをインストールします：

| アダプター | 概要 | 依存エンジン（必要に応じてダウンロード） |
| :--- | :--- | :--- |
| `arcade` | 軽量な AABB（めり込み防止）モード。物理計算なしでプラットフォーマー等に最適。 | なし |
| `box2d` | 本格的な 2D 剛体物理シミュレーション (Box2D v2)。 | `go get github.com/ByteArena/box2d` |
| `box2dgo` | Box2D v3 の Go 移植版。 | `go get github.com/oliverbestmann/box2d-go` |
| `phygo` | 軽量な 2D 物理エンジン。 | `go get github.com/ab-dek/Phygo-2D` |


```bash
# 例1: アーケード(AABB)モードを使う場合（追加の外部ライブラリ不要）
go get github.com/sofia-gros/ebiten/physics/adapters/arcade

# 例2: Box2D (v3) を使う場合（ベースエンジン本体も一緒にダウンロード）
go get github.com/oliverbestmann/box2d-go
go get github.com/sofia-gros/ebiten/physics/adapters/box2dgo
```


---

## 使い方

### クイックスタート

最小限の設定で四角形（`BoxShape`）の物理ボディを作成し、位置座標を取得して描画する基本的な記述方法です。

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
)

type Game struct {
	phys *physics.Manager
	body physics.Body
}

func (g *Game) Init() {
	g.phys = physics.NewManager()
	g.phys.SetWorld(arcade.NewWorld()) // Arcade エンジンをセット

	// 100, 100 の位置に 32x32px の四角形動的ボディを作成
	g.body = g.phys.CreateBody(physics.BodyOptions{
		Type:  physics.BodyTypeDynamic,
		X:     100,
		Y:     100,
		Shape: physics.BoxShape{Width: 32, Height: 32},
	})
}

func (g *Game) Update() error {
	// 移動速度の設定
	g.body.SetVelocity(100, 0)

	// 物理計算の更新 (60FPS)
	g.phys.Update(1.0 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 座標を取り出して手動描画
	x, y := g.body.Position()
	_ = x
	_ = y

	// 当たり判定枠線のデバッグ描画
	g.phys.DrawDebug(screen)
}
```

---

### 全機能の使い方

衝突グループの指定、接触イベントコールバック (`OnCollisionBegin`)、および画像の一括描画制御を行う全機能の使い方です。



```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
)

const (
	GroupPlayer physics.Group = 1
	GroupEnemy  physics.Group = 2
)

type Game struct {
	phys       *physics.Manager
	playerBody physics.Body
	playerImg  *ebiten.Image
}

func (g *Game) Init() {
	g.phys = physics.NewManager()
	g.phys.SetWorld(arcade.NewWorld())

	// 本格スタイル: BodyOptions で詳細設定
	g.playerBody = g.phys.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    100,
		Y:    100,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupEnemy {
				fmt.Println("敵と衝突しました！")
			}
		},
	})

	g.playerBody.SetGroup(GroupPlayer)

	// 一括描画用にスプライト画像とペア登録
	g.phys.AddRenderable(g.playerBody, g.playerImg)
}

func (g *Game) Update() error {
	g.phys.Update(1.0 / 60.0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ペア登録された全オブジェクトを一括描画
	g.phys.Draw(screen)

	// デバッグ用当たり判定の枠描画
	g.phys.DrawDebug(screen)
}
```

---

## 主要 API リファレンス

### `physics.Manager`
- **`SetWorld(world World)`**: 使用する物理エンジンの実体（アダプター）を注入します。（必須）
- **`SetGravity(gx, gy float64)`**: ワールド全体に重力を設定します。（例: `SetGravity(0, 100)`）
- **`CreateBody(options BodyOptions) Body`**: 空間内に物理ボディを生成します。
- **`RemoveBody(body Body)`**: 物体を空間および描画リストから削除します。

- `AddRenderable(body Body, img *ebiten.Image)`: スプライトと物理ボディを一括描画用に登録。
- `Draw(screen *ebiten.Image)`: 登録ペアを一括描画。
- `DrawDebug(screen *ebiten.Image)`: 当たり判定の枠線をデバッグ描画。

### `physics.Body`
- `Position() (x, y float64)` / `SetPosition(x, y float64)`: 座標取得・移動。
- `Velocity() (vx, vy float64)` / `SetVelocity(vx, vy float64)`: 速度取得・設定。
- `ApplyForce(fx, fy float64)`: 力を加える。
- `Group() physics.Group` / `SetGroup(group Group)`: 衝突グループ（レイヤー）の取得・設定。

---

## ライセンス

MIT License
