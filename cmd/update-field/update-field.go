package updatefield

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yanodincov/json-ai-schema-detector/pkg/analyzer"
	"github.com/yanodincov/json-ai-schema-detector/pkg/fieldmanager"
	"github.com/yanodincov/json-ai-schema-detector/pkg/types"
)

var (
	interactive bool
	fieldType   string
	description string
)

// Cmd представляет команду update-field
var Cmd = &cobra.Command{
	Use:   "update-field [schema.json] [json-path] [type]",
	Short: "Обновляет поле в схеме (enum, polymorph, description)",
	Long: `Интерактивно обновляет поле в JSON Schema, позволяя:
- Преобразовать поле в enum тип с выбором значений
- Преобразовать поле в полиморфный тип с вариантами
- Добавить или изменить описание поля
- Изменить тип поля

Примеры использования:
  update-field schema.json "data.0.role" enum
  update-field schema.json "data.0.user" polymorph
  update-field schema.json "data.0.id" description`,
	Args: cobra.MinimumNArgs(2),
	RunE: runUpdateField,
}

func init() {
	Cmd.Flags().BoolVarP(&interactive, "interactive", "i", true, "Интерактивный режим")
	Cmd.Flags().StringVarP(&fieldType, "type", "t", "", "Тип поля (enum, polymorph, description)")
	Cmd.Flags().StringVarP(&description, "description", "d", "", "Описание поля")
}

func runUpdateField(cmd *cobra.Command, args []string) error {
	schemaFile := args[0]
	jsonPath := args[1]

	// Определяем тип операции
	operation := fieldType
	if len(args) >= 3 {
		operation = args[2]
	}

	// Проверяем существование файла схемы
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("файл схемы не найден: %s", schemaFile)
	}

	fmt.Printf("🔧 Обновление поля в схеме\n")
	fmt.Printf("📄 Файл схемы: %s\n", schemaFile)
	fmt.Printf("🎯 Путь к полю: %s\n", jsonPath)
	fmt.Printf("🔄 Операция: %s\n", operation)
	fmt.Println()

	// Загружаем схему
	analyzer := analyzer.New(types.DefaultConfig())
	schema, err := analyzer.LoadSchema(schemaFile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки схемы: %w", err)
	}

	// Создаем менеджер полей
	fieldManager := fieldmanager.New()

	// Выполняем операцию в зависимости от типа
	switch operation {
	case "enum":
		err = handleEnumConversion(fieldManager, schema, jsonPath)
	case "polymorph", "polymorphic":
		err = handlePolymorphicConversion(fieldManager, schema, jsonPath)
	case "description", "desc":
		err = handleDescriptionUpdate(fieldManager, schema, jsonPath)
	case "preserve-default", "preserve":
		err = handlePreserveDefaultUpdate(fieldManager, schema, jsonPath)
	default:
		if interactive {
			operation, err = promptOperation()
			if err != nil {
				return err
			}
			return runUpdateField(cmd, append(args[:2], operation))
		}
		return fmt.Errorf("неподдерживаемая операция: %s. Доступные: enum, polymorph, description, preserve-default", operation)
	}

	if err != nil {
		return fmt.Errorf("ошибка обновления поля: %w", err)
	}

	// Сохраняем обновленную схему
	if err := analyzer.SaveSchema(schema, schemaFile); err != nil {
		return fmt.Errorf("ошибка сохранения схемы: %w", err)
	}

	fmt.Printf("✅ Поле успешно обновлено: %s\n", jsonPath)
	return nil
}

func handleEnumConversion(fm *fieldmanager.FieldManager, schema *types.AnalysisResult, jsonPath string) error {
	fmt.Printf("🎯 Преобразование поля в enum тип\n")
	fmt.Printf("Путь: %s\n", jsonPath)
	fmt.Println()

	// Находим поле по пути
	field, err := fm.FindField(schema.Schema, jsonPath)
	if err != nil {
		return fmt.Errorf("поле не найдено: %w", err)
	}

	if field.Type != "string" {
		return fmt.Errorf("преобразование в enum поддерживается только для string полей, текущий тип: %s", field.Type)
	}

	// Интерактивный ввод значений enum
	fmt.Printf("📝 Введите возможные значения для enum (по одному на строку):\n")
	fmt.Printf("💡 Закончите ввод пустой строкой\n")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	var enumValues []interface{}

	for {
		fmt.Print("Значение: ")
		if !scanner.Scan() {
			break
		}

		value := strings.TrimSpace(scanner.Text())
		if value == "" {
			break
		}

		enumValues = append(enumValues, value)
	}

	if len(enumValues) == 0 {
		return fmt.Errorf("не введено ни одного значения для enum")
	}

	// Обновляем поле
	field.Enum = enumValues

	// Добавляем описание
	if interactive {
		fmt.Print("📝 Описание поля (опционально): ")
		if scanner.Scan() {
			desc := strings.TrimSpace(scanner.Text())
			if desc != "" {
				field.Description = desc
			}
		}
	}

	fmt.Printf("✅ Поле преобразовано в enum с %d значениями\n", len(enumValues))
	fmt.Printf("🎯 Значения: %v\n", enumValues)

	return nil
}

