package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// AnimatorOptions は Animator の初期化オプションです。
type AnimatorOptions struct {
	Speed   float64 // タイムスケール (初期値 1.0)
	Loop    bool    // クリップ全体のデフォルトループ設定上書き
	Reverse bool    // 逆再生フラグ
}

// Animator は単一オブジェクトのアニメーション再生状態・コマ送り・切り替えを管理します。
type Animator struct {
	clips        map[string]*Clip
	currentClip  *Clip
	currentName  string
	frameIndex   int
	frameTimer   float64
	speed        float64
	paused       bool
	stopped      bool
	reverse      bool
	pingPongDir  int // 1: forward, -1: backward
	overrideLoop *bool
}

// NewAnimator は初期 Clip を指定して Animator を作成します。
func NewAnimator(defaultClip *Clip, opts ...AnimatorOptions) *Animator {
	a := &Animator{
		clips:       make(map[string]*Clip),
		speed:       1.0,
		pingPongDir: 1,
	}

	if len(opts) > 0 {
		userOpt := opts[0]
		if userOpt.Speed > 0 {
			a.speed = userOpt.Speed
		}
		a.reverse = userOpt.Reverse
	}

	if defaultClip != nil {
		a.AddClip(defaultClip)
		a.Play(defaultClip.Name())
	}

	return a
}

// AddClip はアニメーション Clip を登録します。
func (a *Animator) AddClip(clip *Clip) *Animator {
	if clip != nil {
		a.clips[clip.Name()] = clip
		if a.currentClip == nil {
			a.Play(clip.Name())
		}
	}
	return a
}

// GetClip は名前指定で Clip を取得します。
func (a *Animator) GetClip(name string) *Clip {
	return a.clips[name]
}

// Play は指定した Clip 名にアニメーションを切り替えて再生します。
// すでに再生中の同名 Clip が指定された場合はリセットせず継続します。
func (a *Animator) Play(name string) {
	if a.currentName == name && !a.stopped {
		a.paused = false
		return
	}
	a.PlayWithReset(name)
}

// PlayWithReset は同じ Clip 名であっても強制的にコマ 0 からリセット再生します。
func (a *Animator) PlayWithReset(name string) {
	clip, ok := a.clips[name]
	if !ok {
		return
	}

	a.currentClip = clip
	a.currentName = name
	a.frameTimer = 0.0
	a.paused = false
	a.stopped = false
	a.pingPongDir = 1

	if clip.IsReverse() || a.reverse {
		a.frameIndex = clip.FrameCount() - 1
		if a.frameIndex < 0 {
			a.frameIndex = 0
		}
	} else {
		a.frameIndex = 0
	}

	// 初期コマ到達イベント発火チェック
	a.triggerFrameHandler(a.frameIndex)
}


// Update は毎フレーム時間 (dt 秒) を進行させ、自動コマ送りとイベント発火を行います。
func (a *Animator) Update(dt float64) {
	if a.paused || a.stopped || a.currentClip == nil {
		return
	}

	count := a.currentClip.FrameCount()
	if count == 0 {
		return
	}

	a.frameTimer += dt * a.speed

	for {
		frameDur := a.currentClip.Duration(a.frameIndex)
		if a.frameTimer < frameDur {
			break
		}
		a.frameTimer -= frameDur



		// 進行方向
		isRev := a.currentClip.IsReverse() || a.reverse
		if a.currentClip.IsPingPong() {
			a.frameIndex += a.pingPongDir
			if a.frameIndex >= count {
				a.frameIndex = count - 2
				a.pingPongDir = -1
				if a.frameIndex < 0 {
					a.frameIndex = 0
				}
			} else if a.frameIndex < 0 {
				a.frameIndex = 1
				a.pingPongDir = 1
				if a.frameIndex >= count {
					a.frameIndex = 0
				}
			}
		} else if isRev {
			a.frameIndex--
			if a.frameIndex < 0 {
				isLoop := a.currentClip.IsLoop()
				if a.overrideLoop != nil {
					isLoop = *a.overrideLoop
				}

				if isLoop {
					a.frameIndex = count - 1
				} else {
					a.frameIndex = 0
					a.stopped = true
					if a.currentClip.onComplete != nil {
						a.currentClip.onComplete()
					}
					break
				}
			}
		} else {
			a.frameIndex++
			if a.frameIndex >= count {
				isLoop := a.currentClip.IsLoop()
				if a.overrideLoop != nil {
					isLoop = *a.overrideLoop
				}


				if isLoop {
					a.frameIndex = 0
				} else {
					a.frameIndex = count - 1
					a.stopped = true
					if a.currentClip.onComplete != nil {
						a.currentClip.onComplete()
					}
					break
				}
			}
		}

		// コマ到達イベント発火
		a.triggerFrameHandler(a.frameIndex)
	}
}

func (a *Animator) triggerFrameHandler(idx int) {
	if a.currentClip != nil && a.currentClip.onFrameHandlers != nil {
		if fn, ok := a.currentClip.onFrameHandlers[idx]; ok && fn != nil {
			fn()
		}
	}
}

// CurrentFrame は現在のアニメーションコマの *ebiten.Image を取得します。
func (a *Animator) CurrentFrame() *ebiten.Image {
	if a.currentClip == nil {
		return nil
	}
	return a.currentClip.FrameImage(a.frameIndex)
}

// --- 個別制御メソッド群 ---

func (a *Animator) Pause()                      { a.paused = true }
func (a *Animator) Resume()                     { a.paused = false }
func (a *Animator) Stop()                       { a.stopped = true; a.frameIndex = 0; a.frameTimer = 0 }
func (a *Animator) Restart()                    { if a.currentClip != nil { a.PlayWithReset(a.currentName) } }
func (a *Animator) IsPaused() bool              { return a.paused }
func (a *Animator) IsStopped() bool             { return a.stopped }
func (a *Animator) CurrentClipName() string     { return a.currentName }
func (a *Animator) CurrentFrameIndex() int      { return a.frameIndex }

func (a *Animator) Speed() float64              { return a.speed }
func (a *Animator) SetSpeed(speed float64)      { a.speed = speed }

func (a *Animator) SetReverse(reverse bool)     { a.reverse = reverse }
func (a *Animator) IsReverse() bool             { return a.reverse }

func (a *Animator) SetLoop(loop bool)           { a.overrideLoop = &loop }

func (a *Animator) GoToFrame(index int) {
	if a.currentClip == nil {
		return
	}
	if index >= 0 && index < a.currentClip.FrameCount() {
		a.frameIndex = index
		a.frameTimer = 0
		a.triggerFrameHandler(index)
	}
}
