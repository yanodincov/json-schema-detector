package fieldmanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yanodincov/json-schema-detector/pkg/types"
)

// FieldManager управляет полями в JSON Schema
type FieldManager struct{}

// New создает новый менеджер полей
func New() *FieldManager {
	return &FieldManager{}
}

// FindField находит поле по JSON Path в схеме
func (fm *FieldManager) FindField(schema *types.JSONSchema, jsonPath string) (*types.Property, error) {
	// Парсим JSON Path
	path, err := fm.parseJSONPath(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга пути: %w", err)
	}

	// Начинаем поиск с корневой схемы
	return fm.findFieldRecursive(schema, path, 0)
}

// parseJSONPath парсит JSON Path в массив сегментов
func (fm *FieldManager) parseJSONPath(jsonPath string) ([]string, error) {
	if jsonPath == "" {
		return nil, fmt.Errorf("пустой путь")
	}

	// Убираем начальную точку если есть
	if strings.HasPrefix(jsonPath, ".") {
		jsonPath = jsonPath[1:]
	}

	// Разбиваем по точкам
	segments := strings.Split(jsonPath, ".")

	// Очищаем пустые сегменты
	var cleanSegments []string
	for _, segment := range segments {
		if segment != "" {
			cleanSegments = append(cleanSegments, segment)
		}
	}

	if len(cleanSegments) == 0 {
		return nil, fmt.Errorf("не найдено валидных сегментов пути")
	}

	return cleanSegments, nil
}

// findFieldRecursive рекурсивно находит поле по пути
func (fm *FieldManager) findFieldRecursive(schema *types.JSONSchema, path []string, index int) (*types.Property, error) {
	if index >= len(path) {
		return nil, fmt.Errorf("достигнут конец пути")
	}

	segment := path[index]

	// Проверяем, является ли сегмент числовым индексом
	if _, err := strconv.Atoi(segment); err == nil {
		// Это индекс массива - нужно найти предыдущее поле (массив) и взять его items
		if index == 0 {
			return nil, fmt.Errorf("числовой индекс не может быть первым сегментом")
		}

		// Получаем предыдущее поле
		prevSegment := path[index-1]
		prevField, err := fm.findFieldInSchema(schema, prevSegment)
		if err != nil {
			return nil, fmt.Errorf("не найдено поле %s: %w", prevSegment, err)
		}

		if prevField.Type != "array" || prevField.Items == nil {
			return nil, fmt.Errorf("поле %s не является массивом", prevSegment)
		}

		// Если это последний сегмент, возвращаем items
		if index == len(path)-1 {
			return prevField.Items, nil
		}

		// Иначе продолжаем поиск в items
		itemSchema := fm.propertyToSchema(prevField.Items)
		return fm.findFieldRecursive(itemSchema, path, index+1)
	}

	// Если это последний сегмент, ищем поле
	if index == len(path)-1 {
		return fm.findFieldInSchema(schema, segment)
	}

	// Если это не последний сегмент, идем глубже
	field, err := fm.findFieldInSchema(schema, segment)
	if err != nil {
		return nil, err
	}

	// Проверяем следующий сегмент - если он числовой, то нам нужно обработать его как индекс массива
	if index+1 < len(path) {
		nextSegment := path[index+1]
		if _, err := strconv.Atoi(nextSegment); err == nil {
			// Следующий сегмент - числовой индекс, поэтому текущее поле должно быть массивом
			if field.Type == "array" && field.Items != nil {
				// Пропускаем индекс и идем к содержимому items
				if index+2 >= len(path) {
					// Если индекс - последний сегмент, возвращаем items
					return field.Items, nil
				}
				// Иначе продолжаем поиск в items, пропуская индекс
				itemSchema := fm.propertyToSchema(field.Items)
				return fm.findFieldRecursive(itemSchema, path, index+2)
			}
			return nil, fmt.Errorf("поле %s должно быть массивом для индекса %s", segment, nextSegment)
		}
	}

	// Если поле это объект, работаем с properties
	if field.Type == "object" && field.Properties != nil {
		// Конвертируем Property в JSONSchema для рекурсии
		objSchema := fm.propertyToSchema(field)
		return fm.findFieldRecursive(objSchema, path, index+1)
	}

	return nil, fmt.Errorf("невозможно перейти глубже по пути %s", segment)
}

