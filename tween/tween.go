package tween

import (
	"sync"
)

// Tween は時間経過に伴う値の補間アニメーションと動的制御を行う構造体です。
type Tween struct {
	mu           sync.RWMutex
	manager      *Manager
	startVal     float64
	endVal       float64
	duration     float64
	ease         EaseFunc
	delay        float64
	loop         int
	yoyo         bool
	onUpdate     func(val float64)
	onRun        func(progress float64)
	onComplete   func()

	// 進行状態変数
	elapsedSec   float64
	delayTimer   float64
	currentLoop  int
	isReversing  bool
	isPlaying    bool
	isPaused     bool
	isFinished   bool
}

// New は Option 構造体またはデフォルト値で待機状態の Tween を作成します。
// 作成後に .Play() を呼ぶことで再生が開始され、マネージャーに自動登録されます。
func New(opts ...Option) *Tween {
	opt := defaultOption(opts)
	return &Tween{
		manager:    DefaultManager,
		startVal:   opt.Start,
		endVal:     opt.End,
		duration:   opt.Duration,
		ease:       opt.Ease,
		delay:      opt.Delay,
		loop:       opt.Loop,
		yoyo:       opt.Yoyo,
		onUpdate:   opt.OnUpdate,
		onRun:      opt.OnRun,
		onComplete: opt.OnComplete,
		delayTimer: opt.Delay,
	}
}

// FromTo は開始値、終了目標値、所要時間(秒)を指定して Tween を作成します。
func FromTo(start, end, duration float64) *Tween {
	return New(Option{
		Start:    start,
		End:      end,
		Duration: duration,
	})
}

// --- メソッドチェーン設定関数 ---

// Ease は使用するイージング関数をセットします。
func (t *Tween) Ease(ease EaseFunc) *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	if ease != nil {
		t.ease = ease
	}
	return t
}

// Delay はアニメーション開始前の遅延時間 (秒) をセットします。
func (t *Tween) Delay(delaySec float64) *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.delay = delaySec
	t.delayTimer = delaySec
	return t
}

// Loop はリピート回数をセットします (-1 で無限リピート)。
func (t *Tween) Loop(count int) *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.loop = count
	return t
}

// Yoyo は往復リピートの ON/OFF をセットします。
func (t *Tween) Yoyo(yoyo bool) *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.yoyo = yoyo
	return t
}

// OnUpdate は毎フレーム補間された最新数値を受け取るコールバックをセットします。
func (t *Tween) OnUpdate(fn func(val float64)) *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onUpdate = fn
	return t
}

// OnRun は毎フレーム補間された進捗率 (0.0 〜 1.0) を受け取るコールバックをセットします。
func (t *Tween) OnRun(fn func(progress float64)) *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onRun = fn
	return t
}

// OnComplete はアニメーション完了時に発火するコールバックをセットします。
func (t *Tween) OnComplete(fn func()) *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onComplete = fn
	return t
}

// --- 動的コントロールメソッド群 ---

// Play はアニメーションの再生を開始し、マネージャーへ登録します。
func (t *Tween) Play() *Tween {
	t.mu.Lock()
	t.isPlaying = true
	t.isPaused = false
	t.isFinished = false
	mgr := t.manager
	t.mu.Unlock()

	if mgr != nil {
		mgr.Add(t)
	}
	return t
}

// Pause はアニメーションを一時停止します。
func (t *Tween) Pause() *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.isPlaying {
		t.isPaused = true
	}
	return t
}

// Resume は一時停止を解除して再生を再開します。
func (t *Tween) Resume() *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.isPlaying && t.isPaused {
		t.isPaused = false
	}
	return t
}

// TogglePause はポーズ状態を反転させます。
func (t *Tween) TogglePause() *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.isPlaying {
		t.isPaused = !t.isPaused
	}
	return t
}

