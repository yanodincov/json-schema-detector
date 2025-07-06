# Правила модульного тестирования Go

## Цель 🎯
- **Проверка корректности бизнес-логики** - убедиться что код работает как задумано
- **Осознанные изменения тестов** - каждое изменение теста должно быть документировано и обосновано

---

## Структура файлов и именование 📁

### Конвенции именования
- **Имя файла**: `<name>_test.go` (например: `get_user.go` → `get_user_test.go`)
- **Пакет**: `filepackage_test` (внешний пакет для изоляции)
- **Функция теста**: `func TestXxx(t *testing.T)`
- **t.Parallel()**: В начале каждого теста и подтеста
- **t.Helper()**: В вспомогательных функциях

### Стили тестирования

#### Язык названий тестов
- **ОБЯЗАТЕЛЬНО АНГЛИЙСКИЙ**: Все названия тест-кейсов только на английском языке
- **Table-driven тесты**: Названия кейсов на английском
- **t.Run кейсы**: Названия на английском  
- **ЗАПРЕЩЕНО**: Русский или любой другой язык кроме английского

```go
// ✅ Правильно - на английском
tests := []struct {
    name     string
    input    string
    expected string
}{
    {
        name:     "success case",
        input:    "test",
        expected: "result",
    },
    {
        name:     "empty input",
        input:    "",
        expected: "",
    },
}

// ❌ Неправильно - на русском
tests := []struct {
    name     string
    input    string
    expected string
}{
    {
        name:     "успешный случай",  // НЕ ДЕЛАЙ ТАК!
        input:    "test",
        expected: "result",
    },
}
```

#### Table-driven тесты (предпочтительный стиль)
```go
func TestUserService_GetUser(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name          string
        userID        string
        mockSetup     func(*testing.T, *mock.MockUserRepository)
        expectedUser  *User
        expectedError string
    }{
        {
            name:   "user exists",
            userID: "user123",
            mockSetup: func(t *testing.T, repo *mock.MockUserRepository) {
                user := &User{ID: "user123", Name: "John"}
                repo.EXPECT().GetByID(gomock.Any(), "user123").Return(user, nil)
            },
            expectedUser: &User{ID: "user123", Name: "John"},
        },
        {
            name:   "user not found",
            userID: "nonexistent",
            mockSetup: func(t *testing.T, repo *mock.MockUserRepository) {
                repo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, ErrUserNotFound)
            },
            expectedError: "user not found",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            f, cleanup := setUp(t)
            defer cleanup()
            
            tt.mockSetup(t, f.userRepo)
            
            user, err := f.userService.GetUser(f.ctx, tt.userID)
            
            if tt.expectedError != "" {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.expectedUser, user)
        })
    }
}
```

#### Простые тесты (для единичных случаев)
```go
func TestUserService_CreateUser_Success(t *testing.T) {
    t.Parallel()
    
    f, cleanup := setUp(t)
    defer cleanup()
    
    req := CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    expectedUser := &User{
        ID:    "user123",
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    f.userRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(expectedUser, nil)
    
    user, err := f.userService.CreateUser(f.ctx, req)
    
    require.NoError(t, err)
    assert.Equal(t, expectedUser, user)
}
```

---

## Fixture pattern (обязательный подход) 🏗️

### Базовая структура fixture
```go
type fixture struct {
    ctx         context.Context
    userService *UserService
    userRepo    *mock.MockUserRepository
    logger      *zap.Logger
}

func setUp(t *testing.T) (*fixture, func()) {
    t.Helper()
    
    // Создание контекста с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    
    // Создание логгера
    logger := zap.NewNop()
    
    // Создание моков
    ctrl := gomock.NewController(t)
    userRepo := mock.NewMockUserRepository(ctrl)
    
    // Создание сервиса
    userService := NewUserService(logger, userRepo)
    
    f := &fixture{
        ctx:         ctx,
        userService: userService,
        userRepo:    userRepo,
        logger:      logger,
    }
    
    cleanup := func() {
        cancel()
        ctrl.Finish()
    }
    
    return f, cleanup
}
```

