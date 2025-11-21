package parser

import (
	"testing"
	"github.com/kxrty/loggerv2/internal/models"
)

func TestCEFParser_Parse(t *testing.T) {
	parser := NewCEFParser()
	
	logLine := "CEF:0|Security|threatmanager|1.0|100|worm successfully stopped|10|src=10.0.0.1 dst=2.1.2.2 spt=1232"
	
	event, err := parser.Parse(logLine)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if event.Description != "worm successfully stopped" {
		t.Errorf("Unexpected description: %s", event.Description)
	}
	
	if event.Severity != models.SeverityCritical {
		t.Errorf("Expected severity КРИТИЧЕСКИЙ, got %s", event.Severity)
	}
	
	if event.AdditionalData["device_vendor"] != "Security" {
		t.Errorf("Expected device_vendor 'Security', got '%v'", event.AdditionalData["device_vendor"])
	}
	
	if event.AdditionalData["device_product"] != "threatmanager" {
		t.Errorf("Expected device_product 'threatmanager', got '%v'", event.AdditionalData["device_product"])
	}
}

func TestCEFParser_WithUser(t *testing.T) {
	parser := NewCEFParser()
	
	logLine := "CEF:0|Vendor|Product|1.0|200|User Login|5|suser=john.doe sdomain=EXAMPLE act=login outcome=success"
	
	event, err := parser.Parse(logLine)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if event.SubjectAccount == nil {
		t.Fatal("Expected SubjectAccount to be set")
	}
	
	if event.SubjectAccount.Username != "john.doe" {
		t.Errorf("Expected username 'john.doe', got '%s'", event.SubjectAccount.Username)
	}
	
	if event.SubjectAccount.Domain != "EXAMPLE" {
		t.Errorf("Expected domain 'EXAMPLE', got '%s'", event.SubjectAccount.Domain)
	}
	
	if event.Result != models.ResultSuccess {
		t.Errorf("Expected result УСПЕХ, got %s", event.Result)
	}
}
