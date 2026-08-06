package physics_test

import (
	"testing"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
)

func TestPhysicsManagerAndArcadeBody(t *testing.T) {
	mgr := physics.NewManager()
	world := arcade.NewWorld()
	mgr.SetWorld(world)

	body1 := mgr.CreateBody(physics.BodyOptions{
		Type:  physics.BodyTypeDynamic,
		X:     10.0,
		Y:     10.0,
		Shape: physics.BoxShape{Width: 16.0, Height: 16.0},
	})

	body2 := mgr.CreateBody(physics.BodyOptions{
		Type:  physics.BodyTypeStatic,
		X:     20.0,
		Y:     10.0,
		Shape: physics.BoxShape{Width: 16.0, Height: 16.0},
	})


	if body1 == nil || body2 == nil {
		t.Fatalf("failed to create physics bodies")
	}

	x, y := body1.Position()
	if x != 10.0 || y != 10.0 {
		t.Errorf("expected position (10,10), got (%.1f, %.1f)", x, y)
	}

	// 速度設定
	body1.SetVelocity(100.0, 0.0)
	vx, vy := body1.Velocity()
	if vx != 100.0 || vy != 0.0 {
		t.Errorf("expected velocity (100,0), got (%.1f, %.1f)", vx, vy)
	}

	// 物理ステップ更新
	mgr.Update(0.1)

	// グループ設定
	body1.SetGroup(1)
	if body1.Group() != 1 {
		t.Errorf("expected group 1, got %d", body1.Group())
	}

	// ユーザーデータ設定
	body1.SetData("player")
	if body1.Data() != "player" {
		t.Errorf("expected data 'player', got %v", body1.Data())
	}
}

