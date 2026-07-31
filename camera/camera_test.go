package camera_test

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
)

func TestCameraBasicAndGetters(t *testing.T) {
	cam := camera.NewWithName("main", 640, 480)

	if cam.Name() != "main" {
		t.Errorf("expected name 'main', got '%s'", cam.Name())
	}

	cam.SetPos(100, 200)
	if cam.X() != 100 || cam.Y() != 200 {
		t.Errorf("expected pos (100, 200), got (%f, %f)", cam.X(), cam.Y())
	}

	x, y := cam.Pos()
	if x != 100 || y != 200 {
		t.Errorf("Pos() mismatch: (%f, %f)", x, y)
	}

	cam.SetZoom(2.0)
	if cam.Zoom() != 2.0 {
		t.Errorf("expected zoom 2.0, got %f", cam.Zoom())
	}

	cam.SetRotation(math.Pi / 2)
	if cam.Rotation() != math.Pi/2 {
		t.Errorf("expected rotation pi/2, got %f", cam.Rotation())
	}

	cam.SetZIndex(10)
	if cam.ZIndex() != 10 {
		t.Errorf("expected zIndex 10, got %d", cam.ZIndex())
	}
}

func TestScreenToWorldAndWorldToScreen(t *testing.T) {
	cam := camera.New(640, 480)
	cam.SetPos(100, 200)
	cam.SetZoom(1.0)

	// 画面中央 (320, 240) はワールド座標 (100, 200) に一致するはず
	worldX, worldY := cam.ScreenToWorld(320, 240)
	if math.Abs(worldX-100) > 0.001 || math.Abs(worldY-200) > 0.001 {
		t.Errorf("expected ScreenToWorld (100, 200), got (%.2f, %.2f)", worldX, worldY)
	}

	// 逆変換 WorldToScreen(100, 200) は (320, 240) になるはず
	screenX, screenY := cam.WorldToScreen(100, 200)
	if math.Abs(screenX-320) > 0.001 || math.Abs(screenY-240) > 0.001 {
		t.Errorf("expected WorldToScreen (320, 240), got (%.2f, %.2f)", screenX, screenY)
	}
}

func TestVisibleBounds(t *testing.T) {
	cam := camera.New(640, 480)
	cam.SetPos(100, 200)
	cam.SetZoom(1.0)

	// 幅640, 高さ480 ➔ 半分は 320, 240
	// 期待可視領域: (100-320, 200-240) 〜 (100+320, 200+240) ➔ (-220, -40) 〜 (420, 440)
	minX, minY, maxX, maxY := cam.VisibleBounds()
	if minX != -220 || minY != -40 || maxX != 420 || maxY != 440 {
		t.Errorf("VisibleBounds mismatch: (%.1f, %.1f) - (%.1f, %.1f)", minX, minY, maxX, maxY)
	}
}

func TestCameraMoveAndBounds(t *testing.T) {
	cam := camera.New(640, 480)
	cam.SetBounds(0, 0, 1000, 1000)

	// 境界制限ありで移動
	cam.SetPos(0, 0)
	// 幅640(半W:320)のため、MinX:0 ➔ カメラ中心は最小 320 にクランプされる
	if cam.X() != 320 {
		t.Errorf("expected clamped X to be 320, got %f", cam.X())
	}
}

func TestShakeAndUpdate(t *testing.T) {
	cam := camera.New(640, 480)
	cam.Shake(10, 0.5)

	if !cam.IsShaking() {
		t.Errorf("expected camera to be shaking")
	}

	cam.Update(0.6)

	if cam.IsShaking() {
		t.Errorf("expected camera shake to finish after 0.6s")
	}
}

func TestGroupZIndexSorting(t *testing.T) {
	cam1 := camera.NewWithName("cam1", 640, 480)
	cam1.SetZIndex(10)

	cam2 := camera.NewWithName("cam2", 640, 480)
	cam2.SetZIndex(0)

	group := camera.NewGroup(cam1, cam2)

	dummyScreen := ebiten.NewImage(640, 480)
	renderedOrder := make([]string, 0)

	group.Render(dummyScreen, func(c *camera.Camera, target *ebiten.Image) {
		renderedOrder = append(renderedOrder, c.Name())
	})

	if len(renderedOrder) != 2 {
		t.Fatalf("expected 2 rendered cameras, got %d", len(renderedOrder))
	}

	// ZIndex の小さい cam2 (0) が先、cam1 (10) が後で描画されるはず
	if renderedOrder[0] != "cam2" || renderedOrder[1] != "cam1" {
		t.Errorf("ZIndex sort failed: order got %v", renderedOrder)
	}
}
