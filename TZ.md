# Техническое задание: JSON AI Schema Detector

## Общее описание

Проект предназначен для автоматического анализа JSON документов и извлечения из них структурированной схемы с возможностью интеллектуального обучения и накопления знаний о структурах данных.

## Основные требования

### Входные данные
- **Формат**: JSON документ
- **Типичная структура**: `{"data": [{"a": 1}, {"a": 2}, {"b": "string"}]}`
- **Поддержка**: Массивы объектов одного типа, вложенные структуры, смешанные типы

### Выходные данные
- **Формат**: Стандартные схемы (JSON Schema, OpenAPI, Proto, TypeScript)
- **Структура**: Схема с типизацией и метаданными
- **Персистентность**: Сохранение и обновление схемы между запусками

## Поддерживаемые форматы схем

### JSON Schema (Рекомендуемый)
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
          }
        }
      }
    }
  },
  "x-analysis-meta": {
    "enum-values": {
      "role": ["admin", "user"],
      "discovered-at": "2024-01-01"
    }
  }
}
```

### OpenAPI 3.0 Components
```yaml
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: number
          description: Уникальный идентификатор
        name:
          type: string
          description: Имя пользователя
        role:
          type: string
          enum: [admin, user]
          description: Роль в системе
      x-analysis-meta:
        enum-values:
          role: [admin, user]
          discovered-at: "2024-01-01"
```

### Protocol Buffers
```protobuf
syntax = "proto3";

// Generated from JSON analysis
message User {
  int32 id = 1;        // Уникальный идентификатор
  string name = 2;     // Имя пользователя
  Role role = 3;       // Роль в системе
}

enum Role {
  ROLE_UNSPECIFIED = 0;
  ROLE_ADMIN = 1;
  ROLE_USER = 2;
}

message UserList {
  repeated User data = 1;
}
```

### TypeScript Definitions
```typescript
/**
 * Generated from JSON analysis
 */
export interface User {
  /** Уникальный идентификатор */
  id: number;
  /** Имя пользователя */
  name: string;
  /** Роль в системе */
  role: "admin" | "user";
  /** Права доступа для админов */
  permissions?: string[];
  /** Статус активности для пользователей */
  active?: boolean;
}

export interface UserResponse {
  data: User[];
}
```

## Функциональные требования

### 1. Анализ типов данных
- **Примитивные типы**: `string`, `number`, `boolean`, `null`
- **Составные типы**: `array`, `object`
- **Специальные типы**: `enum` (для ограниченного набора значений)
- **Автоматическое определение**: Анализ значений для выявления типа

### 2. Система комментариев
- **Автоматические комментарии**: Тип поля определяется автоматически
- **Ручные комментарии**: Пользователь может добавить описание поля
- **Персистентность**: Комментарии сохраняются при повторном анализе
- **Обновление**: Новые поля получают только тип, существующие сохраняют комментарии

### 3. Обработка Enum типов
- **Автоматическое обнаружение**: Поля с ограниченным набором значений
- **Накопление значений**: Новые значения enum добавляются в комментарии
- **Порог определения**: Настраиваемый параметр для определения enum

### 4. Работа с полиморфными объектами
Объекты могут иметь несколько вариантов структуры. Нужно предусмотреть различные подходы к их хранению.

## Работа с полиморфными объектами

### JSON Schema - oneOf/anyOf
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "user": {
      "oneOf": [
        {
          "type": "object",
          "properties": {
            "type": {"const": "admin"},
            "common_field": {"type": "string"},
            "permissions": {"type": "array", "items": {"type": "string"}},
            "admin_level": {"type": "number"}
          },
          "required": ["type", "common_field", "permissions", "admin_level"]
        },
        {
          "type": "object", 
          "properties": {
            "type": {"const": "regular_user"},
            "common_field": {"type": "string"},
            "subscription_type": {"type": "string", "enum": ["free", "premium", "enterprise"]}
          },
          "required": ["type", "common_field", "subscription_type"]
        }
      ]
    }
  }
}
```

