# Основные правила разработки на Go

## Цель 🎯
**Единообразие, читаемость и поддерживаемость кода** - создание качественных решений следуя лучшим практикам Go.

---

## 1. Управление константами и глобальными переменными 📌

### Константы - preferred подход
```go
// ✅ Правильно - группировка в блоки
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
    APIVersion     = "v1"
)

// ✅ Отдельные константы для специфичных значений
const PaymentStatusCompleted = "completed"

// ✅ Типизированные константы
type Status string

const (
    StatusPending   Status = "pending"
    StatusCompleted Status = "completed" 
    StatusFailed    Status = "failed"
)
```

### Глобальные переменные (минимизировать)
```go
// ❌ Избегай глобальных переменных
var GlobalDB *sql.DB

// ✅ Предпочитай внедрение зависимостей
type Service struct {
    db *sql.DB
}

func NewService(db *sql.DB) *Service {
    return &Service{db: db}
}
```

---

## 2. Политика объявления переменных 📝

### Краткие объявления (предпочтительно)
```go
// ✅ Используй := где возможно
user := User{ID: "123", Name: "John"}
timeout := 30 * time.Second
data, err := fetchData()

// ✅ Множественное присваивание
result, err := processRequest(ctx, request)
if err != nil {
    return nil, err
}
```

### Явные объявления var (когда нужно)
```go
// ✅ Zero values
var (
    count    int
    isActive bool
    name     string
)

// ✅ Групповые объявления с начальными значениями
var (
    startTime = time.Now()
    logger    = zap.NewNop()
    cache     = make(map[string]interface{})
)
```

### Правила именования
```go
// ✅ Короткие имена для локальных переменных
for i, v := range items {
    processItem(i, v)
}

// ✅ Полные имена для экспортируемых
type UserRepository interface {
    CreateUser(ctx context.Context, user User) error
}

// ✅ Acronyms в верхнем регистре
type APIClient struct {
    HTTPClient *http.Client
    URLBase    string
}
```

---

## 3. Управление ошибками через errfmt 🚨

### Создание кастомных ошибок
```go
import "github.com/cryptoboyio/payments/pkg/errfmt"

// Определение типов ошибок
var (
    ErrCodeUserNotFound  = errfmt.NewError("failed to get user")
    ErrCodePaymentFailed = errfmt.NewError("payment processing failed")
    ErrCodeValidation    = errfmt.NewError("validation error")
)

func GetUser(ctx context.Context, userID string) (*User, error) {
    if userID == "" {
        return nil, errfmt.Error(ErrCodeValidation, "user id is required")
    }
    
    user, err := repo.GetByID(ctx, userID)
    if err != nil {
        return nil, errfmt.Chain(err, ErrCodeUserNotFound, "failed to get user")
    }
    
    return user, nil
}
```

### Обработка и проверка ошибок
```go
func ProcessPayment(ctx context.Context, paymentID string) error {
    payment, err := GetPayment(ctx, paymentID)
    if err != nil {
        // Проверка типа ошибки
        if errfmt.HasCode(err, ErrCodeUserNotFound) {
            return errfmt.Error(ErrCodePaymentFailed, "user not found for payment")
        }
        return errfmt.Chain(err, ErrCodePaymentFailed, "failed to get payment")
    }
    
    // бизнес логика...
    return nil
}
```

### Структурированные ошибки
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in field %s: %s", e.Field, e.Message)
}

func ValidateUser(user User) error {
    if user.Email == "" {
        return errfmt.Chain(
            ValidationError{Field: "email", Message: "is required"},
            ErrCodeValidation,
            "user validation failed",
        )
    }
    
    return nil
}
```

---

## 4. Правила go generate 🔄

### Обязательное размещение директив в верхней части файла
**ВАЖНО**: Все директивы `//go:generate` должны размещаться в самом начале файла, до объявления `package`.

### Правильное размещение директив
```go
// ✅ ПРАВИЛЬНО: директивы в начале файла
//go:generate mockgen -source=user_repository.go -destination=mock/user_repository.go
//go:generate stringer -type=Status -linecomment

package service

import (
    "context"
)

type Status int

const (
    StatusPending   Status = iota // pending
    StatusCompleted               // completed
    StatusFailed                  // failed
)
```

