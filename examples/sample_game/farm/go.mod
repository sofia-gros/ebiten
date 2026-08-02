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
	golang.org/x/image v0.31.0
)
