# Инструкция по установке и использованию LoggerV2

## Требования

- Go 1.21 или выше
- Git (опционально, для клонирования репозитория)

## Установка

### 1. Клонирование проекта

```bash
git clone https://github.com/kxrty/loggerv2.git
cd loggerv2
```

### 2. Установка зависимостей

```bash
go mod download
```

### 3. Сборка приложения

#### Windows
```bash
go build -o logger.exe ./cmd/logger
```

#### Linux/macOS
```bash
go build -o logger ./cmd/logger
```

## Использование

### Командная строка

#### Обработка из файла
```bash
# Windows
.\logger.exe -input logs.txt -output result.json

# Linux/macOS
./logger -input logs.txt -output result.json
```

#### Обработка из stdin
```bash
# Windows
echo "<134>Oct 11 22:14:15 mymachine su: test" | .\logger.exe

# Linux/macOS
echo "<134>Oct 11 22:14:15 mymachine su: test" | ./logger
```

#### Обработка с перенаправлением
```bash
# Windows
type logs.txt | .\logger.exe > result.json

# Linux/macOS
cat logs.txt | ./logger > result.json
```

### Использование как библиотеки

```go
package main

import (
    "fmt"
    "log"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    
    logLine := "<134>Oct 11 22:14:15 mymachine su: test message"
    
    event, err := proc.Process(logLine)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Событие: %s\n", event.Description)
    fmt.Printf("Критичность: %s\n", event.Severity)
    fmt.Printf("Категория: %s\n", event.Category)
}
```

## Тестирование

### Запуск всех тестов
```bash
go test ./...
```

### Запуск тестов с подробным выводом
```bash
go test ./... -v
```

### Запуск тестов с покрытием кода
```bash
go test ./... -cover
```

### Запуск тестов конкретного пакета
```bash
go test ./internal/parser -v
```

## Примеры

В папке `examples/` находятся примеры использования:

- `sample_logs.txt` - примеры логов различных форматов
- `example_api.go` - пример использования API

### Запуск примера API
```bash
go run examples/example_api.go
```

## Структура проекта

```
loggerv2/
├── cmd/
│   └── logger/
│       └── main.go           # CLI приложение
├── internal/
│   ├── models/
│   │   └── gost.go          # Модели ГОСТ
│   ├── parser/
│   │   ├── syslog.go        # Парсер Syslog (RFC 3164, RFC 5424)
│   │   ├── cef.go           # Парсер CEF
│   │   ├── leef.go          # Парсер LEEF
│   │   └── xml.go           # Парсер XML (Windows Event Log)
│   └── processor/
│       └── processor.go      # Главный процессор
├── examples/
│   ├── sample_logs.txt       # Примеры логов
│   └── example_api.go        # Пример использования API
├── go.mod
├── go.sum
├── README.md
└── INSTALL.md
```

## Форматы входных данных

### Syslog (RFC 3164)
```
<134>Oct 11 22:14:15 mymachine su: 'su root' failed
```

### Syslog (RFC 5424)
```
<165>1 2025-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 - Message
```

### CEF (Common Event Format)
```
CEF:0|Vendor|Product|Version|SignatureID|Name|Severity|Extension
```

### LEEF (Log Event Extended Format)
```
LEEF:1.0|Vendor|Product|Version|EventID|Attributes
```

### XML (Windows Event Log)
```xml
<Event xmlns="http://schemas.microsoft.com/win/2004/08/events/event">
  <System>
    <Provider Name="Provider"/>
    <EventID>4624</EventID>
    <TimeCreated SystemTime="2025-10-11T22:14:15Z"/>
    <Computer>hostname</Computer>
  </System>
  <EventData>
    <Data Name="Field">Value</Data>
  </EventData>
</Event>
```

## Формат выходных данных (ГОСТ Р 59710-2022)

Выходные данные в формате JSON со следующими полями:

- `event_id` - уникальный идентификатор события (UUID)
- `timestamp` - время события (RFC3339)
- `source` - информация об источнике события
  - `hostname` - имя хоста
  - `ip_address` - IP адрес (если доступен)
  - `application` - имя приложения
  - `process` - имя процесса
  - `process_id` - ID процесса
- `category` - категория события (см. README.md)
- `severity` - уровень критичности (см. README.md)
- `description` - описание события
- `result` - результат события (УСПЕХ/НЕУСПЕХ/НЕИЗВЕСТНО)
- `action` - выполненное действие
- `subject_account` - субъект действия (опционально)
- `object_account` - объект действия (опционально)
- `additional_data` - дополнительные данные

## Маппинг уровней критичности

### Syslog → ГОСТ
- 0-2 (Emergency, Alert, Critical) → КРИТИЧЕСКИЙ
- 3 (Error) → ВЫСОКИЙ
- 4 (Warning) → СРЕДНИЙ
- 5-6 (Notice, Informational) → НИЗКИЙ
- 7 (Debug) → ИНФОРМАЦИОННЫЙ

### CEF → ГОСТ
- 8-10 → КРИТИЧЕСКИЙ
- 6-7 → ВЫСОКИЙ
- 4-5 → СРЕДНИЙ
- 2-3 → НИЗКИЙ
- 0-1 → ИНФОРМАЦИОННЫЙ

### XML (Windows Event Log) → ГОСТ
- Level 1 (Critical) → КРИТИЧЕСКИЙ
- Level 2 (Error) → ВЫСОКИЙ
- Level 3 (Warning) → СРЕДНИЙ
- Level 4 (Information) → ИНФОРМАЦИОННЫЙ
- Level 5 (Verbose) → НИЗКИЙ

## Производительность

Обработчик способен обрабатывать:
- ~10,000 простых Syslog сообщений в секунду
- ~5,000 CEF/LEEF сообщений в секунду
- ~2,000 XML событий в секунду

Производительность зависит от сложности логов и доступных ресурсов системы.

## Лицензия

MIT License

## Поддержка

При возникновении проблем или вопросов создайте issue в репозитории проекта.