### Неправильное размещение директив
```go
// ❌ НЕПРАВИЛЬНО: директивы после package
package service

//go:generate mockgen -source=user_repository.go -destination=mock/user_repository.go

// ❌ НЕПРАВИЛЬНО: директивы после imports
package service

import "context"

//go:generate mockgen -source=user_repository.go -destination=mock/user_repository.go
```

### Групповые генерации
```go
//go:generate go run github.com/matryer/moq -out mock_client.go . HTTPClient
//go:generate go run github.com/golang/mock/mockgen -source=interfaces.go -destination=mocks/interfaces.go

package payment

type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}
```

### Отдельный файл для генерации
```go
// generate.go
package service

//go:generate mockgen -source=user_repository.go -destination=mock/user_repository.go
//go:generate mockgen -source=payment_service.go -destination=mock/payment_service.go
//go:generate stringer -type=PaymentStatus
```

---

## 5. Политика комментариев 💬

### Godoc комментарии (обязательно для экспорта)
```go
// User представляет пользователя в системе.
// Содержит основную информацию и методы для работы с пользователем.
type User struct {
    ID    string
    Email string
    Name  string
}

// CreateUser создает нового пользователя в системе.
// Возвращает ошибку если пользователь с таким email уже существует.
func CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // implementation...
}
```

### Внутренние комментарии (только для сложной логики)
```go
func ProcessPayment(ctx context.Context, payment Payment) error {
    // Сложная бизнес-логика требующая объяснения
    if payment.Amount > 10000 {
        // Большие суммы требуют дополнительной верификации
        if err := verifyHighValuePayment(ctx, payment); err != nil {
            return err
        }
    }
    
    // Обычная обработка не требует комментариев
    return s.paymentGateway.Process(ctx, payment)
}
```

### TODO комментарии (временные)
```go
func ProcessRefund(ctx context.Context, refundID string) error {
    // TODO: добавить валидацию refund amount
    // TODO: добавить уведомления пользователя
    
    return s.refundService.Process(ctx, refundID)
}
```

---

## 6. Стиль форматирования 🎨

### Порядок импортов (автоматически через goimports)
```go
package main

import (
    // Стандартная библиотека
    "context"
    "fmt"
    "time"
    
    // Сторонние пакеты
    "github.com/gorilla/mux"
    "go.uber.org/zap"
    
    // Проектные пакеты
    "github.com/cryptoboyio/payments/internal/service"
    "github.com/cryptoboyio/payments/pkg/logger"
)
```

### Группировка полей структур
```go
type User struct {
    // Основные поля
    ID    string
    Email string
    Name  string
    
    // Метаданные
    CreatedAt time.Time
    UpdatedAt time.Time
    
    // Зависимости (последними)
    logger *zap.Logger
}
```

### Выравнивание в блоках
```go
const (
    StatusPending    = "pending"
    StatusCompleted  = "completed"
    StatusFailed     = "failed"
    StatusCancelled  = "cancelled"
)

var (
    ErrUserNotFound     = errors.New("user not found")
    ErrPaymentFailed    = errors.New("payment failed")
    ErrInvalidAmount    = errors.New("invalid amount")
)
```

---

## 7. Контекст как первый параметр 📋

### Обязательное правило
```go
// ✅ Контекст всегда первый параметр
func GetUser(ctx context.Context, userID string) (*User, error) {
    return repo.GetByID(ctx, userID)
}

func ProcessPayment(ctx context.Context, payment Payment, options ProcessOptions) error {
    return gateway.Process(ctx, payment, options)
}

// ✅ В методах контекст тоже первый
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    return s.userRepo.Create(ctx, req.ToUser())
}
```

### Передача контекста
```go
func (s *Service) ComplexOperation(ctx context.Context, data Data) error {
    // Передавай тот же контекст дальше
    if err := s.validateData(ctx, data); err != nil {
        return err
    }
    
    // Создавай дочерний контекст только при необходимости
    processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    return s.processData(processCtx, data)
}
```

---

## 8. DTO для множественных параметров 📦

### Создание Request/Response DTO
```go
// Request DTO для входящих параметров
type CreateUserRequest struct {
    Email     string            `json:"email" validate:"required,email"`
    Name      string            `json:"name" validate:"required"`
    Phone     string            `json:"phone" validate:"phone"`
    Metadata  map[string]string `json:"metadata,omitempty"`
}

func (r CreateUserRequest) ToUser() User {
    return User{
        ID:       uuid.New().String(),
        Email:    r.Email,
        Name:     r.Name,
        Phone:    r.Phone,
        Metadata: r.Metadata,
    }
}

// Response DTO для ответов
type CreateUserResponse struct {
    User      User      `json:"user"`
    CreatedAt time.Time `json:"created_at"`
    Success   bool      `json:"success"`
}
```

