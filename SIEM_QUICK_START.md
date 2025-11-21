# Быстрый старт: Интеграция с SIEM

## ⚡ За 10 минут

### Шаг 1: Выберите метод интеграции

#### **Вариант A: Syslog (Универсальный)**
Работает с любым SIEM. Самый простой способ.

```go
forwarder, _ := siem.NewSyslogForwarder("your-siem.com", 514, "udp")
event, _ := proc.Process(logLine)
forwarder.Forward(event)
```

#### **Вариант B: HTTP API (Splunk/Elastic)**
Для современных SIEM с REST API.

```go
forwarder := siem.NewHTTPForwarder(
    "https://splunk:8088/services/collector/event",
    "YOUR-TOKEN",
    nil,
)
forwarder.Forward(event)
```

#### **Вариант C: Файл (Агентская модель)**
Для SIEM агентов (Filebeat, Universal Forwarder).

```bash
logger -input app.log -output normalized.json
# SIEM агент читает normalized.json
```

### Шаг 2: Настройте SIEM

#### Splunk

**inputs.conf:**
```ini
[udp://514]
sourcetype = gost:json
index = security
```

**Или HEC:**
```bash
curl -k https://splunk:8088/services/collector/event \
  -H "Authorization: Splunk YOUR-TOKEN" \
  -d '{"event": {...}}'
```

#### QRadar

1. Log Sources → Add Log Source
2. Log Source Type: Syslog
3. Protocol: UDP
4. Port: 514

#### Elastic

**filebeat.yml:**
```yaml
filebeat.inputs:
- type: log
  paths: ["/var/log/gost/*.json"]
  json.keys_under_root: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

### Шаг 3: Запустите

```bash
# Демо интеграции
go run examples/siem_integration/main.go

# Реальное использование
tail -f /var/log/app.log | logger | siem-forwarder
```

## 🎯 Готовые решения

### Решение 1: Real-time в Splunk HEC

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
    fwd := siem.NewHTTPForwarder(
        "https://splunk:8088/services/collector/event",
        os.Getenv("SPLUNK_TOKEN"),
        nil,
    )
    
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        event, _ := proc.Process(scanner.Text())
        fwd.Forward(event)
    }
}
```

**Запуск:**
```bash
tail -f /var/log/app.log | go run main.go
```

### Решение 2: Batch в QRadar

```go
package main

import (
    "time"
    "github.com/kxrty/loggerv2/internal/processor"
    "github.com/kxrty/loggerv2/internal/siem"
)

func main() {
    proc := processor.NewProcessor()
    fwd, _ := siem.NewSyslogForwarder("qradar", 514, "tcp")
    defer fwd.Close()
    
    buffer := make([]*models.GOSTEvent, 0, 100)
    ticker := time.NewTicker(5 * time.Second)
    
    for {
        select {
        case logLine := <-logChan:
            event, _ := proc.Process(logLine)
            buffer = append(buffer, event)
            
            if len(buffer) >= 100 {
                fwd.ForwardBatch(buffer)
                buffer = buffer[:0]
            }
            
        case <-ticker.C:
            if len(buffer) > 0 {
                fwd.ForwardBatch(buffer)
                buffer = buffer[:0]
            }
        }
    }
}
```

### Решение 3: Файловая интеграция

```bash
#!/bin/bash
# process-to-file.sh

INPUT="/var/log/app.log"
OUTPUT="/var/log/gost/normalized.json"

tail -f $INPUT | logger -output $OUTPUT
```

**Настройка Filebeat:**
```yaml
filebeat.inputs:
- type: log
  paths: ["/var/log/gost/normalized.json"]
  json.keys_under_root: true
```

## 📊 Архитектуры

### Простая (1-1000 событий/сек)
```
App → LoggerV2 → SIEM
```

### Средняя (1000-10000 событий/сек)
```
Apps → LoggerV2 → File → SIEM Agent → SIEM
```

### Высоконагруженная (10000+ событий/сек)
```
Apps → LoggerV2 → Kafka → LoggerV2 Consumer → SIEM
```

## ✅ Чек-лист интеграции

- [ ] Выбран метод интеграции
- [ ] SIEM настроен для приема логов
- [ ] LoggerV2 обрабатывает логи корректно
- [ ] Тестовые события появляются в SIEM
- [ ] Настроены правила парсинга в SIEM
- [ ] Проверена производительность
- [ ] Настроен мониторинг ошибок
- [ ] Документирована конфигурация

## 🎓 Поддерживаемые SIEM

| SIEM | Метод | Сложность | Рекомендация |
|------|-------|-----------|--------------|
| Splunk | HEC/File | ⭐ Легко | HTTP API |
| Elastic | HTTP/File | ⭐ Легко | Filebeat |
| QRadar | Syslog | ⭐⭐ Средне | TCP Syslog |
| ArcSight | Syslog | ⭐⭐ Средне | Syslog |
| MaxPatrol | Syslog | ⭐⭐ Средне | TCP Syslog |
| Graylog | Syslog/GELF | ⭐ Легко | UDP Syslog |

## 🔧 Устранение проблем

### Логи не появляются в SIEM

1. **Проверьте соединение:**
```bash
nc -zv siem.example.com 514
```

2. **Проверьте формат:**
```bash
echo "test log" | logger | head -5
```

3. **Проверьте SIEM парсер:**
- Убедитесь что sourcetype настроен
- Проверьте что JSON парсится корректно

### Низкая производительность

1. **Используйте batch обработку:**
```go
events, _ := proc.ProcessBatch(logs)
forwarder.ForwardBatch(events)
```

2. **Настройте буферизацию:**
```go
buffer := make([]*models.GOSTEvent, 0, 1000)
```

3. **Используйте TCP вместо UDP для Syslog**

## 📚 Дополнительно

- Полная документация: [SIEM_INTEGRATION.md](SIEM_INTEGRATION.md)
- API документация: [API.md](API.md)
- Примеры: [examples/siem_integration/](examples/siem_integration/)

---

**Нужна помощь?** Создайте [Issue](../../issues) с тегом "siem-integration"
