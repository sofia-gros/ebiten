# ebitenphysics

ebitenphysics は Ebitengine 向けの物理演算ラッパーライブラリです。
RPGやプラットフォーマー向けの軽量な Arcade（めり込み防止AABB）モードと、パズルゲーム等向けの本格的な Box2D（円や剛体物理）モードを、まったく同じAPIで切り替えて利用できます。

## 特徴

- **ピュアな設計**: ebitenphysics 本体には特定の物理エンジンが含まれていません。使いたいエンジン（アダプター）だけをインポートするため、ゲームに余計なコードが混入しません。
- **DIベース**: `arcade` や `box2d` などのエンジンを Manager に渡して動かす設計です。
- **描画は自由**: 自分で `x, y` を取得して自由に描画することもできますし、画像と一緒に登録して一括描画させることもできます。既存のシーン管理やUIと競合しません。

## アップデート情報

Version 1.0.0 2026/07/27
・初期リリース

## インストール

まずは本体をインストールします。

```bash
go get github.com/sofia-gros/ebiten/physics
```

次に、使いたい物理エンジンのアダプターをインストールします。

> **⚠️ 注意: 外部エンジンのインストールについて**
> 本ライブラリ（ebitenphysics）自体は、外部の物理エンジンのコードを一切含んでいません。
> `box2d` などの外部エンジンを利用する場合は、上記のアダプターに加え、**ベースとなる物理エンジン本体もユーザー自身でダウンロード（`go get`）しておく**必要があります。

| アダプター | 概要                                       | 依存エンジン（各自でダウンロードしてください） |
| :--------- | :----------------------------------------- | :--------------------------------------------- |
| `arcade`   | 自作の軽量なAABB（めり込み防止）エンジン。 | なし                                           |
| `box2d`    | 本格的な2D剛体シミュレーション。           | `go get github.com/ByteArena/box2d`            |
| `physix`   | シンプルな2D物理エンジン。                 | `go get github.com/rudransh61/Physix-go`       |
| `phygo`    | 軽量な2D物理エンジン。                     | `go get github.com/ab-dek/Phygo-2D`            |

```bash
# 例: Box2Dエンジンを使う場合
go get github.com/ByteArena/box2d
go get github.com/sofia-gros/ebiten/physics/adapters/box2d
```

## 基本的な使い方

```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/physics"

	// 使いたい物理エンジンをインポート
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
)

// 衝突グループの定義（レイヤー番号みたいなもの）
const (
	GroupPlayer physics.Group = 1
	GroupEnemy  physics.Group = 2
)

type GameScene struct {
	physManager *physics.Manager
	playerBody  physics.Body
	playerImage *ebiten.Image
}

func (s *GameScene) Init() {
	s.physManager = physics.NewManager()
	s.physManager.SetWorld(arcade.NewWorld())

	// Body（物理演算用の実体）を作る
	s.playerBody = s.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    100,
		Y:    100,
		// 基本的な四角形の当たり判定
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		// 衝突コールバック
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupEnemy {
				fmt.Println("敵にヒット！")
			}
		},
	})

	// グループ番号を割り当てる
	s.playerBody.SetGroup(GroupPlayer)

	// 一括描画用に登録
	s.physManager.AddRenderable(s.playerBody, s.playerImage)
}

func (s *GameScene) Update() error {
	s.physManager.Update(1.0 / 60.0)
	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	// 自分で描画する場合
	// x, y := s.playerBody.Position()
	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(x, y)
	// screen.DrawImage(s.playerImage, op)

	// マネージャーに任せる場合
	s.physManager.Draw(screen)

	// デバッグ用当たり判定の枠線を描画する
	s.physManager.DrawDebug(screen)
}
```

## API・機能詳細

本ライブラリが提供するインターフェース（すべての物理エンジン共通）の主な機能一覧です。
※ 摩擦や反発など、どのパラメータが実際に有効になるかは、注入した**物理エンジンのアダプター（のREADME）**に従います。

### 1. `physics.Manager` のメソッド

物理空間全体と描画を統括するオブジェクトです。

