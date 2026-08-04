package animation

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestBasicClipAndAnimator(t *testing.T) {
	img1 := ebiten.NewImage(16, 16)
	img2 := ebiten.NewImage(16, 16)
	img3 := ebiten.NewImage(16, 16)

	frames := []*ebiten.Image{img1, img2, img3}

	clipWalk := NewClip("walk", frames, ClipOptions{
		FPS:  10, // 0.1s per frame
		Loop: true,
	})

	if clipWalk.Name() != "walk" {
		t.Errorf("expected clip name 'walk', got %s", clipWalk.Name())
	}
	if clipWalk.FrameCount() != 3 {
		t.Errorf("expected 3 frames, got %d", clipWalk.FrameCount())
	}

	animator := NewAnimator(clipWalk)
	if animator.CurrentClipName() != "walk" {
		t.Errorf("expected current clip 'walk', got %s", animator.CurrentClipName())
	}
	if animator.CurrentFrameIndex() != 0 {
		t.Errorf("expected initial frame index 0, got %d", animator.CurrentFrameIndex())
	}

	// 0.05秒進行 (まだコマ0のまま)
	animator.Update(0.05)
	if animator.CurrentFrameIndex() != 0 {
		t.Errorf("expected frame index 0, got %d", animator.CurrentFrameIndex())
	}

	// さらに 0.06秒進行 (合計0.11秒 -> コマ1へ移動)
	animator.Update(0.06)
	if animator.CurrentFrameIndex() != 1 {
		t.Errorf("expected frame index 1, got %d", animator.CurrentFrameIndex())
	}
}

func TestCallbacksAndLoop(t *testing.T) {
	img1 := ebiten.NewImage(16, 16)
	img2 := ebiten.NewImage(16, 16)
	frames := []*ebiten.Image{img1, img2}

	hitFrameTriggered := false
	completedTriggered := false

	clipAttack := NewClip("attack", frames, ClipOptions{
		FPS:  10, // 0.1s per frame
		Loop: false,
	})

	clipAttack.OnFrame(1, func() {
		hitFrameTriggered = true
	})

	clipAttack.OnComplete(func() {
		completedTriggered = true
	})

	animator := NewAnimator(clipAttack)

	// コマ0 -> コマ1 へ進行
	animator.Update(0.11)
	if !hitFrameTriggered {
		t.Errorf("expected OnFrame(1) handler to be called")
	}

	// コマ1 -> 再生完了
	animator.Update(0.11)
	if !completedTriggered {
		t.Errorf("expected OnComplete handler to be called")
	}
	if !animator.IsStopped() {
		t.Errorf("expected animator to be stopped after non-looping clip finishes")
	}
}

func TestManagerControl(t *testing.T) {
	mgr := NewManager()

	img := ebiten.NewImage(16, 16)
	// 2コマのアニメーション (FPS: 10 = 1コマ0.1秒)
	clip := NewClip("idle", []*ebiten.Image{img, img}, ClipOptions{FPS: 10, Loop: true})

	anim1 := mgr.CreateAnimator(clip)
	anim2 := mgr.CreateAnimator(clip)

	if anim1 == nil || anim2 == nil {
		t.Fatalf("failed to create animators in manager")
	}

	// 0.1秒進めると コマ0 -> コマ1 へ移行
	mgr.Update(0.1)
	if anim1.CurrentFrameIndex() != 1 || anim2.CurrentFrameIndex() != 1 {
		t.Errorf("expected frame 1 after 0.1s, got anim1=%d, anim2=%d", anim1.CurrentFrameIndex(), anim2.CurrentFrameIndex())
	}

	// 全体一時停止
	mgr.PauseAll()
	mgr.Update(0.5) // ポーズ中なので 0.5秒進めてもコマ1のまま
	if anim1.CurrentFrameIndex() != 1 || anim2.CurrentFrameIndex() != 1 {
		t.Errorf("expected animators to stay paused at frame 1, got anim1=%d", anim1.CurrentFrameIndex())
	}

	// 一括再開
	mgr.ResumeAll()
	mgr.Update(0.1) // 再開して0.1秒進むと コマ1 -> コマ0 へループ
	if anim1.CurrentFrameIndex() != 0 || anim2.CurrentFrameIndex() != 0 {
		t.Errorf("expected animators to resume and loop to frame 0, got anim1=%d", anim1.CurrentFrameIndex())
	}
}










func TestJSONLoader(t *testing.T) {
	baseImg := ebiten.NewImage(64, 64)

	jsonContent := []byte(`{
		"frames": [
			{ "filename": "idle_0.png", "frame": {"x": 0, "y": 0, "w": 16, "h": 16}, "duration": 100 },
			{ "filename": "idle_1.png", "frame": {"x": 16, "y": 0, "w": 16, "h": 16}, "duration": 100 }
		],
		"meta": {
			"frameTags": [
				{ "name": "idle", "from": 0, "to": 1, "direction": "forward" }
			]
		}
	}`)

	mgr := NewManager()
	animator, err := mgr.CreateAnimatorFromJSON(baseImg, jsonContent)
	if err != nil {
		t.Fatalf("failed to load JSON animation: %v", err)
	}

	if animator.CurrentClipName() != "idle" {
		t.Errorf("expected loaded clip name 'idle', got %s", animator.CurrentClipName())
	}
}