### Использование DTO в методах
```go
// ✅ Используй DTO когда параметров больше 2-3
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
    user := req.ToUser()
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return &CreateUserResponse{
        User:      user,
        CreatedAt: time.Now(),
        Success:   true,
    }, nil
}

// ✅ Простые методы могут обойтись без DTO
func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
    return s.userRepo.GetByID(ctx, userID)
}
```

---

## 9. Организация интерфейсов 🔗

### Размещение интерфейсов
**Интерфейсы располагаются над структурой, но под константами и переменными**

```go
package service

// Константы в начале
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// Глобальные переменные после констант
var (
    ErrUserNotFound  = errfmt.NewError("user not found")
    ErrPaymentFailed = errfmt.NewError("payment failed")
)

// Интерфейсы над структурой (приватные по умолчанию)
type userRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user User) error
    Update(ctx context.Context, id string, user User) error
}

type paymentGateway interface {
    Process(ctx context.Context, payment Payment) error
    Refund(ctx context.Context, paymentID string) error
}

// Основная структура после интерфейсов
type Service struct {
    userRepo userRepository
    gateway  paymentGateway
    logger   *zap.Logger
}

// Конструктор
func NewService(userRepo userRepository, gateway paymentGateway, logger *zap.Logger) *Service {
    return &Service{
        userRepo: userRepo,
        gateway:  gateway,
        logger:   logger,
    }
}
```

### Принцип "потребитель определяет интерфейс"
**Сервис сам определяет какую функциональность он требует от зависимостей**

```go
// ✅ Правильно - PaymentService определяет что ему нужно от репозитория
package payment

type userRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    // Только методы, которые нужны PaymentService
}

type PaymentService struct {
    userRepo userRepository // использует свой интерфейс
}

// ❌ Неправильно - использование интерфейса из другого пакета
import "myproject/user"

type PaymentService struct {
    userRepo user.Repository // прямая зависимость на конкретный тип
}
```

### Именование интерфейсов
```go
// ✅ Приватные интерфейсы (по умолчанию)
type userRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
}

type paymentProcessor interface {
    Process(ctx context.Context, payment Payment) error
}

// ✅ Публичные интерфейсы (только когда действительно нужно)
type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id string) (*User, error)
}
```

### Внедрение зависимостей
```go
// main.go или wire.go
func main() {
    // Создаем конкретные реализации
    userRepo := postgres.NewUserRepository(db)
    paymentGW := stripe.NewPaymentGateway(apiKey)
    
    // Внедряем как интерфейсы
    paymentService := payment.NewService(userRepo, paymentGW, logger)
    
    // paymentService знает только о своих интерфейсах,
    // но работает с конкретными реализациями
}
```

### Тестирование с интерфейсами
```go
// payment_service_test.go
type mockUserRepository struct {
    users map[string]*User
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    user, exists := m.users[id]
    if !exists {
        return nil, ErrUserNotFound
    }
    return user, nil
}

func TestPaymentService_Process(t *testing.T) {
    mockRepo := &mockUserRepository{
        users: map[string]*User{
            "123": {ID: "123", Email: "test@example.com"},
        },
    }
    
    service := NewService(mockRepo, nil, logger)
    
    // тестирование...
}
```

### Нет файла types.go
**Интерфейсы размещаются в том же файле, что и структура их использующая**

```go
// ✅ Все в одном файле - payment_service.go
package payment

const DefaultTimeout = 30 * time.Second

// Интерфейсы в этом же файле
type userRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
}

type Service struct {
    userRepo userRepository
}

func (s *Service) ProcessPayment(ctx context.Context, payment Payment) error {
    user, err := s.userRepo.GetByID(ctx, payment.UserID)
    // ...
}
  ```
  
---

## 10. Запрет init функций 🚫

### Почему init запрещены
```go
// ❌ НЕ используй init функции
func init() {
    // Непредсказуемый порядок выполнения
    // Сложно тестировать
    // Скрытые зависимости
    database.Connect()
    logger.SetLevel("debug")
}
```