### Fixture с множественными зависимостями
```go
type fixture struct {
    ctx             context.Context
    paymentService  *PaymentService
    
    // Репозитории
    paymentRepo     *mock.MockPaymentRepository
    userRepo        *mock.MockUserRepository
    
    // Внешние клиенты
    stripeClient    *mock.MockStripeClient
    emailClient     *mock.MockEmailClient
    
    // Инфраструктура
    logger          *zap.Logger
    cache           *mock.MockCache
}

func setUp(t *testing.T) (*fixture, func()) {
    t.Helper()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    logger := zap.NewNop()
    
    // Создание всех моков
    ctrl := gomock.NewController(t)
    paymentRepo := mock.NewMockPaymentRepository(ctrl)
    userRepo := mock.NewMockUserRepository(ctrl)
    stripeClient := mock.NewMockStripeClient(ctrl)
    emailClient := mock.NewMockEmailClient(ctrl)
    cache := mock.NewMockCache(ctrl)
    
    // Создание сервиса со всеми зависимостями
    paymentService := NewPaymentService(PaymentServiceConfig{
        Logger:       logger,
        PaymentRepo:  paymentRepo,
        UserRepo:     userRepo,
        StripeClient: stripeClient,
        EmailClient:  emailClient,
        Cache:        cache,
    })
    
    f := &fixture{
        ctx:            ctx,
        paymentService: paymentService,
        paymentRepo:    paymentRepo,
        userRepo:       userRepo,
        stripeClient:   stripeClient,
        emailClient:    emailClient,
        logger:         logger,
        cache:          cache,
    }
    
    cleanup := func() {
        cancel()
        ctrl.Finish()
    }
    
    return f, cleanup
}
```

---

## Создание моков с gomock 🎭

