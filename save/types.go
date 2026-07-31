package save

// Directory はセーブデータを保存する基準パスを表す型です。
type Directory uint8

const (
	// DirCurrent はカレントディレクトリ (./) を基準とします (デフォルト)。
	// NewManager("save") の場合、 ./save/* にファイルが保存されます。
	DirCurrent Directory = iota

	// DirAppData は OS 標準のアプリケーションデータ領域を基準とします。
	// Windows: %APPDATA%/FolderName/*
	// macOS: ~/Library/Application Support/FolderName/*
	// Linux: ~/.config/FolderName/*
	DirAppData

	// DirCustom はユーザーが CustomDir で指定した任意のカスタムパスを基準とします。
	DirCustom
)

// Format はセーブデータの保存フォーマットを表す型です。
type Format string

const (
	FormatAuto   Format = "auto"   // 拡張子から自動判別
	FormatJSON   Format = "json"   // JSON テキスト形式 (プレーンテキスト/手動編集可能)
	FormatYAML   Format = "yaml"   // YAML テキスト形式
	FormatBinary Format = "binary" // GOB バイナリ形式
)
