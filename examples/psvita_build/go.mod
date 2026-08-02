module ebiten/examples/psvita_build

go 1.26.5

require (
<<<<<<< HEAD
	ebiten/psvita v0.0.0
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/sofia-gros/ebiten/camera v0.0.0
	github.com/sofia-gros/ebiten/physics v0.0.0
	github.com/sofia-gros/ebiten/physics/adapters/arcade v0.0.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace (
	ebiten/psvita => ../../psvita
	github.com/sofia-gros/ebiten/camera => ../../camera
	github.com/sofia-gros/ebiten/physics => ../../physics
	github.com/sofia-gros/ebiten/physics/adapters/arcade => ../../physics/adapters/arcade
=======
	github.com/hajimehoshi/ebiten/v2 v2.8.6
	github.com/sofia-gros/ebiten/camera v0.0.0
	github.com/sofia-gros/ebiten/physics v0.0.0
	github.com/sofia-gros/ebiten/physics/adapters/arcade v0.0.0
	ebiten/psvita v0.0.0
)

replace (
	github.com/sofia-gros/ebiten/camera => "../../camera"
	github.com/sofia-gros/ebiten/physics => "../../physics"
	github.com/sofia-gros/ebiten/physics/adapters/arcade => "../../physics/adapters/arcade"
	ebiten/psvita => "../../psvita"
>>>>>>> 6e4bfbdd6ca2ebbefeddec852bdf6e3739544407
)