### Генерация интерфейсов
```go
//go:generate mockgen -source=user_repository.go -destination=mock/user_repository.go

type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

### Настройка ожиданий моков
```go
func TestPaymentService_ProcessPayment_Success(t *testing.T) {
    f, cleanup := setUp(t)
    defer cleanup()
    
    payment := &Payment{
        ID:     "payment123",
        Amount: 1000,
        Status: "pending",
    }
    
    // Настройка ожиданий в правильном порядке
    f.paymentRepo.EXPECT().
        GetByID(gomock.Any(), "payment123").
        Return(payment, nil)
    
    f.stripeClient.EXPECT().
        ProcessPayment(gomock.Any(), gomock.Any()).
        Return(&StripeResponse{
            ID:     "stripe_123",
            Status: "succeeded",
        }, nil)
    
    f.paymentRepo.EXPECT().
        Update(gomock.Any(), gomock.Any()).
        Do(func(ctx context.Context, p *Payment) {
            // Проверка изменений в платеже
            assert.Equal(t, "completed", p.Status)
            assert.Equal(t, "stripe_123", p.ExternalID)
        }).
        Return(nil)
    
    f.emailClient.EXPECT().
        SendConfirmation(gomock.Any(), gomock.Any()).
        Return(nil)
    
    // Выполнение теста
    err := f.paymentService.ProcessPayment(f.ctx, "payment123")
    
    require.NoError(t, err)
}
```

### Условные моки с gomock.Any()
```go
func TestUserService_UpdateUser_WithValidation(t *testing.T) {
    f, cleanup := setUp(t)
    defer cleanup()
    
    // Мок с условием
    f.userRepo.EXPECT().
        Update(gomock.Any(), gomock.Any()).
        Do(func(ctx context.Context, user *User) {
            // Проверка что email валидный
            assert.Contains(t, user.Email, "@")
            assert.NotEmpty(t, user.Name)
        }).
        Return(nil)
    
    err := f.userService.UpdateUser(f.ctx, &User{
        ID:    "user123",
        Name:  "John Doe",
        Email: "john@example.com",
    })
    
    require.NoError(t, err)
}
```

---

## Тестирование ошибок ❌

### Проверка различных типов ошибок
```go
func TestUserService_GetUser_Errors(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name          string
        userID        string
        mockSetup     func(*mock.MockUserRepository)
        expectedError string
    }{
        {
            name:   "user not found",
            userID: "nonexistent",
            mockSetup: func(repo *mock.MockUserRepository) {
                repo.EXPECT().
                    GetByID(gomock.Any(), "nonexistent").
                    Return(nil, ErrUserNotFound)
            },
            expectedError: "user not found",
        },
        {
            name:   "database error",
            userID: "user123",
            mockSetup: func(repo *mock.MockUserRepository) {
                repo.EXPECT().
                    GetByID(gomock.Any(), "user123").
                    Return(nil, errors.New("connection failed"))
            },
            expectedError: "connection failed",
        },
        {
            name:   "invalid user id",
            userID: "",
            mockSetup: func(repo *mock.MockUserRepository) {
                // Мок не вызывается для невалидных данных
            },
            expectedError: "invalid user id",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            f, cleanup := setUp(t)
            defer cleanup()
            
            tt.mockSetup(f.userRepo)
            
            user, err := f.userService.GetUser(f.ctx, tt.userID)
            
            require.Error(t, err)
            assert.Nil(t, user)
            assert.Contains(t, err.Error(), tt.expectedError)
        })
    }
}
```

### Тестирование цепочки ошибок
```go
func TestPaymentService_ProcessPayment_FailureChain(t *testing.T) {
    f, cleanup := setUp(t)
    defer cleanup()
    
    payment := &Payment{ID: "payment123", Status: "pending"}
    
    // Успешное получение платежа
    f.paymentRepo.EXPECT().
        GetByID(gomock.Any(), "payment123").
        Return(payment, nil)
    
    // Ошибка при обработке в Stripe
    f.stripeClient.EXPECT().
        ProcessPayment(gomock.Any(), gomock.Any()).
        Return(nil, errors.New("card declined"))
    
    // Обновление статуса на failed
    f.paymentRepo.EXPECT().
        Update(gomock.Any(), gomock.Any()).
        Do(func(ctx context.Context, p *Payment) {
            assert.Equal(t, "failed", p.Status)
        }).
        Return(nil)
    
    // Отправка уведомления об ошибке
    f.emailClient.EXPECT().
        SendFailureNotification(gomock.Any(), gomock.Any()).
        Return(nil)
    
    err := f.paymentService.ProcessPayment(f.ctx, "payment123")
    
    require.Error(t, err)
    assert.Contains(t, err.Error(), "card declined")
}
```

---

## Тестирование с внешними зависимостями 🔗

### HTTP клиенты (с gock)
```go
import "github.com/h2non/gock"

func TestAPIClient_GetUser(t *testing.T) {
    t.Parallel()
    
    // Настройка gock
    defer gock.Off()
    defer gock.Clean()
    gock.InterceptClient(http.DefaultClient)
    
    gock.New("https://api.example.com").
        Get("/users/123").
        Reply(200).
        JSON(map[string]interface{}{
            "id":   123,
            "name": "John Doe",
        })
    
    client := NewAPIClient("https://api.example.com")
    user, err := client.GetUser("123")
    
    require.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
    assert.True(t, gock.IsDone())
}
```

### Базы данных (используй моки, не реальные БД)
```go
// ❌ Неправильно - не используй реальную БД в unit тестах
func TestUserRepository_Create_WithRealDB(t *testing.T) {
    db := setupTestDatabase() // НЕ ДЕЛАЙ ТАК!
    // ...
}

// ✅ Правильно - используй моки
func TestUserService_CreateUser_Success(t *testing.T) {
    f, cleanup := setUp(t)
    defer cleanup()
    
    f.userRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(&User{ID: "123"}, nil)
    
    // тест логика...
}
```

---

## Проверки (assertions) 🔍

### Используй require vs assert правильно
```go
func TestExample(t *testing.T) {
    // require - останавливает тест при ошибке
    user, err := service.GetUser("123")
    require.NoError(t, err)        // Если ошибка - тест прекращается
    require.NotNil(t, user)        // Если nil - тест прекращается
    
    // assert - продолжает тест после ошибки
    assert.Equal(t, "John", user.Name)      // Проверяет и продолжает
    assert.Equal(t, "john@example.com", user.Email) // Даже если предыдущая assert failed
}
```

### Полезные проверки
```go
// Проверка содержимого
assert.Contains(t, err.Error(), "user not found")
assert.NotContains(t, response, "password")

