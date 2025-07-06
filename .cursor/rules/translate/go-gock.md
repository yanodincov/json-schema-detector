# HTTP мокирование с gock ⚡ **ТОЛЬКО ДЛЯ GO**

**ПРИМЕНЕНИЕ**: Эти правила применяются исключительно для Go кода и Go проектов с библиотекой gock.

## Цель 🎯
**Детерминированное тестирование HTTP клиентов** - замена реальных HTTP запросов на предсказуемые моки.

---

## Инициализация и очистка 🔧

### Настройка в тестах
```go
import (
    "net/http"
    "github.com/h2non/gock"
)

func setUp(t *testing.T) (*fixture, func()) {
    // Инициализация gock
    gock.InterceptClient(http.DefaultClient)
    
    // Создание fixture...
    f := &fixture{
        client: NewHTTPClient(), // Ваш HTTP клиент
    }
    
    cleanup := func() {
        // Проверка что все моки были использованы
        assert.True(t, gock.IsDone())
        
        // Очистка gock
        gock.Off()
        gock.Clean()
        gock.RestoreClient(http.DefaultClient)
    }
    
    return f, cleanup
}
```

### Перехват HTTP клиента
```go
// Если используете кастомный HTTP клиент
client := &http.Client{
    Timeout: 30 * time.Second,
}

// Включаем перехват для кастомного клиента
gock.InterceptClient(client)
```

---

## Определение моков 📝

### Базовая структура мока
```go
func TestAPICall(t *testing.T) {
    t.Parallel()
    f, cleanup := setUp(t)
    defer cleanup()
    
    // Настройка мока
    gock.New("https://api.example.com").
        Get("/users/123").
        Reply(200).
        JSON(map[string]interface{}{
            "id":   123,
            "name": "John Doe",
        })
    
    // Выполнение тестируемого кода
    user, err := f.client.GetUser("123")
    
    // Проверки
    require.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
}
```

### Мок с заголовками запроса
```go
gock.New("https://api.example.com").
    Post("/payments").
    MatchHeader("Authorization", "Bearer token123").
    MatchHeader("Content-Type", "application/json").
    JSON(map[string]interface{}{
        "amount": 1000,
        "currency": "USD",
    }).
    Reply(201).
    JSON(map[string]interface{}{
        "id": "payment_123",
        "status": "success",
    })
```

### Мок с query параметрами
```go
gock.New("https://api.example.com").
    Get("/users").
    MatchParam("page", "1").
    MatchParam("limit", "10").
    Reply(200).
    JSON(map[string]interface{}{
        "users": []map[string]interface{}{
            {"id": 1, "name": "User 1"},
            {"id": 2, "name": "User 2"},
        },
        "total": 100,
    })
```

---

## Типы ответов 📊

### JSON ответы
```go
// Простой JSON объект
gock.New("https://api.example.com").
    Get("/user/123").
    Reply(200).
    JSON(map[string]interface{}{
        "id": 123,
        "email": "user@example.com",
    })

// JSON массив
gock.New("https://api.example.com").
    Get("/users").
    Reply(200).
    JSON([]map[string]interface{}{
        {"id": 1, "name": "User 1"},
        {"id": 2, "name": "User 2"},
    })

// Использование структуры
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

user := User{ID: 123, Name: "John"}
gock.New("https://api.example.com").
    Get("/user/123").
    Reply(200).
    JSON(user)
```

### Текстовые ответы
```go
// Обычный текст
gock.New("https://api.example.com").
    Get("/health").
    Reply(200).
    BodyString("OK")

// XML ответ
gock.New("https://api.example.com").
    Get("/data.xml").
    Reply(200).
    Body(strings.NewReader(`
        <?xml version="1.0"?>
        <data>
            <id>123</id>
            <name>Test</name>
        </data>
    `))
```

### Ответы с кастомными заголовками
```go
gock.New("https://api.example.com").
    Get("/users").
    Reply(200).
    SetHeader("X-Rate-Limit", "100").
    SetHeader("X-Total-Count", "1000").
    JSON(map[string]interface{}{
        "users": []string{"user1", "user2"},
    })
```

---

## Обработка ошибок 🚨

