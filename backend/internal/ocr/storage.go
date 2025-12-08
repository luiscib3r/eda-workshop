package ocr

import "fmt"

func PageImageKey(fileKey string, pageImageKey string) string {
	return fmt.Sprintf("images/%s/%s.png", fileKey, pageImageKey)
}

func PageImagePrefix(fileKey string) string {
	return fmt.Sprintf("images/%s/", fileKey)
}
