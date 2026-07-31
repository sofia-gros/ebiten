package sound_test

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/sofia-gros/ebiten/sound"
)

// ダミーの 44.1kHz 16bit ステレオ PCM 音声データを作成するヘルパー
func createDummyPCM(seconds float64, sampleRate int) []byte {
	numSamples := int(float64(sampleRate) * seconds)
	buf := new(bytes.Buffer)

	for i := 0; i < numSamples; i++ {
		// 440Hz サイン波
		t := float64(i) / float64(sampleRate)
		sample := int16(math.Sin(2*math.Pi*440*t) * 10000)
		_ = binary.Write(buf, binary.LittleEndian, sample) // Left
		_ = binary.Write(buf, binary.LittleEndian, sample) // Right
	}
	return buf.Bytes()
}

func TestManagerBasicAndTypes(t *testing.T) {
	m := sound.NewManager(44100)

	// 音量セッター・ゲッターテスト
	m.SetMasterVolume(0.8)
	if m.MasterVolume() != 0.8 {
		t.Errorf("expected master volume 0.8, got %f", m.MasterVolume())
	}

	m.SetVolume(sound.TypeSE, 0.9)
	if m.Volume(sound.TypeSE) != 0.9 {
		t.Errorf("expected TypeSE volume 0.9, got %f", m.Volume(sound.TypeSE))
	}

	// iota カスタムタイプ
	const TypeSystem sound.Type = sound.TypeCustom + 1
	m.SetVolume(TypeSystem, 0.5)
	if m.Volume(TypeSystem) != 0.5 {
		t.Errorf("expected TypeSystem volume 0.5, got %f", m.Volume(TypeSystem))
	}
}

func TestPlayAndStop(t *testing.T) {
	m := sound.NewManager(44100)
	dummyPCM := createDummyPCM(0.5, 44100)

	// 再生
	track := m.Play(dummyPCM, sound.Option{Type: sound.TypeSE})
	if track == nil {
		t.Fatalf("failed to play sound track")
	}

	// 即時停止
	m.Stop(track)

	// フェード指定一括停止テスト
	track2 := m.Play(dummyPCM, sound.Option{Type: sound.TypeBGM})
	if track2 == nil {
		t.Fatalf("failed to play BGM track")
	}

	m.StopType(sound.TypeBGM, 1.0) // 1.0秒フェードアウト
	m.Update(0.5)                   // 0.5秒経過

	// まだ停止完了していない
	m.Update(0.6) // 計1.1秒経過で停止完了
}

func TestPositionalAndSetPosition(t *testing.T) {
	m := sound.NewManager(44100)
	dummyPCM := createDummyPCM(0.5, 44100)

	// ポジショナル再生
	track := m.PlayAt(dummyPCM, 100, 0, 0, 0, 500.0, sound.Option{Type: sound.TypeSE})
	if track == nil {
		t.Fatalf("failed to PlayAt sound")
	}

	// 位置更新 (動的追従)
	track.SetPosition(200, 0, 0, 0, 500.0)
}

func TestSampleRate(t *testing.T) {
	m := sound.NewManager(44100)
	if m.SampleRate() != 44100 {
		t.Errorf("expected sample rate 44100, got %d", m.SampleRate())
	}

	m.SetSampleRate(48000)
	if m.SampleRate() != 48000 {
		t.Errorf("expected sample rate 48000, got %d", m.SampleRate())
	}
}
