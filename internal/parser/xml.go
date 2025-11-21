package parser

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/kxrty/loggerv2/internal/models"
	"github.com/google/uuid"
)

// XMLParser парсит XML логи (например Windows Event Log)
type XMLParser struct{}

// Event представляет базовую XML структуру события
type Event struct {
	XMLName xml.Name `xml:"Event"`
	System  System   `xml:"System"`
	EventData EventData `xml:"EventData"`
}

type System struct {
	Provider      Provider `xml:"Provider"`
	EventID       int      `xml:"EventID"`
	Version       int      `xml:"Version"`
	Level         int      `xml:"Level"`
	Task          int      `xml:"Task"`
	Opcode        int      `xml:"Opcode"`
	Keywords      string   `xml:"Keywords"`
	TimeCreated   TimeCreated `xml:"TimeCreated"`
	EventRecordID int      `xml:"EventRecordID"`
	Correlation   Correlation `xml:"Correlation"`
	Execution     Execution `xml:"Execution"`
	Channel       string   `xml:"Channel"`
	Computer      string   `xml:"Computer"`
	Security      Security `xml:"Security"`
}

type Provider struct {
	Name string `xml:"Name,attr"`
	Guid string `xml:"Guid,attr"`
}

type TimeCreated struct {
	SystemTime string `xml:"SystemTime,attr"`
}

type Correlation struct {
	ActivityID string `xml:"ActivityID,attr"`
}

type Execution struct {
	ProcessID int `xml:"ProcessID,attr"`
	ThreadID  int `xml:"ThreadID,attr"`
}

type Security struct {
	UserID string `xml:"UserID,attr"`
}

type EventData struct {
	Data []Data `xml:"Data"`
}

type Data struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

func NewXMLParser() *XMLParser {
	return &XMLParser{}
}

// Parse парсит XML сообщение и возвращает GOSTEvent
func (p *XMLParser) Parse(logLine string) (*models.GOSTEvent, error) {
	var event Event
	
	decoder := xml.NewDecoder(strings.NewReader(logLine))
	if err := decoder.Decode(&event); err != nil {
		return nil, fmt.Errorf("ошибка парсинга XML: %w", err)
	}

	timestamp, err := time.Parse(time.RFC3339Nano, event.System.TimeCreated.SystemTime)
	if err != nil {
		timestamp = time.Now()
	}

	gostEvent := &models.GOSTEvent{
		EventID:     uuid.New().String(),
		Timestamp:   timestamp,
		Description: p.buildDescription(event),
		Source: models.Source{
			Hostname:    event.System.Computer,
			Application: event.System.Provider.Name,
			ProcessID:   event.System.Execution.ProcessID,
		},
		Severity:       p.mapXMLLevelToGOST(event.System.Level),
		Category:       p.categorizeXMLEvent(event),
		Result:         p.determineResult(event),
		AdditionalData: make(map[string]interface{}),
	}

	gostEvent.AdditionalData["xml_event_id"] = event.System.EventID
	gostEvent.AdditionalData["xml_level"] = event.System.Level
	gostEvent.AdditionalData["xml_task"] = event.System.Task
	gostEvent.AdditionalData["xml_opcode"] = event.System.Opcode
	gostEvent.AdditionalData["xml_keywords"] = event.System.Keywords
	gostEvent.AdditionalData["xml_channel"] = event.System.Channel
	gostEvent.AdditionalData["xml_record_id"] = event.System.EventRecordID
	gostEvent.AdditionalData["provider_guid"] = event.System.Provider.Guid

	for _, data := range event.EventData.Data {
		if data.Name != "" {
			gostEvent.AdditionalData["xml_"+data.Name] = data.Value
		}
	}

	if event.System.Security.UserID != "" {
		gostEvent.SubjectAccount = &models.Account{
			UserID: event.System.Security.UserID,
		}
	}

	p.enrichEventFromEventData(gostEvent, event.EventData.Data)

	return gostEvent, nil
}

