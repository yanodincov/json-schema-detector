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

	// Анализируем корневой элемент
	schema, err := a.analyzeValue(data, "", result.Statistics)
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
		return &types.Property{Type: "string"}, nil
	case float64:
		stats.TypeDistribution["number"]++
		return &types.Property{Type: "number"}, nil
	case bool:
		stats.TypeDistribution["boolean"]++
		return &types.Property{Type: "boolean"}, nil
	case nil:
		stats.TypeDistribution["null"]++
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
	// TODO: Реализовать логику объединения схем
	// Пока просто возвращаем новый результат
	return new, nil
}
