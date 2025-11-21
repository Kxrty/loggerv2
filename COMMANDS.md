# Справочник команд LoggerV2

## 🔨 Сборка

### Сборка CLI приложения
```bash
# Windows
go build -o logger.exe ./cmd/logger

# Linux/macOS
go build -o logger ./cmd/logger
```

### Сборка с оптимизацией
```bash
# Уменьшенный размер бинарника
go build -ldflags="-s -w" -o logger.exe ./cmd/logger
```

## 🧪 Тестирование

### Запуск всех тестов
```bash
go test ./...
```

### Тесты с покрытием
```bash
go test ./... -cover
```

### Тесты с подробным выводом
```bash
go test ./... -v
```

### Тесты конкретного пакета
```bash
go test ./internal/parser -v
go test ./internal/processor -v
```

### Тесты с генерацией отчета покрытия
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🚀 Использование CLI

### Обработка из файла в файл
```bash
# Windows
.\logger.exe -input logs.txt -output result.json

# Linux/macOS
./logger -input logs.txt -output result.json
```

### Обработка из stdin в stdout
```bash
# Windows
echo "<134>Oct 11 22:14:15 server su: test" | .\logger.exe

# Linux/macOS
echo "<134>Oct 11 22:14:15 server su: test" | ./logger
```

### Обработка с pipe
```bash
# Windows
type logs.txt | .\logger.exe > output.json

# Linux/macOS
cat logs.txt | ./logger > output.json
```

### Обработка из файла в stdout
```bash
# Windows
.\logger.exe -input logs.txt

# Linux/macOS
./logger -input logs.txt
```

### Обработка с использованием jq
```bash
# Фильтрация только критических событий
cat logs.txt | ./logger | jq 'select(.severity=="КРИТИЧЕСКИЙ")'

# Вывод только описаний
cat logs.txt | ./logger | jq '.description'

# Подсчет событий по категориям
cat logs.txt | ./logger | jq -s 'group_by(.category) | map({category: .[0].category, count: length})'
```

## 💻 Запуск примеров

### API Example
```bash
go run examples/api_example/main.go
```

### HTTP Server
```bash
# Запуск сервера
go run examples/http_server/main.go

# В другом терминале - тестирование
curl http://localhost:8080/health

curl -X POST http://localhost:8080/process-single \
  -H "Content-Type: text/plain" \
  -d "<134>Oct 11 22:14:15 server su: test"

curl -X POST http://localhost:8080/process \
  -H "Content-Type: application/json" \
  -d '{"logs": ["<134>Oct 11 22:14:15 server su: test1", "CEF:0|Vendor|Product|1.0|100|Event|5|src=1.1.1.1"]}'
```

## 📦 Управление зависимостями

### Установка зависимостей
```bash
go mod download
```

### Обновление зависимостей
```bash
go get -u ./...
go mod tidy
```

### Проверка зависимостей
```bash
go list -m all
```

### Очистка кеша модулей
```bash
go clean -modcache
```

## 🔍 Проверка кода

### Форматирование кода
```bash
go fmt ./...
```

### Линтинг (требует golangci-lint)
```bash
golangci-lint run
```

### Проверка на ошибки
```bash
go vet ./...
```

## 📊 Бенчмарки

### Создание бенчмарка
```go
// internal/parser/syslog_bench_test.go
func BenchmarkSyslogParser_Parse(b *testing.B) {
    parser := NewSyslogParser()
    logLine := "<134>Oct 11 22:14:15 mymachine su: test"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        parser.Parse(logLine)
    }
}
```

### Запуск бенчмарков
```bash
go test -bench=. ./internal/parser
go test -bench=. -benchmem ./internal/parser
```

## 🛠️ Разработка

### Добавление нового парсера

1. Создайте файл парсера:
```bash
# Windows
New-Item internal/parser/myformat.go

# Linux/macOS
touch internal/parser/myformat.go
```

2. Создайте тестовый файл:
```bash
# Windows
New-Item internal/parser/myformat_test.go

# Linux/macOS
touch internal/parser/myformat_test.go
```

3. Обновите processor.go для поддержки нового формата

4. Запустите тесты:
```bash
go test ./internal/parser -v
```

## 📝 Генерация документации

### GoDoc локально
```bash
godoc -http=:6060
# Откройте http://localhost:6060/pkg/github.com/kxrty/loggerv2/
```

## 🧹 Очистка

### Удаление бинарников
```bash
# Windows
Remove-Item logger.exe -ErrorAction SilentlyContinue

# Linux/macOS
rm -f logger
```

### Полная очистка
```bash
go clean -cache -testcache -modcache
```

## 📋 Создание релиза

### Сборка для разных платформ
```bash
# Windows amd64
GOOS=windows GOARCH=amd64 go build -o logger-windows-amd64.exe ./cmd/logger

# Linux amd64
GOOS=linux GOARCH=amd64 go build -o logger-linux-amd64 ./cmd/logger

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o logger-darwin-amd64 ./cmd/logger

# macOS arm64 (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o logger-darwin-arm64 ./cmd/logger
```

### Создание архива
```bash
# Windows
Compress-Archive -Path logger.exe,README.md,LICENSE -DestinationPath logger-v1.0.0-windows-amd64.zip

# Linux/macOS
tar -czf logger-v1.0.0-linux-amd64.tar.gz logger README.md LICENSE
```

## 🐛 Отладка

### Запуск с race detector
```bash
go test -race ./...
go build -race -o logger ./cmd/logger
```

### Профилирование
```bash
# CPU профиль
go test -cpuprofile=cpu.prof -bench=. ./internal/parser
go tool pprof cpu.prof

# Memory профиль
go test -memprofile=mem.prof -bench=. ./internal/parser
go tool pprof mem.prof
```

## 📊 Статистика кода

### Подсчет строк кода
```bash
# Linux/macOS
find . -name '*.go' -not -path './vendor/*' | xargs wc -l

# Windows PowerShell
Get-ChildItem -Recurse -Filter *.go | Measure-Object -Line | Select-Object Lines
```

### Сложность кода (требует gocyclo)
```bash
gocyclo -over 15 .
```

## 🔐 Безопасность

### Проверка уязвимостей
```bash
go list -json -m all | nancy sleuth
```

### Аудит зависимостей
```bash
go mod verify
```

## 📤 Публикация

### Создание тега версии
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### Публикация на pkg.go.dev
```bash
# После push тега, модуль автоматически появится на pkg.go.dev
# Принудительное обновление:
GOPROXY=proxy.golang.org go list -m github.com/kxrty/loggerv2@v1.0.0
```

## 🆘 Помощь

### Справка по командам
```bash
go help
go help build
go help test
```

### Версия Go
```bash
go version
```

### Информация о окружении
```bash
go env
```

---

**Полезные ссылки:**
- [README.md](README.md) - Основная документация
- [API.md](API.md) - API справочник
- [INSTALL.md](INSTALL.md) - Установка
- [QUICKSTART.md](QUICKSTART.md) - Быстрый старт
