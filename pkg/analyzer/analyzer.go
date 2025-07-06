package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yanodincov/json-ai-schema-detector/pkg/types"
)

// Analyzer представляет анализатор JSON структур
type Analyzer struct {
	config *types.Config
}

// New создает новый анализатор
func New(config *types.Config) *Analyzer {
	return &Analyzer{
		config: config,
	}
}

// AnalyzeFile анализирует JSON файл и возвращает результат
func (a *Analyzer) AnalyzeFile(filename string) (*types.AnalysisResult, error) {
	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	// Парсим JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	// Анализируем структуру
	return a.analyzeData(jsonData)
}

// analyzeData анализирует JSON данные
func (a *Analyzer) analyzeData(data interface{}) (*types.AnalysisResult, error) {
	// Создаем результат
	result := &types.AnalysisResult{
		Metadata: &types.AnalysisMetadata{
			GeneratedAt: time.Now(),
			Version:     "1.0.0",
		},
		Statistics: &types.AnalysisStatistics{
			FieldFrequency:   make(map[string]int),
			TypeDistribution: make(map[string]int),
			EnumCandidates:   make(map[string][]interface{}),
		},
	}

	// Определяем тип корневого элемента
	var schema *types.Property
	var err error

	switch v := data.(type) {
	case map[string]interface{}:
		// Проверяем, есть ли поле 'data' - массив
		if dataField, exists := v["data"]; exists {
			if _, ok := dataField.([]interface{}); ok {
				// Это структура с массивом данных
				schema, err = a.analyzeValue(data, "", result.Statistics)
			} else {
				// Поле 'data' существует, но не массив - анализируем как обычный объект
				schema, err = a.analyzeValue(data, "", result.Statistics)
			}
		} else {
			// Нет поля 'data' - считаем за один объект
			schema, err = a.analyzeValue(data, "", result.Statistics)
		}
	default:
		// Анализируем как есть
		schema, err = a.analyzeValue(data, "", result.Statistics)
	}

	if err != nil {
		return nil, err
	}

	// Создаем JSON Schema
	result.Schema = &types.JSONSchema{
		Schema:      "http://json-schema.org/draft-07/schema#",
		Type:        schema.Type,
		Properties:  schema.Properties,
		Items:       schema.Items,
		Required:    schema.Required,
		Default:     schema.Default,
		Description: "Generated JSON Schema",
	}

	return result, nil
}

// analyzeValue анализирует JSON значение
func (a *Analyzer) analyzeValue(value interface{}, path string, stats *types.AnalysisStatistics) (*types.Property, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return a.analyzeObject(v, path, stats)
	case []interface{}:
		return a.analyzeArray(v, path, stats)
	case string:
		stats.TypeDistribution["string"]++
		property := &types.Property{Type: "string"}
		if v != "" { // Заполняем default только если строка не пустая
			property.Default = v
		}
		return property, nil
	case float64:
		stats.TypeDistribution["number"]++
		property := &types.Property{Type: "number"}
		if v != 0 { // Заполняем default только если число не равно 0
			property.Default = v
		}
		return property, nil
	case bool:
		stats.TypeDistribution["boolean"]++
		property := &types.Property{Type: "boolean"}
		// Для boolean всегда заполняем default
		property.Default = v
		return property, nil
	case nil:
		stats.TypeDistribution["null"]++
		// Для null не заполняем default
		return &types.Property{Type: "null"}, nil
	default:
		return nil, fmt.Errorf("неподдерживаемый тип данных: %T", v)
	}
}

// analyzeObject анализирует объект
func (a *Analyzer) analyzeObject(obj map[string]interface{}, path string, stats *types.AnalysisStatistics) (*types.Property, error) {
	stats.TypeDistribution["object"]++
	stats.TotalObjects++

	property := &types.Property{
		Type:       "object",
		Properties: make(map[string]*types.Property),
		Required:   make([]string, 0),
	}

	// Анализируем каждое поле
	for key, value := range obj {
		fieldPath := path + "." + key
		stats.FieldFrequency[key]++

		fieldProperty, err := a.analyzeValue(value, fieldPath, stats)
		if err != nil {
			return nil, err
		}

		property.Properties[key] = fieldProperty
		property.Required = append(property.Required, key)
	}

	return property, nil
}

