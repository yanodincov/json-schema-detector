# JSON AI Schema Detector

Инструмент для автоматического анализа JSON документов и генерации структурированных схем с поддержкой JSON Schema стандарта.

## Возможности

- 🔍 **Автоматический анализ типов данных** - определение примитивных и составных типов
- 📋 **Генерация JSON Schema** - создание стандартных схем JSON Schema
- 🔄 **Обновление схем** - слияние новых данных с существующими схемами
- ✅ **Валидация** - проверка JSON данных против схем
- 📊 **Статистика** - подробная аналитика по структурам данных
- 🎯 **Enum типы** - интерактивное преобразование полей в enum с выбором значений
- 🔗 **Полиморфные типы** - создание полиморфных объектов с oneOf/anyOf
- 🛠️ **Интерактивное управление полями** - изменение типов и описаний через команды
- 📍 **JSON Path навигация** - точная адресация полей в сложных схемах

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

### Интерактивное управление полями

```bash
# Просмотр всех полей в схеме
json-schema-detector list-fields user_schema.json

# Просмотр полей с типами
json-schema-detector list-fields user_schema.json --types

# Подробный просмотр полей
json-schema-detector list-fields user_schema.json --verbose

# Преобразование поля в enum тип
json-schema-detector update-field user_schema.json "data.0.role" enum

# Создание полиморфного типа
json-schema-detector update-field user_schema.json "data.0.user" polymorph

# Обновление описания поля
json-schema-detector update-field user_schema.json "data.0.id" description

# Интерактивный режим (выбор операции)
json-schema-detector update-field user_schema.json "data.0.status"
```

### JSON Path навигация

Для работы с полями в сложных схемах используется JSON Path синтаксис:

```bash
# Простые поля
data.name           # поле name в объекте data
data.id             # поле id в объекте data

# Массивы
data.0.name         # поле name в первом элементе массива data
users.0.profile.age # поле age в профиле первого пользователя

# Вложенные объекты
user.profile.settings.theme    # глубоко вложенное поле
config.database.connection.host # поле в конфигурации

# Примеры команд
json-schema-detector list-fields schema.json
json-schema-detector update-field schema.json "data.0.role" enum
json-schema-detector update-field schema.json "users.0.profile.type" polymorph
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

## Примеры работы

### Интерактивное управление полями

```bash
# Просмотр всех полей в схеме
$ json-schema-detector list-fields examples/sample_data.schema.json

🔍 Поля в схеме: examples/sample_data.schema.json
├── data (array)
│   ├── 0.active (boolean)
│   ├── 0.created_at (string)
│   ├── 0.id (number)
│   ├── 0.name (string)
│   ├── 0.role (string - enum: admin, user, manager)
│   └── 0.permissions (array)

# Преобразование поля в enum тип
$ json-schema-detector update-field examples/sample_data.schema.json "data.0.role" enum

🔧 Обновление поля в схеме
📄 Файл схемы: examples/sample_data.schema.json
🎯 Путь к полю: data.0.role
🔄 Операция: enum

📝 Введите возможные значения для enum (по одному на строку):
💡 Закончите ввод пустой строкой

Значение: admin
Значение: user
Значение: manager
Значение: 

📝 Описание поля (опционально): Роль пользователя в системе
✅ Поле преобразовано в enum с 3 значениями
🎯 Значения: [admin user manager]
✅ Поле успешно обновлено: data.0.role
```

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
│   ├── validate/          # Команда валидации
│   ├── update-field/      # Интерактивное управление полями
│   └── list-fields/       # Просмотр полей схемы
├── pkg/                   # Основные пакеты
│   ├── types/             # Типы данных
│   ├── analyzer/          # Анализатор JSON
│   ├── validator/         # Валидатор схем
│   └── fieldmanager/      # Менеджер полей схемы
├── examples/              # Примеры данных
└── schemas/               # Сгенерированные схемы
```

## Планы развития

### В разработке
- 🔄 **Полиморфные типы** - создание oneOf/anyOf схем для разных вариантов объектов
- 🧪 **Расширенное тестирование** - автоматические тесты для всех компонентов
- 📈 **Статистика использования** - аналитика по полям и типам

### Планируется
- 🌐 **Web интерфейс** - графический интерфейс для управления схемами
- 🔌 **API интерфейс** - REST API для интеграции с другими системами
- 📊 **Расширенная аналитика** - детальные отчеты по структурам данных
- 🎨 **Кастомизация схем** - темы и настройки вывода
- 🔍 **Поиск по схемам** - быстрый поиск полей и типов

### Запуск тестов

```bash
go test ./...
```

## Лицензия

MIT License 