// Restart は経過時間を 0 にリセットし、最初から再生をリスタートします。
func (t *Tween) Restart() *Tween {
	t.mu.Lock()
	t.elapsedSec = 0
	t.delayTimer = t.delay
	t.currentLoop = 0
	t.isReversing = false
	t.isPlaying = true
	t.isPaused = false
	t.isFinished = false
	mgr := t.manager
	t.mu.Unlock()

	if mgr != nil {
		mgr.Add(t)
	}
	return t
}

// Reset は進行時間を 0 にリセットします (再生/ポーズ状態は維持)。
func (t *Tween) Reset() *Tween {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.elapsedSec = 0
	t.delayTimer = t.delay
	t.currentLoop = 0
	t.isReversing = false
	return t
}

// Stop はアニメーションを強制停止し、マネージャーから除去します。
func (t *Tween) Stop() {
	t.stopInternal()
	if mgr := t.manager; mgr != nil {
		mgr.Remove(t)
	}
}

func (t *Tween) stopInternal() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.isPlaying = false
	t.isFinished = true
}


// --- クエリ・ステータス表示 ---

// IsPlaying は現在再生中であるかを返します。
func (t *Tween) IsPlaying() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isPlaying && !t.isPaused && !t.isFinished
}

// IsPaused は一時停止中であるかを返します。
func (t *Tween) IsPaused() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isPaused
}

// IsFinished は完了したかを返します。
func (t *Tween) IsFinished() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isFinished
}

// Progress は現在のイージング適用前の進捗率 (0.0 〜 1.0) を返します。
func (t *Tween) Progress() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.duration <= 0 {
		return 1.0
	}
	p := t.elapsedSec / t.duration
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	return p
}

// CurrentValue は現在のイージング適用後の補間値を返します。
func (t *Tween) CurrentValue() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.calculateCurrentVal()
}

// --- 内部更新処理 ---

// updateInternal はマネージャーから毎フレーム呼び出される内部タイマー更新処理です。
func (t *Tween) updateInternal(dt float64) bool {
	t.mu.Lock()

	if !t.isPlaying || t.isPaused || t.isFinished {
		t.mu.Unlock()
		return t.isFinished
	}

	// 1. 遅延タイマーの消化
	if t.delayTimer > 0 {
		t.delayTimer -= dt
		if t.delayTimer > 0 {
			t.mu.Unlock()
			return false
		}
		// 遅延時間を超過した分を経過時間に充当
		dt = -t.delayTimer
		t.delayTimer = 0
	}

	// 2. 経過時間の進捗
	t.elapsedSec += dt

	completedThisCycle := false
	if t.elapsedSec >= t.duration {
		t.elapsedSec = t.duration
		completedThisCycle = true
	}

	// 3. 補間値と進捗率の算出
	currVal := t.calculateCurrentVal()
	progress := t.elapsedSec / t.duration

	onUpdate := t.onUpdate
	onRun := t.onRun
	onComplete := t.onComplete

	// 4. リピート・Yoyo判定
	if completedThisCycle {
		if t.yoyo {
			t.isReversing = !t.isReversing
		}

		if t.loop < 0 || t.currentLoop < t.loop {
			t.currentLoop++
			t.elapsedSec = 0
		} else {
			t.isPlaying = false
			t.isFinished = true
		}
	}

	isFinished := t.isFinished
	t.mu.Unlock()

	// コールバックの発火 (ロック外で安全に実行)
	if onUpdate != nil {
		onUpdate(currVal)
	}
	if onRun != nil {
		onRun(progress)
	}

	if completedThisCycle && isFinished && onComplete != nil {
		onComplete()
	}

	return isFinished
}

func (t *Tween) calculateCurrentVal() float64 {
	p := t.elapsedSec / t.duration
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}

	if t.isReversing {
		p = 1.0 - p
	}

	easedP := p
	if t.ease != nil {
		easedP = t.ease(p)
	}

	return t.startVal + (t.endVal-t.startVal)*easedP
}
