package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kxrty/loggerv2/internal/processor"
	"github.com/kxrty/loggerv2/internal/siem"
)

func main() {
	fmt.Println("=== Пример интеграции с SIEM ===\n")

	// Создаем процессор логов
	proc := processor.NewProcessor()

	// Примеры логов для обработки
	logs := []string{
		"<134>Oct 11 22:14:15 server su: authentication failed for user admin",
		"CEF:0|Security|IDS|1.0|100|Port scan detected|8|src=10.0.0.1 dst=192.168.1.1",
		"LEEF:1.0|IBM|QRadar|7.3|Login|usrName=admin\tresult=success\tsev=3",
	}

	fmt.Println("--- Сценарий 1: Отправка в SIEM через Syslog (UDP) ---\n")
	
	// ВАЖНО: Замените на адрес вашего SIEM
	// Пример для локального тестирования (если SIEM недоступен):
	// forwarder, err := siem.NewSyslogForwarder("localhost", 514, "udp")
	
	// Для демонстрации покажем, как это работало бы:
	fmt.Println("Подключение к SIEM через Syslog...")
	fmt.Println("  Адрес: siem.example.com:514 (UDP)")
	fmt.Println("  Формат: RFC 5424 с JSON payload\n")

	// Обработка и форматирование для отправки
	for i, logLine := range logs {
		event, err := proc.Process(logLine)
		if err != nil {
			log.Printf("Ошибка обработки лога %d: %v\n", i+1, err)
			continue
		}

		fmt.Printf("Лог %d обработан:\n", i+1)
		fmt.Printf("  Категория: %s\n", event.Category)
		fmt.Printf("  Критичность: %s\n", event.Severity)
		fmt.Printf("  Описание: %s\n", event.Description)
		fmt.Println("  -> Готов к отправке в SIEM")
		fmt.Println()
	}

	fmt.Println("--- Сценарий 2: Отправка через HTTP API (Splunk HEC, Elastic) ---\n")
	
	// Пример для Splunk HTTP Event Collector
	fmt.Println("Конфигурация для Splunk HEC:")
	fmt.Println("  URL: https://splunk.example.com:8088/services/collector/event")
	fmt.Println("  Токен: ********-****-****-****-************")
	fmt.Println("  Формат: JSON с ГОСТ полями\n")

	_ = siem.NewHTTPForwarder(
		"https://splunk.example.com:8088/services/collector/event",
		"YOUR-HEC-TOKEN",
		map[string]string{
			"X-Custom-Header": "LoggerV2",
		},
	)

	// Обработка логов
	events, _ := proc.ProcessBatch(logs)

	fmt.Printf("Подготовлено %d событий для отправки\n", len(events))
	fmt.Println("В реальной ситуации здесь произойдет отправка через HTTP POST\n")

	// В реальном использовании:
	// err = httpForwarder.ForwardBatch(events)

	fmt.Println("--- Сценарий 3: Интеграция через File Monitoring ---\n")
	
	fmt.Println("Запись нормализованных логов в файл для SIEM:")
	outputFile := "normalized_logs.json"
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, logLine := range logs {
		event, err := proc.Process(logLine)
		if err != nil {
			continue
		}

		jsonOutput, err := proc.ConvertToJSON(event)
		if err != nil {
			continue
		}

		writer.WriteString(jsonOutput + "\n")
		fmt.Printf("Событие %d записано в %s\n", i+1, outputFile)
	}
	writer.Flush()

	fmt.Printf("\nФайл %s готов для мониторинга SIEM агентом\n", outputFile)
	fmt.Println("(Splunk Universal Forwarder, Filebeat, NXLog, etc.)\n")

	fmt.Println("--- Поддерживаемые SIEM системы ---\n")
	fmt.Println("✓ Splunk (через HEC или файлы)")
	fmt.Println("✓ Elastic Stack (через HTTP API)")
	fmt.Println("✓ IBM QRadar (через Syslog)")
	fmt.Println("✓ ArcSight (через Syslog или CEF)")
	fmt.Println("✓ LogRhythm (через Syslog)")
	fmt.Println("✓ AlienVault OSSIM (через Syslog)")
	fmt.Println("✓ Graylog (через GELF или Syslog)")
	fmt.Println("✓ MaxPatrol SIEM (через Syslog)")
	fmt.Println("✓ R-Vision SOAR (через REST API)")
	fmt.Println()

	fmt.Println("=== Демонстрация завершена ===")
}
