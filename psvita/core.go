package psvita

import "log"

// ClockFrequency は CPU / GPU の動作周波数設定を表します。
type ClockFrequency int

const (
	ClockFrequencyDefault ClockFrequency = 333 // 333MHz (標準)
	ClockFrequencyMax     ClockFrequency = 444 // 444MHz (ブースト)
	ClockFrequencyEco     ClockFrequency = 166 // 166MHz (省電力)
)

// MemoryInfo は PSVita のメモリ使用状況構造体です。
type MemoryInfo struct {
	TotalMainMemoryBytes uint64
	FreeMainMemoryBytes  uint64
	TotalCDRAMMemoryBytes uint64
	FreeCDRAMMemoryBytes  uint64
}

// SetCPUClock は PSVita Core 層の CPU 動作周波数を設定します。
func SetCPUClock(freq ClockFrequency) error {
	return setCPUClockImpl(freq)
}

// GetMemoryInfo は Core 層から現在のメモリ空き状況を取得します。
func GetMemoryInfo() (MemoryInfo, error) {
	return getMemoryInfoImpl()
}

func setCPUClockImpl(freq ClockFrequency) error {
	// PC環境ではログ出力のみで安全にスキップ
	log.Printf("[psvita core] CPU Clock frequency set to %d MHz (PC Emulated)", freq)
	return nil
}

func getMemoryInfoImpl() (MemoryInfo, error) {
	// PSVita メモリ仕様 (メインメモリ 512MB, CDRAM 128MB) のサンプル値
	return MemoryInfo{
		TotalMainMemoryBytes: 512 * 1024 * 1024,
		FreeMainMemoryBytes:  256 * 1024 * 1024,
		TotalCDRAMMemoryBytes: 128 * 1024 * 1024,
		FreeCDRAMMemoryBytes:  64 * 1024 * 1024,
	}, nil
}
