# Интеграция логирования в Go ⚡ **ТОЛЬКО ДЛЯ GO**

**НАЗНАЧЕНИЕ**: **ПРАВИЛА ТОЛЬКО ДЛЯ GO КОДА** - интеграция ZAP логгера с обязательным маскированием чувствительных данных.

---

## 1. Политика интеграции логгера 📝

### Основной логгер
**Обязательное использование ZAP Logger** - единственный разрешенный логгер в проекте.

### Интеграция в структуры
```go
// Обязательное поле логгера в структуре
type Service struct {
    logger *zap.Logger  // Обязательное поле с именем "logger"
    repo   Repository
}

// Инъекция через конструктор - включай во все конструкторы
func NewService(logger *zap.Logger, repo Repository) *Service {
    return &Service{
        logger: logger,
        repo:   repo,
    }
}
```

**Правила интеграции:**
- **Поле в структуре**: Обязательное поле `*zap.Logger` в каждой структуре
- **Инъекция через конструктор**: Логгер передается во все конструкторы  
- **Имя поля**: `logger` (единообразное именование)

### Статические методы
```go
func ProcessPayment(ctx context.Context) error {
    logger := logger.NewEntry().Zap()
    logger.Info(ctx, "processing payment started")
    
    return nil
}
```

### Методы контекстного логирования
```go
// Доступные методы с контекстом
s.logger.Info(ctx, "message", fields...)
s.logger.Warn(ctx, "message", fields...)
s.logger.Error(ctx, "message", fields...)
s.logger.Debug(ctx, "message", fields...)
```

### Справочник типов ZAP полей
```go
// Строки
zap.String("key", val)
zap.ByteString("key", []byte)

// Числа
zap.Int("key", val)
zap.Int64("key", val)
zap.Uint64("key", val)
zap.Float64("key", val)
zap.Float32("key", val)

// Временные и логические данные
zap.Bool("key", val)
zap.Time("key", val)
zap.Duration("key", val)

// Специальные типы
zap.Error(err)
zap.Any("key", val)

// Массивы
zap.Strings("key", []string{})
zap.Ints("key", []int{})
```

**Приоритет**: Используй типизированные поля вместо `zap.Any` всегда когда возможно.

---

## 2. Политика логирования ошибок 🚨

### Выбор уровня логирования
- **Неожиданные ошибки**: ERROR уровень
- **Бизнес-логические ошибки**: WARN уровень

### Моменты логирования ошибок
- **При создании новой ошибки**: Логируй при инициализации ошибки
- **При получении ошибки**: Логируй при получении от внешнего сервиса

### Содержание лога ошибки
```go
if err != nil {
    s.logger.Error(ctx, "payment processing failed", 
        zap.Error(err),                    // Включай объект ошибки
        zap.String("payment_id", id),      // Контекст
    )
    
    return errfmt.Error(err, "payment processing failed")
}
```

**Обязательно**: Включай объект ошибки в лог-запись для полного контекста.

---

## 3. Политика логирования успешных операций ✅

### Обязательное логирование успеха
**INFO уровень обязателен** для успешного завершения операций.

### Области покрытия
- **Методы сервисов**: Логируй успешное завершение
- **Клиентские операции**: Логируй успешные операции
- **Внешние запросы**: Логируй успешные вызовы внешних API
- **Обработка запросов**: Логируй успешную обработку

```go
s.logger.Info(ctx, "payment processed successfully", 
    zap.String("payment_id", payment.ID),
    zap.String("status", payment.Status),
)
```

---

## 4. Политика включения данных в логи 📊

### Вспомогательные данные
**Включай полезные контекстные данные** для понимания операций.

### Определение чувствительных данных

#### 🚫 **Запрещенные пользовательские данные:**
- **Личная информация**: имена, адреса, даты рождения, документы, карты, счета, телефоны, email
- **Биометрические данные**: отпечатки, геолокация, секреты

#### 🚫 **Запрещенные сырые данные:**
- **JSON строки**: json_strings (потенциально содержат чувствительную информацию)
- **XML строки**: xml_strings  
- **Base64 строки**: base64_strings
- **Непарсированные данные**: unparsed_data

#### ✅ **Разрешенные технические данные:**
- **Системные идентификаторы**: system_ids, external_ids
- **Статусы**: statuses, operation_types
- **Финансовые общие данные**: currencies, amounts, timestamps
- **Версии**: api_versions

### Ограничения
- **Максимум 5 свойств на лог** - не перегружай логи
- **Структуры с чувствительными данными**: Используй только с методом `Mask()`

---

## 5. Политика маскирования данных 🔒