### OpenAPI - discriminator
```yaml
components:
  schemas:
    User:
      discriminator:
        propertyName: type
        mapping:
          admin: '#/components/schemas/AdminUser'
          regular_user: '#/components/schemas/RegularUser'
      oneOf:
        - $ref: '#/components/schemas/AdminUser'
        - $ref: '#/components/schemas/RegularUser'
    
    AdminUser:
      type: object
      properties:
        type:
          type: string
          enum: [admin]
        common_field:
          type: string
        permissions:
          type: array
          items:
            type: string
        admin_level:
          type: number
      required: [type, common_field, permissions, admin_level]
    
    RegularUser:
      type: object
      properties:
        type:
          type: string
          enum: [regular_user]
        common_field:
          type: string
        subscription_type:
          type: string
          enum: [free, premium, enterprise]
      required: [type, common_field, subscription_type]
```

### Protocol Buffers - oneof
```protobuf
syntax = "proto3";

message User {
  string common_field = 1;
  
  oneof user_type {
    AdminUser admin = 2;
    RegularUser regular_user = 3;
  }
}

message AdminUser {
  repeated string permissions = 1;
  int32 admin_level = 2;
}

message RegularUser {
  SubscriptionType subscription_type = 1;
}

enum SubscriptionType {
  SUBSCRIPTION_TYPE_UNSPECIFIED = 0;
  SUBSCRIPTION_TYPE_FREE = 1;
  SUBSCRIPTION_TYPE_PREMIUM = 2;
  SUBSCRIPTION_TYPE_ENTERPRISE = 3;
}
```

### TypeScript - Union Types
```typescript
interface BaseUser {
  common_field: string;
}

interface AdminUser extends BaseUser {
  type: "admin";
  permissions: string[];
  admin_level: number;
}

interface RegularUser extends BaseUser {
  type: "regular_user";
  subscription_type: "free" | "premium" | "enterprise";
}

type User = AdminUser | RegularUser;

// Type guards для работы с union types
function isAdminUser(user: User): user is AdminUser {
  return user.type === "admin";
}

function isRegularUser(user: User): user is RegularUser {
  return user.type === "regular_user";
}
```

## Технические требования

### Архитектура
- **Язык**: Go
- **Парсер**: Стандартный `encoding/json`
- **Хранение**: Файловая система (JSONC файлы)
- **Конфигурация**: YAML/JSON файл настроек

### Модули
1. **Parser** - анализ входного JSON
2. **Schema Analyzer** - определение типов и структур
3. **Metadata Manager** - управление описаниями и метаданными
4. **Format Converter** - конвертация между форматами схем
5. **Persistence** - сохранение и загрузка схем
6. **Merger** - слияние новых данных с существующими схемами
7. **Validator** - валидация JSON против схем

### Настройки
- Порог определения enum (по умолчанию: 10 различных значений)
- Стратегия обработки полиморфных объектов
- Путь к файлам схем
- Правила именования полей

## Интерфейс взаимодействия

### CLI команды
```bash
# Анализ JSON файла с выводом в JSON Schema (по умолчанию)
json-schema-detector analyze input.json -o schema.json

# Анализ с выбором формата выходной схемы
json-schema-detector analyze input.json -o schema.yaml --format openapi
json-schema-detector analyze input.json -o schema.proto --format protobuf  
json-schema-detector analyze input.json -o types.ts --format typescript

# Обновление существующей схемы
json-schema-detector update schema.json -i new_data.json

# Валидация JSON против схемы
json-schema-detector validate data.json schema.json

# Конвертация между форматами
json-schema-detector convert schema.json --from json-schema --to openapi -o api.yaml
json-schema-detector convert schema.json --from json-schema --to protobuf -o schema.proto
```

