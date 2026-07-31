package asset

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// DataAs は指定したキーのアセットデータを指定した構造体型 T の値にバインドして返します。
// JSON, YAML, TOML などの形式から型安全にデータをパース・変換できます。
func DataAs[T any](m *Manager, key string) (T, error) {
	var zero T
	if m == nil {
		return zero, fmt.Errorf("manager is nil")
	}

	rawData, err := m.Data(key)
	if err != nil {
		return zero, err
	}

	// 1. すでに T 型である場合
	if val, ok := rawData.(T); ok {
		return val, nil
	}

	// 2. []byte データを JSON/YAML 経由で型 T にバインド
	if byteData, ok := rawData.([]byte); ok {
		var result T
		if err := json.Unmarshal(byteData, &result); err == nil {
			return result, nil
		}
		if err := yaml.Unmarshal(byteData, &result); err == nil {
			return result, nil
		}
		return zero, fmt.Errorf("failed to unmarshal byte data into target type %T", zero)
	}

	// 3. 汎用 any (map/slice) データを JSON 再シリアライズ経由で T にバインド
	jsonBytes, err := json.Marshal(rawData)
	if err != nil {
		return zero, fmt.Errorf("failed to marshal raw data: %w", err)
	}

	var result T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return zero, fmt.Errorf("failed to unmarshal into target type %T: %w", zero, err)
	}

	return result, nil
}
