package listfields

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yanodincov/json-ai-schema-detector/pkg/analyzer"
	"github.com/yanodincov/json-ai-schema-detector/pkg/fieldmanager"
)

var (
	showTypes bool
	verbose   bool
)

// Cmd представляет команду list-fields
var Cmd = &cobra.Command{
	Use:   "list-fields [schema.json]",
	Short: "Показывает список всех полей в схеме",
	Long: `Отображает список всех доступных полей в JSON Schema с их путями.
Полезно для определения правильного JSON Path для команды update-field.

Примеры использования:
  list-fields schema.json
  list-fields schema.json --types
  list-fields schema.json --verbose`,
	Args: cobra.ExactArgs(1),
	RunE: runListFields,
}

func init() {
	Cmd.Flags().BoolVarP(&showTypes, "types", "t", false, "Показать типы полей")
	Cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Подробный вывод")
}

func runListFields(cmd *cobra.Command, args []string) error {
	schemaFile := args[0]

	// Проверяем существование файла схемы
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("файл схемы не найден: %s", schemaFile)
	}

	fmt.Printf("📋 Список полей в схеме: %s\n", schemaFile)
	fmt.Println()

	// Загружаем схему
	analyzer := analyzer.New()
	schema, err := analyzer.LoadSchema(schemaFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки схемы: %w", err)
	}

	// Создаем менеджер полей
	fieldManager := fieldmanager.New()

	// Получаем список полей
	fields := fieldManager.ListFields(schema.Schema)

	if len(fields) == 0 {
		fmt.Println("⚠️ Поля не найдены в схеме")
		return nil
	}

	// Сортируем поля для удобства
	sort.Strings(fields)

	fmt.Printf("🎯 Найдено полей: %d\n", len(fields))
	fmt.Println()

	// Выводим список полей
	for i, fieldPath := range fields {
		fmt.Printf("%3d. %s", i+1, fieldPath)

		if showTypes || verbose {
			// Получаем информацию о поле
			field, err := fieldManager.FindField(schema.Schema, fieldPath)
			if err == nil {
				fmt.Printf(" (%s)", field.Type)

				if verbose {
					// Дополнительная информация
					if field.Description != "" {
						fmt.Printf(" - %s", field.Description)
					}

					if field.Enum != nil {
						fmt.Printf(" [enum: %v]", field.Enum)
					}

					if field.OneOf != nil {
						fmt.Printf(" [polymorphic: %d variants]", len(field.OneOf))
					}
				}
			}
		}

		fmt.Println()
	}

	fmt.Println()
	fmt.Printf("💡 Используйте пути из списка с командой update-field:\n")
	fmt.Printf("   ./json-schema-detector update-field %s \"<path>\" <operation>\n", schemaFile)
	fmt.Println()

	return nil
}
