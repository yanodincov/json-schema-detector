# Правила параллельности в Go ⚡ **ТОЛЬКО ДЛЯ GO**

**ПРИМЕНЕНИЕ**: Эти правила применяются исключительно для Go кода и Go проектов с горутинами.

## Цель 🎯
**Безопасная работа с конкурентностью** - обеспечение корректной синхронизации при параллельном выполнении кода.

---

## Защита разделяемых данных 🔒

### RWMutex для читаемых данных
**Используй когда много читателей, мало писателей**

```go
type UserCache struct {
    mu    sync.RWMutex
    users map[string]*User
}

// Чтение (множественные горутины могут читать одновременно)
func (uc *UserCache) GetUser(id string) (*User, bool) {
    uc.mu.RLock()
    defer uc.mu.RUnlock()
    
    user, exists := uc.users[id]
    return user, exists
}

// Запись (эксклюзивный доступ)
func (uc *UserCache) SetUser(id string, user *User) {
    uc.mu.Lock()
    defer uc.mu.Unlock()
    
    uc.users[id] = user
}
```

### Mutex для простой синхронизации
**Используй когда читателей и писателей примерно поровну**

```go
type Counter struct {
    mu    sync.Mutex
    value int64
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.value++
}

func (c *Counter) Value() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    return c.value
}
```

### Atomic для простых операций
**Используй для простых числовых операций**

```go
type AtomicCounter struct {
    value int64
}

func (ac *AtomicCounter) Increment() {
    atomic.AddInt64(&ac.value, 1)
}

func (ac *AtomicCounter) Value() int64 {
    return atomic.LoadInt64(&ac.value)
}

func (ac *AtomicCounter) CompareAndSwap(old, new int64) bool {
    return atomic.CompareAndSwapInt64(&ac.value, old, new)
}
```

---

## Параллельное выполнение 🚀

### WaitGroup для ожидания завершения
**Используй когда нужно дождаться завершения всех горутин**

```go
func ProcessUsers(users []User) {
    var wg sync.WaitGroup
    
    for _, user := range users {
        wg.Add(1)
        
        go func(u User) {
            defer wg.Done()
            
            // Обработка пользователя
            processUser(u)
        }(user)
    }
    
    wg.Wait() // Ждем завершения всех горутин
}
```

### errgroup для параллельной обработки с ошибками
**Используй когда нужно обработать ошибки от горутин**

```go
import "golang.org/x/sync/errgroup"

func ProcessUsersWithErrors(ctx context.Context, users []User) error {
    g, ctx := errgroup.WithContext(ctx)
    
    for _, user := range users {
        user := user // Важно! Захват переменной
        
        g.Go(func() error {
            return processUserWithError(ctx, user)
        })
    }
    
    return g.Wait() // Возвращает первую ошибку или nil
}
```

### errgroup с ограничением на количество горутин
**Используй для контроля нагрузки**

```go
func ProcessUsersWithLimit(ctx context.Context, users []User, maxConcurrency int) error {
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(maxConcurrency) // Ограничиваем количество параллельных горутин
    
    for _, user := range users {
        user := user
        
        g.Go(func() error {
            return processUserWithError(ctx, user)
        })
    }
    
    return g.Wait()
}
```

---

## Асинхронный сбор данных 📊

### Предопределенные слайсы для результатов
**Создавай слайс нужного размера заранее**

```go
func FetchUsersAsync(ctx context.Context, userIDs []string) ([]*User, error) {
    // Предопределяем слайс с правильным размером
    users := make([]*User, len(userIDs))
    g, ctx := errgroup.WithContext(ctx)
    
    for i, userID := range userIDs {
        i, userID := i, userID // Захват переменных
        
        g.Go(func() error {
            user, err := fetchUser(ctx, userID)
            if err != nil {
                return err
            }
            
            users[i] = user // Записываем в правильную позицию
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    
    return users, nil
}
```

### Сбор результатов с каналами
**Используй когда порядок результатов не важен**

```go
func FetchUsersWithChannel(ctx context.Context, userIDs []string) ([]*User, error) {
    results := make(chan *User, len(userIDs))
    g, ctx := errgroup.WithContext(ctx)
    
    for _, userID := range userIDs {
        userID := userID
        
        g.Go(func() error {
            user, err := fetchUser(ctx, userID)
            if err != nil {
                return err
            }
            
            select {
            case results <- user:
                return nil
            case <-ctx.Done():
                return ctx.Err()
            }
        })
    }
    
    go func() {
        g.Wait()
        close(results)
    }()
    
    var users []*User
    for user := range results {
        users = append(users, user)
    }
    
    return users, g.Wait()
}
```

