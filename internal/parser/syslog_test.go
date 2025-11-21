package parser

import (
	"testing"
	"github.com/kxrty/loggerv2/internal/models"
)

func TestSyslogParser_RFC3164(t *testing.T) {
	parser := NewSyslogParser()
	
	logLine := "<134>Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8"
	
	event, err := parser.Parse(logLine)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if event.Source.Hostname != "mymachine" {
		t.Errorf("Expected hostname 'mymachine', got '%s'", event.Source.Hostname)
	}
	
	if event.Source.Application != "su" {
		t.Errorf("Expected application 'su', got '%s'", event.Source.Application)
	}
	
	if event.Description != "'su root' failed for lonvick on /dev/pts/8" {
		t.Errorf("Unexpected description: %s", event.Description)
	}
	
	if event.Severity != models.SeverityLow {
		t.Errorf("Expected severity НИЗКИЙ, got %s", event.Severity)
	}
}

func TestSyslogParser_RFC5424(t *testing.T) {
	parser := NewSyslogParser()
	
	logLine := "<165>1 2023-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 - User login successful"
	
	event, err := parser.Parse(logLine)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if event.Source.Hostname != "mymachine.example.com" {
		t.Errorf("Expected hostname 'mymachine.example.com', got '%s'", event.Source.Hostname)
	}
	
	if event.Source.Application != "evntslog" {
		t.Errorf("Expected application 'evntslog', got '%s'", event.Source.Application)
	}
}
