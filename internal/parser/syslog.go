package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kxrty/loggerv2/internal/models"
	"github.com/google/uuid"
)

// SyslogParser парсит RFC 3164 и RFC 5424 syslog сообщения
type SyslogParser struct{}

// RFC3164Pattern - паттерн для RFC 3164 формата
var RFC3164Pattern = regexp.MustCompile(`^<(\d+)>(\w+\s+\d+\s+\d+:\d+:\d+)\s+(\S+)\s+(\S+?)(\[(\d+)\])?:\s+(.*)$`)

// RFC5424Pattern - паттерн для RFC 5424 формата
var RFC5424Pattern = regexp.MustCompile(`^<(\d+)>(\d+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(.*)$`)

func NewSyslogParser() *SyslogParser {
	return &SyslogParser{}
}

// Parse парсит syslog сообщение и возвращает GOSTEvent
func (p *SyslogParser) Parse(logLine string) (*models.GOSTEvent, error) {
	if match := RFC5424Pattern.FindStringSubmatch(logLine); match != nil {
		return p.parseRFC5424(match)
	}
	
	if match := RFC3164Pattern.FindStringSubmatch(logLine); match != nil {
		return p.parseRFC3164(match)
	}
	
	return nil, fmt.Errorf("неподдерживаемый формат syslog")
}

func (p *SyslogParser) parseRFC3164(match []string) (*models.GOSTEvent, error) {
	priority, _ := strconv.Atoi(match[1])
	timestamp := match[2]
	hostname := match[3]
	appName := match[4]
	processID := 0
	if match[6] != "" {
		processID, _ = strconv.Atoi(match[6])
	}
	message := match[7]

	parsedTime, err := time.Parse("Jan 2 15:04:05", timestamp)
	if err != nil {
		parsedTime = time.Now()
	} else {
		parsedTime = parsedTime.AddDate(time.Now().Year(), 0, 0)
	}

	event := &models.GOSTEvent{
		EventID:     uuid.New().String(),
		Timestamp:   parsedTime,
		Description: message,
		Source: models.Source{
			Hostname:    hostname,
			Application: appName,
			ProcessID:   processID,
		},
		Severity:       p.mapSyslogSeverityToGOST(priority),
		Category:       p.categorizeSyslogMessage(message),
		Result:         models.ResultUnknown,
		AdditionalData: make(map[string]interface{}),
	}

	event.AdditionalData["syslog_priority"] = priority
	event.AdditionalData["syslog_facility"] = priority / 8
	event.AdditionalData["syslog_severity"] = priority % 8

	return event, nil
}

func (p *SyslogParser) parseRFC5424(match []string) (*models.GOSTEvent, error) {
	priority, _ := strconv.Atoi(match[1])
	version := match[2]
	timestamp := match[3]
	hostname := match[4]
	appName := match[5]
	procID := match[6]
	msgID := match[7]
	message := match[9]

	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		parsedTime = time.Now()
	}

	processID := 0
	if procID != "-" {
		processID, _ = strconv.Atoi(procID)
	}

	event := &models.GOSTEvent{
		EventID:     uuid.New().String(),
		Timestamp:   parsedTime,
		Description: message,
		Source: models.Source{
			Hostname:    hostname,
			Application: appName,
			ProcessID:   processID,
		},
		Severity:       p.mapSyslogSeverityToGOST(priority),
		Category:       p.categorizeSyslogMessage(message),
		Result:         models.ResultUnknown,
		AdditionalData: make(map[string]interface{}),
	}

	event.AdditionalData["syslog_priority"] = priority
	event.AdditionalData["syslog_facility"] = priority / 8
	event.AdditionalData["syslog_severity"] = priority % 8
	event.AdditionalData["syslog_version"] = version
	event.AdditionalData["syslog_msgid"] = msgID

	return event, nil
}

func (p *SyslogParser) mapSyslogSeverityToGOST(priority int) string {
	severity := priority % 8
	
	switch severity {
	case 0, 1, 2:
		return models.SeverityCritical
	case 3:
		return models.SeverityHigh
	case 4:
		return models.SeverityMedium
	case 5, 6:
		return models.SeverityLow
	case 7:
		return models.SeverityInfo
	default:
		return models.SeverityInfo
	}
}

func (p *SyslogParser) categorizeSyslogMessage(message string) string {
	messageLower := strings.ToLower(message)
	
	if strings.Contains(messageLower, "login") || strings.Contains(messageLower, "auth") {
		return models.CategoryAuthentication
	}
	if strings.Contains(messageLower, "access") || strings.Contains(messageLower, "denied") {
		return models.CategoryAccess
	}
	if strings.Contains(messageLower, "network") || strings.Contains(messageLower, "connection") {
		return models.CategoryNetworkEvent
	}
	if strings.Contains(messageLower, "security") || strings.Contains(messageLower, "breach") {
		return models.CategorySecurityEvent
	}
	
	return models.CategorySystemEvent
}
