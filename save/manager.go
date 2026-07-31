package save

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)


// SlotManager は特定のスロット番号 (1, 2, 3...) に固定された操作用ハンドルです。
type SlotManager struct {
	manager *Manager
	slotNum int
}

// Save はこのスロット番号のサフィックスを付与してデータを保存します (例: ./save/data_1.json)。
func (sm *SlotManager) Save(name, ext string, data any, opts ...Option) error {
	if sm == nil || sm.manager == nil {
		return fmt.Errorf("slot manager is nil")
	}
	return sm.manager.saveSlot(sm.slotNum, name, ext, data, opts...)
}

// Load はこのスロット番号のファイルからデータを読み込み構造体へ反映します。
func (sm *SlotManager) Load(name, ext string, target any, opts ...Option) error {
	if sm == nil || sm.manager == nil {
		return fmt.Errorf("slot manager is nil")
	}
	return sm.manager.loadSlot(sm.slotNum, name, ext, target, opts...)
}

// Exists はこのスロットのファイルが存在するか判定します。
func (sm *SlotManager) Exists(name, ext string, opts ...Option) bool {
	if sm == nil || sm.manager == nil {
		return false
	}
	return sm.manager.existsSlot(sm.slotNum, name, ext, opts...)
}

// Delete はこのスロットのファイルを削除します。
func (sm *SlotManager) Delete(name, ext string, opts ...Option) error {
	if sm == nil || sm.manager == nil {
		return fmt.Errorf("slot manager is nil")
	}
	return sm.manager.deleteSlot(sm.slotNum, name, ext, opts...)
}

// Manager はセーブデータの保存・復元・スロット・暗号化を一元管理するメインマネージャーです。
type Manager struct {
	mu         sync.RWMutex
	folderName string
	defaultOpt Option
	storage    *Storage
}

// NewManager は保存フォルダ名 (例: "save", "my_game_save") とオプションを指定してマネージャーを作成します。
func NewManager(folderName string, opts ...Option) *Manager {
	opt := defaultOption(opts)
	storage, err := newStorage(folderName, opt)
	if err != nil {
		// ディレクトリ作成失敗時フォールバック
		storage, _ = newStorage(folderName, Option{Dir: DirCurrent})
	}

	return &Manager{
		folderName: folderName,
		defaultOpt: opt,
		storage:    storage,
	}
}

// Slot は指定したスロット番号 (1, 2, 3...) の操作ハンドルを取得します。
func (m *Manager) Slot(num int) *SlotManager {
	return &SlotManager{
		manager: m,
		slotNum: num,
	}
}

// Save はスロットサフィックスなしの単一ファイルとしてデータを保存します (例: ./save/data.json)。
func (m *Manager) Save(name, ext string, data any, opts ...Option) error {
	return m.saveSlot(0, name, ext, data, opts...)
}

// SaveNewSlot は欠番を含む未使用の最小スロット番号 (1, 2, 3...) を自動検出して保存し、割当スロット番号を返します。
func (m *Manager) SaveNewSlot(name, ext string, data any, opts ...Option) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	slotNum := m.storage.ScanUnusedSlot(name, ext)
	err := m.saveSlotUnlocked(slotNum, name, ext, data, opts...)
	if err != nil {
		return 0, err
	}
	return slotNum, nil
}

// Load はスロットサフィックスなしのファイルからデータを読み込み構造体へバインドします。
func (m *Manager) Load(name, ext string, target any, opts ...Option) error {
	return m.loadSlot(0, name, ext, target, opts...)
}

// Exists はスロットサフィックスなしのファイルが存在するか判定します。
func (m *Manager) Exists(name, ext string, opts ...Option) bool {
	return m.existsSlot(0, name, ext, opts...)
}

// Delete はスロットサフィックスなしのファイルを削除します。
func (m *Manager) Delete(name, ext string, opts ...Option) error {
	return m.deleteSlot(0, name, ext, opts...)
}

// --- 内部処理 ---

func (m *Manager) saveSlot(slot int, name, ext string, data any, opts ...Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saveSlotUnlocked(slot, name, ext, data, opts...)
}

