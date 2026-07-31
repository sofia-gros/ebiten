package asset

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
)

// decodeImage は画像ファイルデータを *ebiten.Image にデコードします。
func decodeImage(r io.Reader) (*ebiten.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return ebiten.NewImageFromImage(img), nil
}

// detectFormat はファイルパスの拡張子から DataFormat を推定します。
func detectFormat(path string) DataFormat {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return FormatJSON
	case ".yaml", ".yml":
		return FormatYAML
	case ".toml":
		return FormatTOML
	case ".csv":
		return FormatCSV
	case ".txt":
		return FormatText
	default:
		return FormatRaw
	}
}

// decodeStructuredData は JSON, YAML, CSV などの構造化データをパースします。
func decodeStructuredData(data []byte, format DataFormat) (any, error) {
	switch format {
	case FormatJSON:
		var result any
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("json unmarshal failed: %w", err)
		}
		return result, nil

	case FormatYAML:
		var result any
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("yaml unmarshal failed: %w", err)
		}
		return result, nil

	case FormatCSV:
		r := csv.NewReader(bytes.NewReader(data))
		records, err := r.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("csv read failed: %w", err)
		}
		return records, nil

	case FormatText:
		return string(data), nil

	case FormatRaw:
		return data, nil

	default:
		// デフォルトは JSON 試行、ダメなら YAML 試行
		var result any
		if err := json.Unmarshal(data, &result); err == nil {
			return result, nil
		}
		if err := yaml.Unmarshal(data, &result); err == nil {
			return result, nil
		}
		return data, nil
	}
}
