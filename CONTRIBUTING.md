# Руководство для контрибьюторов

Спасибо за интерес к проекту LoggerV2! Это руководство поможет вам внести свой вклад.

## 📋 Содержание

- [Как начать](#как-начать)
- [Стандарты кодирования](#стандарты-кодирования)
- [Добавление нового парсера](#добавление-нового-парсера)
- [Тестирование](#тестирование)
- [Pull Request процесс](#pull-request-процесс)
- [Сообщение об ошибках](#сообщение-об-ошибках)

## 🚀 Как начать

### 1. Форк и клонирование

```bash
# Форкните репозиторий через GitHub UI
# Затем клонируйте свой форк
git clone https://github.com/YOUR_USERNAME/loggerv2.git
cd loggerv2
```

### 2. Настройка окружения

```bash
# Установите зависимости
go mod download

# Убедитесь что тесты проходят
go test ./...
```

### 3. Создание ветки

```bash
# Создайте ветку для вашей работы
git checkout -b feature/my-new-feature
# или
git checkout -b fix/bug-description
```

## 📝 Стандарты кодирования

### Форматирование

```bash
# Всегда форматируйте код перед коммитом
go fmt ./...
```

### Именование

- **Переменные**: camelCase для локальных, PascalCase для экспортируемых
- **Функции**: PascalCase для экспортируемых, camelCase для приватных
- **Константы**: PascalCase с префиксом категории
- **Файлы**: snake_case.go

### Комментарии

```go
// Экспортируемые функции должны иметь документацию
// Parse парсит CEF сообщение и возвращает GOSTEvent
func Parse(logLine string) (*models.GOSTEvent, error) {
    // Внутренние комментарии для сложной логики
    if match := pattern.FindStringSubmatch(logLine); match != nil {
        return parseMatch(match)
    }
    return nil, errors.New("invalid format")
}
```

### Обработка ошибок

```go
// Хорошо
event, err := parser.Parse(logLine)
if err != nil {
    return nil, fmt.Errorf("failed to parse log: %w", err)
}

// Плохо
event, _ := parser.Parse(logLine)  // Игнорирование ошибок
```

## 🔧 Добавление нового парсера

### Шаг 1: Создайте файл парсера

```go
// internal/parser/myformat.go
package parser

import (
    "github.com/kxrty/loggerv2/internal/models"
    "github.com/google/uuid"
)

type MyFormatParser struct{}

func NewMyFormatParser() *MyFormatParser {
    return &MyFormatParser{}
}

// Parse парсит MyFormat сообщение и возвращает GOSTEvent
func (p *MyFormatParser) Parse(logLine string) (*models.GOSTEvent, error) {
    // Ваша логика парсинга
    
    event := &models.GOSTEvent{
        EventID:     uuid.New().String(),
        Timestamp:   time.Now(),
        Description: "parsed message",
        Source: models.Source{
            Hostname: "hostname",
        },
        Severity:       models.SeverityInfo,
        Category:       models.CategorySystemEvent,
        Result:         models.ResultUnknown,
        AdditionalData: make(map[string]interface{}),
    }
    
    return event, nil
}
```

### Шаг 2: Создайте тесты

```go
// internal/parser/myformat_test.go
package parser

import (
    "testing"
)

func TestMyFormatParser_Parse(t *testing.T) {
    parser := NewMyFormatParser()
    
    tests := []struct {
        name    string
        logLine string
        wantErr bool
    }{
        {
            name:    "Valid log",
            logLine: "valid log format here",
            wantErr: false,
        },
        {
            name:    "Invalid log",
            logLine: "invalid format",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            event, err := parser.Parse(tt.logLine)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && event == nil {
                t.Error("Expected non-nil event")
            }
        })
    }
}
```

### Шаг 3: Интегрируйте в процессор

```go
// internal/processor/processor.go

// Добавьте новый тип лога
const (
    LogTypeUnknown LogType = iota
    LogTypeSyslog
    LogTypeCEF
    LogTypeLEEF
    LogTypeXML
    LogTypeMyFormat  // <- Новый тип
)

// Добавьте парсер в структуру
type Processor struct {
    syslogParser   *parser.SyslogParser
    cefParser      *parser.CEFParser
    leefParser     *parser.LEEFParser
    xmlParser      *parser.XMLParser
    myFormatParser *parser.MyFormatParser  // <- Новый парсер
}

// Инициализируйте парсер
func NewProcessor() *Processor {
    return &Processor{
        syslogParser:   parser.NewSyslogParser(),
        cefParser:      parser.NewCEFParser(),
        leefParser:     parser.NewLEEFParser(),
        xmlParser:      parser.NewXMLParser(),
        myFormatParser: parser.NewMyFormatParser(),  // <- Инициализация
    }
}

// Добавьте обработку в DetectLogType
func (p *Processor) DetectLogType(logLine string) LogType {
    logLine = strings.TrimSpace(logLine)
    
    // Добавьте проверку вашего формата
    if strings.HasPrefix(logLine, "MYFORMAT:") {
        return LogTypeMyFormat
    }
    
    // ... остальные проверки
}

// Добавьте case в Process
func (p *Processor) Process(logLine string) (*models.GOSTEvent, error) {
    logType := p.DetectLogType(logLine)
    
    switch logType {
    case LogTypeSyslog:
        return p.syslogParser.Parse(logLine)
    case LogTypeCEF:
        return p.cefParser.Parse(logLine)
    case LogTypeLEEF:
        return p.leefParser.Parse(logLine)
    case LogTypeXML:
        return p.xmlParser.Parse(logLine)
    case LogTypeMyFormat:  // <- Новый case
        return p.myFormatParser.Parse(logLine)
    default:
        return nil, fmt.Errorf("неизвестный тип лога")
    }
}
```

### Шаг 4: Обновите документацию

Добавьте информацию о новом формате в:
- `README.md` - в раздел "Поддерживаемые форматы"
- `INSTALL.md` - примеры использования
- `CHANGELOG.md` - в раздел "Added"

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
go test ./...

# Конкретный пакет
go test ./internal/parser -v

# С покрытием
go test ./... -cover
```

### Требования к тестам

- Минимальное покрытие: 60%
- Тесты для всех публичных функций
- Тесты граничных случаев
- Тесты ошибочных ситуаций

### Пример хорошего теста

```go
func TestParser_EdgeCases(t *testing.T) {
    parser := NewMyParser()
    
    tests := []struct {
        name        string
        input       string
        wantErr     bool
        checkResult func(*testing.T, *models.GOSTEvent)
    }{
        {
            name:    "Empty string",
            input:   "",
            wantErr: true,
        },
        {
            name:    "Very long message",
            input:   strings.Repeat("a", 10000),
            wantErr: false,
            checkResult: func(t *testing.T, event *models.GOSTEvent) {
                if len(event.Description) == 0 {
                    t.Error("Description should not be empty")
                }
            },
        },
        {
            name:    "Special characters",
            input:   "test\x00\xff",
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            event, err := parser.Parse(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if tt.checkResult != nil && event != nil {
                tt.checkResult(t, event)
            }
        })
    }
}
```

## 📨 Pull Request процесс

### 1. Убедитесь что все проходит

```bash
# Форматирование
go fmt ./...

# Тесты
go test ./...

# Vet
go vet ./...
```

### 2. Коммит изменений

```bash
git add .
git commit -m "feat: добавлен парсер MyFormat"
```

### Стиль коммитов

Используйте [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - новая функциональность
- `fix:` - исправление бага
- `docs:` - изменения в документации
- `test:` - добавление/изменение тестов
- `refactor:` - рефакторинг кода
- `perf:` - улучшение производительности
- `chore:` - рутинные задачи

Примеры:
```
feat: добавлен парсер JSON логов
fix: исправлена обработка пустых строк в CEF
docs: обновлена документация API
test: добавлены тесты для XML парсера
```

### 3. Push и создание PR

```bash
git push origin feature/my-new-feature
```

Создайте Pull Request через GitHub UI и заполните шаблон:

```markdown
## Описание
Краткое описание изменений

## Тип изменения
- [ ] Новая функциональность
- [ ] Исправление бага
- [ ] Улучшение документации
- [ ] Рефакторинг

## Тестирование
- [ ] Добавлены unit-тесты
- [ ] Все тесты проходят
- [ ] Покрытие не уменьшилось

## Чеклист
- [ ] Код отформатирован (go fmt)
- [ ] Нет ошибок (go vet)
- [ ] Документация обновлена
- [ ] CHANGELOG.md обновлен
```

## 🐛 Сообщение об ошибках

### Хороший отчет об ошибке включает:

1. **Описание проблемы**
   - Что произошло
   - Что ожидалось

2. **Шаги для воспроизведения**
   ```
   1. Запустите logger.exe
   2. Передайте лог формата X
   3. Наблюдайте ошибку Y
   ```

3. **Окружение**
   - ОС: Windows 10 / Linux Ubuntu 20.04
   - Go версия: 1.21.0
   - Версия проекта: 1.0.0

4. **Минимальный пример**
   ```go
   // Код который воспроизводит проблему
   parser := NewParser()
   event, err := parser.Parse("...")
   ```

5. **Логи/Output**
   ```
   Вывод ошибки
   ```

## 📚 Дополнительные ресурсы

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [ГОСТ Р 59710-2021](https://docs.cntd.ru/document/1200179861) - стандарт логирования

## 💬 Связь

- GitHub Issues - для багов и feature requests
- GitHub Discussions - для вопросов и обсуждений

## 🙏 Благодарности

Спасибо за ваш вклад в проект!

---

**Важно:** 
- Следуйте стандартам кодирования Go
- Пишите тесты для нового кода
- Обновляйте документацию
- Будьте вежливы и конструктивны в обсуждениях
