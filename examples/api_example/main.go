package main

import (
	"fmt"
	"log"

	"github.com/kxrty/loggerv2/internal/processor"
)

func main() {
	proc := processor.NewProcessor()

	logs := []string{
		"<134>Oct 11 22:14:15 mymachine su: 'su root' failed",
		"CEF:0|Security|IDS|1.0|100|Attack detected|10|src=10.0.0.1 dst=192.168.1.1",
		"LEEF:1.0|Microsoft|MSExchange|4.0|15345|src=10.0.0.1\tdst=172.50.123.1\tsev=5",
	}

	for i, logLine := range logs {
		event, err := proc.Process(logLine)
		if err != nil {
			log.Printf("Ошибка обработки лога %d: %v\n", i+1, err)
			continue
		}

		fmt.Printf("=== Лог %d ===\n", i+1)
		fmt.Printf("ID события: %s\n", event.EventID)
		fmt.Printf("Время: %s\n", event.Timestamp)
		fmt.Printf("Источник: %s\n", event.Source.Hostname)
		fmt.Printf("Приложение: %s\n", event.Source.Application)
		fmt.Printf("Категория: %s\n", event.Category)
		fmt.Printf("Критичность: %s\n", event.Severity)
		fmt.Printf("Описание: %s\n", event.Description)
		fmt.Printf("Результат: %s\n", event.Result)
		fmt.Println()

		jsonOutput, err := proc.ConvertToJSON(event)
		if err != nil {
			log.Printf("Ошибка преобразования в JSON: %v\n", err)
			continue
		}
		fmt.Println("JSON:")
		fmt.Println(jsonOutput)
		fmt.Println()
	}

	events, errors := proc.ProcessBatch(logs)
	fmt.Printf("Обработано событий: %d\n", len(events))
	fmt.Printf("Ошибок: %d\n", len(errors))
}
