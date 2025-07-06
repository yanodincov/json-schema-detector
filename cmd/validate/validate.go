package validate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yanodincov/json-ai-schema-detector/pkg/validator"
)

var (
	verbose bool
	strict  bool
)

// Cmd представляет команду validate
var Cmd = &cobra.Command{
	Use:   "validate [data.json] [schema.json]",
	Short: "Валидирует JSON файл против схемы",
	Long: `Валидирует JSON файл против JSON Schema и выводит результат валидации 
с подробным описанием ошибок.`,
	Args: cobra.ExactArgs(2),
	RunE: runValidate,
}

func init() {
	Cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Подробный вывод")
	Cmd.Flags().BoolVarP(&strict, "strict", "s", false, "Строгая валидация")
}

func runValidate(cmd *cobra.Command, args []string) error {
	dataFile := args[0]
	schemaFile := args[1]

	// Проверяем существование файлов
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return fmt.Errorf("файл данных не найден: %s", dataFile)
	}

	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("файл схемы не найден: %s", schemaFile)
	}

	fmt.Printf("Валидация данных: %s\n", dataFile)
	fmt.Printf("Против схемы: %s\n", schemaFile)

	// Создаем валидатор
	validator := validator.New(strict)

	// Выполняем валидацию
	result, err := validator.ValidateFile(dataFile, schemaFile)
	if err != nil {
		return fmt.Errorf("ошибка валидации: %w", err)
	}

	// Выводим результат
	if result.Valid {
		fmt.Printf("✅ Валидация прошла успешно\n")
		if verbose {
			fmt.Printf("Проверено полей: %d\n", result.ValidatedFields)
			fmt.Printf("Время валидации: %s\n", result.Duration)
		}
	} else {
		fmt.Printf("❌ Валидация не пройдена\n")
		fmt.Printf("Найдено ошибок: %d\n", len(result.Errors))

		for i, err := range result.Errors {
			fmt.Printf("  %d. %s\n", i+1, err.Description)
			if verbose {
				fmt.Printf("     Путь: %s\n", err.Field)
				fmt.Printf("     Тип: %s\n", err.Type)
			}
		}

		// Возвращаем код ошибки для CI/CD
		os.Exit(1)
	}

	return nil
}
