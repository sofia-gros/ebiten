package emit

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// ListenerID はリスナーを一意に識別するためのIDです。
// Off(e, id) を使用してリスナーの登録解除を行う際に使用します。
type ListenerID uint64

// listenerWrapper はイベントハンドラとそのメタデータを保持する内部構造体です。
type listenerWrapper struct {
	id   ListenerID
	once bool
	fn   any // func(T)
}

// queuedEvent は Queue で保持される一括処理用イベント構造体です。
type queuedEvent struct {
	eventType reflect.Type
	eventData any
}

// Emitter は型安全なイベントの発行・購読・キューイングを管理するイベントマネージャーです。
// 他のパッケージに一切依存せず単体で動作します。
type Emitter struct {
	mu          sync.RWMutex
	listeners   map[reflect.Type][]listenerWrapper
	nextID      atomic.Uint64
	queueMu     sync.Mutex
	eventQueue  []queuedEvent
	isEmitting  bool
}

// New は新しい Emitter インスタンスを作成します。
func New() *Emitter {
	return &Emitter{
		listeners:  make(map[reflect.Type][]listenerWrapper),
		eventQueue: make([]queuedEvent, 0),
	}
}

// nextListenerID は新しいユニークな ListenerID を発行します。
func (e *Emitter) nextListenerID() ListenerID {
	return ListenerID(e.nextID.Add(1))
}

// addListener は内部的にハンドラを登録し、ListenerID を返します。
func (e *Emitter) addListener(t reflect.Type, fn any, once bool) ListenerID {
	e.mu.Lock()
	defer e.mu.Unlock()

	id := e.nextListenerID()
	wrapper := listenerWrapper{
		id:   id,
		once: once,
		fn:   fn,
	}

	e.listeners[t] = append(e.listeners[t], wrapper)
	return id
}

// removeListener は指定された ID のリスナーを登録解除（削除）します。
func (e *Emitter) removeListener(id ListenerID) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	for t, list := range e.listeners {
		for i, l := range list {
			if l.id == id {
				// スライスから要素を削除
				e.listeners[t] = append(list[:i], list[i+1:]...)
				if len(e.listeners[t]) == 0 {
					delete(e.listeners, t)
				}
				return true
			}
		}
	}
	return false
}

// removeListenersByType は指定されたイベント型のすべてのリスナーを登録解除（削除）します。
func (e *Emitter) removeListenersByType(t reflect.Type) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.listeners, t)
}

// Reset は Emitter 内のすべてのリスナーおよびキューに入っている未処理イベントを全消去します。
func (e *Emitter) Reset() {
	e.mu.Lock()
	e.listeners = make(map[reflect.Type][]listenerWrapper)
	e.mu.Unlock()

	e.queueMu.Lock()
	e.eventQueue = e.eventQueue[:0]
	e.queueMu.Unlock()
}
