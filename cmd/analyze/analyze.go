package analyze

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yanodincov/json-ai-schema-detector/pkg/analyzer"
	"github.com/yanodincov/json-ai-schema-detector/pkg/types"
)

var (
	outputFile string
	configFile string
)

// Cmd представляет команду analyze
var Cmd = &cobra.Command{
	Use:   "analyze [input.json]",
	Short: "Анализирует JSON файл и создает схему",
	Long: `Анализирует структуру JSON файла и генерирует соответствующую 
JSON Schema с автоматическим определением типов и структур.`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	Cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Выходной файл для схемы")
	Cmd.Flags().StringVarP(&configFile, "config", "c", "", "Файл конфигурации")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Проверяем существование входного файла
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("входной файл не найден: %s", inputFile)
	}

	// Если выходной файл не указан, создаем его на основе входного
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		outputFile = inputFile[:len(inputFile)-len(ext)] + ".schema.json"
	}

	fmt.Printf("Анализ файла: %s\n", inputFile)
	fmt.Printf("Выходной файл: %s\n", outputFile)

	// Загружаем конфигурацию
	config := types.DefaultConfig()
	if configFile != "" {
		// TODO: Загрузить конфигурацию из файла
		fmt.Printf("Использование конфигурации: %s\n", configFile)
	}

	// Создаем анализатор
	analyzer := analyzer.New(config)

	// Анализируем файл
	result, err := analyzer.AnalyzeFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка анализа: %w", err)
	}

	// Сохраняем результат
	if err := analyzer.SaveSchema(result, outputFile); err != nil {
		return fmt.Errorf("ошибка сохранения схемы: %w", err)
	}

	fmt.Printf("Схема успешно создана: %s\n", outputFile)
	fmt.Printf("Проанализировано объектов: %d\n", result.Statistics.TotalObjects)
	fmt.Printf("Уникальных структур: %d\n", result.Statistics.UniqueStructures)

	return nil
}
