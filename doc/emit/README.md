# ebitenemit

[English](./README_en.md)

ebitenemit は Go のジェネリクスを活用した、型安全かつ軽量なイベントディスパッチャー／エミッターライブラリです。
他のライブラリに一切依存せず単体で動作し、コンポーネント間の疎結合なイベント通知やアニメーション・ゲームロジックのイベント駆動をシンプルに記述できます。

---

## 特徴

- **型安全なジェネリクス (`Emitter[T]`)**:
  - イベントデータ `T` の型をコンパイル時に検証。文字列キーによるキャストバグを防止。
- **直感的なハンドラー登録 (`On` / `Once` / `Off`)**:
  - イベントの継続リスニング、1回限りの発火、ハンドラー解除を簡単に管理。
- **2通りのイベント発行モード (`Emit` vs `Queue` / `Flush`)**:
  - `Emit`: 呼出時に即座にハンドラーを実行（同期）。
  - `Queue` + `Flush`: イベントをキューに蓄積し、ゲームループの指定タイミングでまとめて一括処理（遅延キュー）。

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/emit
```

---

## 使い方

### クイックスタート

型安全なイベントを発行し、リスナー側で即座に受け取る基本的なコード例です。


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/emit"
)

// イベント構造体の定義
type PlayerDamageEvent struct {
	PlayerID int
	Damage   int
}

func main() {
	// イベントエミッターの作成
	e := emit.New[PlayerDamageEvent]()

	// リスナーの登録 (On)
	unsubscribe := e.On(func(ev PlayerDamageEvent) {
		fmt.Printf("Player %d took %d damage!\n", ev.PlayerID, ev.Damage)
	})

	// イベントの即時発行 (Emit)
	e.Emit(PlayerDamageEvent{PlayerID: 1, Damage: 25})

	// 登録解除
	unsubscribe()
}
```

---

### 全機能の使い方

ゲームループの指定フレームタイミングでまとめてイベントを一括処理する、遅延キュー（`Queue` + `Flush`）を行う全機能の使い方です。


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/emit"
)

type ItemPickedEvent struct {
	ItemName string
	Score    int
}

type Game struct {
	itemEmitter *emit.Emitter[ItemPickedEvent]
}

func (g *Game) Init() {
	g.itemEmitter = emit.New[ItemPickedEvent]()

	// 1回だけ発火するリスナー (Once)
	g.itemEmitter.Once(func(ev ItemPickedEvent) {
		fmt.Println("First item picked up!", ev.ItemName)
	})

	// 通常の継続リスナー
	g.itemEmitter.On(func(ev ItemPickedEvent) {
		fmt.Printf("Picked: %s (+%d pts)\n", ev.ItemName, ev.Score)
	})
}

func (g *Game) OnItemCollide(name string, pts int) {
	// 即時処理せずキューに溜める (Queue)
	g.itemEmitter.Queue(ItemPickedEvent{ItemName: name, Score: pts})
}

func (g *Game) Update() error {
	// ゲームループのこのタイミングでキューを一括処理・フラッシュ！
	g.itemEmitter.Flush()
	return nil
}
```

---

## 主要 API リファレンス

### `emit.Emitter[T]`
- `New[T]()`: 型 `T` のエミッターを作成。
- `On(handler func(T)) func()`: イベントリスナーを登録。戻り値の関数で登録解除可能。
- `Once(handler func(T))`: 1回発火したら自動解除されるリスナーを登録。
- `Off(handlerID)` / `RemoveAll()`: 特定ハンドラーまたは全リスナーの解除。
- `Emit(data T)`: イベントを即座に同期発行。
- `Queue(data T)`: イベントを一時キューに蓄積。
- `Flush()`: キューに溜まった全イベントを一括処理。
- `Reset()`: キューと全リスナーの完全クリア。

---

## ライセンス

MIT License
