package asset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)


// assetType はアセットの分類を表す内部型です。
type assetType int

const (
	typeImage assetType = iota
	typeSprite
	typeAudio
	typeTilemap
	typeAnimation
	typeData
)

// assetEntry は登録されたアセットのメタデータとキャッシュを保持する構造体です。
type entry struct {
	key       string
	path      string
	assetType assetType
	option    Option
	loaded    bool
	val       any
	err       error
}

// Manager はアセットの登録、一括ロード、型安全な取得、キャッシュ解除を行うメインマネージャーです。
type Manager struct {
	mu           sync.RWMutex
	fileSystem   fs.FS
	entries      map[string]*entry
	keysInOrder  []string

	onStart    func(total int)
	onProgress func(progress float64, key string)
	onComplete func()
	onError    func(key string, err error)
}

// NewManager は新しいアセットマネージャーを作成します。
func NewManager() *Manager {
	return &Manager{
		entries:     make(map[string]*entry),
		keysInOrder: make([]string, 0),
	}
}

// SetFS は読み込みに使用するファイルシステム (`fs.FS` または `embed.FS`) を設定します。
func (m *Manager) SetFS(fileSys fs.FS) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fileSystem = fileSys
}

// OnStart は一括ロード開始時に呼ばれるコールバックを登録します。
func (m *Manager) OnStart(fn func(total int)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onStart = fn
}

