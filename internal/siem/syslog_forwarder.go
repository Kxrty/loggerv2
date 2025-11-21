package siem

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/kxrty/loggerv2/internal/models"
)

// SyslogForwarder отправляет нормализованные события в SIEM через Syslog
type SyslogForwarder struct {
	conn     net.Conn
	host     string
	port     int
	protocol string
}

// NewSyslogForwarder создает новый форвардер
func NewSyslogForwarder(host string, port int, protocol string) (*SyslogForwarder, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к SIEM: %w", err)
	}

	return &SyslogForwarder{
		conn:     conn,
		host:     host,
		port:     port,
		protocol: protocol,
	}, nil
}

// Forward отправляет событие в SIEM
func (f *SyslogForwarder) Forward(event *models.GOSTEvent) error {
	message, err := f.formatMessage(event)
	if err != nil {
		return fmt.Errorf("ошибка форматирования сообщения: %w", err)
	}

	_, err = f.conn.Write([]byte(message + "\n"))
	if err != nil {
		return fmt.Errorf("ошибка отправки в SIEM: %w", err)
	}

	return nil
}

// ForwardBatch отправляет несколько событий
func (f *SyslogForwarder) ForwardBatch(events []*models.GOSTEvent) []error {
	var errors []error

	for i, event := range events {
		if err := f.Forward(event); err != nil {
			errors = append(errors, fmt.Errorf("событие %d: %w", i, err))
		}
	}

	return errors
}

// formatMessage форматирует событие ГОСТ в Syslog формат для SIEM
func (f *SyslogForwarder) formatMessage(event *models.GOSTEvent) (string, error) {
	// Преобразуем в JSON для структурированного лога
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	// RFC 5424 формат: <PRI>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID STRUCTURED-DATA MSG
	priority := f.calculatePriority(event.Severity)
	timestamp := event.Timestamp.Format(time.RFC3339)
	hostname := event.Source.Hostname
	if hostname == "" {
		hostname = "unknown"
	}
	appName := event.Source.Application
	if appName == "" {
		appName = "loggerv2"
	}

	// Формируем сообщение
	message := fmt.Sprintf("<%d>1 %s %s %s - - - GOST: %s",
		priority,
		timestamp,
		hostname,
		appName,
		string(eventJSON))

	return message, nil
}

// calculatePriority вычисляет Syslog priority на основе ГОСТ критичности
func (f *SyslogForwarder) calculatePriority(severity string) int {
	facility := 16 // local0
	var level int

	switch severity {
	case models.SeverityCritical:
		level = 2 // critical
	case models.SeverityHigh:
		level = 3 // error
	case models.SeverityMedium:
		level = 4 // warning
	case models.SeverityLow:
		level = 5 // notice
	case models.SeverityInfo:
		level = 6 // informational
	default:
		level = 6
	}

	return facility*8 + level
}

// Close закрывает соединение с SIEM
func (f *SyslogForwarder) Close() error {
	if f.conn != nil {
		return f.conn.Close()
	}
	return nil
}
