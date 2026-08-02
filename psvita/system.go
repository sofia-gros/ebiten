package psvita

// BatteryLevel はバッテリー残量をパーセンテージ (0 〜 100) で返します。
func BatteryLevel() int {
	return getBatteryLevelImpl()
}

// IsCharging は現在充電器が接続されているかを返します。
func IsCharging() bool {
	return isChargingImpl()
}

// PowerMode は PSVita の省電力モード状態を示します。
type PowerMode int

const (
	PowerModeNormal PowerMode = iota
	PowerModeSave
)

// CurrentPowerMode は現在の電源モードを返します。
func CurrentPowerMode() PowerMode {
	return getPowerModeImpl()
}

func getBatteryLevelImpl() int {
	// PC環境では 100% 固定フォールバック
	return 100
}

func isChargingImpl() bool {
	return true
}

func getPowerModeImpl() PowerMode {
	return PowerModeNormal
}
