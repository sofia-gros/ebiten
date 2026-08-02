package psvita

import (
	"testing"
)

func TestPSVitaConstants(t *testing.T) {
	if ScreenWidth != 960 {
		t.Errorf("expected ScreenWidth to be 960, got %d", ScreenWidth)
	}
	if ScreenHeight != 544 {
		t.Errorf("expected ScreenHeight to be 544, got %d", ScreenHeight)
	}
}

func TestPSVitaSystem(t *testing.T) {
	bat := BatteryLevel()
	if bat < 0 || bat > 100 {
		t.Errorf("invalid battery level: %d", bat)
	}
}

func TestPSVitaCore(t *testing.T) {
	if err := SetCPUClock(ClockFrequencyMax); err != nil {
		t.Errorf("SetCPUClock failed: %v", err)
	}

	mem, err := GetMemoryInfo()
	if err != nil {
		t.Errorf("GetMemoryInfo failed: %v", err)
	}
	if mem.TotalMainMemoryBytes == 0 {
		t.Errorf("expected non-zero total main memory")
	}
}
