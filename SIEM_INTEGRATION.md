# Интеграция LoggerV2 с SIEM системами

## 📋 Содержание

- [Обзор](#обзор)
- [Способы интеграции](#способы-интеграции)
- [Поддерживаемые SIEM](#поддерживаемые-siem)
- [Примеры интеграции](#примеры-интеграции)
- [Лучшие практики](#лучшие-практики)

## 🎯 Обзор

LoggerV2 идеально подходит для интеграции с SIEM системами благодаря:

✅ **Нормализация данных** - приведение к единому формату ГОСТ  
✅ **Структурированный вывод** - JSON с полной информацией  
✅ **Автоматическое обогащение** - извлечение пользователей, IP, действий  
✅ **Категоризация** - автоматическое определение типа события  
✅ **Маппинг критичности** - преобразование в стандартные уровни  

## 🔌 Способы интеграции

### 1. Syslog Forwarder (Рекомендуется)

**Преимущества:**
- Стандартный протокол, поддерживаемый всеми SIEM
- Высокая надежность
- Низкие накладные расходы

**Как работает:**
```go
package main

import (
    "github.com/kxrty/loggerv2/internal/processor"
    "github.com/kxrty/loggerv2/internal/siem"
)

func main() {
    // Создаем процессор
    proc := processor.NewProcessor()
    
    // Подключаемся к SIEM
    forwarder, err := siem.NewSyslogForwarder("siem.example.com", 514, "udp")
    if err != nil {
        panic(err)
    }
    defer forwarder.Close()
    
    // Обрабатываем и отправляем лог
    event, _ := proc.Process("<134>Oct 11 22:14:15 server su: failed login")
    forwarder.Forward(event)
}
```

**Конфигурация SIEM:**
- **Порт:** 514 (UDP) или 6514 (TCP/TLS)
- **Формат:** RFC 5424 с JSON payload
- **Парсер:** JSON в structured data

### 2. HTTP/HTTPS API

**Используется для:**
- Splunk HTTP Event Collector (HEC)
- Elastic Stack
- Современные cloud SIEM

**Пример:**
```go
package main

import (
    "github.com/kxrty/loggerv2/internal/processor"
    "github.com/kxrty/loggerv2/internal/siem"
)

func main() {
    proc := processor.NewProcessor()
    
    // Для Splunk HEC
    forwarder := siem.NewHTTPForwarder(
        "https://splunk.example.com:8088/services/collector/event",
        "YOUR-HEC-TOKEN",
        map[string]string{
            "X-Splunk-Request-Channel": "unique-guid",
        },
    )
    
    event, _ := proc.Process(logLine)
    forwarder.Forward(event)
}
```

### 3. File Monitoring

**Подходит для:**
- Legacy SIEM без API
- Агентская модель (Filebeat, Universal Forwarder)
- Оффлайн обработка

**Схема работы:**
```
Логи → LoggerV2 → normalized_logs.json → SIEM Agent → SIEM
```

**Пример:**
```bash
# Непрерывная обработка логов в файл
tail -f /var/log/app.log | logger | tee normalized.json

# SIEM агент мониторит normalized.json
```

### 4. Kafka/Message Queue

**Для высоконагруженных систем:**
```go
// Отправка в Kafka
producer.Send("gost-events", event)
```

## 🏢 Поддерживаемые SIEM

### 1. Splunk

**Метод 1: HTTP Event Collector (HEC)**

```go
forwarder := siem.NewHTTPForwarder(
    "https://splunk.example.com:8088/services/collector/event",
    "YOUR-HEC-TOKEN",
    nil,
)
```

**Метод 2: Universal Forwarder**
```bash
# LoggerV2 пишет в файл
logger -input app.log -output /var/log/gost/normalized.json

# Splunk UF мониторит файл
[monitor:///var/log/gost/normalized.json]
sourcetype = gost:json
index = security
```

**Метод 3: Syslog**
```bash
# inputs.conf
[udp://514]
sourcetype = syslog
```

### 2. Elastic Stack (ELK)

**Через Filebeat:**

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/gost/*.json
  json.keys_under_root: true
  json.add_error_key: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  index: "gost-events-%{+yyyy.MM.dd}"
```

**Через HTTP API:**
```go
forwarder := siem.NewHTTPForwarder(
    "https://elasticsearch:9200/gost-events/_doc",
    "",
    map[string]string{
        "Content-Type": "application/json",
    },
)
```

### 3. IBM QRadar

**Через Syslog:**
```go
forwarder, _ := siem.NewSyslogForwarder("qradar.example.com", 514, "udp")
```

**QRadar Log Source:**
- Protocol: Syslog
- Log Source Type: Custom JSON
- Parser: JSON

### 4. ArcSight

**Через Syslog/CEF:**
```go
// LoggerV2 может отправлять в CEF формате обратно
forwarder, _ := siem.NewSyslogForwarder("arcsight.example.com", 514, "tcp")
```

**ArcSight Connector:**
- Type: Syslog File
- Format: CEF or JSON

### 5. MaxPatrol SIEM (Positive Technologies)

**Через Syslog:**
```go
forwarder, _ := siem.NewSyslogForwarder("maxpatrol.example.com", 514, "tcp")
```

**Настройка источника:**
- Тип: Syslog
- Формат: JSON в structured data
- Парсер: Пользовательский JSON парсер

### 6. R-Vision SOAR

**Через REST API:**
```go
forwarder := siem.NewHTTPForwarder(
    "https://rvision.example.com/api/events",
    "API-TOKEN",
    map[string]string{
        "X-API-Version": "1.0",
    },
)
```

### 7. Graylog

**Через GELF:**
```bash
# Конвертация JSON в GELF
logger -input app.log | jq -c '{version:"1.1",host:.source.hostname,short_message:.description,full_message:.,level:1}' | nc graylog 12201
```

**Через Syslog:**
```go
forwarder, _ := siem.NewSyslogForwarder("graylog.example.com", 514, "udp")
```

## 📝 Примеры интеграции

### Пример 1: Real-time обработка с отправкой в Splunk

```go
package main

import (
    "bufio"
    "log"
    "os"
    
    "github.com/kxrty/loggerv2/internal/processor"
    "github.com/kxrty/loggerv2/internal/siem"
)

func main() {
    proc := processor.NewProcessor()
    
    forwarder := siem.NewHTTPForwarder(
        "https://splunk:8088/services/collector/event",
        os.Getenv("SPLUNK_HEC_TOKEN"),
        nil,
    )
    
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        event, err := proc.Process(scanner.Text())
        if err != nil {
            log.Printf("Ошибка: %v\n", err)
            continue
        }
        
        if err := forwarder.Forward(event); err != nil {
            log.Printf("Ошибка отправки: %v\n", err)
        }
    }
}
```

### Пример 2: Batch обработка с отправкой в QRadar

```go
package main

import (
    "bufio"
    "os"
    
    "github.com/kxrty/loggerv2/internal/processor"
    "github.com/kxrty/loggerv2/internal/siem"
)

func main() {
    proc := processor.NewProcessor()
    
    forwarder, _ := siem.NewSyslogForwarder("qradar", 514, "tcp")
    defer forwarder.Close()
    
    file, _ := os.Open("logs.txt")
    defer file.Close()
    
    var logs []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        logs = append(logs, scanner.Text())
        
        // Отправляем пакетами по 100
        if len(logs) >= 100 {
            events, _ := proc.ProcessBatch(logs)
            forwarder.ForwardBatch(events)
            logs = logs[:0]
        }
    }
    
    // Отправляем остаток
    if len(logs) > 0 {
        events, _ := proc.ProcessBatch(logs)
        forwarder.ForwardBatch(events)
    }
}
```

### Пример 3: Микросервис для обработки логов

```go
package main

import (
    "encoding/json"
    "net/http"
    
    "github.com/kxrty/loggerv2/internal/processor"
    "github.com/kxrty/loggerv2/internal/siem"
)

var (
    proc      = processor.NewProcessor()
    forwarder *siem.SyslogForwarder
)

func init() {
    forwarder, _ = siem.NewSyslogForwarder("siem", 514, "udp")
}

func handleLog(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Log string `json:"log"`
    }
    
    json.NewDecoder(r.Body).Decode(&req)
    
    event, err := proc.Process(req.Log)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    
    if err := forwarder.Forward(event); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    json.NewEncoder(w).Encode(event)
}

func main() {
    http.HandleFunc("/log", handleLog)
    http.ListenAndServe(":8080", nil)
}
```

## 🎯 Лучшие практики

### 1. Надежность

```go
// Используйте retry механизм
for i := 0; i < 3; i++ {
    if err := forwarder.Forward(event); err == nil {
        break
    }
    time.Sleep(time.Second * time.Duration(i+1))
}
```

### 2. Производительность

```go
// Используйте буферизацию
buffer := make([]*models.GOSTEvent, 0, 100)
for event := range eventChan {
    buffer = append(buffer, event)
    if len(buffer) >= 100 {
        forwarder.ForwardBatch(buffer)
        buffer = buffer[:0]
    }
}
```

### 3. Мониторинг

```go
// Логируйте статистику
var (
    processed int64
    errors    int64
)

// Периодически выводите метрики
go func() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        log.Printf("Обработано: %d, Ошибок: %d\n", processed, errors)
    }
}()
```

### 4. Безопасность

```go
// Используйте TLS для Syslog
forwarder, _ := siem.NewSyslogForwarder("siem", 6514, "tcp-tls")

// Храните токены в переменных окружения
token := os.Getenv("SIEM_TOKEN")
forwarder := siem.NewHTTPForwarder(url, token, nil)
```

## 📊 Архитектуры развертывания

### Архитектура 1: Прямая интеграция

```
Источники логов → LoggerV2 → SIEM
```

Подходит для: небольших объемов, простых сценариев

### Архитектура 2: Агентская модель

```
Источники логов → LoggerV2 → Файл → SIEM Agent → SIEM
```

Подходит для: distributed deployments, fault tolerance

### Архитектура 3: Message Queue

```
Источники логов → LoggerV2 → Kafka → SIEM Consumer → SIEM
```

Подходит для: высокие нагрузки, горизонтальное масштабирование

### Архитектура 4: API Gateway

```
Источники логов → LoggerV2 HTTP Service → Load Balancer → SIEM
```

Подходит для: централизованная обработка, множество источников

## 🚀 Запуск примера

```bash
# Запустить пример интеграции
go run examples/siem_integration/main.go

# Построить микросервис
go build -o siem-forwarder examples/siem_integration/main.go
./siem-forwarder
```

## 📞 Поддержка

Для помощи с интеграцией:
- См. [API.md](API.md) для деталей API
- Создайте [Issue](../../issues) для вопросов
- Проверьте [примеры](examples/) для кода

---

**Совместимость:** LoggerV2 совместим со всеми SIEM, поддерживающими Syslog, HTTP API или file monitoring.