func (m *Manager) saveSlotUnlocked(slot int, name, ext string, data any, opts ...Option) error {
	opt := m.mergeOption(opts)
	filePath := m.storage.BuildFilePath(name, ext, slot)

	// 1. データシリアライズ
	serialized, err := m.serialize(data, ext, opt.Format)
	if err != nil {
		return fmt.Errorf("failed to serialize save data: %w", err)
	}

	// 2. チェックサム付加 (オプション)
	if opt.Checksum {
		serialized = attachChecksum(serialized)
	}

	// 3. AES 暗号化 (オプション)
	if opt.EncryptKey != "" {
		encrypted, err := encryptAES(serialized, opt.EncryptKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt save data: %w", err)
		}
		serialized = encrypted
	}

	// 4. ディレクトリ書き込み (バックアップ退避付き)
	return m.storage.WriteFile(filePath, serialized, opt.Backup)
}

func (m *Manager) loadSlot(slot int, name, ext string, target any, opts ...Option) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	opt := m.mergeOption(opts)
	filePath := m.storage.BuildFilePath(name, ext, slot)

	// 1. メインファイルからの読み込み＆デシリアライズ試行
	err := m.loadFromFile(filePath, ext, opt, target)
	if err == nil {
		return nil
	}

	// 2. メインファイルが破損/失敗しており、かつ Backup が有効な場合、.bak からの復旧試行
	if opt.Backup {
		bakPath := filePath + ".bak"
		if m.storage.Exists(bakPath) {
			bakErr := m.loadFromFile(bakPath, ext, opt, target)
			if bakErr == nil {
				// バックアップからの読み込みに成功した場合、壊れたメインファイルをバックアップデータで上書き修復
				if bakData, readBakErr := m.storage.ReadFile(bakPath, false); readBakErr == nil {
					_ = m.storage.WriteFile(filePath, bakData, false)
				}
				return nil
			}
		}
	}

	return fmt.Errorf("failed to load save data from %s: %w", filePath, err)
}

func (m *Manager) loadFromFile(filePath, ext string, opt Option, target any) error {
	rawData, err := m.storage.ReadFile(filePath, false)
	if err != nil {
		return err
	}

	if opt.EncryptKey != "" {
		decrypted, err := decryptAES(rawData, opt.EncryptKey)
		if err != nil {
			return err
		}
		rawData = decrypted
	}

	if opt.Checksum {
		verified, err := verifyAndStripChecksum(rawData)
		if err != nil {
			return err
		}
		rawData = verified
	}

	return m.deserialize(rawData, ext, opt.Format, target)
}


func (m *Manager) existsSlot(slot int, name, ext string, opts ...Option) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	filePath := m.storage.BuildFilePath(name, ext, slot)
	return m.storage.Exists(filePath)
}

func (m *Manager) deleteSlot(slot int, name, ext string, opts ...Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := m.storage.BuildFilePath(name, ext, slot)
	return m.storage.Remove(filePath)
}

func (m *Manager) mergeOption(opts []Option) Option {
	if len(opts) > 0 {
		return opts[0]
	}
	return m.defaultOpt
}

func (m *Manager) serialize(data any, ext string, fmtType Format) ([]byte, error) {
	if fmtType == FormatAuto || fmtType == "" {
		fmtType = detectFormatFromExt(ext)
	}

	switch fmtType {
	case FormatJSON:
		return json.MarshalIndent(data, "", "  ")

	case FormatYAML:
		return yaml.Marshal(data)

	case FormatBinary:
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(data); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil

	default:
		return json.MarshalIndent(data, "", "  ")
	}
}

func (m *Manager) deserialize(rawData []byte, ext string, fmtType Format, target any) error {
	if fmtType == FormatAuto || fmtType == "" {
		fmtType = detectFormatFromExt(ext)
	}

	switch fmtType {
	case FormatJSON:
		return json.Unmarshal(rawData, target)

	case FormatYAML:
		return yaml.Unmarshal(rawData, target)

	case FormatBinary:
		buf := bytes.NewBuffer(rawData)
		dec := gob.NewDecoder(buf)
		return dec.Decode(target)

	default:
		return json.Unmarshal(rawData, target)
	}
}

func detectFormatFromExt(ext string) Format {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	switch ext {
	case "json":
		return FormatJSON
	case "yaml", "yml":
		return FormatYAML
	case "dat", "bin", "gob":
		return FormatBinary
	default:
		return FormatJSON
	}
}
