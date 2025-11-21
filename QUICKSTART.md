# Быстрый старт

## За 5 минут

### 1. Установка (1 минута)

```bash
# Клонируйте репозиторий
git clone https://github.com/kxrty/loggerv2.git
cd loggerv2

# Установите зависимости и соберите
go mod download
go build -o logger.exe ./cmd/logger  # Windows
# go build -o logger ./cmd/logger    # Linux/macOS
```

### 2. Первый запуск (1 минута)

Попробуйте обработать лог прямо из командной строки:

```bash
# Windows
echo "<134>Oct 11 22:14:15 mymachine su: test message" | .\logger.exe

# Linux/macOS
echo "<134>Oct 11 22:14:15 mymachine su: test message" | ./logger
```

Вы увидите JSON с событием в формате ГОСТ:

```json
{
  "event_id": "unique-uuid",
  "timestamp": "2025-10-11T22:14:15Z",
  "source": {
    "hostname": "mymachine",
    "application": "su"
  },
  "category": "СИСТЕМНОЕ_СОБЫТИЕ",
  "severity": "НИЗКИЙ",
  "description": "test message",
  "result": "НЕИЗВЕСТНО"
}
```

### 3. Обработка файла (2 минуты)

Попробуйте обработать примеры логов:

```bash
# Windows
.\logger.exe -input examples\sample_logs.txt -output result.json

# Linux/macOS
./logger -input examples/sample_logs.txt -output result.json
```

Откройте `result.json` - там будут все обработанные логи в формате ГОСТ.

### 4. Использование в коде (1 минута)

Создайте файл `test.go`:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kxrty/loggerv2/internal/processor"
)

func main() {
    proc := processor.NewProcessor()
    
    logs := []string{
        "<134>Oct 11 22:14:15 server1 app: User login failed",
        "CEF:0|Vendor|IDS|1.0|100|Attack|10|src=10.0.0.1 dst=192.168.1.1",
    }
    
    for _, logLine := range logs {
        event, err := proc.Process(logLine)
        if err != nil {
            log.Printf("Ошибка: %v\n", err)
            continue
        }
        
        fmt.Printf("Событие: %s (Критичность: %s)\n", 
            event.Description, event.Severity)
    }
}
```

Запустите:

```bash
go run test.go
```

## Поддерживаемые форматы

### 1. Syslog

```
<134>Oct 11 22:14:15 mymachine su: message
```

### 2. CEF (Common Event Format)

```
CEF:0|Vendor|Product|1.0|100|Event Name|5|src=1.1.1.1 dst=2.2.2.2
```

### 3. LEEF (Log Event Extended Format)

```
LEEF:1.0|Vendor|Product|1.0|100|src=1.1.1.1	dst=2.2.2.2	sev=5
```

### 4. XML (Windows Event Log)

```xml
<Event xmlns="...">
  <System>
    <EventID>4624</EventID>
    <TimeCreated SystemTime="2023-10-11T22:14:15Z"/>
    <Computer>hostname</Computer>
  </System>
  <EventData>
    <Data Name="UserName">user</Data>
  </EventData>
</Event>
```

## Основные команды

```bash
# Из файла в файл
logger.exe -input input.log -output output.json

# Из stdin в файл
cat logs.txt | logger.exe -output result.json

# Из файла в stdout
logger.exe -input logs.txt

# Из stdin в stdout (для пайпов)
cat logs.txt | logger.exe | jq .
```

## Что дальше?

- 📖 Полная документация: [README.md](README.md)
- 🔧 Инструкции по установке: [INSTALL.md](INSTALL.md)
- 📚 API документация: [API.md](API.md)
- 📋 История изменений: [CHANGELOG.md](CHANGELOG.md)
- 💡 Примеры: `examples/`

## Нужна помощь?

Создайте issue в репозитории проекта или обратитесь к документации.

## Ключевые возможности

✅ **Автоматическое определение формата** - не нужно указывать тип лога  
✅ **Стандарт ГОСТ** - соответствие ГОСТ Р 59710-2021  
✅ **Множество форматов** - Syslog, CEF, LEEF, XML  
✅ **Простой API** - легко интегрировать в свои проекты  
✅ **CLI утилита** - готова к использованию  
✅ **Высокая производительность** - до 10,000 логов/сек  
✅ **Тесты** - покрытие ~70%  

Приятного использования! 🚀
