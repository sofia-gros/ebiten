package sound

// Option はサウンド再生時の各種オプション（タイプ、個別音量、ループ、パン、フェードインなど）を設定する構造体です。
type Option struct {
	// Type はサウンドの分類指定です。省略時は TypeSE になります。
	Type Type

	// Volume はこの再生の個別音量倍率 (0.0 〜 1.0) です。デフォルトは 1.0 です。
	Volume float64

	// Loop はループ再生を行うかのフラグです。
	Loop bool

	// Pan は左右の音量バランス（パンニング）です。
	// -1.0 (完全左) 〜 0.0 (中央) 〜 1.0 (完全右)。デフォルトは 0.0 です。
	Pan float64

	// Pitch は再生速度/高低（ピッチ倍率）です。デフォルトは 1.0 です。
	Pitch float64

	// FadeInDuration は再生開始時のフェードイン時間(秒)です。0 の場合は即座に規定音量で再生されます。
	FadeInDuration float64
}

// defaultOption は省略されたオプション項目にデフォルト値を適用します。
func defaultOption(opts []Option) Option {
	opt := Option{
		Type:     TypeSE,
		Volume:   1.0,
		Loop:     false,
		Pan:      0.0,
		Pitch:    1.0,
		FadeInDuration: 0.0,
	}

	if len(opts) > 0 {
		userOpt := opts[0]
		opt.Type = userOpt.Type
		if userOpt.Volume > 0 || userOpt.Volume < 0 {
			opt.Volume = userOpt.Volume
		}
		opt.Loop = userOpt.Loop
		opt.Pan = userOpt.Pan
		if userOpt.Pitch > 0 {
			opt.Pitch = userOpt.Pitch
		}
		opt.FadeInDuration = userOpt.FadeInDuration
	}

	return opt
}
