package asset_test

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/sofia-gros/ebiten/asset"
)

// テスト用画像ファイルをヘルパー作成
func createDummyPNG(t *testing.T, path string, width, height int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create dummy png: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("failed to encode dummy png: %v", err)
	}
}

func TestImageAndSpriteSheet(t *testing.T) {
	tmpDir := t.TempDir()
	imgFilePath := filepath.Join(tmpDir, "sheet.png")
	createDummyPNG(t, imgFilePath, 64, 64)

	m := asset.NewManager()

	// 単体画像登録
	m.AddImage("single_img", imgFilePath)

	// 32x32のコマ割り画像登録
	m.AddSprite("hero_sheet", imgFilePath, asset.Option{
		FrameWidth:  32,
		FrameHeight: 32,
	})

	m.Load()

	// 単体画像の検証
	img, err := m.Image("single_img")
	if err != nil || img == nil {
		t.Fatalf("failed to get single image: %v", err)
	}

	// コマ割りスプライトの検証 (64x64 を 32x32 で切るので 4 コマになるはず)
	sheet, err := m.Sprite("hero_sheet")
	if err != nil || sheet == nil {
		t.Fatalf("failed to get SpriteSheet: %v", err)
	}

	if sheet.Count() != 4 {
		t.Errorf("expected 4 frames, got %d", sheet.Count())
	}

	frame0 := sheet.Frame(0)
	if frame0 == nil {
		t.Errorf("frame 0 should not be nil")
	}

	allFrames := sheet.Frames()
	if len(allFrames) != 4 {
		t.Errorf("expected 4 frames in slice, got %d", len(allFrames))
	}
}

func TestJSONAndDataAs(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "config.json")

	jsonData := []byte(`{
		"framerate": 12.0,
		"loop": true,
		"tags": {
			"walk": {
				"name": "walk",
				"frames": [0, 1, 2, 1],
				"loop": true
			}
		}
	}`)
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		t.Fatalf("failed to write json: %v", err)
	}

	m := asset.NewManager()
	m.AddData("config_json", jsonPath)
	m.AddAnimation("hero_anim", jsonPath)

	m.Load()

	// 汎用 DataAs で型安全取得
	anim, err := asset.DataAs[asset.AnimationData](m, "config_json")
	if err != nil {
		t.Fatalf("DataAs failed: %v", err)
	}

	if anim.FrameRate != 12.0 {
		t.Errorf("expected framerate 12.0, got %f", anim.FrameRate)
	}

	walkTag, ok := anim.Tags["walk"]
	if !ok {
		t.Fatalf("expected 'walk' tag in AnimationData")
	}

	if len(walkTag.Frames) != 4 {
		t.Errorf("expected 4 frames in walk tag, got %d", len(walkTag.Frames))
	}
}

func TestUnloadAndClear(t *testing.T) {
	tmpDir := t.TempDir()
	imgFilePath := filepath.Join(tmpDir, "test.png")
	createDummyPNG(t, imgFilePath, 32, 32)

	m := asset.NewManager()
	m.AddImage("img1", imgFilePath)
	m.AddImage("img2", imgFilePath)

	m.Load()

	// 解放前の確認
	_, err1 := m.Image("img1")
	_, err2 := m.Image("img2")
	if err1 != nil || err2 != nil {
		t.Fatalf("failed initial image load")
	}

	// img1 を個別アンロード
	m.Unload("img1")

	_, err1After := m.Image("img1")
	if err1After == nil {
		t.Errorf("expected img1 to be unloaded and return error")
	}

	// 全クリア
	m.Clear()

	_, err2After := m.Image("img2")
	if err2After == nil {
		t.Errorf("expected img2 to be cleared and return error")
	}
}

func TestProgressEvents(t *testing.T) {
	tmpDir := t.TempDir()
	imgFilePath := filepath.Join(tmpDir, "test.png")
	createDummyPNG(t, imgFilePath, 16, 16)

	m := asset.NewManager()

	var startedTotal int
	var progressCount int
	var completed bool

	m.OnStart(func(total int) {
		startedTotal = total
	})
	m.OnProgress(func(progress float64, key string) {
		progressCount++
	})
	m.OnComplete(func() {
		completed = true
	})

	m.AddImage("img1", imgFilePath)
	m.AddImage("img2", imgFilePath)

	m.Load()

	if startedTotal != 2 {
		t.Errorf("expected OnStart total 2, got %d", startedTotal)
	}
	if progressCount != 2 {
		t.Errorf("expected OnProgress called 2 times, got %d", progressCount)
	}
	if !completed {
		t.Errorf("expected OnComplete to be called")
	}
}
