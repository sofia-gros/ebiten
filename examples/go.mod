module github.com/sofia-gros/ebiten/examples

go 1.25.1

replace github.com/sofia-gros/ebiten/physics => ../physics

replace github.com/sofia-gros/ebiten/pad => ../pad

require (
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/sofia-gros/ebiten/pad v0.0.0-00010101000000-000000000000
	github.com/sofia-gros/ebiten/physics v0.0.0
	github.com/sofia-gros/ebiten/physics/adapters/arcade v0.0.0-20260727103936-064daed58be2
	github.com/sofia-gros/ebiten/physics/adapters/box2d v0.0.0-20260727103936-064daed58be2
	github.com/sofia-gros/ebiten/physics/adapters/phygo v0.0.0-20260727103936-064daed58be2
)

require (
	github.com/ByteArena/box2d v1.0.2 // indirect
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)
