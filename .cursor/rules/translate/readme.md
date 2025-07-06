# Стандарт документации README

## Философия 🎯
**Живой документ - главный справочник разработчика** с акцентом на ясность, точность и полноту информации.

- **Язык**: ОБЯЗАТЕЛЬНО английский
- **Цель**: Комплексная практическая документация сервиса
- **Подход**: Первоисточник информации для разработчиков

---

## Обязательная структура разделов (строгий порядок) 📋

### 1. Заголовок сервиса (H1)
```markdown
# <service_name> service
```

**Содержание**: 
- Назначение сервиса
- Основной функционал
- Главные обязанности
- Случаи использования

### 2. Запускаемые компоненты (H2)
```markdown
## Running components
```

**Формат**: Список имен бинарных файлов (только имена)
```markdown
- api
- worker
- callback-processor
```

### 3. Зависимости (H2)
```markdown
## Dependencies
```

#### Инфраструктура (H3)
```markdown
### Infrastructure
```
Базы данных, очереди, кеши, внешние системы

#### Внутренние сервисы (H3)
```markdown
### Internal services
```
**Формат**: `- service-name (PROTOCOL)`
```markdown
- user-service (GRPC)
- notification-service (HTTP)
```

#### Внешние API (H3)
```markdown
### External API
```
**Формат**: `- provider-name, mandatory/optional`
```markdown
- stripe-api, mandatory
- paypal-api, optional
```

### 4. Критические метрики (H2)
```markdown
## Critical metrics
```

**Формат**: `* \`metric_name{labels="values"}\` Description with interpretation guidance`

```markdown
* `payment_transactions_total{status="success|failed"}` Monitor for >5% failure rate indicating provider issues.
* `queue_messages{queue="callbacks"}` Values >100 for >10min indicate processing bottleneck.
```

**Содержание**:
- **Только критические метрики** - не полный список
- **Пороговые значения** - когда начинать беспокоиться
- **Индикаторы проблем** - что означают аномальные значения

### 5. Ответственные за проект (H2)
```markdown
## Project maintainers
```

**Формат**:
```markdown
Code changes to the project should be reviewed by:

* Name <email>
* Name <email>

In case of inaccessibility ask Unit-X members.
```

**Содержание**:
- Список ревьюеров кода
- Информация для эскалации при недоступности

---

## Опциональные разделы (разрешены) 📝

### Архитектура и поток данных
```markdown
## Architecture & Flow
```

### Руководство по устранению неисправностей
```markdown
## Troubleshooting Guide
```

### Любые другие разделы
По мере необходимости - README может быть расширен специфичными для проекта разделами.

---

## Требования к форматированию 🎨

### Markdown стандарт
- **Основа**: Стандартные расширения GitHub
- **Таблицы**: Разрешены где уместно
- **Заголовки**: Строгое соответствие названиям и уровням (для автоматической обработки)

### Блоки кода
Обязательно указывай язык для подсветки синтаксиса:
```markdown
```json
{
  "config": "value"
}
```

```yaml
version: "3.8"
services:
  app:
    image: app:latest
