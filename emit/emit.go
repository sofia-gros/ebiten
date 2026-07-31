package emit

import (
	"reflect"
)

// On は指定したイベント型 T のハンドラを登録します。
// 登録解除（削除）に使用できる ListenerID を返します。
func On[T any](e *Emitter, handler func(T)) ListenerID {
	if e == nil || handler == nil {
		return 0
	}
	var zero T
	t := reflect.TypeOf(zero)
	return e.addListener(t, handler, false)
}

// Once は指定したイベント型 T のハンドラを1回限り実行されるよう登録します。
// イベントが1回発火してハンドラが呼び出されると、自動的に登録解除（削除）されます。
func Once[T any](e *Emitter, handler func(T)) ListenerID {
	if e == nil || handler == nil {
		return 0
	}
	var zero T
	t := reflect.TypeOf(zero)
	return e.addListener(t, handler, true)
}

// Off は指定した ListenerID のリスナーを登録解除（削除）します。
// 以降、そのリスナーがイベントを受信して実行されることは一切なくなります。
func Off(e *Emitter, id ListenerID) bool {
	if e == nil || id == 0 {
		return false
	}
	return e.removeListener(id)
}

// RemoveAll は指定したイベント型 T のすべてのリスナーを一括して登録解除（削除）します。
func RemoveAll[T any](e *Emitter) {
	if e == nil {
		return
	}
	var zero T
	t := reflect.TypeOf(zero)
	e.removeListenersByType(t)
}

// Emit はイベント event を即座に発行し、登録されているすべてのハンドラを呼び出します。
// イベント発火中の安全な登録解除・追加をサポートするため、ハンドラ一覧のコピーに対して実行されます。
func Emit[T any](e *Emitter, event T) {
	if e == nil {
		return
	}

	var zero T
	t := reflect.TypeOf(zero)

	// ハンドラスライスの安全なスナップショットを作成
	e.mu.RLock()
	rawList, ok := e.listeners[t]
	if !ok || len(rawList) == 0 {
		e.mu.RUnlock()
		return
	}

	// 呼び出し対象のスナップショットを複製
	targets := make([]listenerWrapper, len(rawList))
	copy(targets, rawList)
	e.mu.RUnlock()

	// 1回限り(Once)リスナーの解除用IDリスト
	var onceIDs []ListenerID

	// ハンドラ実行
	for _, l := range targets {
		if fn, ok := l.fn.(func(T)); ok {
			fn(event)
		}
		if l.once {
			onceIDs = append(onceIDs, l.id)
		}
	}

	// Once リスナーの自動解除
	for _, id := range onceIDs {
		e.removeListener(id)
	}
}

// Queue はイベント event を直ちに実行せず、キューに蓄積します。
// 蓄積されたイベントは Flush() を呼び出すことでまとめて配送・実行されます。
func Queue[T any](e *Emitter, event T) {
	if e == nil {
		return
	}

	var zero T
	t := reflect.TypeOf(zero)

	e.queueMu.Lock()
	e.eventQueue = append(e.eventQueue, queuedEvent{
		eventType: t,
		eventData: event,
	})
	e.queueMu.Unlock()
}

// Flush は Queue() で溜められたイベントをすべて順番に発行し、キューをクリアします。
// ゲームループの Update 内など、特定のタイミングで安全に一括処理したい場合に呼び出します。
func (e *Emitter) Flush() {
	if e == nil {
		return
	}

	e.queueMu.Lock()
	if len(e.eventQueue) == 0 {
		e.queueMu.Unlock()
		return
	}

	// キューの安全な退避
	pending := make([]queuedEvent, len(e.eventQueue))
	copy(pending, e.eventQueue)
	e.eventQueue = e.eventQueue[:0]
	e.queueMu.Unlock()

	// 退避したイベントを順次ディスパッチ
	for _, qItem := range pending {
		e.dispatchDynamic(qItem.eventType, qItem.eventData)
	}
}

// dispatchDynamic は動的な型情報を元にイベントを発行する内部ヘルパーです。
func (e *Emitter) dispatchDynamic(t reflect.Type, eventData any) {
	e.mu.RLock()
	rawList, ok := e.listeners[t]
	if !ok || len(rawList) == 0 {
		e.mu.RUnlock()
		return
	}

	targets := make([]listenerWrapper, len(rawList))
	copy(targets, rawList)
	e.mu.RUnlock()

	var onceIDs []ListenerID

	for _, l := range targets {
		// リフレクションによる関数呼び出し
		v := reflect.ValueOf(l.fn)
		if v.IsValid() {
			v.Call([]reflect.Value{reflect.ValueOf(eventData)})
		}
		if l.once {
			onceIDs = append(onceIDs, l.id)
		}
	}

	for _, id := range onceIDs {
		e.removeListener(id)
	}
}
