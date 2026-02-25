package imageutil

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

func GenerateBarcode(data string) ([]byte, error) {
	bc, err := code128.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode barcode: %w", err)
	}
	scaled, err := barcode.Scale(bc, 400, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to scale barcode: %w", err)
	}
	var buf bytes.Buffer
	if err = png.Encode(&buf, scaled); err != nil {
		return nil, fmt.Errorf("failed to encode barcode image: %w", err)
	}
	return buf.Bytes(), nil
}

func GenerateBarcodeBase64(data string) (string, error) {
	dataBytes, err := GenerateBarcode(data)
	if err != nil {
		return "", fmt.Errorf("failed to encode barcode image: %w", err)
	}
	base64Str := base64.StdEncoding.EncodeToString(dataBytes)
	return fmt.Sprintf("data:image/png;base64,%s", base64Str), nil
}
