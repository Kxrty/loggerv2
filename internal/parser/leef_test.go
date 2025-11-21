package parser

import (
	"testing"
)

func TestLEEFParser_Parse(t *testing.T) {
	parser := NewLEEFParser()
	
	logLine := "LEEF:1.0|Microsoft|MSExchange|4.0 SP1|15345|src=10.0.0.1\tdst=172.50.123.1\tsev=5\tcat=anomaly\tsrcPort=81\tdstPort=21"
	
	event, err := parser.Parse(logLine)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if event.AdditionalData["vendor"] != "Microsoft" {
		t.Errorf("Expected vendor 'Microsoft', got '%v'", event.AdditionalData["vendor"])
	}
	
	if event.AdditionalData["product"] != "MSExchange" {
		t.Errorf("Expected product 'MSExchange', got '%v'", event.AdditionalData["product"])
	}
	
	if event.Source.IPAddress == "" {
		t.Error("Expected IPAddress to be set")
	}
}

func TestLEEFParser_Version2(t *testing.T) {
	parser := NewLEEFParser()
	
	logLine := "LEEF:2.0|Vendor|Product|1.0|EventID|usrName=admin^action=login^result=success"
	
	event, err := parser.Parse(logLine)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if event.SubjectAccount == nil {
		t.Fatal("Expected SubjectAccount to be set")
	}
	
	if event.SubjectAccount.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", event.SubjectAccount.Username)
	}
}
