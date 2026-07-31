package save_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sofia-gros/ebiten/save"
)

type TestGameData struct {
	PlayerName string `json:"player_name" yaml:"player_name"`
	Level      int    `json:"level"       yaml:"level"`
	Gold       int    `json:"gold"        yaml:"gold"`
}

func TestBasicSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	m := save.NewManager("save_test", save.Option{
		Dir:       save.DirCustom,
		CustomDir: tmpDir,
	})

	data := TestGameData{PlayerName: "Hero", Level: 5, Gold: 1000}

	// 1. スロットなし単一保存 (./save_test/data.json)
	err := m.Save("data", "json", data)
	if err != nil {
		t.Fatalf("failed to Save: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "save_test", "data.json")
	if !m.Exists("data", "json") {
		t.Errorf("expected file %s to exist", expectedPath)
	}

	var loaded TestGameData
	err = m.Load("data", "json", &loaded)
	if err != nil {
		t.Fatalf("failed to Load: %v", err)
	}

	if loaded.PlayerName != "Hero" || loaded.Level != 5 || loaded.Gold != 1000 {
		t.Errorf("loaded data mismatch: %+v", loaded)
	}
}

func TestSlotSaveAndAutoReuse(t *testing.T) {
	tmpDir := t.TempDir()

	m := save.NewManager("save_slots", save.Option{
		Dir:       save.DirCustom,
		CustomDir: tmpDir,
	})

	data := TestGameData{PlayerName: "Hero", Level: 1, Gold: 100}

	// スロット 1, 2, 3 を作成
	m.Slot(1).Save("data", "json", data)
	m.Slot(2).Save("data", "json", data)
	m.Slot(3).Save("data", "json", data)

	// スロット 2 を削除 (欠番発生: ./save_slots/data_2.json)
	err := m.Slot(2).Delete("data", "json")
	if err != nil {
		t.Fatalf("failed to delete slot 2: %v", err)
	}

	if m.Slot(2).Exists("data", "json") {
		t.Errorf("slot 2 should be deleted")
	}

	// SaveNewSlot 実行 ➔ 欠番になったスロット 2 が自動選択されて穴埋め保存されるはず！
	newSlotNum, err := m.SaveNewSlot("data", "json", data)
	if err != nil {
		t.Fatalf("SaveNewSlot failed: %v", err)
	}

	if newSlotNum != 2 {
		t.Errorf("expected SaveNewSlot to pick reused slot 2, got %d", newSlotNum)
	}

	// 再度 SaveNewSlot 実行 ➔ 次の空きであるスロット 4 が割り当てられるはず！
	nextSlotNum, err := m.SaveNewSlot("data", "json", data)
	if err != nil {
		t.Fatalf("SaveNewSlot 2 failed: %v", err)
	}

	if nextSlotNum != 4 {
		t.Errorf("expected SaveNewSlot to pick next slot 4, got %d", nextSlotNum)
	}
}

func TestEncryptionAndChecksum(t *testing.T) {
	tmpDir := t.TempDir()

	m := save.NewManager("save_crypto", save.Option{
		Dir:       save.DirCustom,
		CustomDir: tmpDir,
	})

	data := TestGameData{PlayerName: "SecretAgent", Level: 99, Gold: 999999}
	secretKey := "my_pass_12345"

	// 暗号化 + チェックサム保存
	err := m.Save("secret", "dat", data, save.Option{
		Format:     save.FormatBinary,
		EncryptKey: secretKey,
		Checksum:   true,
	})
	if err != nil {
		t.Fatalf("failed to save encrypted data: %v", err)
	}

	// 正しいキーで復号
	var loaded TestGameData
	err = m.Load("secret", "dat", &loaded, save.Option{
		Format:     save.FormatBinary,
		EncryptKey: secretKey,
		Checksum:   true,
	})
	if err != nil {
		t.Fatalf("failed to load encrypted data: %v", err)
	}

	if loaded.PlayerName != "SecretAgent" || loaded.Level != 99 {
		t.Errorf("loaded encrypted data mismatch: %+v", loaded)
	}

	// 誤ったキーで復号するとエラーになるかのテスト
	var loadedWrong TestGameData
	errWrong := m.Load("secret", "dat", &loadedWrong, save.Option{
		Format:     save.FormatBinary,
		EncryptKey: "wrong_password",
		Checksum:   true,
	})
	if errWrong == nil {
		t.Errorf("expected error when decrypting with wrong key")
	}
}

func TestBackupRecovery(t *testing.T) {
	tmpDir := t.TempDir()

	m := save.NewManager("save_backup", save.Option{
		Dir:       save.DirCustom,
		CustomDir: tmpDir,
	})

	data1 := TestGameData{PlayerName: "ValidData", Level: 1, Gold: 100}
	m.Save("data", "json", data1, save.Option{Backup: true})

	// 2回目の保存 (data.json.bak に valid 状態が保存される)
	data2 := TestGameData{PlayerName: "ValidData2", Level: 2, Gold: 200}
	m.Save("data", "json", data2, save.Option{Backup: true})

	filePath := filepath.Join(tmpDir, "save_backup", "data.json")

	// メインの data.json ファイルの内容を破壊 (壊れたJSON化)
	err := os.WriteFile(filePath, []byte("CORRUPTED_JSON_DATA!!!"), 0644)
	if err != nil {
		t.Fatalf("failed to corrupt file: %v", err)
	}

	// Load 実行 ➔ 自動的に .bak から復旧されてロード成功するはず！
	var loaded TestGameData
	err = m.Load("data", "json", &loaded, save.Option{Backup: true})
	if err != nil {
		t.Fatalf("failed to recover from backup: %v", err)
	}

	if loaded.PlayerName != "ValidData" && loaded.PlayerName != "ValidData2" {
		t.Errorf("recovered data mismatch: %+v", loaded)
	}
}
