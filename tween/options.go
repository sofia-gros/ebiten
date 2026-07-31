package tween

// Option は Tween の初期化・設定パラメータを保持する構造体です。
type Option struct {
	// Start は開始値です (デフォルト: 0.0)。
	Start float64

	// End は終了目標値です (デフォルト: 1.0)。
	End float64

	// Duration は補間アニメーションの所要時間 (秒) です (デフォルト: 1.0)。
	Duration float64

	// Ease は適用するイージング関数です (デフォルト: EaseLinear)。
	Ease EaseFunc

	// Delay はアニメーション開始前の遅延時間 (秒) です。
	Delay float64

	// Loop はリピート回数です (0: 1回のみ再生, n > 0: n回リピート, -1: 無限リピート)。
	Loop int

	// Yoyo はリピート時に逆再生往復を行うかのフラグです。
	Yoyo bool

	// OnUpdate は毎フレーム補間された最新値 (val float64) を受け取るコールバックです。
	OnUpdate func(val float64)

	// OnRun は毎フレーム補間された進捗率 (progress float64: 0.0 〜 1.0) を受け取るコールバックです。
	OnRun func(progress float64)

	// OnComplete はアニメーションが正常完了した際に発火するコールバックです。
	OnComplete func()
}

// defaultOption は省略された設定項目を補填します。
func defaultOption(opts []Option) Option {
	opt := Option{
		Start:    0.0,
		End:      1.0,
		Duration: 1.0,
		Ease:     EaseLinear,
		Loop:     0,
		Yoyo:     false,
	}

	if len(opts) > 0 {
		userOpt := opts[0]
		opt.Start = userOpt.Start
		opt.End = userOpt.End
		if userOpt.Duration > 0 {
			opt.Duration = userOpt.Duration
		}
		if userOpt.Ease != nil {
			opt.Ease = userOpt.Ease
		}
		opt.Delay = userOpt.Delay
		opt.Loop = userOpt.Loop
		opt.Yoyo = userOpt.Yoyo
		opt.OnUpdate = userOpt.OnUpdate
		opt.OnRun = userOpt.OnRun
		opt.OnComplete = userOpt.OnComplete
	}

	return opt
}
