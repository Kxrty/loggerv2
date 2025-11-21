package parser

import (
	"testing"
	"github.com/kxrty/loggerv2/internal/models"
)

func TestXMLParser_Parse(t *testing.T) {
	parser := NewXMLParser()
	
	xmlLog := `<Event xmlns="http://schemas.microsoft.com/win/2004/08/events/event">
  <System>
    <Provider Name="Microsoft-Windows-Security-Auditing" Guid="{54849625-5478-4994-A5BA-3E3B0328C30D}"/>
    <EventID>4624</EventID>
    <Level>0</Level>
    <TimeCreated SystemTime="2023-10-11T22:14:15.123456Z"/>
    <Computer>workstation.example.com</Computer>
    <Execution ProcessID="500" ThreadID="600"/>
  </System>
  <EventData>
    <Data Name="TargetUserName">john.doe</Data>
    <Data Name="TargetDomainName">EXAMPLE</Data>
  </EventData>
</Event>`
	
	event, err := parser.Parse(xmlLog)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if event.Source.Hostname != "workstation.example.com" {
		t.Errorf("Expected hostname 'workstation.example.com', got '%s'", event.Source.Hostname)
	}
	
	if event.Source.Application != "Microsoft-Windows-Security-Auditing" {
		t.Errorf("Expected application 'Microsoft-Windows-Security-Auditing', got '%s'", event.Source.Application)
	}
	
	if event.Source.ProcessID != 500 {
		t.Errorf("Expected ProcessID 500, got %d", event.Source.ProcessID)
	}
	
	if event.Category != models.CategoryAuthentication {
		t.Errorf("Expected category АУТЕНТИФИКАЦИЯ, got %s", event.Category)
	}
	
	if event.SubjectAccount == nil {
		t.Fatal("Expected SubjectAccount to be set")
	}
	
	if event.SubjectAccount.Username != "john.doe" {
		t.Errorf("Expected username 'john.doe', got '%s'", event.SubjectAccount.Username)
	}
}
