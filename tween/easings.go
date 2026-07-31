package tween

import "math"

// EaseFunc は 0.0 〜 1.0 の進捗率 t を受け取り、補間された進捗率 (0.0 〜 1.0 前後) を返すイージング関数の型です。
type EaseFunc func(t float64) float64

// --- 標準イージング関数群 ---

// EaseLinear (直線的)
func EaseLinear(t float64) float64 { return t }

// --- Quad ---
func EaseInQuad(t float64) float64  { return t * t }
func EaseOutQuad(t float64) float64 { return t * (2 - t) }
func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// --- Cubic ---
func EaseInCubic(t float64) float64  { return t * t * t }
func EaseOutCubic(t float64) float64 { t--; return t*t*t + 1 }
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return (t-1)*(2*t-2)*(2*t-2) + 1
}

// --- Quart ---
func EaseInQuart(t float64) float64  { return t * t * t * t }
func EaseOutQuart(t float64) float64 { t--; return 1 - t*t*t*t }
func EaseInOutQuart(t float64) float64 {
	if t < 0.5 {
		return 8 * t * t * t * t
	}
	t--
	return 1 - 8*t*t*t*t
}

// --- Quint ---
func EaseInQuint(t float64) float64  { return t * t * t * t * t }
func EaseOutQuint(t float64) float64 { t--; return t*t*t*t*t + 1 }
func EaseInOutQuint(t float64) float64 {
	if t < 0.5 {
		return 16 * t * t * t * t * t
	}
	t--
	return 16*t*t*t*t*t + 1
}

// --- Sine ---
func EaseInSine(t float64) float64  { return 1 - math.Cos((t*math.Pi)/2) }
func EaseOutSine(t float64) float64 { return math.Sin((t * math.Pi) / 2) }
func EaseInOutSine(t float64) float64 {
	return -(math.Cos(math.Pi*t) - 1) / 2
}

// --- Expo ---
func EaseInExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	return math.Pow(2, 10*(t-1))
}
func EaseOutExpo(t float64) float64 {
	if t == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}
func EaseInOutExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	if t < 0.5 {
		return math.Pow(2, 20*t-10) / 2
	}
	return (2 - math.Pow(2, -20*t+10)) / 2
}

// --- Circ ---
func EaseInCirc(t float64) float64  { return 1 - math.Sqrt(1-t*t) }
func EaseOutCirc(t float64) float64 { t--; return math.Sqrt(1 - t*t) }
func EaseInOutCirc(t float64) float64 {
	if t < 0.5 {
		return (1 - math.Sqrt(1-4*t*t)) / 2
	}
	t = 2*t - 2
	return (math.Sqrt(1-t*t) + 1) / 2
}

// --- Back ---
func EaseInBack(t float64) float64 {
	s := 1.70158
	return t * t * ((s+1)*t - s)
}
func EaseOutBack(t float64) float64 {
	s := 1.70158
	t--
	return t*t*((s+1)*t+s) + 1
}
func EaseInOutBack(t float64) float64 {
	s := 1.70158 * 1.525
	t *= 2
	if t < 1 {
		return (t * t * ((s+1)*t - s)) / 2
	}
	t -= 2
	return (t*t*((s+1)*t+s) + 2) / 2
}

// --- Elastic ---
func EaseInElastic(t float64) float64 {
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	return -math.Pow(2, 10*(t-1)) * math.Sin((t-1.1)*(2*math.Pi)/0.3)
}
func EaseOutElastic(t float64) float64 {
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	return math.Pow(2, -10*t)*math.Sin((t-0.1)*(2*math.Pi)/0.3) + 1
}
func EaseInOutElastic(t float64) float64 {
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	t *= 2
	if t < 1 {
		return -0.5 * math.Pow(2, 10*(t-1)) * math.Sin((t-1.1)*(2*math.Pi)/0.45)
	}
	return 0.5*math.Pow(2, -10*(t-1))*math.Sin((t-1.1)*(2*math.Pi)/0.45) + 1
}

// --- Bounce ---
func EaseOutBounce(t float64) float64 {
	n1 := 7.5625
	d1 := 2.75

	if t < 1/d1 {
		return n1 * t * t
	} else if t < 2/d1 {
		t -= 1.5 / d1
		return n1*t*t + 0.75
	} else if t < 2.5/d1 {
		t -= 2.25 / d1
		return n1*t*t + 0.9375
	} else {
		t -= 2.625 / d1
		return n1*t*t + 0.984375
	}
}

func EaseInBounce(t float64) float64 {
	return 1 - EaseOutBounce(1-t)
}

func EaseInOutBounce(t float64) float64 {
	if t < 0.5 {
		return (1 - EaseOutBounce(1-2*t)) / 2
	}
	return (1 + EaseOutBounce(2*t-1)) / 2
}