### HTTP ошибки
```go
// 404 Not Found
gock.New("https://api.example.com").
    Get("/users/999").
    Reply(404).
    JSON(map[string]interface{}{
        "error": "User not found",
        "code": "USER_NOT_FOUND",
    })

// 500 Internal Server Error
gock.New("https://api.example.com").
    Post("/payments").
    Reply(500).
    JSON(map[string]interface{}{
        "error": "Internal server error",
        "message": "Database connection failed",
    })

// 429 Rate Limit
gock.New("https://api.example.com").
    Get("/api/data").
    Reply(429).
    SetHeader("Retry-After", "60").
    JSON(map[string]interface{}{
        "error": "Rate limit exceeded",
    })
```

### Сетевые ошибки
```go
// Таймаут соединения
gock.New("https://api.example.com").
    Get("/slow-endpoint").
    ReplyError(errors.New("context deadline exceeded"))

// Ошибка соединения
gock.New("https://api.example.com").
    Post("/data").
    ReplyError(errors.New("connection refused"))
```

---

## Примеры использования 💡

### Тестирование успешного сценария
```go
func TestPaymentService_CreatePayment_Success(t *testing.T) {
    t.Parallel()
    f, cleanup := setUp(t)
    defer cleanup()
    
    // Мок успешного ответа от платежного API
    gock.New("https://payment-api.com").
        Post("/v1/payments").
        MatchHeader("Authorization", "Bearer test-token").
        JSON(map[string]interface{}{
            "amount": 1000,
            "currency": "USD",
            "card_token": "card_123",
        }).
        Reply(201).
        JSON(map[string]interface{}{
            "id": "payment_456",
            "status": "completed",
            "amount": 1000,
            "currency": "USD",
            "created_at": "2023-01-01T00:00:00Z",
        })
    
    // Выполнение тестируемого кода
    payment, err := f.paymentService.CreatePayment(f.ctx, CreatePaymentRequest{
        Amount:    1000,
        Currency:  "USD",
        CardToken: "card_123",
    })
    
    // Проверки
    require.NoError(t, err)
    assert.Equal(t, "payment_456", payment.ID)
    assert.Equal(t, "completed", payment.Status)
    assert.Equal(t, 1000, payment.Amount)
}
```

### Тестирование ошибок от внешнего API
```go
func TestPaymentService_CreatePayment_APIError(t *testing.T) {
    t.Parallel()
    f, cleanup := setUp(t)
    defer cleanup()
    
    // Мок ошибки от платежного API
    gock.New("https://payment-api.com").
        Post("/v1/payments").
        Reply(422).
        JSON(map[string]interface{}{
            "error": "invalid_card",
            "message": "Card number is invalid",
            "code": "INVALID_CARD_NUMBER",
        })
    
    // Выполнение тестируемого кода
    payment, err := f.paymentService.CreatePayment(f.ctx, CreatePaymentRequest{
        Amount:    1000,
        Currency:  "USD",
        CardToken: "invalid_card",
    })
    
    // Проверки
    require.Error(t, err)
    assert.Nil(t, payment)
    assert.Contains(t, err.Error(), "invalid_card")
}
```

### Тестирование множественных API вызовов
```go
func TestUserService_GetUserWithPayments_Success(t *testing.T) {
    t.Parallel()
    f, cleanup := setUp(t)
    defer cleanup()
    
    userID := "user_123"
    
    // Мок для получения пользователя
    gock.New("https://user-api.com").
        Get("/v1/users/" + userID).
        Reply(200).
        JSON(map[string]interface{}{
            "id": userID,
            "name": "John Doe",
            "email": "john@example.com",
        })
    
    // Мок для получения платежей пользователя
    gock.New("https://payment-api.com").
        Get("/v1/payments").
        MatchParam("user_id", userID).
        Reply(200).
        JSON(map[string]interface{}{
            "payments": []map[string]interface{}{
                {
                    "id": "payment_1",
                    "amount": 1000,
                    "status": "completed",
                },
                {
                    "id": "payment_2", 
                    "amount": 2000,
                    "status": "pending",
                },
            },
        })
    
    // Выполнение тестируемого кода
    result, err := f.userService.GetUserWithPayments(f.ctx, userID)
    
    // Проверки
    require.NoError(t, err)
    assert.Equal(t, "John Doe", result.User.Name)
    assert.Len(t, result.Payments, 2)
    assert.Equal(t, "payment_1", result.Payments[0].ID)
}
```

---

