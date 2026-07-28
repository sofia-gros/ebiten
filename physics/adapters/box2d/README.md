# box2d (ebitenphysics adapter)

[English](./README_en.md)

`github.com/ByteArena/box2d` を裏で動かす本格的な剛体物理エンジン用アダプターです。
複雑な衝突判定や、摩擦・反発などの物理シミュレーションをフルサポートします。

## サポートしている機能

`BodyOptions` の全機能に対応しています。

### `BodyOptions` の対応状況

| プロパティ | 対応状況 | 備考 |
| :--- | :--- | :--- |
| `Type` | ◯ | `BodyTypeStatic` (StaticBody), `BodyTypeDynamic` (DynamicBody) |
| `X`, `Y` | ◯ | |
| `Shapes` | ◯ | BoxShape, CircleShape をサポート。複数渡すことで複合形状（複数Fixture）に対応。 |
| `Friction` | ◯ | 摩擦係数として機能します。 |
| `Restitution`| ◯ | 反発係数（バウンス）として機能します。 |
| `Density` | ◯ | 密度として機能します（質量計算に使用）。 |
| `IsSensor` | ◯ | Fixtureをセンサー化します。 |
| `OnCollisionBegin` | ◯ | ContactListenerで実現 |
| `OnCollisionEnd` | ◯ | ContactListenerで実現 |
| `OnOverlap` | ◯ | PreSolve または ContactListener で実現 |

## インストール手順

Box2D本体も必要になります。

```bash
go get github.com/ByteArena/box2d
go get github.com/sofia-gros/ebiten/physics/adapters/box2d
```
