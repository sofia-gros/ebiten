# ebitenemit

[日本語](./README.md)

`ebitenemit` is a lightweight, type-safe event dispatcher and emitter library designed for Ebitengine (ebiten) and general Go applications.
It functions independently without any external dependencies, allowing you to easily write event publishing, subscription, and unsubscription for game loops and animations (such as Tweens).

## Features

- **Type-Safe Event Handling**: Utilizes Go Generics (`[T any]`) for casting-free, reflection-free, and type-safe event handling.
- **Intuitive Subscription & Unsubscription (`On` / `Off`)**: `On` returns a `ListenerID`, which can be passed to `Off(e, id)` at any time to remove unwanted listeners.
- **Temporary & Batch Operations**:
  - `Once`: Automatically unsubscribes after firing once.
  - `RemoveAll[T]`: Removes all listeners for a specific event type `T`.
  - `Reset`: Clears all listeners and queued events in the Emitter.
- **Immediate (`Emit`) & Queued (`Queue` / `Flush`) Modes**:
  - `Emit`: Immediately invokes all registered handlers.
  - `Queue` & `Flush`: Queues events for batch processing during specific phases, such as frame updates.
- **Concurrent Safety**: Safe execution even when unsubscribing (`Off`) or subscribing (`On`) inside event handlers.

## Installation

```bash
go get github.com/sofia-gros/ebiten/emit
```

## Usage

### Quick Start

```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/emit"
)

type PlayerDamageEvent struct {
	PlayerID int
	Damage   int
}

func main() {
	e := emit.New()

	// Subscribe (On)
	id := emit.On(e, func(ev PlayerDamageEvent) {
		fmt.Printf("Player %d took %d damage!\n", ev.PlayerID, ev.Damage)
	})

	// Emit immediately
	emit.Emit(e, PlayerDamageEvent{PlayerID: 1, Damage: 35})

	// Unsubscribe (Off)
	emit.Off(e, id)
}
```

---

### Full Usage


```go
type TweenUpdateEvent struct {
	Progress float64
}

type TweenManager struct {
	emitter *emit.Emitter
}

func (t *TweenManager) StartTween() {
	var listenerID emit.ListenerID

	listenerID = emit.On(t.emitter, func(ev TweenUpdateEvent) {
		fmt.Printf("Animating: %.1f%%\n", ev.Progress*100)

		if ev.Progress >= 1.0 {
			emit.Off(t.emitter, listenerID)
			fmt.Println("Tween finished and listener removed.")
		}
	})
}
```

## License

MIT License