## Продвинутые возможности 🎛️

### Мокирование с условиями
```go
// Разные ответы для разных условий
gock.New("https://api.example.com").
    Get("/users").
    MatchParam("role", "admin").
    Reply(200).
    JSON([]map[string]interface{}{
        {"id": 1, "name": "Admin User", "role": "admin"},
    })

gock.New("https://api.example.com").
    Get("/users").
    MatchParam("role", "user").
    Reply(200).
    JSON([]map[string]interface{}{
        {"id": 2, "name": "Regular User", "role": "user"},
    })
```

### Проверка тела запроса
```go
gock.New("https://api.example.com").
    Post("/users").
    MatchType("json").
    JSON(map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
    }).
    Reply(201).
    JSON(map[string]interface{}{
        "id": 123,
        "name": "John Doe",
    })
```

### Множественные использования мока
```go
// Мок будет срабатывать несколько раз
gock.New("https://api.example.com").
    Get("/health").
    Times(3). // Сработает 3 раза
    Reply(200).
    BodyString("OK")

// Бесконечное количество раз
gock.New("https://api.example.com").
    Get("/health").
    Persist(). // Будет срабатывать всегда
    Reply(200).
    BodyString("OK")
```

---

## Лучшие практики ✨

### Всегда проверяй что моки использованы
```go
func TestExample(t *testing.T) {
    // Настройка моков...
    
    // Выполнение тестируемого кода...
    
    // ОБЯЗАТЕЛЬНО: Проверка что все моки были использованы
    assert.True(t, gock.IsDone(), "Not all HTTP mocks were used")
}
```

### Очистка после каждого теста
```go
func TestExample(t *testing.T) {
    // Настройка gock
    gock.InterceptClient(http.DefaultClient)
    
    defer func() {
        gock.Off()          // Отключить перехват
        gock.Clean()        // Очистить все моки
        gock.RestoreClient(http.DefaultClient) // Восстановить клиент
    }()
    
    // Тест логика...
}
```

### Использование в table-driven тестах
```go
func TestAPIClient_GetUser(t *testing.T) {
    tests := []struct {
        name           string
        userID         string
        mockSetup      func()
        expectedUser   *User
        expectedError  string
    }{
        {
            name:   "success",
            userID: "123",
            mockSetup: func() {
                gock.New("https://api.example.com").
                    Get("/users/123").
                    Reply(200).
                    JSON(map[string]interface{}{
                        "id": 123,
                        "name": "John",
                    })
            },
            expectedUser: &User{ID: 123, Name: "John"},
        },
        {
            name:   "not found",
            userID: "999",
            mockSetup: func() {
                gock.New("https://api.example.com").
                    Get("/users/999").
                    Reply(404).
                    JSON(map[string]interface{}{
                        "error": "Not found",
                    })
            },
            expectedError: "user not found",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            // Очистка для каждого подтеста
            gock.Clean()
            gock.InterceptClient(http.DefaultClient)
            defer func() {
                gock.Off()
                gock.Clean()
                gock.RestoreClient(http.DefaultClient)
            }()
            
            // Настройка мока
            tt.mockSetup()
            
            // Выполнение теста
            user, err := client.GetUser(tt.userID)
            
            // Проверки
            if tt.expectedError != "" {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.expectedUser, user)
            }
            
            // Проверка использования моков
            assert.True(t, gock.IsDone())
        })
    }
}
```

---

## Контрольный список ✅

### Настройка:
- [ ] Используй `gock.InterceptClient()` для перехвата HTTP клиента
- [ ] Всегда очищай gock в cleanup функции
- [ ] Восстанавливай оригинальный клиент через `gock.RestoreClient()`

### Моки:
- [ ] Настраивай специфичные URL, методы и заголовки
- [ ] Используй реалистичные ответы похожие на настоящие API
- [ ] Покрывай как успешные сценарии, так и ошибки

### Проверки:
- [ ] Всегда проверяй `gock.IsDone()` в конце теста
- [ ] Проверяй что тестируемый код правильно обрабатывает ответы
- [ ] Тестируй различные HTTP статус коды

### Производительность:
- [ ] Используй `t.Parallel()` в тестах с gock
- [ ] Очищай моки между тестами для изоляции
- [ ] Группируй связанные моки логически

gock делает HTTP тестирование предсказуемым и быстрым! 🚀