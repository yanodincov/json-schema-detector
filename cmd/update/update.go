package update

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanodincov/json-ai-schema-detector/pkg/analyzer"
	"github.com/yanodincov/json-ai-schema-detector/pkg/types"
)

var (
	inputFile  string
	configFile string
)

// Cmd представляет команду update
var Cmd = &cobra.Command{
	Use:   "update [schema.json]",
	Short: "Обновляет существующую схему новыми данными",
	Long: `Обновляет существующую JSON Schema новыми данными из JSON файла,
сохраняя существующие описания и комментарии.`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

func init() {
	Cmd.Flags().StringVarP(&inputFile, "input", "i", "", "JSON файл с новыми данными")
	Cmd.Flags().StringVarP(&configFile, "config", "c", "", "Файл конфигурации")
	Cmd.MarkFlagRequired("input")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	schemaFile := args[0]

	// Проверяем существование файлов
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("файл схемы не найден: %s", schemaFile)
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("входной файл не найден: %s", inputFile)
	}

	fmt.Printf("Обновление схемы: %s\n", schemaFile)
	fmt.Printf("Новые данные: %s\n", inputFile)

	// Загружаем конфигурацию
	config := types.DefaultConfig()
	if configFile != "" {
		// TODO: Загрузить конфигурацию из файла
		fmt.Printf("Использование конфигурации: %s\n", configFile)
	}

	// Создаем анализатор
	analyzer := analyzer.New(config)

	// Загружаем существующую схему
	existingSchema, err := analyzer.LoadSchema(schemaFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки схемы: %w", err)
	}

	// Анализируем новые данные
	newResult, err := analyzer.AnalyzeFile(inputFile)
	if err != nil {
		return fmt.Errorf("ошибка анализа новых данных: %w", err)
	}

	// Объединяем схемы
	mergedResult, err := analyzer.MergeResults(existingSchema, newResult)
	if err != nil {
		return fmt.Errorf("ошибка объединения схем: %w", err)
	}

	// Сохраняем обновленную схему
	if err := analyzer.SaveSchema(mergedResult, schemaFile); err != nil {
		return fmt.Errorf("ошибка сохранения схемы: %w", err)
	}

	fmt.Printf("Схема успешно обновлена: %s\n", schemaFile)
	fmt.Printf("Добавлено новых объектов: %d\n", newResult.Statistics.TotalObjects)

	return nil
}
