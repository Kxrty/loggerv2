package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kxrty/loggerv2/internal/models"
	"github.com/google/uuid"
)

// LEEFParser парсит Log Event Extended Format (LEEF) логи
type LEEFParser struct{}

// LEEF Format: LEEF:Version|Vendor|Product|Version|EventID|Attributes
var leefPattern = regexp.MustCompile(`^LEEF:([\d.]+)\|([^|]*)\|([^|]*)\|([^|]*)\|([^|]*)\|(.*)$`)

func NewLEEFParser() *LEEFParser {
	return &LEEFParser{}
}

// Parse парсит LEEF сообщение и возвращает GOSTEvent
func (p *LEEFParser) Parse(logLine string) (*models.GOSTEvent, error) {
	match := leefPattern.FindStringSubmatch(logLine)
	if match == nil {
		return nil, fmt.Errorf("неверный формат LEEF")
	}

	version := match[1]
	vendor := match[2]
	product := match[3]
	productVersion := match[4]
	eventID := match[5]
	attributes := match[6]

	attrs := p.parseAttributes(attributes, version)

	timestamp := time.Now()
	if devTime, ok := attrs["devTime"]; ok {
		if t, err := p.parseTimestamp(devTime); err == nil {
			timestamp = t
		}
	}

	event := &models.GOSTEvent{
		EventID:     uuid.New().String(),
		Timestamp:   timestamp,
		Description: p.getAttributeValue(attrs, "usrName", "msg", "eventId"),
		Source: models.Source{
			Hostname:    p.getAttributeValue(attrs, "devName", "srcHostName", "dstHostName"),
			Application: fmt.Sprintf("%s %s", vendor, product),
			IPAddress:   p.getAttributeValue(attrs, "src", "dst"),
		},
		Severity:       p.mapLEEFSeverityToGOST(attrs),
		Category:       p.categorizeLEEFEvent(eventID, attrs),
		Result:         p.determineResult(attrs),
		AdditionalData: make(map[string]interface{}),
	}

	event.AdditionalData["leef_version"] = version
	event.AdditionalData["vendor"] = vendor
	event.AdditionalData["product"] = product
	event.AdditionalData["product_version"] = productVersion
	event.AdditionalData["event_id"] = eventID

	if srcUser, ok := attrs["srcUser"]; ok {
		event.SubjectAccount = &models.Account{
			Username: srcUser,
			Domain:   attrs["srcDomain"],
		}
	} else if usrName, ok := attrs["usrName"]; ok {
		event.SubjectAccount = &models.Account{
			Username: usrName,
		}
	}

	if dstUser, ok := attrs["dstUser"]; ok {
		event.ObjectAccount = &models.Account{
			Username: dstUser,
			Domain:   attrs["dstDomain"],
		}
	}

	for k, v := range attrs {
		event.AdditionalData["leef_"+k] = v
	}

	if action, ok := attrs["action"]; ok {
		event.Action = action
	} else if cat, ok := attrs["cat"]; ok {
		event.Action = cat
	}

	return event, nil
}

func (p *LEEFParser) parseAttributes(attributes string, version string) map[string]string {
	attrs := make(map[string]string)
	
	delimiter := "\t"
	if version == "2.0" {
		if strings.Contains(attributes, "x09") {
			delimiter = "x09"
		} else if strings.Contains(attributes, "^") {
			delimiter = "^"
		}
	}
	
	var parts []string
	if delimiter == "x09" {
		parts = strings.Split(strings.ReplaceAll(attributes, "x09", "\t"), "\t")
	} else if delimiter == "^" {
		parts = strings.Split(attributes, "^")
	} else {
		parts = strings.Split(attributes, delimiter)
	}
	
	for _, part := range parts {
		if idx := strings.Index(part, "="); idx != -1 {
			key := strings.TrimSpace(part[:idx])
			value := strings.TrimSpace(part[idx+1:])
			attrs[key] = value
		}
	}
	
	return attrs
}

func (p *LEEFParser) parseTimestamp(ts string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"Jan 02 2006 15:04:05",
		"MMM dd yyyy HH:mm:ss",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("невозможно распарсить время: %s", ts)
}

func (p *LEEFParser) mapLEEFSeverityToGOST(attrs map[string]string) string {
	if sev, ok := attrs["sev"]; ok {
		sevLower := strings.ToLower(sev)
		switch {
		case strings.Contains(sevLower, "critical") || strings.Contains(sevLower, "fatal") || sev == "10":
			return models.SeverityCritical
		case strings.Contains(sevLower, "high") || strings.Contains(sevLower, "error") || sev == "8" || sev == "7":
			return models.SeverityHigh
		case strings.Contains(sevLower, "medium") || strings.Contains(sevLower, "warn") || sev == "5" || sev == "6":
			return models.SeverityMedium
		case strings.Contains(sevLower, "low") || sev == "3" || sev == "4":
			return models.SeverityLow
		case strings.Contains(sevLower, "info") || sev == "1" || sev == "2":
			return models.SeverityInfo
		}
	}
	
	return models.SeverityInfo
}

func (p *LEEFParser) categorizeLEEFEvent(eventID string, attrs map[string]string) string {
	eventIDLower := strings.ToLower(eventID)
	
	if cat, ok := attrs["cat"]; ok {
		catLower := strings.ToLower(cat)
		if strings.Contains(catLower, "auth") {
			return models.CategoryAuthentication
		}
		if strings.Contains(catLower, "access") {
			return models.CategoryAccess
		}
		if strings.Contains(catLower, "network") {
			return models.CategoryNetworkEvent
		}
	}
	
	if strings.Contains(eventIDLower, "login") || strings.Contains(eventIDLower, "auth") {
		return models.CategoryAuthentication
	}
	if strings.Contains(eventIDLower, "access") || strings.Contains(eventIDLower, "permission") {
		return models.CategoryAccess
	}
	if strings.Contains(eventIDLower, "modify") || strings.Contains(eventIDLower, "change") {
		return models.CategoryDataModification
	}
	if strings.Contains(eventIDLower, "network") || strings.Contains(eventIDLower, "connection") {
		return models.CategoryNetworkEvent
	}
	if strings.Contains(eventIDLower, "security") || strings.Contains(eventIDLower, "threat") {
		return models.CategorySecurityEvent
	}
	
	return models.CategorySystemEvent
}

func (p *LEEFParser) determineResult(attrs map[string]string) string {
	if result, ok := attrs["result"]; ok {
		resultLower := strings.ToLower(result)
		if strings.Contains(resultLower, "success") || strings.Contains(resultLower, "allow") {
			return models.ResultSuccess
		}
		if strings.Contains(resultLower, "fail") || strings.Contains(resultLower, "deny") {
			return models.ResultFailure
		}
	}
	
	if action, ok := attrs["action"]; ok {
		actionLower := strings.ToLower(action)
		if strings.Contains(actionLower, "allow") || strings.Contains(actionLower, "permit") {
			return models.ResultSuccess
		}
		if strings.Contains(actionLower, "block") || strings.Contains(actionLower, "deny") {
			return models.ResultFailure
		}
	}
	
	return models.ResultUnknown
}

func (p *LEEFParser) getAttributeValue(attrs map[string]string, keys ...string) string {
	for _, key := range keys {
		if val, ok := attrs[key]; ok && val != "" {
			return val
		}
	}
	return ""
}
