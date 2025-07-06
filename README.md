# JSON AI Schema Detector

Инструмент для автоматического анализа JSON документов и генерации структурированных схем с поддержкой JSON Schema стандарта.

## Возможности

- 🔍 **Автоматический анализ типов данных** - определение примитивных и составных типов
- 📋 **Генерация JSON Schema** - создание стандартных схем JSON Schema
- 🔄 **Обновление схем** - слияние новых данных с существующими схемами
- ✅ **Валидация** - проверка JSON данных против схем
- 📊 **Статистика** - подробная аналитика по структурам данных
- 🎯 **Enum detection** - автоматическое определение перечислений
- 🔗 **Полиморфизм** - поддержка объектов с несколькими вариантами структуры

## Установка

```bash
go install github.com/yanodincov/json-ai-schema-detector/cmd@latest
```

## Использование

### Анализ JSON файла

```bash
# Базовый анализ
json-schema-detector analyze examples/sample_data.json

# Анализ с указанием выходного файла
json-schema-detector analyze examples/sample_data.json -o user_schema.json

# Анализ с конфигурацией
json-schema-detector analyze examples/sample_data.json -c config.yaml
```

### Обновление схемы

```bash
# Обновление существующей схемы новыми данными
json-schema-detector update user_schema.json -i new_data.json
```

### Валидация данных

```bash
# Базовая валидация
json-schema-detector validate data.json user_schema.json

# Подробная валидация
json-schema-detector validate data.json user_schema.json -v

# Строгая валидация
json-schema-detector validate data.json user_schema.json -s
```

## Конфигурация

Создайте файл `config.yaml`:

```yaml
# Порог для определения enum типов
enum_threshold: 10

# Формат выходного файла
output_format: "json-schema"

# Директория для сохранения схем
schemas_directory: "schemas"

# Сохранять комментарии при обновлении
preserve_comments: true

# Автоматическое определение полиморфных объектов
detect_polymorphic: true

# Строгая валидация
strict_validation: false
```

## Пример работы

### Входные данные (sample_data.json)

```json
{
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "role": "admin",
      "permissions": ["read", "write", "delete"],
      "active": true
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "role": "user",
      "active": true
    }
  ]
}
```

### Генерируемая схема

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "data": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "id": {
            "type": "number",
            "description": "Уникальный идентификатор"
          },
          "name": {
            "type": "string",
            "description": "Имя пользователя"
          },
          "role": {
            "type": "string",
            "enum": ["admin", "user"],
            "description": "Роль в системе"
          },
          "permissions": {
            "type": "array",
            "items": {"type": "string"},
            "description": "Права доступа"
          },
          "active": {
            "type": "boolean",
            "description": "Статус активности"
          }
        },
        "required": ["id", "name", "role", "active"]
      }
    }
  },
  "required": ["data"]
}
```

## Сборка из исходников

```bash
git clone https://github.com/yanodincov/json-ai-schema-detector.git
cd json-ai-schema-detector
go mod tidy
go build -o json-schema-detector cmd/main.go
```

## Разработка

### Структура проекта

```
├── cmd/                    # CLI команды
│   ├── main.go            # Точка входа
│   ├── root/              # Корневая команда
│   ├── analyze/           # Команда анализа
│   ├── update/            # Команда обновления
│   └── validate/          # Команда валидации
├── pkg/                   # Основные пакеты
│   ├── types/             # Типы данных
│   ├── analyzer/          # Анализатор JSON
│   └── validator/         # Валидатор схем
├── examples/              # Примеры данных
└── schemas/               # Сгенерированные схемы
```

### Запуск тестов

```bash
go test ./...
```

## Лицензия

MIT License 