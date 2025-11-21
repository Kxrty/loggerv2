package processor

import (
	"testing"
)

func TestProcessor_DetectLogType(t *testing.T) {
	proc := NewProcessor()
	
	tests := []struct {
		name     string
		logLine  string
		expected LogType
	}{
		{
			name:     "CEF format",
			logLine:  "CEF:0|Vendor|Product|1.0|100|Test|5|src=1.1.1.1",
			expected: LogTypeCEF,
		},
		{
			name:     "LEEF format",
			logLine:  "LEEF:1.0|Vendor|Product|1.0|100|test=value",
			expected: LogTypeLEEF,
		},
		{
			name:     "XML format",
			logLine:  "<Event><System><EventID>100</EventID></System></Event>",
			expected: LogTypeXML,
		},
		{
			name:     "Syslog format",
			logLine:  "<134>Oct 11 22:14:15 mymachine test: message",
			expected: LogTypeSyslog,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := proc.DetectLogType(tt.logLine)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestProcessor_Process(t *testing.T) {
	proc := NewProcessor()
	
	tests := []struct {
		name    string
		logLine string
		wantErr bool
	}{
		{
			name:    "Valid Syslog",
			logLine: "<134>Oct 11 22:14:15 mymachine su: test message",
			wantErr: false,
		},
		{
			name:    "Valid CEF",
			logLine: "CEF:0|Vendor|Product|1.0|100|Test|5|src=1.1.1.1",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			logLine: "This is not a valid log format",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := proc.Process(tt.logLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && event == nil {
				t.Error("Expected event to be non-nil")
			}
		})
	}
}

func TestProcessor_ConvertToJSON(t *testing.T) {
	proc := NewProcessor()
	
	logLine := "<134>Oct 11 22:14:15 mymachine su: test message"
	event, err := proc.Process(logLine)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	
	json, err := proc.ConvertToJSON(event)
	if err != nil {
		t.Fatalf("ConvertToJSON failed: %v", err)
	}
	
	if json == "" {
		t.Error("Expected non-empty JSON output")
	}
	
	if !containsString(json, "event_id") {
		t.Error("JSON should contain event_id field")
	}
	
	if !containsString(json, "timestamp") {
		t.Error("JSON should contain timestamp field")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