// Проверка длины и пустоты
assert.Len(t, users, 3)
assert.Empty(t, errors)
assert.NotEmpty(t, user.ID)

// Проверка типов
assert.IsType(t, &ValidationError{}, err)

// Проверка времени (с допуском)
assert.WithinDuration(t, expectedTime, actualTime, time.Second)

// Проверка JSON равенства
assert.JSONEq(t, expectedJSON, actualJSON)
```

---

## Тестирование контекста и таймаутов ⏰

### Тестирование отмены контекста
```go
func TestUserService_GetUser_ContextCancellation(t *testing.T) {
    f, cleanup := setUp(t)
    defer cleanup()
    
    // Создание контекста с отменой
    ctx, cancel := context.WithCancel(f.ctx)
    
    // Мок будет ждать отмены
    f.userRepo.EXPECT().
        GetByID(gomock.Any(), "user123").
        DoAndReturn(func(ctx context.Context, id string) (*User, error) {
            // Ожидание отмены контекста
            <-ctx.Done()
            return nil, ctx.Err()
        })
    
    // Запуск горутины которая отменит контекст
    go func() {
        time.Sleep(10 * time.Millisecond)
        cancel()
    }()
    
    user, err := f.userService.GetUser(ctx, "user123")
    
    require.Error(t, err)
    assert.Nil(t, user)
    assert.Equal(t, context.Canceled, err)
}
```

### Тестирование таймаутов
```go
func TestUserService_GetUser_Timeout(t *testing.T) {
    f, cleanup := setUp(t)
    defer cleanup()
    
    // Контекст с коротким таймаутом
    ctx, cancel := context.WithTimeout(f.ctx, 10*time.Millisecond)
    defer cancel()
    
    // Мок имитирует медленную операцию
    f.userRepo.EXPECT().
        GetByID(gomock.Any(), "user123").
        DoAndReturn(func(ctx context.Context, id string) (*User, error) {
            time.Sleep(20 * time.Millisecond) // Дольше чем таймаут
            return &User{ID: id}, nil
        })
    
    user, err := f.userService.GetUser(ctx, "user123")
    
    require.Error(t, err)
    assert.Nil(t, user)
    assert.Equal(t, context.DeadlineExceeded, err)
}
```

---

## Вспомогательные функции 🛠️

### Хелперы для создания тестовых данных
```go
func createTestUser(id string) *User {
    return &User{
        ID:        id,
        Name:      "Test User " + id,
        Email:     fmt.Sprintf("user%s@example.com", id),
        CreatedAt: time.Now(),
    }
}

func createTestPayment(userID string, amount int64) *Payment {
    return &Payment{
        ID:        fmt.Sprintf("payment_%s_%d", userID, amount),
        UserID:    userID,
        Amount:    amount,
        Currency:  "USD",
        Status:    "pending",
        CreatedAt: time.Now(),
    }
}

// Использование в тестах
func TestExample(t *testing.T) {
    user := createTestUser("123")
    payment := createTestPayment(user.ID, 1000)
    
    // тест логика...
}
```

### Хелперы для проверок
```go
func assertUserEqual(t *testing.T, expected, actual *User) {
    t.Helper()
    
    assert.Equal(t, expected.ID, actual.ID)
    assert.Equal(t, expected.Name, actual.Name)
    assert.Equal(t, expected.Email, actual.Email)
    // Время проверяется с допуском
    assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
}

