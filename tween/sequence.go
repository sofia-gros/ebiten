package tween

// Sequence は複数の Tween を順番に連続実行するシーケンス構造体です。
type Sequence struct {
	tweens []func() *Tween
}

// NewSequence は新しい Sequence を作成します。
func NewSequence() *Sequence {
	return &Sequence{
		tweens: make([]func() *Tween, 0),
	}
}

// Append はシーケンスの末尾に実行する Tween の生成関数を追加します。
func (s *Sequence) Append(tweenFunc func() *Tween) *Sequence {
	s.tweens = append(s.tweens, tweenFunc)
	return s
}

// Play はシーケンスを最初から順番に再生開始します。
func (s *Sequence) Play() {
	if len(s.tweens) == 0 {
		return
	}
	s.playStep(0)
}

func (s *Sequence) playStep(index int) {
	if index >= len(s.tweens) {
		return
	}

	tw := s.tweens[index]()
	if tw == nil {
		s.playStep(index + 1)
		return
	}

	origComplete := tw.onComplete
	tw.OnComplete(func() {
		if origComplete != nil {
			origComplete()
		}
		s.playStep(index + 1)
	})

	tw.Play()
}
