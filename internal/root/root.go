package root

import (
	"github.com/spf13/cobra"
	"github.com/yanodincov/json-schema-detector/internal/analyze"
	listfields "github.com/yanodincov/json-schema-detector/internal/list-fields"
	"github.com/yanodincov/json-schema-detector/internal/update"
	updatefield "github.com/yanodincov/json-schema-detector/internal/update-field"
	"github.com/yanodincov/json-schema-detector/internal/validate"
)

var rootCmd = &cobra.Command{
	Use:   "json-schema-detector",
	Short: "Инструмент для анализа JSON структур и генерации схем",
	Long: `JSON AI Schema Detector - инструмент для автоматического анализа JSON документов
и генерации структурированных схем с поддержкой JSON Schema стандарта.`,
}

func init() {
	// Добавляем подкоманды
	rootCmd.AddCommand(analyze.Cmd)
	rootCmd.AddCommand(listfields.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(updatefield.Cmd)
	rootCmd.AddCommand(validate.Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}
