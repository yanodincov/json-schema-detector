package types

import (
	"time"
)

// AnalysisResult представляет результат анализа JSON структуры
type AnalysisResult struct {
	Schema     *JSONSchema         `json:"schema"`
	Metadata   *AnalysisMetadata   `json:"metadata"`
	Statistics *AnalysisStatistics `json:"statistics"`
}

// JSONSchema представляет JSON Schema
type JSONSchema struct {
	Schema      string                 `json:"$schema"`
	Type        string                 `json:"type"`
	Properties  map[string]*Property   `json:"properties,omitempty"`
	Items       *Property              `json:"items,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Enum        []interface{}          `json:"enum,omitempty"`
	OneOf       []*JSONSchema          `json:"oneOf,omitempty"`
	AnyOf       []*JSONSchema          `json:"anyOf,omitempty"`
	Description string                 `json:"description,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Extensions  map[string]interface{} `json:"-"`
}

// Property представляет свойство в JSON Schema
type Property struct {
	Type        string                 `json:"type"`
	Properties  map[string]*Property   `json:"properties,omitempty"`
	Items       *Property              `json:"items,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Enum        []interface{}          `json:"enum,omitempty"`
	OneOf       []*JSONSchema          `json:"oneOf,omitempty"`
	AnyOf       []*JSONSchema          `json:"anyOf,omitempty"`
	Description string                 `json:"description,omitempty"`
	Default     interface{}            `json:"default,omitempty"`
	Extensions  map[string]interface{} `json:"-"`

	// Дополнительные поля для управления поведением
	PreserveDefault bool `json:"x-preserve-default,omitempty"` // Защита от перезатирания default
}

// AnalysisMetadata содержит метаданные анализа
type AnalysisMetadata struct {
	EnumValues        map[string][]interface{} `json:"enum_values,omitempty"`
	OptionalFields    []string                 `json:"optional_fields,omitempty"`
	PolymorphicFields map[string][]string      `json:"polymorphic_patterns,omitempty"`
	GeneratedAt       time.Time                `json:"generated_at"`
	Version           string                   `json:"version"`
}

// AnalysisStatistics содержит статистику анализа
type AnalysisStatistics struct {
	TotalObjects     int                      `json:"total_objects"`
	UniqueStructures int                      `json:"unique_structures"`
	FieldFrequency   map[string]int           `json:"field_frequency"`
	TypeDistribution map[string]int           `json:"type_distribution"`
	EnumCandidates   map[string][]interface{} `json:"enum_candidates"`
}

// JSONType представляет тип JSON значения
type JSONType string

const (
	TypeString  JSONType = "string"
	TypeNumber  JSONType = "number"
	TypeBoolean JSONType = "boolean"
	TypeObject  JSONType = "object"
	TypeArray   JSONType = "array"
	TypeNull    JSONType = "null"
)
