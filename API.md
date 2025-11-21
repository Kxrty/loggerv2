# API документация LoggerV2

## Пакет processor

### NewProcessor()

Создает новый экземпляр процессора логов.

```go
proc := processor.NewProcessor()
```

### Process(logLine string) (*models.GOSTEvent, error)

Обрабатывает одну строку лога и возвращает событие в формате ГОСТ.

**Параметры:**
- `logLine` - строка лога в одном из поддерживаемых форматов

**Возвращает:**
- `*models.GOSTEvent` - событие в формате ГОСТ
- `error` - ошибка, если лог не удалось распарсить

**Пример:**
```go
event, err := proc.Process("<134>Oct 11 22:14:15 mymachine su: test")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Событие: %s\n", event.Description)
```

### ProcessBatch(logLines []string) ([]*models.GOSTEvent, []error)

Обрабатывает массив строк логов.

**Параметры:**
- `logLines` - массив строк логов

**Возвращает:**
- `[]*models.GOSTEvent` - массив успешно обработанных событий
- `[]error` - массив ошибок для строк, которые не удалось обработать

**Пример:**
```go
logs := []string{
    "<134>Oct 11 22:14:15 mymachine su: test1",
    "CEF:0|Vendor|Product|1.0|100|Event|5|src=1.1.1.1",
}
events, errors := proc.ProcessBatch(logs)
fmt.Printf("Обработано: %d, Ошибок: %d\n", len(events), len(errors))
```

### DetectLogType(logLine string) LogType

Автоматически определяет тип лога.

**Параметры:**
- `logLine` - строка лога

**Возвращает:**
- `LogType` - тип лога (LogTypeSyslog, LogTypeCEF, LogTypeLEEF, LogTypeXML, LogTypeUnknown)

**Пример:**
```go
logType := proc.DetectLogType("<134>Oct 11 22:14:15 mymachine su: test")
if logType == processor.LogTypeSyslog {
    fmt.Println("Это Syslog")
}
```

### ConvertToJSON(event *models.GOSTEvent) (string, error)

Преобразует событие ГОСТ в JSON строку.

**Параметры:**
- `event` - событие в формате ГОСТ

**Возвращает:**
- `string` - JSON строка
- `error` - ошибка сериализации

**Пример:**
```go
json, err := proc.ConvertToJSON(event)
if err != nil {
    log.Fatal(err)
}
fmt.Println(json)
```

### ConvertBatchToJSON(events []*models.GOSTEvent) (string, error)

Преобразует массив событий ГОСТ в JSON строку.

**Параметры:**
- `events` - массив событий в формате ГОСТ

**Возвращает:**
- `string` - JSON строка
- `error` - ошибка сериализации

## Пакет parser

### SyslogParser

Парсер для Syslog сообщений (RFC 3164 и RFC 5424).

```go
parser := parser.NewSyslogParser()
event, err := parser.Parse("<134>Oct 11 22:14:15 mymachine su: test")
```

### CEFParser

Парсер для CEF (Common Event Format) логов.

```go
parser := parser.NewCEFParser()
event, err := parser.Parse("CEF:0|Vendor|Product|1.0|100|Event|5|src=1.1.1.1")
```

### LEEFParser

Парсер для LEEF (Log Event Extended Format) логов.

```go
parser := parser.NewLEEFParser()
event, err := parser.Parse("LEEF:1.0|Vendor|Product|1.0|100|key=value")
```

### XMLParser

Парсер для XML логов (например, Windows Event Log).

```go
parser := parser.NewXMLParser()
event, err := parser.Parse("<Event>...</Event>")
```

## Пакет models

### GOSTEvent

Структура события в формате ГОСТ Р 59710-2022.

```go
type GOSTEvent struct {
    EventID          string                 // Уникальный ID события
    Timestamp        time.Time              // Время события
    Source           Source                 // Источник события
    Category         string                 // Категория события
    Severity         string                 // Уровень критичности
    Description      string                 // Описание события
    AdditionalData   map[string]interface{} // Дополнительные данные
    SubjectAccount   *Account               // Субъект действия
    ObjectAccount    *Account               // Объект действия
    Result           string                 // Результат (УСПЕХ/НЕУСПЕХ/НЕИЗВЕСТНО)
    Action           string                 // Выполненное действие
}
```

### Source

Информация об источнике события.

