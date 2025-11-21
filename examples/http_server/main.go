package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/kxrty/loggerv2/internal/processor"
)

type LogRequest struct {
	Logs []string `json:"logs"`
}

type LogResponse struct {
	Success int                    `json:"success"`
	Errors  int                    `json:"errors"`
	Events  []interface{}          `json:"events"`
	ErrorMessages []string          `json:"error_messages,omitempty"`
}

func main() {
	proc := processor.NewProcessor()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "ok"}`)
	})

	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req LogRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		events, errors := proc.ProcessBatch(req.Logs)

		var errorMessages []string
		for _, err := range errors {
			errorMessages = append(errorMessages, err.Error())
		}

		var eventsJSON []interface{}
		for _, event := range events {
			eventMap := map[string]interface{}{
				"event_id":    event.EventID,
				"timestamp":   event.Timestamp,
				"source":      event.Source,
				"category":    event.Category,
				"severity":    event.Severity,
				"description": event.Description,
				"result":      event.Result,
				"action":      event.Action,
			}
			if event.SubjectAccount != nil {
				eventMap["subject_account"] = event.SubjectAccount
			}
			if event.ObjectAccount != nil {
				eventMap["object_account"] = event.ObjectAccount
			}
			if len(event.AdditionalData) > 0 {
				eventMap["additional_data"] = event.AdditionalData
			}
			eventsJSON = append(eventsJSON, eventMap)
		}

		response := LogResponse{
			Success:       len(events),
			Errors:        len(errors),
			Events:        eventsJSON,
			ErrorMessages: errorMessages,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/process-single", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		logLine := strings.TrimSpace(string(body))
		event, err := proc.Process(logLine)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to process log: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(event)
	})

	http.HandleFunc("/detect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		logLine := strings.TrimSpace(string(body))
		logType := proc.DetectLogType(logLine)

		typeNames := map[processor.LogType]string{
			processor.LogTypeSyslog: "Syslog",
			processor.LogTypeCEF:    "CEF",
			processor.LogTypeLEEF:   "LEEF",
			processor.LogTypeXML:    "XML",
			processor.LogTypeUnknown: "Unknown",
		}

		response := map[string]string{
			"log_type": typeNames[logType],
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	port := ":8080"
	log.Printf("Запуск HTTP сервера на порту %s\n", port)
	log.Printf("Endpoints:\n")
	log.Printf("  GET  /health         - Проверка здоровья сервиса\n")
	log.Printf("  POST /process        - Обработка массива логов (JSON)\n")
	log.Printf("  POST /process-single - Обработка одного лога (text)\n")
	log.Printf("  POST /detect         - Определение типа лога (text)\n")
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