### Приоритет использования методов Mask
```go
// ✅ Правильно - используй весь объект с методом Mask
s.logger.Info(ctx, "user data processed", zap.Any("user", user.Mask()))

// ❌ Неправильно - не маскируй отдельные поля если есть метод Mask
s.logger.Info(ctx, "user data processed",
    zap.String("user_id", user.ID),
    zap.String("email", mask.Email(user.Email))) // НЕ ДЕЛАЙ ТАК
```

### Правила генерации метода Mask

#### Для внутренних структур - добавь метод Mask
```go
type User struct {
    ID    string
    Email string
    Phone string
}

func (u User) Mask() User {
    return User{
        ID:    u.ID,              // ID не маскируется
        Email: mask.Email(u.Email),
        Phone: mask.Phone(u.Phone),
    }
}
```

#### Для внешних структур - создай статическую функцию маскирования
```go
func maskExternalUser(user external.User) map[string]interface{} {
    return map[string]interface{}{
        "id":    user.ID,
        "email": mask.Email(user.Email),
        "phone": mask.Phone(user.Phone),
    }
}

s.logger.Info(ctx, "external user processed", 
    zap.Any("user", maskExternalUser(externalUser)))
```

#### Для вложенных структур - вызывай Mask для вложенных структур
```go
func (o Order) Mask() Order {
    return Order{
        ID:   o.ID,
        User: o.User.Mask(),  // Вызови Mask для вложенной структуры
        Items: o.Items,
    }
}
```

**Правило возврата**: Метод Mask должен возвращать тот же тип, что и получатель.

### Доступные функции маскирования
```go
// Полное маскирование
mask.Full(val)     // Полностью скрывает значение
mask.Empty(val)    // Заменяет пустой строкой
mask.Fixed(val)    // Фиксированная замена

// Частичное маскирование
mask.Holder(val)   // Показывает только держателя
mask.Email(val)    // Маскирует email (us***@example.com)
mask.Phone(val)    // Маскирует телефон (+7***1234)

// Финансовые данные
mask.Pan(val)      // Маскирует номер карты
mask.Cvv(val)      // Маскирует CVV
mask.Secret(val)   // Маскирует секретные данные

// Компоненты дат
mask.Month(val)    // Маскирует месяц
mask.Year(val)     // Маскирует год
```

---

## 6. Примеры использования 💡

### Структура с логгером
```go
type Service struct {
    logger *zap.Logger
    repo   Repository
}

func NewService(logger *zap.Logger, repo Repository) *Service {
    return &Service{logger: logger, repo: repo}
}
```

### Логирование в статических методах
```go
func ProcessPayment(ctx context.Context) error {
    logger := logger.NewEntry().Zap()
    logger.Info(ctx, "processing payment started")
    
    return nil
}
```

### Логирование ошибок
```go
if err != nil {
    s.logger.Error(ctx, "payment processing failed", 
        zap.Error(err), 
        zap.String("payment_id", id))
    
    return errfmt.Error(err, "payment processing failed")
}
```

### Логирование успешных операций
```go
s.logger.Info(ctx, "payment processed successfully", 
    zap.String("payment_id", payment.ID),
    zap.String("status", payment.Status))
```

### Правильное маскированное логирование
```go
s.logger.Info(ctx, "user data processed", zap.Any("user", user.Mask()))
```

### Неправильное маскирование (НЕ ДЕЛАЙ)
```go
// НИКОГДА НЕ ДЕЛАЙ ТАК ЕСЛИ ЕСТЬ МЕТОД MASK:
s.logger.Info(ctx, "user data processed",
    zap.String("user_id", user.ID),
    zap.String("email", mask.Email(user.Email)))
```

### Маскирование внешних структур
```go
func maskExternalUser(user external.User) map[string]interface{} {
    return map[string]interface{}{
        "id":    user.ID,
        "email": mask.Email(user.Email),
    }
}

s.logger.Info(ctx, "external user processed", 
    zap.Any("user", maskExternalUser(externalUser)))
```

---

## 📋 Чек-лист для внедрения

### При создании новой структуры:
- [ ] Добавил поле `logger *zap.Logger`
- [ ] Включил логгер в конструктор
- [ ] Добавил метод `Mask()` если структура содержит чувствительные данные

### При логировании:
- [ ] Использую типизированные ZAP поля вместо `zap.Any`
- [ ] Логирую ошибки на правильном уровне (ERROR/WARN)
- [ ] Логирую успешные операции на INFO уровне
- [ ] Не превышаю 5 свойств на лог
- [ ] Маскирую чувствительные данные
- [ ] Использую весь объект с `Mask()` вместо маскирования отдельных полей

### При работе с внешними структурами:
- [ ] Создал статическую функцию маскирования
- [ ] Не логирую сырые JSON/XML/Base64 данные
- [ ] Проверил что не логирую личную информацию