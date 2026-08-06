# box2dgo (ebitenphysics adapter)

[English](./README_en.md)

`github.com/oliverbestmann/box2d-go` (Box2D v3のGo言語移植版) を裏で動かす剛体物理エンジン用アダプターです。
従来の Box2D(v2) アダプターと比較して、より高速でモダンな内部設計のエンジンを利用することができます。複雑な衝突判定や摩擦・反発などの物理シミュレーションをサポートします。

## サポートしている機能

`BodyOptions` の全機能に対応しています。

### `BodyOptions` の対応状況

| プロパティ | 対応状況 | 備考 |
| :--- | :--- | :--- |
| `Type` | ◯ | `BodyTypeStatic` (StaticBody), `BodyTypeDynamic` (DynamicBody) |
| `X`, `Y` | ◯ | |
| `Shapes` | ◯ | BoxShape, CircleShape をサポート。複数渡すことで複合形状（複数Shape）に対応。 |
| `Friction` | ◯ | 摩擦係数として機能します。 |
| `Restitution`| ◯ | 反発係数（バウンス）として機能します。 |
| `Density` | ◯ | 密度として機能します（質量計算に使用）。 |
| `IsSensor` | ◯ | 形状をセンサー化します。 |
| `OnCollisionBegin` | ◯ | `ContactEvents.BeginEvents` で実現 |
| `OnCollisionEnd` | ◯ | `ContactEvents.EndEvents` で実現 |
| `OnOverlap` | △ | （現在は未実装：将来のアップデートでセンサーイベント等にてサポート予定） |

## インストール手順

Box2D-Go 本体も必要になります。

```bash
go get github.com/oliverbestmann/box2d-go
go get github.com/sofia-gros/ebiten/physics/adapters/box2dgo
```
