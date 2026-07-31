package sound

// Type はサウンドの分類・チャンネル（SE, BGM, Voice, Env 等）を表す uint8 定数型です。
// Go 標準の iota により、ユーザーが自由に追加拡張できます。
type Type uint8

const (
	// TypeSE は効果音（デフォルト）を表します。
	TypeSE Type = iota

	// TypeBGM はバックグラウンドミュージックを表します。
	TypeBGM

	// TypeVoice はキャラクターボイス・セリフ音声を表します。
	TypeVoice

	// TypeEnv は環境音・アンビエントサウンドを表します。
	TypeEnv

	// TypeCustom はユーザーが独自のカスタムタイプを追加する際の基準値です。
	// 例: const TypeSystem sound.Type = sound.TypeCustom + iota
	TypeCustom
)