// findFieldInSchema находит поле в конкретной схеме
func (fm *FieldManager) findFieldInSchema(schema *types.JSONSchema, fieldName string) (*types.Property, error) {
	// Ищем поле по имени
	if schema.Properties != nil {
		if field, exists := schema.Properties[fieldName]; exists {
			return field, nil
		}
	}

	// Если не найдено в основной схеме, проверяем oneOf/anyOf
	if schema.OneOf != nil {
		for _, variant := range schema.OneOf {
			if field, err := fm.findFieldInSchema(variant, fieldName); err == nil {
				return field, nil
			}
		}
	}

	if schema.AnyOf != nil {
		for _, variant := range schema.AnyOf {
			if field, err := fm.findFieldInSchema(variant, fieldName); err == nil {
				return field, nil
			}
		}
	}

	return nil, fmt.Errorf("поле %s не найдено", fieldName)
}

// propertyToSchema конвертирует Property в JSONSchema
func (fm *FieldManager) propertyToSchema(prop *types.Property) *types.JSONSchema {
	schema := &types.JSONSchema{
		Type:        prop.Type,
		Properties:  prop.Properties,
		Required:    prop.Required,
		Enum:        prop.Enum,
		OneOf:       prop.OneOf,
		AnyOf:       prop.AnyOf,
		Description: prop.Description,
	}

	if prop.Items != nil {
		schema.Items = prop.Items
	}

	return schema
}

// schemaToProperty конвертирует JSONSchema в Property
func (fm *FieldManager) schemaToProperty(schema *types.JSONSchema) *types.Property {
	prop := &types.Property{
		Type:        schema.Type,
		Properties:  schema.Properties,
		Required:    schema.Required,
		Enum:        schema.Enum,
		OneOf:       schema.OneOf,
		AnyOf:       schema.AnyOf,
		Description: schema.Description,
	}

	if schema.Items != nil {
		prop.Items = schema.Items
	}

	return prop
}

// ListFields возвращает список всех полей в схеме
func (fm *FieldManager) ListFields(schema *types.JSONSchema) []string {
	var fields []string
	fm.listFieldsRecursive(schema, "", &fields)
	return fields
}

// listFieldsRecursive рекурсивно собирает все поля
func (fm *FieldManager) listFieldsRecursive(schema *types.JSONSchema, prefix string, fields *[]string) {
	if schema.Properties != nil {
		for fieldName, field := range schema.Properties {
			fullPath := fieldName
			if prefix != "" {
				fullPath = prefix + "." + fieldName
			}

			*fields = append(*fields, fullPath)

			// Рекурсивно обрабатываем вложенные объекты
			if field.Type == "object" && field.Properties != nil {
				subSchema := fm.propertyToSchema(field)
				fm.listFieldsRecursive(subSchema, fullPath, fields)
			}

			// Рекурсивно обрабатываем массивы
			if field.Type == "array" && field.Items != nil {
				subSchema := fm.propertyToSchema(field.Items)
				fm.listFieldsRecursive(subSchema, fullPath+".0", fields)
			}
		}
	}

	// Обрабатываем oneOf/anyOf
	if schema.OneOf != nil {
		for i, variant := range schema.OneOf {
			variantPrefix := prefix
			if prefix != "" {
				variantPrefix = prefix + fmt.Sprintf(".oneOf[%d]", i)
			} else {
				variantPrefix = fmt.Sprintf("oneOf[%d]", i)
			}
			fm.listFieldsRecursive(variant, variantPrefix, fields)
		}
	}

	if schema.AnyOf != nil {
		for i, variant := range schema.AnyOf {
			variantPrefix := prefix
			if prefix != "" {
				variantPrefix = prefix + fmt.Sprintf(".anyOf[%d]", i)
			} else {
				variantPrefix = fmt.Sprintf("anyOf[%d]", i)
			}
			fm.listFieldsRecursive(variant, variantPrefix, fields)
		}
	}
}

// UpdateField обновляет поле в схеме
func (fm *FieldManager) UpdateField(schema *types.JSONSchema, jsonPath string, updater func(*types.Property) error) error {
	field, err := fm.FindField(schema, jsonPath)
	if err != nil {
		return err
	}

	return updater(field)
}

// ValidateJSONPath проверяет, существует ли поле по указанному пути
func (fm *FieldManager) ValidateJSONPath(schema *types.JSONSchema, jsonPath string) error {
	_, err := fm.FindField(schema, jsonPath)
	return err
}