func handlePolymorphicConversion(fm *fieldmanager.FieldManager, schema *types.AnalysisResult, jsonPath string) error {
	fmt.Printf("🎯 Преобразование поля в полиморфный тип\n")
	fmt.Printf("Путь: %s\n", jsonPath)
	fmt.Println()

	// Находим поле по пути
	field, err := fm.FindField(schema.Schema, jsonPath)
	if err != nil {
		return fmt.Errorf("поле не найдено: %w", err)
	}

	if field.Type != "object" {
		return fmt.Errorf("преобразование в полиморфный тип поддерживается только для object полей, текущий тип: %s", field.Type)
	}

	fmt.Printf("📝 Создание полиморфного типа\n")
	fmt.Printf("💡 Введите варианты полиморфного типа\n")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	var variants []*types.JSONSchema

	for {
		fmt.Print("Название варианта (или пустая строка для завершения): ")
		if !scanner.Scan() {
			break
		}

		variantName := strings.TrimSpace(scanner.Text())
		if variantName == "" {
			break
		}

		// Создаем базовый вариант
		variant := &types.JSONSchema{
			Type:        "object",
			Properties:  make(map[string]*types.Property),
			Description: fmt.Sprintf("Вариант %s", variantName),
		}

		// Добавляем дискриминатор
		variant.Properties["type"] = &types.Property{
			Type: "string",
			Enum: []interface{}{variantName},
		}

		variants = append(variants, variant)
		fmt.Printf("✅ Добавлен вариант: %s\n", variantName)
	}

	if len(variants) == 0 {
		return fmt.Errorf("не создано ни одного варианта")
	}

	// Обновляем поле как oneOf
	field.OneOf = variants
	field.Type = "" // Убираем базовый тип

	fmt.Printf("✅ Поле преобразовано в полиморфный тип с %d вариантами\n", len(variants))

	return nil
}

func handlePreserveDefaultUpdate(fm *fieldmanager.FieldManager, schema *types.AnalysisResult, jsonPath string) error {
	fmt.Printf("🔒 Защита default значения от перезатирания\n")
	fmt.Printf("Путь: %s\n", jsonPath)
	fmt.Println()

	// Находим поле по пути
	field, err := fm.FindField(schema.Schema, jsonPath)
	if err != nil {
		return fmt.Errorf("поле не найдено: %w", err)
	}

	// Устанавливаем защиту от перезатирания
	field.PreserveDefault = true

	if field.Default != nil {
		fmt.Printf("✅ Default значение защищено: %v\n", field.Default)
	} else {
		fmt.Printf("⚠️ Default значение отсутствует, но защита установлена\n")
		fmt.Printf("💡 При следующем анализе default будет заполнен и защищен\n")
	}

	fmt.Printf("✅ Поле защищено от перезатирания default: %s\n", jsonPath)
	return nil
}

func handleDescriptionUpdate(fm *fieldmanager.FieldManager, schema *types.AnalysisResult, jsonPath string) error {
	fmt.Printf("🎯 Обновление описания поля\n")
	fmt.Printf("Путь: %s\n", jsonPath)
	fmt.Println()

	// Находим поле по пути
	field, err := fm.FindField(schema.Schema, jsonPath)
	if err != nil {
		return fmt.Errorf("поле не найдено: %w", err)
	}

	// Показываем текущее описание
	if field.Description != "" {
		fmt.Printf("📄 Текущее описание: %s\n", field.Description)
	} else {
		fmt.Printf("📄 Текущее описание: отсутствует\n")
	}

	// Интерактивный ввод нового описания
	fmt.Print("📝 Новое описание: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		newDesc := strings.TrimSpace(scanner.Text())
		if newDesc != "" {
			field.Description = newDesc
			fmt.Printf("✅ Описание обновлено: %s\n", newDesc)
		} else {
			fmt.Printf("⚠️ Пустое описание, изменения не внесены\n")
		}
	}

	return nil
}

func promptOperation() (string, error) {
	fmt.Printf("🎯 Выберите операцию:\n")
	fmt.Printf("1. enum - преобразовать в enum тип\n")
	fmt.Printf("2. polymorph - преобразовать в полиморфный тип\n")
	fmt.Printf("3. description - обновить описание\n")
	fmt.Printf("4. preserve-default - защитить default от перезатирания\n")
	fmt.Print("Ваш выбор (1-4): ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		choice := strings.TrimSpace(scanner.Text())
		switch choice {
		case "1":
			return "enum", nil
		case "2":
			return "polymorph", nil
		case "3":
			return "description", nil
		case "4":
			return "preserve-default", nil
		default:
			return "", fmt.Errorf("неверный выбор: %s", choice)
		}
	}

	return "", fmt.Errorf("ошибка ввода")
}
