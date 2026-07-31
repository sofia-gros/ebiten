package save

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"strings"
)

// Storage は保存先パスの解決と実際のファイルIOを行う抽象化インターフェースです。
type Storage struct {
	baseDir string
}

// newStorage は指定された基準ディレクトリ `Dir` とフォルダ名 `folderName` から物理ディレクトリパスを解決します。
func newStorage(folderName string, opt Option) (*Storage, error) {
	var root string

	switch opt.Dir {
	case DirCurrent:
		cwd, err := os.Getwd()
		if err != nil {
			root = "."
		} else {
			root = cwd
		}

	case DirAppData:
		root = getAppDataDir()

	case DirCustom:
		if opt.CustomDir != "" {
			root = opt.CustomDir
		} else {
			root = "."
		}

	default:
		root = "."
	}

	targetDir := filepath.Join(root, folderName)

	// 保存先ディレクトリが存在しない場合は自動作成
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create save directory %s: %w", targetDir, err)
	}

	return &Storage{baseDir: targetDir}, nil
}

// BaseDir は解決された物理ディレクトリパスを返します。
func (s *Storage) BaseDir() string {
	return s.baseDir
}

// BuildFilePath はフラットなファイルパスを組み立てます。
// slot == 0 の場合: {baseDir}/{name}.{ext}
// slot > 0 の場合:  {baseDir}/{name}_{slot}.{ext}
func (s *Storage) BuildFilePath(name, ext string, slot int) string {
	ext = strings.TrimPrefix(ext, ".")
	if slot <= 0 {
		return filepath.Join(s.baseDir, fmt.Sprintf("%s.%s", name, ext))
	}
	return filepath.Join(s.baseDir, fmt.Sprintf("%s_%d.%s", name, slot, ext))
}

// WriteFile は指定パスにアトミックにデータを保存します。Backup が有効な場合は旧ファイルを .bak へ退避します。
func (s *Storage) WriteFile(filePath string, data []byte, backup bool) error {
	if backup {
		if _, err := os.Stat(filePath); err == nil {
			bakPath := filePath + ".bak"
			_ = os.Remove(bakPath)
			_ = os.Rename(filePath, bakPath)
		}
	}

	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp file %s: %w", tmpPath, err)
	}

	if err := os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("failed to replace target file %s: %w", filePath, err)
	}

	return nil
}

// ReadFile は指定パスからデータを読み込みます。壊れている場合は .bak からの自動復旧を試みます。
func (s *Storage) ReadFile(filePath string, backup bool) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err == nil {
		return data, nil
	}

	// メインファイルのロードに失敗した場合、バックアップからの復旧を試行
	if backup {
		bakPath := filePath + ".bak"
		bakData, bakErr := os.ReadFile(bakPath)
		if bakErr == nil {
			// バックアップをメインファイルに復元
			_ = os.WriteFile(filePath, bakData, 0644)
			return bakData, nil
		}
	}

	return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
}

// Exists は指定パスのファイルが存在するかを判定します。
func (s *Storage) Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// Remove は指定パスのファイルを削除します。
func (s *Storage) Remove(filePath string) error {
	_ = os.Remove(filePath + ".bak")
	return os.Remove(filePath)
}

// ScanUnusedSlot は既存ファイル名サフィックス ({name}_{slot}.{ext}) を検索し、
// 欠番を含む最小の未使用スロット番号 (1, 2, 3...) を検出して返します。
func (s *Storage) ScanUnusedSlot(name, ext string) int {
	ext = strings.TrimPrefix(ext, ".")
	usedSlots := make(map[int]bool)

	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return 1
	}

	prefix := name + "_"
	suffix := "." + ext

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if strings.HasPrefix(filename, prefix) && strings.HasSuffix(filename, suffix) {
			numStr := filename[len(prefix) : len(filename)-len(suffix)]
			if num, err := strconv.Atoi(numStr); err == nil && num > 0 {
				usedSlots[num] = true
			}
		}
	}

	// 最小の未使用番号 (欠番補填) を検索
	for slot := 1; ; slot++ {
		if !usedSlots[slot] {
			return slot
		}
	}
}

// getAppDataDir は OS ごとの標準 AppData ディレクトリパスを取得します。
func getAppDataDir() string {
	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return appData
		}
	case "darwin":
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, "Library", "Application Support")
		}
	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return xdg
		}
		if home := os.Getenv("HOME"); home != "" {
			return filepath.Join(home, ".config")
		}
	}
	return "."
}