```go
type Source struct {
    Hostname    string // Имя хоста
    IPAddress   string // IP адрес
    Application string // Имя приложения
    Process     string // Имя процесса
    ProcessID   int    // ID процесса
}
```

### Account

Информация об учетной записи.

```go
type Account struct {
    Username string // Имя пользователя
    Domain   string // Домен
    UserID   string // ID пользователя
}
```

## Константы

### Уровни критичности (Severity)

```go
const (
    SeverityCritical = "КРИТИЧЕСКИЙ"      // Критические события
    SeverityHigh     = "ВЫСОКИЙ"          // Высокий уровень
    SeverityMedium   = "СРЕДНИЙ"          // Средний уровень
    SeverityLow      = "НИЗКИЙ"           // Низкий уровень
    SeverityInfo     = "ИНФОРМАЦИОННЫЙ"   // Информационные события
)
```

### Результаты событий (Result)

```go
const (
    ResultSuccess = "УСПЕХ"      // Успешное выполнение
    ResultFailure = "НЕУСПЕХ"    // Неуспешное выполнение
    ResultUnknown = "НЕИЗВЕСТНО" // Результат неизвестен
)
```

### Категории событий (Category)

```go
const (
    CategoryAuthentication   = "АУТЕНТИФИКАЦИЯ"      // События входа/аутентификации
    CategoryAuthorization    = "АВТОРИЗАЦИЯ"         // События авторизации
    CategoryAccess           = "ДОСТУП"              // События доступа
    CategoryDataModification = "ИЗМЕНЕНИЕ_ДАННЫХ"    // События изменения данных
    CategorySystemEvent      = "СИСТЕМНОЕ_СОБЫТИЕ"   // Системные события
    CategorySecurityEvent    = "СОБЫТИЕ_БЕЗОПАСНОСТИ" // События безопасности
    CategoryNetworkEvent     = "СЕТЕВОЕ_СОБЫТИЕ"     // Сетевые события
)
```

## Полный пример использования

```go
package main

import (
    "fmt"
    "log"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    // Создаем процессор
    proc := processor.NewProcessor()
    
    // Обрабатываем один лог
    logLine := "<134>Oct 11 22:14:15 mymachine su: 'su root' failed"
    event, err := proc.Process(logLine)
    if err != nil {
        log.Fatal(err)
    }
    
    // Выводим информацию о событии
    fmt.Printf("ID: %s\n", event.EventID)
    fmt.Printf("Время: %s\n", event.Timestamp)
    fmt.Printf("Источник: %s\n", event.Source.Hostname)
    fmt.Printf("Категория: %s\n", event.Category)
    fmt.Printf("Критичность: %s\n", event.Severity)
    fmt.Printf("Описание: %s\n", event.Description)
    
    // Преобразуем в JSON
    jsonOutput, err := proc.ConvertToJSON(event)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(jsonOutput)
    
    // Обрабатываем несколько логов
    logs := []string{
        "<134>Oct 11 22:14:15 mymachine su: test1",
        "CEF:0|Vendor|Product|1.0|100|Event|5|src=1.1.1.1",
        "LEEF:1.0|Vendor|Product|1.0|100|key=value",
    }
    
    events, errors := proc.ProcessBatch(logs)
    fmt.Printf("Обработано: %d событий\n", len(events))
    fmt.Printf("Ошибок: %d\n", len(errors))
    
    // Преобразуем все события в JSON
    batchJSON, err := proc.ConvertBatchToJSON(events)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(batchJSON)
}
```

## Обработка ошибок

```go
event, err := proc.Process(logLine)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "неподдерживаемый формат"):
        fmt.Println("Неизвестный формат лога")
    case strings.Contains(err.Error(), "неверный формат"):
        fmt.Println("Лог поврежден или имеет неправильную структуру")
    default:
        fmt.Printf("Ошибка: %v\n", err)
    }
}
```

## Расширение функционала

### Добавление нового парсера

1. Создайте новый файл в `internal/parser/`
2. Реализуйте интерфейс `Parse(logLine string) (*models.GOSTEvent, error)`
3. Добавьте новый `LogType` в `processor.go`
4. Обновите метод `DetectLogType()` для определения нового формата
5. Добавьте обработку нового типа в метод `Process()`

### Добавление новых полей в GOSTEvent

1. Обновите структуру `GOSTEvent` в `internal/models/gost.go`
2. Обновите все парсеры для заполнения новых полей
3. Добавьте соответствующие тесты