func assertPaymentCompleted(t *testing.T, payment *Payment) {
    t.Helper()
    
    assert.Equal(t, "completed", payment.Status)
    assert.NotEmpty(t, payment.ExternalID)
    assert.NotZero(t, payment.CompletedAt)
}
```

---

## Тестирование конкурентности 🔄

### Тестирование горутин
```go
func TestUserService_BulkUpdate_Concurrent(t *testing.T) {
    f, cleanup := setUp(t)
    defer cleanup()
    
    users := []*User{
        createTestUser("1"),
        createTestUser("2"),
        createTestUser("3"),
    }
    
    // Моки должны быть thread-safe
    for _, user := range users {
        f.userRepo.EXPECT().
            Update(gomock.Any(), user).
            Return(nil).
            AnyTimes() // Может вызываться любое количество раз
    }
    
    // Запуск обновления в горутинах
    var wg sync.WaitGroup
    for _, user := range users {
        wg.Add(1)
        go func(u *User) {
            defer wg.Done()
            err := f.userService.UpdateUser(f.ctx, u)
            assert.NoError(t, err)
        }(user)
    }
    
    wg.Wait()
}
```

### Тестирование race conditions
```go
func TestCounter_Increment_RaceCondition(t *testing.T) {
    counter := NewCounter()
    
    // Запуск множественных горутин
    const numGoroutines = 100
    const incrementsPerGoroutine = 100
    
    var wg sync.WaitGroup
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < incrementsPerGoroutine; j++ {
                counter.Increment()
            }
        }()
    }
    
    wg.Wait()
    
    expected := numGoroutines * incrementsPerGoroutine
    assert.Equal(t, expected, counter.Value())
}
```

---

## Использование gotestsum для запуска тестов 🚀

### Основные принципы
- **Предпочтительный инструмент**: `gotestsum` вместо стандартного `go test`
- **Команда**: `gotestsum --format pkgname --packages="$(go list ./...)"`
- **Резервная команда**: `go test ./...` если gotestsum недоступен

### Преимущества gotestsum
- **Человекочитаемый вывод**: Цветная подсветка и лучшее форматирование
- **Названия пакетов**: Видимость какой пакет сейчас тестируется  
- **Поддержка JUnit XML**: Интеграция с CI/CD системами
- **Улучшенное форматирование ошибок**: Легче найти проблемы
- **Индикатор прогресса**: Показывает количество пройденных тестов

### Примеры использования
```bash
# Основная команда для всех тестов
gotestsum --format pkgname --packages="$(go list ./...)"

# Запуск конкретного пакета
gotestsum --format pkgname ./internal/service

# С дополнительными флагами
gotestsum --format pkgname --packages="$(go list ./...)" -- -v -count=1

# Генерация JUnit XML для CI
gotestsum --format pkgname --junitfile tests.xml --packages="$(go list ./...)"
```

### Интеграция с автоматизацией
- **Обязательное использование**: При автоматическом запуске тестов в пайплайне
- **Замена go test**: В скриптах и CI/CD конфигурациях
- **Обработка ошибок**: Тот же код возврата что и у `go test`

---

## Контрольный список качественных тестов ✅

### Обязательные элементы:
- [ ] Каждый тест использует `t.Parallel()`
- [ ] Используется fixture pattern с `setUp()`
- [ ] Все названия тестов на английском языке
- [ ] Используется table-driven подход где уместно
- [ ] Тестируются как успешные случаи, так и ошибки

### Моки и зависимости:
- [ ] Все внешние зависимости замокированы
- [ ] Моки настроены с правильными ожиданиями
- [ ] Используется `gomock.Controller` и `ctrl.Finish()`
- [ ] Не используются реальные БД или сетевые вызовы

### Проверки:
- [ ] Используется `require` для критических проверок
- [ ] Используется `assert` для дополнительных проверок
- [ ] Проверяются все важные поля результата
- [ ] Ошибки проверяются на содержание сообщения

### Читаемость:
- [ ] Тесты легко понять без изучения implementation
- [ ] Используются говорящие названия переменных
- [ ] Каждый тест проверяет одну конкретную вещь
- [ ] Есть комментарии для сложной логики

Качественные тесты - это ваша страховка от багов! 🛡️