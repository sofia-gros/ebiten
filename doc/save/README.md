# ebitensave

[English](./README_en.md)

ebitensave は Ebitengine および汎用 Go 開発向けのセーブデータ管理ライブラリです。
平文テキスト (JSON/YAML) やバイナリ形式の保存・復元、安全なディレクトリ管理、AES 暗号化、改ざん検出、削除済み欠番スロットの自動穴埋め保存、破損時のバックアップ自動復旧を一元管理できます。

---

## 特徴

- **安全な保存ディレクトリ管理 (`Dir`)**:
  - `DirCurrent` (デフォルト): カレント配下の専用フォルダ (例: `./save/*`) 内にファイルを安全に隔離保存。
  - `DirAppData`: OS標準のアプリケーションデータ領域 (%APPDATA%/FolderName/ や ~/Library/Application Support/FolderName/) に保存。
  - `DirCustom`: ユーザー独自のカスタムパスに保存。
- **フラットなファイル命名規則**:
  - スロットなし単一保存: `./save/data.json`
  - スロット指定保存: `./save/data_1.json`
- **欠番スロットの優先自動穴埋め保存 (`SaveNewSlot`)**:
  - スロット 2 が削除された場合、`SaveNewSlot` を呼び出すと空いた最小番号である **2** を自動検出して優先的に再利用保存。
- **暗号化と改ざん防止 (`EncryptKey`, `Checksum`)**:
  - メモ帳等でプレイヤーが手動編集できる平文 JSON/YAML 保存と、改ざんを防ぐ AES-GCM 暗号化バイナリ保存の両方に対応。
- **破損防止バックアップ (`Backup`)**:
  - 保存直前に旧ファイルを `.bak` へ退避。ファイル破壊を検知すると自動的に `.bak` から復旧。

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/save
```

---

## 使い方

### クイックスタート

`./save/data.json` としてプレーンテキスト保存・ロードする最もシンプルな方法です。


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/save"
)

type GameData struct {
	PlayerName string `json:"player_name"`
	Level      int    `json:"level"`
	Gold       int    `json:"gold"`
}

func main() {
	// カレントディレクトリ配下の "save" フォルダ (./save/*) 内に保存するマネージャーを作成
	m := save.NewManager("save", save.Option{
		Dir: save.DirCurrent, // ./
	})

	data := GameData{PlayerName: "Hero", Level: 10, Gold: 5000}

	// 1. スロットなし単一保存 (./save/data.json にプレーンテキストで保存)
	err := m.Save("data", "json", data, save.Option{
		Format: save.FormatJSON,
	})
	if err != nil {
		panic(err)
	}

	// 2. 単一データのロード
	var loaded GameData
	if err := m.Load("data", "json", &loaded); err == nil {
		fmt.Printf("Loaded: Level %d, Gold %d\n", loaded.Level, loaded.Gold)
	}
}
```

---

### 全機能の使い方

スロット 1 への指定上書き保存、`SaveNewSlot` による欠番スロットの自動穴埋め保存、AES 暗号化、およびデータ削除を行う全機能の使い方です。


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/save"
)

type GameData struct {
	PlayerName string `json:"player_name"`
	Level      int    `json:"level"`
}

func main() {
	m := save.NewManager("save", save.Option{
		Dir: save.DirCurrent,
	})

	data := GameData{PlayerName: "Hero", Level: 5}

	// 1. 特定スロットへの上書き保存 (./save/data_1.json)
	m.Slot(1).Save("data", "json", data)

	// 2. 新規スロットの自動穴埋め保存 (新規セーブ作成用)
	// (もしスロット 2 が過去に削除されていた場合、自動的に 2 番が選択されて ./save/data_2.json として保存される)
	slotNum, err := m.SaveNewSlot("data", "json", data, save.Option{
		Format: save.FormatJSON,
	})
	if err == nil {
		fmt.Printf("Saved into Slot %d!\n", slotNum)
	}

	// 3. 暗号化 + 改ざん検出付きバイナリ保存 (./save/secret.dat)
	m.Save("secret", "dat", data, save.Option{
		Format:     save.FormatBinary,
		EncryptKey: "secret_pass_123",
		Checksum:   true,
		Backup:     true, // クラッシュ防止バックアップON
	})

	// 4. スロット 2 のデータを削除 (./save/data_2.json を削除)
	m.Slot(2).Delete("data", "json")
}
```

---

## 主要 API リファレンス

### `save.Manager`
- `NewManager(folderName, opts...)`: 保存マネージャーを作成。
- `Save(name, ext, data, opts...)`: スロットサフィックスなしの単一ファイル保存 (例: `./save/data.json`)。
- `Load(name, ext, target, opts...)`: 単一ファイルをロードして構造体へバインド。
- `SaveNewSlot(name, ext, data, opts...)`: 欠番を含む最小のスロット番号を自動検出して保存し、割当スロット番号を返却。
- `Slot(slotNum)`: 指定スロット (`SlotManager`) の操作ハンドルを取得。
- `Exists(name, ext, opts...)` / `Delete(name, ext, opts...)`: ファイルの存在確認・削除。

### `save.SlotManager`
- `Slot(num).Save(name, ext, data, opts...)`: スロット番号を付与して保存 (例: `./save/data_1.json`)。
- `Slot(num).Load(name, ext, target, opts...)`: 特定スロットのファイルからロード。
- `Slot(num).Delete(name, ext, opts...)`: 特定スロットのファイルを削除。

---

## ライセンス

MIT License
