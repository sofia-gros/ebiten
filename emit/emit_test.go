package emit_test

import (
	"testing"

	"github.com/sofia-gros/ebiten/emit"
)

type ScoreEvent struct {
	Points int
}

type PlayerDamageEvent struct {
	Damage int
}

func TestEmitAndOn(t *testing.T) {
	e := emit.New()

	var receivedScore int
	emit.On(e, func(ev ScoreEvent) {
		receivedScore += ev.Points
	})

	emit.Emit(e, ScoreEvent{Points: 100})
	emit.Emit(e, ScoreEvent{Points: 50})

	if receivedScore != 150 {
		t.Errorf("expected score 150, got %d", receivedScore)
	}
}

func TestOff(t *testing.T) {
	e := emit.New()

	var count int
	id := emit.On(e, func(ev ScoreEvent) {
		count++
	})

	emit.Emit(e, ScoreEvent{Points: 10})
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// 登録解除
	ok := emit.Off(e, id)
	if !ok {
		t.Errorf("expected Off to return true")
	}

	// 以降は呼ばれないこと
	emit.Emit(e, ScoreEvent{Points: 10})
	if count != 1 {
		t.Errorf("expected count 1 after Off, got %d", count)
	}
}

func TestOnce(t *testing.T) {
	e := emit.New()

	var count int
	emit.Once(e, func(ev ScoreEvent) {
		count++
	})

	emit.Emit(e, ScoreEvent{Points: 10})
	emit.Emit(e, ScoreEvent{Points: 10})
	emit.Emit(e, ScoreEvent{Points: 10})

	if count != 1 {
		t.Errorf("expected Once listener to be called 1 time, got %d", count)
	}
}

func TestRemoveAll(t *testing.T) {
	e := emit.New()

	var count1, count2, countDamage int

	emit.On(e, func(ev ScoreEvent) {
		count1++
	})
	emit.On(e, func(ev ScoreEvent) {
		count2++
	})
	emit.On(e, func(ev PlayerDamageEvent) {
		countDamage++
	})

	emit.Emit(e, ScoreEvent{Points: 10})
	emit.Emit(e, PlayerDamageEvent{Damage: 20})

	if count1 != 1 || count2 != 1 || countDamage != 1 {
		t.Errorf("initial emit failed")
	}

	// ScoreEvent だけを一括解除
	emit.RemoveAll[ScoreEvent](e)

	emit.Emit(e, ScoreEvent{Points: 10})
	emit.Emit(e, PlayerDamageEvent{Damage: 20})

	if count1 != 1 || count2 != 1 {
		t.Errorf("ScoreEvent listeners should not be called after RemoveAll")
	}
	if countDamage != 2 {
		t.Errorf("PlayerDamageEvent listener should still be called, expected 2, got %d", countDamage)
	}
}

func TestQueueAndFlush(t *testing.T) {
	e := emit.New()

	var totalScore int
	emit.On(e, func(ev ScoreEvent) {
		totalScore += ev.Points
	})

	// キューに積む (この時点ではハンドラは実行されない)
	emit.Queue(e, ScoreEvent{Points: 10})
	emit.Queue(e, ScoreEvent{Points: 20})

	if totalScore != 0 {
		t.Errorf("expected 0 before Flush, got %d", totalScore)
	}

	// 一括処理
	e.Flush()

	if totalScore != 30 {
		t.Errorf("expected 30 after Flush, got %d", totalScore)
	}

	// 2回目のFlushでは何も起きない
	e.Flush()
	if totalScore != 30 {
		t.Errorf("expected 30 on second Flush, got %d", totalScore)
	}
}

func TestReset(t *testing.T) {
	e := emit.New()

	var count int
	emit.On(e, func(ev ScoreEvent) {
		count++
	})
	emit.Queue(e, ScoreEvent{Points: 10})

	e.Reset()

	emit.Emit(e, ScoreEvent{Points: 10})
	e.Flush()

	if count != 0 {
		t.Errorf("expected count 0 after Reset, got %d", count)
	}
}

func TestOffInsideHandler(t *testing.T) {
	e := emit.New()

	var id emit.ListenerID
	var count int

	id = emit.On(e, func(ev ScoreEvent) {
		count++
		emit.Off(e, id) // ハンドラ実行中に自分自身を解除
	})

	emit.Emit(e, ScoreEvent{Points: 10})
	emit.Emit(e, ScoreEvent{Points: 10})

	if count != 1 {
		t.Errorf("expected count 1 when Off is called inside handler, got %d", count)
	}
}
