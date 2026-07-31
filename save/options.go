package save

// Option はセーブ/ロード処理のオプション設定構造体です。
type Option struct {
	// Dir は保存先フォルダを生成・配置する基準ディレクトリです。デフォルトは DirCurrent (./) です。
	Dir Directory

	// CustomDir は DirCustom を選択した際の基準カスタムディレクトリパスです。
	CustomDir string

	// Format は保存形式 (JSON, YAML, Binary) を明示的に指定します。
	Format Format

	// EncryptKey は暗号化に使用するキー文字列です。空文字でない場合 AES 暗号化が行われます。
	EncryptKey string

	// Checksum はデータ末尾に SHA256 チェックサムを付加し、改ざん検出を行うかのフラグです。
	Checksum bool

	// Backup は書き込み直前に旧ファイルを .bak として退避し、ロード破損時に自動復旧するかのフラグです。
	Backup bool
}

// defaultOption は省略されたオプション項目にデフォルト値を補填します。
func defaultOption(opts []Option) Option {
	opt := Option{
		Dir:        DirCurrent,
		Format:     FormatAuto,
		EncryptKey: "",
		Checksum:   false,
		Backup:     true, // デフォルトで破損防止バックアップを有効化
	}

	if len(opts) > 0 {
		userOpt := opts[0]
		opt.Dir = userOpt.Dir
		if userOpt.CustomDir != "" {
			opt.CustomDir = userOpt.CustomDir
		}
		if userOpt.Format != "" {
			opt.Format = userOpt.Format
		}
		opt.EncryptKey = userOpt.EncryptKey
		opt.Checksum = userOpt.Checksum
		opt.Backup = userOpt.Backup
	}

	return opt
}