func (p *XMLParser) buildDescription(event Event) string {
	if len(event.EventData.Data) > 0 {
		var parts []string
		for _, data := range event.EventData.Data {
			if data.Value != "" {
				parts = append(parts, fmt.Sprintf("%s: %s", data.Name, data.Value))
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "; ")
		}
	}
	
	return fmt.Sprintf("Event ID %d from %s", event.System.EventID, event.System.Provider.Name)
}

func (p *XMLParser) mapXMLLevelToGOST(level int) string {
	switch level {
	case 1:
		return models.SeverityCritical
	case 2:
		return models.SeverityHigh
	case 3:
		return models.SeverityMedium
	case 4:
		return models.SeverityInfo
	case 5:
		return models.SeverityLow
	default:
		return models.SeverityInfo
	}
}

func (p *XMLParser) categorizeXMLEvent(event Event) string {
	eventID := event.System.EventID
	channel := strings.ToLower(event.System.Channel)
	providerName := strings.ToLower(event.System.Provider.Name)

	switch {
	case eventID >= 4624 && eventID <= 4634:
		return models.CategoryAuthentication
	case eventID >= 4720 && eventID <= 4767:
		return models.CategoryAuthorization
	case eventID >= 4660 && eventID <= 4663:
		return models.CategoryAccess
	case eventID >= 4670 && eventID <= 4690:
		return models.CategoryDataModification
	}

	if strings.Contains(channel, "security") || strings.Contains(providerName, "security") {
		return models.CategorySecurityEvent
	}

	if strings.Contains(channel, "system") || strings.Contains(providerName, "system") {
		return models.CategorySystemEvent
	}

	if strings.Contains(channel, "application") {
		return models.CategorySystemEvent
	}

	return models.CategorySystemEvent
}

func (p *XMLParser) determineResult(event Event) string {
	for _, data := range event.EventData.Data {
		nameLower := strings.ToLower(data.Name)
		valueLower := strings.ToLower(data.Value)
		
		if strings.Contains(nameLower, "status") || strings.Contains(nameLower, "result") {
			if strings.Contains(valueLower, "success") || valueLower == "0" || valueLower == "0x0" {
				return models.ResultSuccess
			}
			if strings.Contains(valueLower, "fail") || strings.Contains(valueLower, "error") {
				return models.ResultFailure
			}
		}
	}

	if event.System.Level == 1 || event.System.Level == 2 {
		return models.ResultFailure
	}

	eventID := event.System.EventID
	if (eventID >= 4624 && eventID <= 4625) || (eventID >= 4634 && eventID <= 4647) {
		if eventID == 4625 {
			return models.ResultFailure
		}
		return models.ResultSuccess
	}

	return models.ResultUnknown
}

func (p *XMLParser) enrichEventFromEventData(gostEvent *models.GOSTEvent, data []Data) {
	for _, d := range data {
		nameLower := strings.ToLower(d.Name)
		
		switch {
		case strings.Contains(nameLower, "targetusername") || strings.Contains(nameLower, "subjectusername"):
			if gostEvent.SubjectAccount == nil {
				gostEvent.SubjectAccount = &models.Account{}
			}
			gostEvent.SubjectAccount.Username = d.Value
			
		case strings.Contains(nameLower, "targetdomainname") || strings.Contains(nameLower, "subjectdomainname"):
			if gostEvent.SubjectAccount == nil {
				gostEvent.SubjectAccount = &models.Account{}
			}
			gostEvent.SubjectAccount.Domain = d.Value
			
		case strings.Contains(nameLower, "ipaddress") || strings.Contains(nameLower, "workstationname"):
			if gostEvent.Source.IPAddress == "" {
				gostEvent.Source.IPAddress = d.Value
			}
			
		case strings.Contains(nameLower, "processname"):
			gostEvent.Source.Process = d.Value
		}
	}
}