### Альтернативы init
```go
// ✅ Используй явную инициализацию
func NewService(config Config) (*Service, error) {
    db, err := connectToDatabase(config.DatabaseURL)
    if err != nil {
        return nil, err
    }
    
    logger := setupLogger(config.LogLevel)
    
    return &Service{
        db:     db,
        logger: logger,
    }, nil
}

// ✅ Используй sync.Once для одноразовой инициализации
type SingletonService struct {
    instance *Service
    once     sync.Once
}

func (s *SingletonService) GetInstance() *Service {
    s.once.Do(func() {
        s.instance = &Service{
            // инициализация...
        }
    })
    return s.instance
}
```

---

## 11. Организация моделей (3 варианта) 📁

### Вариант 1: DTO над методом (рекомендуется для небольших методов)
```go
// create_user.go

// CreateUserRequest представляет запрос на создание пользователя
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
}

// CreateUserResponse представляет ответ создания пользователя
type CreateUserResponse struct {
    User    User   `json:"user"`
    Success bool   `json:"success"`
}

// CreateUser создает нового пользователя в системе
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
    // implementation...
}
```

### Вариант 2: Общие DTO в model.go (рекомендуется для больших сервисов)
```go
// model.go

// User представляет пользователя в системе
type User struct {
    ID        string            `json:"id"`
    Email     string            `json:"email"`
    Name      string            `json:"name"`
    CreatedAt time.Time         `json:"created_at"`
    Metadata  map[string]string `json:"metadata"`
}

// CreateUserRequest для создания пользователя
type CreateUserRequest struct {
    Email    string            `json:"email" validate:"required,email"`
    Name     string            `json:"name" validate:"required"`
    Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateUserRequest для обновления пользователя
type UpdateUserRequest struct {
    Name     *string           `json:"name,omitempty"`
    Metadata map[string]string `json:"metadata,omitempty"`
}
```

### Вариант 3: Отдельные файлы для каждой модели (для сложных доменов)
```go
// user.go
type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
}

// payment.go  
type Payment struct {
    ID       string
    UserID   string
    Amount   int64
    Currency string
    Status   PaymentStatus
}

// order.go
type Order struct {
    ID       string
    UserID   string
    Items    []OrderItem
    Total    int64
    Status   OrderStatus
}
```

---

## 12. Размещение ошибок, констант и приватных методов 📍

### Ошибки в начале файла
```go
package service

import (...)

// Ошибки сразу после импортов
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrEmailExists      = errors.New("email already exists")
    ErrInvalidPassword  = errors.New("invalid password")
)

// Константы после ошибок
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// Основная структура сервиса
type Service struct {
    userRepo UserRepository
    logger   *zap.Logger
}

// Публичные методы
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) error {
    return s.createUserInternal(ctx, req)
}

// Приватные методы в конце файла
func (s *Service) createUserInternal(ctx context.Context, req CreateUserRequest) error {
    // implementation...
}

func (s *Service) validateUser(user User) error {
    // implementation...
}
```

### Организация больших файлов
```go
// service.go - основной файл
package service

// Errors
var (
    ErrUserNotFound    = errors.New("user not found")
)

// Constants
const (
    DefaultTimeout = 30 * time.Second
)

// Main struct
type Service struct {
    userRepo UserRepository
    logger   *zap.Logger
}

// Constructor
func NewService(userRepo UserRepository, logger *zap.Logger) *Service {
    return &Service{
        userRepo: userRepo,
        logger:   logger,
    }
}

// Public methods
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) error {
    // implementation...
}

// Private methods (в конце)
func (s *Service) validateUser(user User) error {
    // implementation...
}
```

---

## 13. Организация файлов для больших структур 📂

### Разделение по функциональности
```
service/
├── service.go          # Основная структура и конструктор
├── user.go            # Методы работы с пользователями
├── payment.go         # Методы работы с платежами
├── notification.go    # Методы уведомлений
├── validation.go      # Приватные методы валидации
└── helpers.go         # Приватные вспомогательные методы
```

### service.go - основной файл
```go
// service.go
package service

// Все ошибки в основном файле
var (
    ErrUserNotFound    = errors.New("user not found")
)

// Все константы в основном файле
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// Основная структура
type Service struct {
    userRepo    UserRepository
    paymentRepo PaymentRepository
    logger      *zap.Logger
}

// Конструктор
func NewService(
    userRepo UserRepository,
    paymentRepo PaymentRepository,
    logger *zap.Logger,
) *Service {
    return &Service{
        userRepo:    userRepo,
        paymentRepo: paymentRepo,
        logger:      logger,
    }
}
```