### Файловая структура
```
schemas/
├── user.json           # JSON Schema пользователя
├── user.yaml           # OpenAPI схема пользователя
├── user.proto          # Protocol Buffers схема
├── user.ts             # TypeScript определения
├── product.json        # JSON Schema продукта
└── config.yaml         # Конфигурация
```

## Примеры использования

### Входной JSON
```json
{
  "data": [
    {"id": 1, "name": "John", "role": "admin", "permissions": ["read", "write"]},
    {"id": 2, "name": "Jane", "role": "user", "active": true},
    {"id": 3, "name": "Bob", "role": "admin", "permissions": ["read"]}
  ]
}
```

### Выходная схема (JSON Schema)
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
            "description": "Права доступа для админов"
          },
          "active": {
            "type": "boolean",
            "description": "Статус активности для пользователей"
          }
        },
        "required": ["id", "name", "role"]
      }
    }
  },
  "required": ["data"],
  "x-analysis-meta": {
    "enum-values": {
      "role": ["admin", "user"]
    },
    "optional-fields": ["permissions", "active"],
    "polymorphic-patterns": {
      "admin": ["id", "name", "role", "permissions"],
      "user": ["id", "name", "role", "active"]
    }
  }
}
```

## Дополнительные возможности

### Фазы развития
1. **MVP**: Базовый анализ типов и сохранение в JSONC
2. **Enum detection**: Автоматическое определение перечислений
3. **Polymorphic objects**: Поддержка полиморфных объектов
4. **Advanced analytics**: Статистика, валидация, миграции схем

### Возможные расширения
- Интеграция с системами валидации
- Генерация Go структур из схем
- Web интерфейс для просмотра схем
- Автоматическое создание документации
- Интеграция с IDE (VS Code, IntelliJ)
- Генерация тестовых данных на основе схем

## Сравнение форматов схем

| Критерий | JSON Schema | OpenAPI | Protocol Buffers | TypeScript |
|----------|-------------|---------|------------------|------------|
| **Стандартизация** | ✅ RFC стандарт | ✅ OpenAPI Spec | ✅ Google стандарт | ✅ Microsoft стандарт |
| **Валидация JSON** | ✅ Нативная | ✅ Через JSON Schema | ❌ Только бинарные данные | ❌ Только типы |
| **Документация** | ✅ description поля | ✅ Богатая документация | ✅ Комментарии | ✅ JSDoc комментарии |
| **Полиморфизм** | ✅ oneOf/anyOf | ✅ discriminator | ✅ oneof | ✅ Union types |
| **Тулинг** | ✅ Широкая поддержка | ✅ Swagger/OpenAPI | ✅ protoc, buf | ✅ TypeScript compiler |
| **Читаемость** | ✅ Хорошая | ✅ Отличная | ✅ Хорошая | ✅ Отличная |
| **Размер файла** | ✅ Компактный | ⚠️ Объемный | ✅ Компактный | ✅ Компактный |
| **Метаданные** | ✅ x-* extensions | ✅ Нативные | ✅ Опции | ✅ Декораторы |

## Рекомендации по выбору формата

### JSON Schema - По умолчанию
- **Когда использовать**: Для большинства случаев
- **Преимущества**: Стандарт валидации JSON, широкая поддержка
- **Недостатки**: Менее читаемый чем OpenAPI

### OpenAPI - Для API документации
- **Когда использовать**: Когда нужна подробная документация API
- **Преимущества**: Отличная документация, Swagger UI
- **Недостатки**: Избыточность для простых схем

### Protocol Buffers - Для производительности
- **Когда использовать**: Для gRPC сервисов и бинарной сериализации
- **Преимущества**: Высокая производительность, строгая типизация
- **Недостатки**: Не подходит для JSON валидации

### TypeScript - Для фронтенда
- **Когда использовать**: Для генерации типов для фронтенд приложений
- **Преимущества**: Нативная поддержка в TypeScript/JavaScript
- **Недостатки**: Не валидирует данные в runtime 