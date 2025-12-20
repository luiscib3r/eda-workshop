package main

import (
	"backend/cmd/api"
	"backend/cmd/ocr"
	ocrllm "backend/cmd/ocr-llm"
	"backend/cmd/telegram"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ocr-system",
	Short: "OCR System",
}

func init() {
	rootCmd.AddCommand(api.ApiCmd)
	rootCmd.AddCommand(ocr.OcrCmd)
	rootCmd.AddCommand(ocrllm.OcrLlmCmd)
	rootCmd.AddCommand(telegram.TelegramCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