- **`SetWorld(world World)`**: 使用する物理エンジンの実体（アダプター）を注入します。（必須）
- **`CreateBody(options BodyOptions) Body`**: 空間内に物体を生成します。
- **`RemoveBody(body Body)`**: 物体を空間および描画リストから完全に削除します。
- **`Update(dt float64)`**: 物理シミュレーションを `dt` 秒だけ進めます。
- **`AddRenderable(body Body, img *ebiten.Image)`**: 一括描画対象として、画像と物理ボディを紐付けて登録します。
- **`Draw(screen *ebiten.Image)`**: `AddRenderable` で登録されたペアを一括で描画します。
- **`DrawDebug(screen *ebiten.Image)`**: 各ボディの当たり判定（矩形や円）の枠線を画面上に可視化します。

### 2. `physics.Body` のメソッド

`CreateBody` で返される物理オブジェクトです。いつでも状態の取得・変更が可能です。

#### 状態の取得（Get）

- **`Position() (x, y float64)`**: 現在のX, Y座標を取得します。
- **`Rotation() float64`**: 現在の回転角度（ラジアン）を取得します。
- **`Velocity() (vx, vy float64)`**: 現在のX, Y方向の速度を取得します。
- **`AngularVelocity() float64`**: 現在の回転速度を取得します。
- **`Group() physics.Group`**: 衝突レイヤーや種類の判定に使うためのグループ番号（数値）を取得します。
- **`Data() interface{}`**: ボディに紐付けられた任意のユーザーデータ（独自の構造体ポインタなど）を取得します。

#### 状態の変更（Set）

- **`SetPosition(x, y float64)`**: 強制的に座標を移動させます（ワープ）。
- **`SetRotation(angle float64)`**: 強制的に角度を変更します。
- **`SetVelocity(vx, vy float64)`**: 指定した速度（ピクセル/秒など）を直接セットします。ジャンプやダッシュなどの実装に便利です。
- **`SetAngularVelocity(omega float64)`**: 回転速度をセットします。
- **`ApplyForce(fx, fy float64)`**: 現在の速度に対して徐々に力を加えます（風や重力などの持続的な力向き）。
- **`SetGroup(group physics.Group)`**: オブジェクトの種類（Player=1, Enemy=2 など）を数値で割り当てます。文字列比較より高速です。
- **`SetData(data interface{})`**: 衝突判定時などで詳細な処理をするために、独自のデータ（ステータス構造体など）を持たせることができます。

### 3. `physics.BodyOptions` (生成時の設定)

- **`Type`**: `physics.BodyTypeDynamic`（動く物体）か `physics.BodyTypeStatic`（動かない壁など）を指定。
- **`Shapes []ShapeDef`**: オブジェクトを構成する形状。複数渡すことで、腕や足が独立した「複合形状」を1つのボディとして作れます。
- **`Friction`**: 摩擦係数。（滑りにくさ）
- **`Restitution`**: 反発係数。（バウンス・跳ね返りやすさ）
- **`Density`**: 密度。質量計算に使われます。
- **`IsSensor`**: `true` にすると、衝突判定（コールバック）は発生しますが、物理的に押し返されることはなくなります（アイテム取得エリアや毒沼などに最適）。

#### 衝突コールバック

- **`OnCollisionBegin func(other Body)`**: 他の物体と「接触した瞬間」に1回だけ呼ばれます。
- **`OnOverlap func(other Body)`**: 他の物体と「重なっている（または接触し続けている）間」毎フレーム呼ばれ続けます。
- **`OnCollisionEnd func(other Body)`**: 他の物体から「離れた瞬間」に1回だけ呼ばれます。

### 4. 形状（Shape）の定義

`ShapeDef` を使うことで、ボディの中心座標から相対的にズラした位置に形状を配置できます。

```go
// 中心の円の横に、四角形をくっつけた形の例
Shapes: []physics.ShapeDef{
	{Shape: physics.CircleShape{Radius: 16}, OffsetX: 0, OffsetY: 0},
	{Shape: physics.BoxShape{Width: 32, Height: 8}, OffsetX: 32, OffsetY: 0},
}
```

## ライセンス

MIT License
