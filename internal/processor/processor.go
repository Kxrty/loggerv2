package processor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kxrty/loggerv2/internal/models"
	"github.com/kxrty/loggerv2/internal/parser"
)

// LogType представляет тип лога
type LogType int

const (
	LogTypeUnknown LogType = iota
	LogTypeSyslog
	LogTypeCEF
	LogTypeLEEF
	LogTypeXML
)

// Processor обрабатывает логи различных форматов
type Processor struct {
	syslogParser *parser.SyslogParser
	cefParser    *parser.CEFParser
	leefParser   *parser.LEEFParser
	xmlParser    *parser.XMLParser
}

// NewProcessor создает новый процессор логов
func NewProcessor() *Processor {
	return &Processor{
		syslogParser: parser.NewSyslogParser(),
		cefParser:    parser.NewCEFParser(),
		leefParser:   parser.NewLEEFParser(),
		xmlParser:    parser.NewXMLParser(),
	}
}

// Process обрабатывает лог и преобразует его в формат ГОСТ
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
	default:
		return nil, fmt.Errorf("неизвестный тип лога")
	}
}

// ProcessBatch обрабатывает массив логов
func (p *Processor) ProcessBatch(logLines []string) ([]*models.GOSTEvent, []error) {
	events := make([]*models.GOSTEvent, 0, len(logLines))
	errors := make([]error, 0)
	
	for i, logLine := range logLines {
		event, err := p.Process(logLine)
		if err != nil {
			errors = append(errors, fmt.Errorf("строка %d: %w", i+1, err))
			continue
		}
		events = append(events, event)
	}
	
	return events, errors
}

// DetectLogType автоматически определяет тип лога
func (p *Processor) DetectLogType(logLine string) LogType {
	logLine = strings.TrimSpace(logLine)
	
	if strings.HasPrefix(logLine, "CEF:") {
		return LogTypeCEF
	}
	
	if strings.HasPrefix(logLine, "LEEF:") {
		return LogTypeLEEF
	}
	
	if strings.HasPrefix(logLine, "<") && strings.Contains(logLine, "<?xml") || 
	   strings.HasPrefix(logLine, "<?xml") {
		return LogTypeXML
	}
	
	if strings.HasPrefix(logLine, "<Event") {
		return LogTypeXML
	}
	
	if strings.HasPrefix(logLine, "<") && strings.Contains(logLine, ">") {
		priorityEnd := strings.Index(logLine, ">")
		if priorityEnd > 0 && priorityEnd < 10 {
			priority := logLine[1:priorityEnd]
			if isNumeric(priority) {
				return LogTypeSyslog
			}
		}
	}
	
	return LogTypeUnknown
}

// ConvertToJSON преобразует GOSTEvent в JSON
func (p *Processor) ConvertToJSON(event *models.GOSTEvent) (string, error) {
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации в JSON: %w", err)
	}
	return string(data), nil
}

// ConvertBatchToJSON преобразует массив GOSTEvent в JSON
func (p *Processor) ConvertBatchToJSON(events []*models.GOSTEvent) (string, error) {
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации в JSON: %w", err)
	}
	return string(data), nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
