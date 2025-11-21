package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/kxrty/loggerv2/internal/processor"
)

func main() {
	inputFile := flag.String("input", "", "Входной файл с логами")
	outputFile := flag.String("output", "", "Выходной файл для результатов (по умолчанию stdout)")
	flag.Parse()

	proc := processor.NewProcessor()

	var scanner *bufio.Scanner
	
	if *inputFile != "" {
		file, err := os.Open(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	var output *os.File
	if *outputFile != "" {
		var err error
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка создания выходного файла: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	lineNum := 0
	successCount := 0
	errorCount := 0

	for scanner.Scan() {
		lineNum++
		logLine := scanner.Text()
		
		if logLine == "" {
			continue
		}

		event, err := proc.Process(logLine)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка обработки строки %d: %v\n", lineNum, err)
			errorCount++
			continue
		}

		jsonOutput, err := proc.ConvertToJSON(event)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка преобразования в JSON строки %d: %v\n", lineNum, err)
			errorCount++
			continue
		}

		fmt.Fprintln(output, jsonOutput)
		successCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения входных данных: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "\nОбработка завершена:\n")
	fmt.Fprintf(os.Stderr, "  Успешно: %d\n", successCount)
	fmt.Fprintf(os.Stderr, "  Ошибок: %d\n", errorCount)
	fmt.Fprintf(os.Stderr, "  Всего строк: %d\n", lineNum)
}
