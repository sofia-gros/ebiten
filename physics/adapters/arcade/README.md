# arcade (ebitenphysics adapter)

自作の軽量なAABB（Axis-Aligned Bounding Box / めり込み防止）物理エンジンアダプターです。
外部の依存ライブラリを一切持たず、RPGやプラットフォーマーなど、単純な四角形や円形の押し出し・重なり検知に特化しています。

## サポートしている機能

このエンジンは軽量さを優先しているため、高度な物理シミュレーション機能（摩擦や反発）には対応していません。

### `BodyOptions` の対応状況

| プロパティ | 対応状況 | 備考 |
| :--- | :--- | :--- |
| `Type` | ◯ | `BodyTypeStatic` と `BodyTypeDynamic` のみサポート。 |
| `X`, `Y` | ◯ | |
| `Shapes` | ◯ | `BoxShape` と `CircleShape` をサポート。オフセット付きの複合形状にも対応。 |
| `Friction` | ✕ | 無視されます。 |
| `Restitution`| ✕ | 無視されます。 |
| `Density` | ✕ | 無視されます。（すべてのDynamic物体は等しい押し返し力を持ちます） |
| `IsSensor` | ◯ | 物理的な押し返しは行わず、コールバックのみ発火します。 |
| `OnCollisionBegin` | ◯ | |
| `OnCollisionEnd` | ◯ | |
| `OnOverlap` | ◯ | |

## インストール

```bash
go get github.com/sofia-gros/ebiten/physics/adapters/arcade
```
