package asset

// DataFormat はアセットデータのフォーマット（拡張子判別等）を表します。
type DataFormat string

const (
	FormatAuto DataFormat = "auto"
	FormatJSON DataFormat = "json"
	FormatYAML DataFormat = "yaml"
	FormatTOML DataFormat = "toml"
	FormatCSV  DataFormat = "csv"
	FormatText DataFormat = "text"
	FormatRaw  DataFormat = "raw"
)

// Option はアセット登録時の詳細オプションを設定する構造体です。
type Option struct {
	// FrameWidth は SpriteSheet (コマ割り画像) の1コマの横幅(px)です。
	FrameWidth int

	// FrameHeight は SpriteSheet (コマ割り画像) の1コマの縦幅(px)です。
	FrameHeight int

	// FrameMargin は SpriteSheet の外周マージン(px)です。
	FrameMargin int

	// FrameSpacing は SpriteSheet コマ間の間隔(px)です。
	FrameSpacing int

	// Format はデータファイルのフォーマットを明示的に指定します。
	Format DataFormat
}