// OnProgress は個別アセットのロード成功時に呼ばれるコールバックを登録します。
// progress は 0.0 〜 1.0 の進捗率を示します。
func (m *Manager) OnProgress(fn func(progress float64, key string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onProgress = fn
}

// OnComplete は全アセットのロード完了時に呼ばれるコールバックを登録します。
func (m *Manager) OnComplete(fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onComplete = fn
}

// OnError はアセットロードエラー発生時に呼ばれるコールバックを登録します。
func (m *Manager) OnError(fn func(key string, err error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onError = fn
}

// --- 登録メソッド (Add...) ---

func (m *Manager) addEntry(key, path string, aType assetType, opts ...Option) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var opt Option
	if len(opts) > 0 {
		opt = opts[0]
	}

	if _, exists := m.entries[key]; !exists {
		m.keysInOrder = append(m.keysInOrder, key)
	}

	m.entries[key] = &entry{
		key:       key,
		path:      path,
		assetType: aType,
		option:    opt,
		loaded:    false,
	}
}

// AddImage は単体画像アセット (`*ebiten.Image`) を登録します。
func (m *Manager) AddImage(key, path string, opts ...Option) {
	m.addEntry(key, path, typeImage, opts...)
}

// AddSprite はコマ割り画像アセット (`*SpriteSheet`) を登録します。
func (m *Manager) AddSprite(key, path string, opts ...Option) {
	m.addEntry(key, path, typeSprite, opts...)
}

// AddAudio は音声データアセットを登録します。
func (m *Manager) AddAudio(key, path string, opts ...Option) {
	m.addEntry(key, path, typeAudio, opts...)
}

// AddTilemap は Tilemap マップデータアセット (`*Tilemap`) を登録します。
func (m *Manager) AddTilemap(key, path string, opts ...Option) {
	m.addEntry(key, path, typeTilemap, opts...)
}

// AddAnimation はスプライトアニメーション定義データ (`*AnimationData`) を登録します。
func (m *Manager) AddAnimation(key, path string, opts ...Option) {
	m.addEntry(key, path, typeAnimation, opts...)
}

// AddData は JSON/YAML/TOML/CSV などの構造化データを登録します。
func (m *Manager) AddData(key, path string, opts ...Option) {
	m.addEntry(key, path, typeData, opts...)
}

// --- ロード実行 ---

// Load は登録されている未ロードのアセットを一括して順次ロード・パースします。
func (m *Manager) Load() {
	m.mu.Lock()
	total := len(m.entries)
	onStart := m.onStart
	onProgress := m.onProgress
	onComplete := m.onComplete
	onError := m.onError
	m.mu.Unlock()

	if onStart != nil {
		onStart(total)
	}

	loadedCount := 0

	for _, key := range m.keysInOrder {
		m.mu.RLock()
		e, exists := m.entries[key]
		m.mu.RUnlock()

		if !exists || e.loaded {
			loadedCount++
			continue
		}

		err := m.loadEntry(e)
		m.mu.Lock()
		if err != nil {
			e.err = err
			if onError != nil {
				onError(key, err)
			}
		} else {
			e.loaded = true
		}
		m.mu.Unlock()

		loadedCount++
		if onProgress != nil {
			progress := float64(loadedCount) / float64(total)
			onProgress(progress, key)
		}
	}

	if onComplete != nil {
		onComplete()
	}
}

// loadEntry は1つのエントリを読み込む内部処理です。
func (m *Manager) loadEntry(e *entry) error {
	dataBytes, err := m.readFile(e.path)
	if err != nil {
		return err
	}

	switch e.assetType {
	case typeImage:
		img, err := decodeImage(bytes.NewReader(dataBytes))
		if err != nil {
			return err
		}
		e.val = img

	case typeSprite:
		img, err := decodeImage(bytes.NewReader(dataBytes))
		if err != nil {
			return err
		}
		sheet := NewSpriteSheet(img, e.option.FrameWidth, e.option.FrameHeight, e.option.FrameMargin, e.option.FrameSpacing)
		e.val = sheet

	case typeAudio, typeData:
		fmt := e.option.Format
		if fmt == FormatAuto || fmt == "" {
			fmt = detectFormat(e.path)
		}
		parsed, err := decodeStructuredData(dataBytes, fmt)
		if err != nil {
			return err
		}
		e.val = parsed

	case typeTilemap:
		fmt := e.option.Format
		if fmt == FormatAuto || fmt == "" {
			fmt = detectFormat(e.path)
		}
		parsed, err := decodeStructuredData(dataBytes, fmt)
		if err != nil {
			return err
		}
		// JSON/YAML 等から Tilemap 構造体へ変換
		jsonBytes, err := json.Marshal(parsed)
		if err == nil {
			var tm Tilemap
			if err := json.Unmarshal(jsonBytes, &tm); err == nil {
				e.val = &tm
				return nil
			}
		}
		e.val = parsed

	case typeAnimation:
		fmt := e.option.Format
		if fmt == FormatAuto || fmt == "" {
			fmt = detectFormat(e.path)
		}
		parsed, err := decodeStructuredData(dataBytes, fmt)
		if err != nil {
			return err
		}
		// JSON/YAML 等から AnimationData 構造体へ変換
		jsonBytes, err := json.Marshal(parsed)
		if err == nil {
			var anim AnimationData
			if err := json.Unmarshal(jsonBytes, &anim); err == nil {
				e.val = &anim
				return nil
			}
		}
		e.val = parsed
	}

	return nil
}


func (m *Manager) readFile(path string) ([]byte, error) {
	if m.fileSystem != nil {
		file, err := m.fileSystem.Open(path)
		if err != nil {
			return nil, fmt.Errorf("fs.Open failed for %s: %w", path, err)
		}
		defer file.Close()
		return io.ReadAll(file)
	}

	return os.ReadFile(path)
}

// --- 取得メソッド (Image, Sprite, Tilemap, Animation, Data) ---

// Image は指定キーの単体画像 (`*ebiten.Image`) を取得します。
// 未ロードの場合は同期ロードを試行します。
func (m *Manager) Image(key string) (*ebiten.Image, error) {
	val, err := m.getOrLoad(key)
	if err != nil {
		return nil, err
	}
	if img, ok := val.(*ebiten.Image); ok {
		return img, nil
	}
	return nil, fmt.Errorf("asset '%s' is not an image", key)
}

// Sprite は指定キーのコマ割りスプライト (`*SpriteSheet`) を取得します。
func (m *Manager) Sprite(key string) (*SpriteSheet, error) {
	val, err := m.getOrLoad(key)
	if err != nil {
		return nil, err
	}
	if sheet, ok := val.(*SpriteSheet); ok {
		return sheet, nil
	}
	return nil, fmt.Errorf("asset '%s' is not a SpriteSheet", key)
}

// Tilemap は指定キーの Tilemap マップデータを取得します。
func (m *Manager) Tilemap(key string) (*Tilemap, error) {
	val, err := m.getOrLoad(key)
	if err != nil {
		return nil, err
	}
	if tm, ok := val.(*Tilemap); ok {
		return tm, nil
	}
	// バインド再試行
	tm, err := DataAs[Tilemap](m, key)
	if err == nil {
		return &tm, nil
	}
	return nil, fmt.Errorf("asset '%s' is not a Tilemap", key)
}

// Animation は指定キーのアニメーション定義データを取得します。
func (m *Manager) Animation(key string) (*AnimationData, error) {
	val, err := m.getOrLoad(key)
	if err != nil {
		return nil, err
	}
	if anim, ok := val.(*AnimationData); ok {
		return anim, nil
	}
	anim, err := DataAs[AnimationData](m, key)
	if err == nil {
		return &anim, nil
	}
	return nil, fmt.Errorf("asset '%s' is not AnimationData", key)
}

// Data は指定キーのパース済み構造化データ (any) を取得します。
func (m *Manager) Data(key string) (any, error) {
	return m.getOrLoad(key)
}

func (m *Manager) getOrLoad(key string) (any, error) {
	m.mu.RLock()
	e, exists := m.entries[key]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("asset '%s' not found", key)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if !e.loaded {
		if err := m.loadEntry(e); err != nil {
			e.err = err
			return nil, err
		}
		e.loaded = true
	}

	return e.val, e.err
}

// --- クリーンアップ・メモリ解放 ---

// Unload は指定したキーのアセットをキャッシュおよびマネージャーから削除しメモリを解放します。
func (m *Manager) Unload(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.entries, key)
	for i, k := range m.keysInOrder {
		if k == key {
			m.keysInOrder = append(m.keysInOrder[:i], m.keysInOrder[i+1:]...)
			break
		}
	}
}

// Clear はすべての読み込み済みアセットおよび未処理登録を一括クリアして初期化します。
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = make(map[string]*entry)
	m.keysInOrder = m.keysInOrder[:0]
}
