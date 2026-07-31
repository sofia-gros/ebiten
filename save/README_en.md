# ebitensave

[日本語](./README.md)

`ebitensave` is a save data management library for Ebitengine and general Go applications.
It handles JSON, YAML, or binary save files, directory resolution, AES encryption, checksum tampering verification, automatic slot index reuse, and backup recovery.

---

## Features

- **Safe Directory Isolation (`Dir`)**:
  - `DirCurrent` (Default): Saves inside a designated folder relative to working dir (`./save/*`).
  - `DirAppData`: Saves in OS standard application data paths (%APPDATA%/FolderName/).
  - `DirCustom`: Saves in any custom directory.
- **Flat File Naming**:
  - Non-slotted single save: `./save/data.json`
  - Slotted save: `./save/data_1.json`
- **Automatic Slot Reuse (`SaveNewSlot`)**:
  - Automatically scans existing files and fills deleted/missing slot indices first.
- **Encryption & Verification (`EncryptKey`, `Checksum`)**:
  - Supports both human-readable text (JSON/YAML) and encrypted binary data.
- **Corruption Recovery (`Backup`)**:
  - Safely backs up `.bak` files before writing and recovers automatically if file corruption occurs.

---

## Installation

```bash
go get github.com/sofia-gros/ebiten/save
```

---

## Usage

### Quick Start

Save and load plain JSON data without slots.


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/save"
)

type GameData struct {
	PlayerName string `json:"player_name"`
	Level      int    `json:"level"`
}

func main() {
	m := save.NewManager("save", save.Option{
		Dir: save.DirCurrent,
	})

	data := GameData{PlayerName: "Hero", Level: 10}

	// Save single file to ./save/data.json
	_ = m.Save("data", "json", data, save.Option{
		Format: save.FormatJSON,
	})

	// Load single file
	var loaded GameData
	if err := m.Load("data", "json", &loaded); err == nil {
		fmt.Println("Loaded Level:", loaded.Level)
	}
}
```

---

### Full Usage

Slotted saves, automatic index reuse for deleted saves, and AES encryption.


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/save"
)

type GameData struct {
	PlayerName string `json:"player_name"`
	Level      int    `json:"level"`
}

func main() {
	m := save.NewManager("save", save.Option{
		Dir: save.DirCurrent,
	})

	data := GameData{PlayerName: "Hero", Level: 5}

	// Overwrite slot 1 (./save/data_1.json)
	m.Slot(1).Save("data", "json", data)

	// Save to next available slot (fills deleted indices first)
	slotNum, err := m.SaveNewSlot("data", "json", data, save.Option{
		Format: save.FormatJSON,
	})
	if err == nil {
		fmt.Println("Saved into slot:", slotNum)
	}

	// Encrypted binary save (./save/secret.dat)
	m.Save("secret", "dat", data, save.Option{
		Format:     save.FormatBinary,
		EncryptKey: "secret_pass_123",
		Checksum:   true,
		Backup:     true,
	})

	// Delete slot 2
	m.Slot(2).Delete("data", "json")
}
```

---

## License

MIT License
