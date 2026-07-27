package scene_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/scene"
)

// モックシーン
type MockScene struct {
	UpdateFunc func(ctx *scene.Context) error
	DrawFunc   func(screen *ebiten.Image)
}

func (m *MockScene) Update(ctx *scene.Context) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx)
	}
	return nil
}

func (m *MockScene) Draw(screen *ebiten.Image) {
	if m.DrawFunc != nil {
		m.DrawFunc(screen)
	}
}

type SceneA struct {
	MockScene
}

type SceneB struct {
	MockScene
}

func TestManager_Start(t *testing.T) {
	m := scene.NewManager(800, 600)

	var updated bool
	sceneA := &SceneA{
		MockScene: MockScene{
			UpdateFunc: func(ctx *scene.Context) error {
				updated = true
				return nil
			},
		},
	}

	// 初期状態ではシーンがないためUpdateを呼んでも何も起きない
	m.Update()

	// Contextを取り出してStart
	ctx := m.Context()
	ctx.Start(sceneA)

	// まだコマンドがflushされていないため、updatedはfalse
	if updated {
		t.Errorf("expected updated to be false before flush")
	}

	// Updateを呼ぶとコマンドがflushされてSceneAのUpdateが実行される
	m.Update()

	if !updated {
		t.Errorf("expected updated to be true after flush")
	}

	// Overlayのテスト
	var updatedB bool
	sceneB := &SceneB{
		MockScene: MockScene{
			UpdateFunc: func(ctx *scene.Context) error {
				updatedB = true
				return nil
			},
		},
	}
	ctx.Overlay(sceneB)
	updated = false // リセット
	m.Update()

	if !updated || !updatedB {
		t.Errorf("expected both scenes to be updated in Overlay state: A=%v, B=%v", updated, updatedB)
	}

	// Hideのテスト (SceneAを隠す)
	ctx.Hide(sceneA)
	updated = false
	updatedB = false
	m.Update()

	if updated {
		t.Errorf("expected sceneA to not update after Hide")
	}
	if !updatedB {
		t.Errorf("expected sceneB to update")
	}
}