### user.go - методы пользователей
```go
// user.go
package service

// Только методы работы с пользователями
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    if err := s.validateUserRequest(req); err != nil {
        return nil, err
    }
    
    return s.userRepo.Create(ctx, req.ToUser())
}

func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
    return s.userRepo.GetByID(ctx, userID)
}

func (s *Service) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error {
    return s.userRepo.Update(ctx, userID, req)
}
```

### validation.go - приватные методы
```go
// validation.go
package service

// Приватные методы валидации
func (s *Service) validateUserRequest(req CreateUserRequest) error {
    if req.Email == "" {
        return ErrInvalidEmail
    }
    
    return nil
}

func (s *Service) validatePaymentRequest(req PaymentRequest) error {
    if req.Amount <= 0 {
        return ErrInvalidAmount
    }
    
    return nil
}
```

---

## 14. Правила форматирования кода 🎨

### Пустые строки перед return

**Основное правило**: Перед каждым `return` должна быть пустая строка.

**Исключение**: Если `return` является первой инструкцией в области видимости (начало метода, if, for, switch и т.д.).

#### ✅ Правильно - return в начале области видимости
```go
// Начало метода
func GetUser() User {
    return defaultUser // OK - первая инструкция в методе
}

// Начало if блока
if condition {
    return value // OK - первая инструкция в if
}

// Начало for цикла
for i := 0; i < len(items); i++ {
    return items[i] // OK - первая инструкция в for
}

// Начало switch case
switch status {
case "active":
    return processActive() // OK - первая инструкция в case
default:
    return defaultResult() // OK - первая инструкция в default
}

// Начало select case
select {
case result := <-ch:
    return result // OK - первая инструкция в case
default:
    return nil // OK - первая инструкция в default
}

// Начало функции-литерала
handler := func() error {
    return processRequest() // OK - первая инструкция в функции
}
```

#### ✅ Правильно - return с пустой строкой
```go
// После объявления переменной
func ProcessData() int {
    result := calculateResult()
    
    return result // Пустая строка обязательна
}

// После вызова метода
func SaveUser(user User) error {
    s.processData()
    
    return nil // Пустая строка обязательна
}

// После присваивания
func UpdateUser(user *User) *User {
    user.Status = "active"
    user.UpdatedAt = time.Now()
    
    return user // Пустая строка обязательна
}

// После цикла
func ProcessItems(items []Item) error {
    for _, item := range items {
        processItem(item)
    }
    
    return nil // Пустая строка обязательна
}
```

#### ❌ Неправильно - отсутствие пустой строки
```go
func ProcessPayment() error {
    amount := calculateAmount()
    return processAmount(amount) // ❌ Нужна пустая строка выше

    s.logOperation()
    return nil // ❌ Нужна пустая строка выше

    user.Status = "processed"
    return user // ❌ Нужна пустая строка выше
}
```

### Обоснование правила

1. **Читаемость**: Пустая строка визуально отделяет логику от возврата значения
2. **Понимание потока**: Легче увидеть, где заканчивается обработка и начинается возврат
3. **Единообразие**: Стандартизация форматирования для всей команды
4. **Сопровождение**: Упрощает понимание кода при ревью и отладке

---

## Контрольный список качественного Go кода ✅

### Структура и организация:
- [ ] Константы и ошибки в начале файла
- [ ] Публичные методы перед приватными
- [ ] Логическая группировка методов по файлам
- [ ] Правильные импорты с группировкой

### Обработка ошибок:
- [ ] Использование errfmt для кастомных ошибок
- [ ] Проверка всех возвращаемых ошибок
- [ ] Информативные сообщения об ошибках
- [ ] Правильная передача ошибок вверх

### Стиль кодирования:
- [ ] Контекст всегда первый параметр
- [ ] DTO для сложных параметров
- [ ] Отсутствие init функций
- [ ] Правильное именование переменных

### Документация:
- [ ] Godoc комментарии для экспортируемых элементов
- [ ] Комментарии только для сложной логики
- [ ] Информативные названия функций и переменных
- [ ] Примеры использования в комментариях

Следование этим правилам обеспечивает поддерживаемый и читаемый код! 🚀