---

## Ограничения использования каналов ⚠️

### Избегай каналы для простой синхронизации
```go
// ❌ Неправильно - излишне сложно
func BadExample() {
    done := make(chan bool)
    
    go func() {
        doWork()
        done <- true
    }()
    
    <-done
}

// ✅ Правильно - проще и эффективнее
func GoodExample() {
    var wg sync.WaitGroup
    wg.Add(1)
    
    go func() {
        defer wg.Done()
        doWork()
    }()
    
    wg.Wait()
}
```

### Используй каналы для передачи данных
```go
// ✅ Правильно - канал для передачи данных
func ProducerConsumer() {
    jobs := make(chan Job, 100)
    
    // Producer
    go func() {
        defer close(jobs)
        for i := 0; i < 1000; i++ {
            jobs <- Job{ID: i}
        }
    }()
    
    // Consumer
    for job := range jobs {
        processJob(job)
    }
}
```

---

## Паттерны конкурентности 🎨

### Worker Pool
```go
type WorkerPool struct {
    jobs    chan Job
    results chan Result
    workers int
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        jobs:    make(chan Job, 100),
        results: make(chan Result, 100),
        workers: workers,
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        go wp.worker(ctx)
    }
}

func (wp *WorkerPool) worker(ctx context.Context) {
    for {
        select {
        case job, ok := <-wp.jobs:
            if !ok {
                return
            }
            
            result := processJob(job)
            
            select {
            case wp.results <- result:
            case <-ctx.Done():
                return
            }
            
        case <-ctx.Done():
            return
        }
    }
}
```

### Fan-Out/Fan-In
```go
func FanOutFanIn(ctx context.Context, input <-chan Data) <-chan Result {
    numWorkers := runtime.NumCPU()
    
    // Fan-out: распределяем работу между воркерами
    workers := make([]<-chan Result, numWorkers)
    for i := 0; i < numWorkers; i++ {
        workers[i] = worker(ctx, input)
    }
    
    // Fan-in: собираем результаты от всех воркеров
    return fanIn(ctx, workers...)
}

func worker(ctx context.Context, input <-chan Data) <-chan Result {
    output := make(chan Result)
    
    go func() {
        defer close(output)
        
        for {
            select {
            case data, ok := <-input:
                if !ok {
                    return
                }
                
                result := processData(data)
                
                select {
                case output <- result:
                case <-ctx.Done():
                    return
                }
                
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return output
}

func fanIn(ctx context.Context, inputs ...<-chan Result) <-chan Result {
    output := make(chan Result)
    var wg sync.WaitGroup
    
    for _, input := range inputs {
        wg.Add(1)
        
        go func(ch <-chan Result) {
            defer wg.Done()
            
            for {
                select {
                case result, ok := <-ch:
                    if !ok {
                        return
                    }
                    
                    select {
                    case output <- result:
                    case <-ctx.Done():
                        return
                    }
                    
                case <-ctx.Done():
                    return
                }
            }
        }(input)
    }
    
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}
```

---

## Контрольный список безопасности ✅

### Мьютексы:
- [ ] Используй `defer mutex.Unlock()` сразу после `mutex.Lock()`
- [ ] RWMutex для случаев с частым чтением
- [ ] Mutex для сбалансированного чтения/записи
- [ ] Atomic для простых числовых операций

### Горутины:
- [ ] Всегда захватывай переменные цикла в горутинах
- [ ] Используй `errgroup` для обработки ошибок
- [ ] Ограничивай количество горутин при необходимости
- [ ] Контролируй время жизни горутин через контекст

### Каналы:
- [ ] Закрывай каналы только на стороне отправителя
- [ ] Используй буферизованные каналы для производительности
- [ ] Всегда проверяй `ok` при получении из канала
- [ ] Используй `select` с `ctx.Done()` для отмены

### Общие принципы:
- [ ] Избегай гонок данных (data races)
- [ ] Предпочитай сообщения блокировкам (каналы vs мьютексы)
- [ ] Проектируй для отмены операций через контекст
- [ ] Тестируй конкурентный код с флагом `-race`

Правильная работа с конкурентностью критически важна для надежности системы! 🛡️