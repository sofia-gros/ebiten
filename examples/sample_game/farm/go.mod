module github.com/sofia-gros/ebiten/examples/sample_game/farm

go 1.26.5

replace (
	github.com/sofia-gros/ebiten/asset => ../../../asset
	github.com/sofia-gros/ebiten/camera => ../../../camera
	github.com/sofia-gros/ebiten/emit => ../../../emit
	github.com/sofia-gros/ebiten/pad => ../../../pad
	github.com/sofia-gros/ebiten/physics => ../../../physics
	github.com/sofia-gros/ebiten/physics/adapters/arcade => ../../../physics/adapters/arcade
	github.com/sofia-gros/ebiten/save => ../../../save
	github.com/sofia-gros/ebiten/scene => ../../../scene
	github.com/sofia-gros/ebiten/sound => ../../../sound
	github.com/sofia-gros/ebiten/tween => ../../../tween
	github.com/sofia-gros/ebiten/ui => ../../../ui
)

require (
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/sofia-gros/ebiten/asset v0.0.0
	github.com/sofia-gros/ebiten/camera v0.0.0
	github.com/sofia-gros/ebiten/emit v0.0.0
	github.com/sofia-gros/ebiten/pad v0.0.0
	github.com/sofia-gros/ebiten/physics v0.0.0
	github.com/sofia-gros/ebiten/physics/adapters/arcade v0.0.0
	github.com/sofia-gros/ebiten/save v0.0.0
	github.com/sofia-gros/ebiten/scene v0.0.0
	github.com/sofia-gros/ebiten/sound v0.0.0
	github.com/sofia-gros/ebiten/tween v0.0.0
	github.com/sofia-gros/ebiten/ui v0.0.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.4.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/image v0.31.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
