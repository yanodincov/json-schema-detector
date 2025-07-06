package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

// Validator представляет валидатор JSON схем
type Validator struct {
	strict bool
}

// ValidationResult представляет результат валидации
type ValidationResult struct {
	Valid           bool              `json:"valid"`
	Errors          []ValidationError `json:"errors,omitempty"`
	ValidatedFields int               `json:"validated_fields"`
	Duration        time.Duration     `json:"duration"`
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field       string      `json:"field"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Value       interface{} `json:"value,omitempty"`
}

// New создает новый валидатор
func New(strict bool) *Validator {
	return &Validator{
		strict: strict,
	}
}

// ValidateFile валидирует JSON файл против схемы
func (v *Validator) ValidateFile(dataFile, schemaFile string) (*ValidationResult, error) {
	start := time.Now()

	// Читаем файл данных
	dataBytes, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла данных: %w", err)
	}

	// Читаем файл схемы
	schemaBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла схемы: %w", err)
	}

	// Валидируем
	result, err := v.ValidateBytes(dataBytes, schemaBytes)
	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(start)
	return result, nil
}

// ValidateBytes валидирует JSON данные против схемы
func (v *Validator) ValidateBytes(data, schema []byte) (*ValidationResult, error) {
	// Создаем загрузчики для gojsonschema
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(data)

	// Выполняем валидацию
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("ошибка валидации: %w", err)
	}

	// Преобразуем результат
	validationResult := &ValidationResult{
		Valid:  result.Valid(),
		Errors: make([]ValidationError, 0),
	}

	// Если есть ошибки, преобразуем их
	if !result.Valid() {
		for _, desc := range result.Errors() {
			validationResult.Errors = append(validationResult.Errors, ValidationError{
				Field:       desc.Field(),
				Type:        desc.Type(),
				Description: desc.Description(),
				Value:       desc.Value(),
			})
		}
	}

	// Подсчитываем количество проверенных полей
	validationResult.ValidatedFields = v.countFields(data)

	return validationResult, nil
}

// countFields подсчитывает количество полей в JSON
func (v *Validator) countFields(data []byte) int {
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return 0
	}

	return v.countFieldsRecursive(jsonData)
}

// countFieldsRecursive рекурсивно подсчитывает поля
func (v *Validator) countFieldsRecursive(data interface{}) int {
	count := 0

	switch val := data.(type) {
	case map[string]interface{}:
		for _, value := range val {
			count++                                // Считаем само поле
			count += v.countFieldsRecursive(value) // Рекурсивно считаем вложенные поля
		}
	case []interface{}:
		for _, item := range val {
			count += v.countFieldsRecursive(item)
		}
	default:
		// Примитивные типы не добавляют к счетчику
	}

	return count
}