```

```bash
make build
docker-compose up -d
```

```go
func main() {
    fmt.Println("Hello, World!")
}
```

### Акценты и выделения
- **Жирный текст**: `**text**` для ключевых терминов, команд, имен файлов
- **Код**: `` `text` `` для инлайн-кода, путей, переменных окружения, JSON-ключей

```markdown
Run **make build** to compile the project.
Set `DATABASE_URL` environment variable.
The config is stored in `/etc/app/config.yaml`.
```

---

## Организация документации 📂

### Корневой README
**Обязательный** - должен быть в корне проекта и следовать стандарту.

### README в подпапках
**Приветствуется** - детальная документация для сложных компонентов.

### Общая документация
**Отдельная папка docs/** - для диаграмм последовательности, архитектурных документов, технических спецификаций.

---

## Процесс ревью и обновления 🔄

### Ревью документации
- **Изменения документации**: Ревьюируются как код через PR
- **Деплой не требуется** - документация автоматически доступна
- **Изменения maintainer'ов**: Требуют одобрения Tech Lead'а команды Go
- **Связь с JIRA**: Опциональна для изменений документации

### Политика обновления
**Обновляй вместе с значительными изменениями сервиса** - README должен отражать текущее состояние.

---

## Обеспечение качества ✅

### Диаграммы
- **Используй где уместно**: PlantUML, D2 предпочтительны
- **Актуальность**: Диаграммы должны соответствовать реальной архитектуре

### Внешние ссылки
- **Confluence документация**: Проверяй валидность ссылок
- **Внешние ресурсы**: Убедись что ссылки рабочие

### Информация о maintainer'ах
- **Актуальность**: Контакты должны быть актуальными
- **Ясные пути эскалации**: Четкие инструкции кому обращаться

---

## Шаблон README 📋

```markdown
# Payment Gateway service

Payment processing service that handles transactions, manages payment methods, and integrates with multiple payment providers. Supports card payments, bank transfers, and digital wallets with real-time fraud detection.

## Running components
- api
- worker  
- callback-processor

## Dependencies

### Infrastructure
- postgresql
- redis
- rabbitmq

### Internal services
- user-service (GRPC)
- notification-service (HTTP)
- fraud-detection-service (HTTP)

### External API
- stripe-api, mandatory
- paypal-api, optional
- fraud-check-api, mandatory

## Critical metrics

* `payment_transactions_total{status="success|failed"}` Monitor for >5% failure rate indicating provider issues.
* `queue_messages{queue="callbacks"}` Values >100 for >10min indicate processing bottleneck.
* `external_api_duration{provider="stripe"}` Response times >5s indicate provider degradation.
* `database_connections_active` Values near max pool size indicate connection leaks.

## Project maintainers

Code changes to the project should be reviewed by:

* John Doe <j.doe@company.com>
* Jane Smith <j.smith@company.com>

In case of inaccessibility ask Unit-1 members.
```

---

## Дополнительные рекомендации 💡

### Для команд
- **Живой документ**: README должен обновляться при каждом значимом изменении
- **Onboarding**: Новые разработчики должны смочь развернуть проект используя только README
- **Troubleshooting**: Включай решения частых проблем

### Для maintainer'ов
- **Регулярная проверка**: Периодически проверяй актуальность всех разделов
- **Feedback**: Собирай обратную связь от команды об удобстве документации
- **Автоматизация**: Рассматривай возможность автоматического обновления некоторых разделов

### Для безопасности
- **Не включай секреты**: Никаких паролей, ключей, токенов
- **Ограничь внутреннюю информацию**: README может быть доступен широкому кругу
- **Ссылки на конфиденциальную документацию**: Используй внутренние системы

---

## Контрольный список качества ✅

### Обязательные элементы:
- [ ] Заголовок в формате `# <service_name> service`
- [ ] Раздел "Running components" со списком бинарников
- [ ] Раздел "Dependencies" с подразделами Infrastructure/Internal/External
- [ ] Раздел "Critical metrics" с интерпретацией
- [ ] Раздел "Project maintainers" с контактами

### Качество содержания:
- [ ] Описание сервиса понятно новичку
- [ ] Все зависимости перечислены с протоколами
- [ ] Критические метрики включают пороговые значения
- [ ] Контакты maintainer'ов актуальны

### Форматирование:
- [ ] Весь текст на английском языке
- [ ] Блоки кода имеют указание языка
- [ ] Использованы правильные Markdown заголовки
- [ ] Ссылки проверены и рабочие

### Поддержка:
- [ ] README обновляется с изменениями проекта
- [ ] Процесс ревью документации настроен
- [ ] Feedback от команды учитывается

Качественный README экономит время всей команды и упрощает onboarding новых разработчиков! 🚀 