// analyzeArray анализирует массив
func (a *Analyzer) analyzeArray(arr []interface{}, path string, stats *types.AnalysisStatistics) (*types.Property, error) {
	stats.TypeDistribution["array"]++

	property := &types.Property{
		Type: "array",
	}

	if len(arr) == 0 {
		return property, nil
	}

	// Анализируем первый элемент для определения типа элементов
	itemProperty, err := a.analyzeValue(arr[0], path+"[0]", stats)
	if err != nil {
		return nil, err
	}

	property.Items = itemProperty
	return property, nil
}

// SaveSchema сохраняет схему в файл
func (a *Analyzer) SaveSchema(result *types.AnalysisResult, filename string) error {
	// Создаем JSON Schema с метаданными
	schema := result.Schema
	if schema.Extensions == nil {
		schema.Extensions = make(map[string]interface{})
	}
	schema.Extensions["x-analysis-meta"] = result.Metadata

	// Сериализуем в JSON
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации схемы: %w", err)
	}

	// Записываем в файл
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

// LoadSchema загружает схему из файла
func (a *Analyzer) LoadSchema(filename string) (*types.AnalysisResult, error) {
	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	// Парсим JSON Schema
	var schema types.JSONSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("ошибка парсинга схемы: %w", err)
	}

	// Извлекаем метаданные
	result := &types.AnalysisResult{
		Schema: &schema,
	}

	// TODO: Извлечь метаданные из extensions
	result.Metadata = &types.AnalysisMetadata{
		GeneratedAt: time.Now(),
		Version:     "1.0.0",
	}

	return result, nil
}

// MergeResults объединяет результаты анализа
func (a *Analyzer) MergeResults(existing, new *types.AnalysisResult) (*types.AnalysisResult, error) {
	// Обновляем схему с учетом новых данных
	a.mergeProperties(existing.Schema.Properties, new.Schema.Properties, "")

	// Обновляем статистики
	if existing.Statistics != nil && new.Statistics != nil {
		for key, count := range new.Statistics.FieldFrequency {
			existing.Statistics.FieldFrequency[key] += count
		}
		for key, count := range new.Statistics.TypeDistribution {
			existing.Statistics.TypeDistribution[key] += count
		}
		existing.Statistics.TotalObjects += new.Statistics.TotalObjects
	}

	return existing, nil
}

// mergeProperties рекурсивно объединяет свойства схем
func (a *Analyzer) mergeProperties(existing, new map[string]*types.Property, path string) {
	for key, newProp := range new {
		currentPath := path + "." + key
		if currentPath[0] == '.' {
			currentPath = currentPath[1:]
		}

		if existingProp, exists := existing[key]; exists {
			// Поле уже существует - обновляем
			a.mergeProperty(existingProp, newProp, currentPath)
		} else {
			// Новое поле - добавляем
			existing[key] = newProp
		}
	}
}

// mergeProperty объединяет два свойства
func (a *Analyzer) mergeProperty(existing, new *types.Property, path string) {
	// Обновляем default значения
	if !existing.PreserveDefault {
		a.updateDefaultValue(existing, new)
	}

	// Рекурсивно обновляем вложенные свойства
	if existing.Type == "object" && new.Type == "object" {
		if existing.Properties == nil {
			existing.Properties = make(map[string]*types.Property)
		}
		if new.Properties != nil {
			a.mergeProperties(existing.Properties, new.Properties, path)
		}
	}

	// Для массивов обновляем items
	if existing.Type == "array" && new.Type == "array" {
		if existing.Items != nil && new.Items != nil {
			a.mergeProperty(existing.Items, new.Items, path+"[0]")
		}
	}
}

// updateDefaultValue обновляет default значение согласно правилам
func (a *Analyzer) updateDefaultValue(existing, new *types.Property) {
	// Если у существующего свойства нет default, устанавливаем из нового
	if existing.Default == nil && new.Default != nil {
		existing.Default = new.Default
		return
	}

	// Если у существующего есть default, а у нового другое значение - обнуляем default
	if existing.Default != nil && new.Default != nil {
		if !a.isEqualValue(existing.Default, new.Default) {
			existing.Default = nil
		}
	}
}

// isEqualValue сравнивает два значения
func (a *Analyzer) isEqualValue(a1, a2 interface{}) bool {
	// Простое сравнение значений
	return fmt.Sprintf("%v", a1) == fmt.Sprintf("%v", a2)
}
