# Примеры использования LoggerV2

## Файлы и папки

- `sample_logs.txt` - Примеры логов различных форматов
- `api_example/main.go` - Пример использования API в коде
- `http_server/main.go` - HTTP сервер для обработки логов через REST API

## 1. Использование CLI

### Базовое использование

```bash
# Обработка файла с логами
logger.exe -input sample_logs.txt -output result.json

# Обработка из stdin
echo "<134>Oct 11 22:14:15 server su: test" | logger.exe

# С использованием pipe
cat sample_logs.txt | logger.exe > output.json
```

## 2. Использование API (api_example/main.go)

```bash
# Запуск примера
go run api_example/main.go
```

Этот пример демонстрирует:
- Создание процессора
- Обработку отдельных логов
- Вывод информации о событиях
- Преобразование в JSON
- Пакетную обработку

## 3. HTTP Сервер (http_server/main.go)

### Запуск сервера

```bash
go run http_server/main.go
```

Сервер запустится на `http://localhost:8080`

### Endpoints

#### GET /health
Проверка работоспособности сервиса.

```bash
curl http://localhost:8080/health
```

Ответ:
```json
{"status": "ok"}
```

#### POST /process
Обработка массива логов.

```bash
curl -X POST http://localhost:8080/process \
  -H "Content-Type: application/json" \
  -d '{
    "logs": [
      "<134>Oct 11 22:14:15 mymachine su: test1",
      "CEF:0|Vendor|Product|1.0|100|Event|5|src=1.1.1.1"
    ]
  }'
```

Ответ:
```json
{
  "success": 2,
  "errors": 0,
  "events": [
    {
      "event_id": "uuid",
      "timestamp": "2025-10-11T22:14:15Z",
      "source": {...},
      "category": "СИСТЕМНОЕ_СОБЫТИЕ",
      "severity": "НИЗКИЙ",
      "description": "test1",
      ...
    }
  ]
}
```

#### POST /process-single
Обработка одного лога (text/plain).

```bash
curl -X POST http://localhost:8080/process-single \
  -H "Content-Type: text/plain" \
  -d "<134>Oct 11 22:14:15 mymachine su: test message"
```

#### POST /detect
Определение типа лога.

```bash
curl -X POST http://localhost:8080/detect \
  -H "Content-Type: text/plain" \
  -d "CEF:0|Vendor|Product|1.0|100|Event|5|src=1.1.1.1"
```

Ответ:
```json
{
  "log_type": "CEF"
}
```

## 4. Интеграция в свой проект

### Установка

```bash
go get github.com/kxrty/loggerv2
```

### Базовый пример

```go
package main

import (
    "fmt"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    event, _ := proc.Process("<134>Oct 11 22:14:15 server su: test")
    fmt.Printf("Событие: %s\n", event.Description)
}
```

### Обработка с проверкой ошибок

```go
package main

import (
    "fmt"
    "log"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    
    logLine := "<134>Oct 11 22:14:15 mymachine su: failed login"
    event, err := proc.Process(logLine)
    if err != nil {
        log.Fatalf("Ошибка обработки: %v", err)
    }
    
    fmt.Printf("ID: %s\n", event.EventID)
    fmt.Printf("Время: %s\n", event.Timestamp)
    fmt.Printf("Хост: %s\n", event.Source.Hostname)
    fmt.Printf("Приложение: %s\n", event.Source.Application)
    fmt.Printf("Категория: %s\n", event.Category)
    fmt.Printf("Критичность: %s\n", event.Severity)
    fmt.Printf("Описание: %s\n", event.Description)
}
```

### Пакетная обработка

```go
package main

import (
    "fmt"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    
    logs := []string{
        "<134>Oct 11 22:14:15 server1 app: event1",
        "CEF:0|Vendor|IDS|1.0|100|Attack|10|src=10.0.0.1",
        "LEEF:1.0|Vendor|Product|1.0|100|src=1.1.1.1\tdst=2.2.2.2",
    }
    
    events, errors := proc.ProcessBatch(logs)
    
    fmt.Printf("Обработано: %d событий\n", len(events))
    fmt.Printf("Ошибок: %d\n", len(errors))
    
    for i, event := range events {
        fmt.Printf("Событие %d: %s (%s)\n", 
            i+1, event.Description, event.Severity)
    }
}
```

### Определение типа лога

```go
package main

import (
    "fmt"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    
    logs := map[string]string{
        "Syslog": "<134>Oct 11 22:14:15 server su: test",
        "CEF":    "CEF:0|Vendor|Product|1.0|100|Event|5|src=1.1.1.1",
        "LEEF":   "LEEF:1.0|Vendor|Product|1.0|100|key=value",
        "XML":    "<Event><System><EventID>100</EventID></System></Event>",
    }
    
    for name, logLine := range logs {
        logType := proc.DetectLogType(logLine)
        fmt.Printf("%s: тип %v\n", name, logType)
    }
}
```

## 5. Тестирование примеров логов

### Syslog RFC 3164
```bash
echo "<134>Oct 11 22:14:15 mymachine su: 'su root' failed for user" | logger.exe
```

### Syslog RFC 5424
```bash
echo "<165>1 2023-10-11T22:14:15.003Z mymachine.example.com app - ID47 - Message" | logger.exe
```

### CEF
```bash
echo "CEF:0|Security|IDS|1.0|100|Attack detected|10|src=10.0.0.1 dst=192.168.1.1 act=blocked" | logger.exe
```

### LEEF
```bash
echo "LEEF:1.0|Microsoft|MSExchange|4.0|15345|src=10.0.0.1	dst=172.50.123.1	sev=5	cat=anomaly" | logger.exe
```

## 6. Производительность

### Бенчмарк обработки

```go
package main

import (
    "fmt"
    "time"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    
    logLine := "<134>Oct 11 22:14:15 mymachine su: test"
    iterations := 10000
    
    start := time.Now()
    for i := 0; i < iterations; i++ {
        proc.Process(logLine)
    }
    duration := time.Since(start)
    
    fmt.Printf("Обработано %d логов за %v\n", iterations, duration)
    fmt.Printf("Скорость: %.0f логов/сек\n", 
        float64(iterations)/duration.Seconds())
}
```

## 7. Интеграция с файлами

### Чтение из файла построчно

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    
    file, err := os.Open("sample_logs.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    successCount := 0
    errorCount := 0
    
    for scanner.Scan() {
        logLine := scanner.Text()
        _, err := proc.Process(logLine)
        if err != nil {
            errorCount++
            fmt.Printf("Ошибка: %v\n", err)
        } else {
            successCount++
        }
    }
    
    fmt.Printf("Успешно: %d, Ошибок: %d\n", successCount, errorCount)
}
```

## Поддержка

Для получения дополнительной информации см.:
- [README.md](../README.md) - Общая документация
- [API.md](../API.md) - API документация
- [INSTALL.md](../INSTALL.md) - Инструкции по установке
