package tween_test

import (
	"math"
	"testing"

	"github.com/sofia-gros/ebiten/tween"
)

func TestBasicTweenAndPlay(t *testing.T) {
	mgr := tween.NewManager()

	var lastVal float64
	completed := false

	tw := mgr.New(tween.Option{
		Start:    0.0,
		End:      100.0,
		Duration: 1.0,
		Ease:     tween.EaseLinear,
	}).OnUpdate(func(val float64) {
		lastVal = val
	}).OnComplete(func() {
		completed = true
	})

	if tw.IsPlaying() {
		t.Errorf("tween should not be playing before Play()")
	}

	tw.Play()

	if !tw.IsPlaying() {
		t.Errorf("tween should be playing after Play()")
	}

	// 0.5 秒進行 (50% 補間)
	mgr.Update(0.5)
	if math.Abs(lastVal-50.0) > 0.001 {
		t.Errorf("expected val 50.0 at 0.5s, got %f", lastVal)
	}

	// 残り 0.5 秒進行 (100% 補間・完了)
	mgr.Update(0.5)
	if math.Abs(lastVal-100.0) > 0.001 {
		t.Errorf("expected val 100.0 at 1.0s, got %f", lastVal)
	}

	if !completed {
		t.Errorf("expected OnComplete to be called")
	}

	if mgr.Count() != 0 {
		t.Errorf("expected finished tween to be removed from manager, got count %d", mgr.Count())
	}
}

func TestPauseResumeRestart(t *testing.T) {
	mgr := tween.NewManager()

	var lastVal float64
	tw := mgr.New(tween.Option{
		Start:    0.0,
		End:      100.0,
		Duration: 1.0,
		Ease:     tween.EaseLinear,
	}).OnUpdate(func(val float64) {
		lastVal = val
	}).Play()

	mgr.Update(0.2) // 20.0
	tw.Pause()

	mgr.Update(0.5) // ポーズ中なので進まない
	if math.Abs(lastVal-20.0) > 0.001 {
		t.Errorf("expected val to remain 20.0 while paused, got %f", lastVal)
	}

	tw.Resume()
	mgr.Update(0.3) // 50.0 まで進む
	if math.Abs(lastVal-50.0) > 0.001 {
		t.Errorf("expected val to be 50.0 after resume, got %f", lastVal)
	}

	tw.Restart() // 最初からリスタート (0.0)
	if math.Abs(lastVal-50.0) > 0.001 {
		// Update で評価されるまで値は残る
	}
	mgr.Update(0.1) // 10.0
	if math.Abs(lastVal-10.0) > 0.001 {
		t.Errorf("expected val to restart at 10.0, got %f", lastVal)
	}
}

func TestGroupOperations(t *testing.T) {
	group := tween.NewGroup()

	var val1, val2 float64

	group.New(tween.Option{
		Start:    0.0,
		End:      100.0,
		Duration: 1.0,
	}).OnUpdate(func(v float64) { val1 = v }).Play()

	group.New(tween.Option{
		Start:    0.0,
		End:      200.0,
		Duration: 1.0,
	}).OnUpdate(func(v float64) { val2 = v }).Play()

	if group.Count() != 2 {
		t.Fatalf("expected 2 tweens in group, got %d", group.Count())
	}

	group.Update(0.5)
	if math.Abs(val1-50.0) > 0.001 || math.Abs(val2-100.0) > 0.001 {
		t.Errorf("group update mismatch: val1=%f, val2=%f", val1, val2)
	}

	group.PauseAll()
	group.Update(0.5)
	// ポーズ中なので数値は進まない
	if math.Abs(val1-50.0) > 0.001 {
		t.Errorf("group PauseAll failed")
	}

	group.Clear()
	if group.Count() != 0 {
		t.Errorf("group Clear failed, count=%d", group.Count())
	}
}

func TestYoyoAndLoop(t *testing.T) {
	mgr := tween.NewManager()

	var lastVal float64
	tw := mgr.New(tween.Option{
		Start:    0.0,
		End:      100.0,
		Duration: 1.0,
		Loop:     1, // 1回リピート (計2サイクル)
		Yoyo:     true,
	}).OnUpdate(func(v float64) {
		lastVal = v
	}).Play()

	// 1サイクル目 0.5s ➔ 50.0
	mgr.Update(0.5)
	if math.Abs(lastVal-50.0) > 0.001 {
		t.Errorf("expected 50.0, got %f", lastVal)
	}

	// 1サイクル目 1.0s ➔ 100.0 (完) ➔ 2サイクル目 (Yoyo 逆再生開始)
	mgr.Update(0.5)

	// 2サイクル目 0.5s (折り返して半分の50.0)
	mgr.Update(0.5)
	if math.Abs(lastVal-50.0) > 0.001 {
		t.Errorf("expected yoyo reverse 50.0, got %f", lastVal)
	}

	// 2サイクル目 1.0s (原点0.0に到着・完了)
	mgr.Update(0.5)
	if math.Abs(lastVal-0.0) > 0.001 {
		t.Errorf("expected yoyo reverse 0.0, got %f", lastVal)
	}

	if !tw.IsFinished() {
		t.Errorf("expected tween to finish after yoyo loop")
	}
}
