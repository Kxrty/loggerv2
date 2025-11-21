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

// CEFParser парсит Common Event Format (CEF) логи
type CEFParser struct{}

// CEF Format: CEF:Version|Device Vendor|Device Product|Device Version|Signature ID|Name|Severity|Extension
var cefPattern = regexp.MustCompile(`^CEF:(\d+)\|([^|]*)\|([^|]*)\|([^|]*)\|([^|]*)\|([^|]*)\|([^|]*)\|(.*)$`)

func NewCEFParser() *CEFParser {
	return &CEFParser{}
}

// Parse парсит CEF сообщение и возвращает GOSTEvent
func (p *CEFParser) Parse(logLine string) (*models.GOSTEvent, error) {
	match := cefPattern.FindStringSubmatch(logLine)
	if match == nil {
		return nil, fmt.Errorf("неверный формат CEF")
	}

	version := match[1]
	deviceVendor := match[2]
	deviceProduct := match[3]
	deviceVersion := match[4]
	signatureID := match[5]
	name := match[6]
	severity := match[7]
	extension := match[8]

	extensions := p.parseExtensions(extension)

	timestamp := time.Now()
	if rt, ok := extensions["rt"]; ok {
		if t, err := p.parseTimestamp(rt); err == nil {
			timestamp = t
		}
	} else if end, ok := extensions["end"]; ok {
		if t, err := p.parseTimestamp(end); err == nil {
			timestamp = t
		}
	}

	event := &models.GOSTEvent{
		EventID:     uuid.New().String(),
		Timestamp:   timestamp,
		Description: name,
		Source: models.Source{
			Hostname:    p.getExtensionValue(extensions, "dvc", "shost", "dvchost"),
			Application: fmt.Sprintf("%s %s", deviceVendor, deviceProduct),
			IPAddress:   p.getExtensionValue(extensions, "src", "dst"),
		},
		Severity:       p.mapCEFSeverityToGOST(severity),
		Category:       p.categorizeCEFEvent(signatureID, name, extensions),
		Result:         p.determineResult(extensions),
		AdditionalData: make(map[string]interface{}),
	}

	event.AdditionalData["cef_version"] = version
	event.AdditionalData["device_vendor"] = deviceVendor
	event.AdditionalData["device_product"] = deviceProduct
	event.AdditionalData["device_version"] = deviceVersion
	event.AdditionalData["signature_id"] = signatureID
	event.AdditionalData["cef_severity"] = severity

	if suser, ok := extensions["suser"]; ok {
		event.SubjectAccount = &models.Account{
			Username: suser,
			Domain:   extensions["sdomain"],
		}
	}

	if duser, ok := extensions["duser"]; ok {
		event.ObjectAccount = &models.Account{
			Username: duser,
			Domain:   extensions["ddomain"],
		}
	}

	for k, v := range extensions {
		event.AdditionalData["cef_"+k] = v
	}

	if act, ok := extensions["act"]; ok {
		event.Action = act
	}

	return event, nil
}

func (p *CEFParser) parseExtensions(extension string) map[string]string {
	extensions := make(map[string]string)
	
	parts := strings.Split(extension, " ")
	for _, part := range parts {
		if idx := strings.Index(part, "="); idx != -1 {
			key := part[:idx]
			value := part[idx+1:]
			extensions[key] = value
		}
	}
	
	return extensions
}

func (p *CEFParser) parseTimestamp(ts string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"Jan 02 2006 15:04:05",
		"2006-01-02 15:04:05",
		"1136239445000",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, ts); err == nil {
			return t, nil
		}
	}
	
	if msec, err := strconv.ParseInt(ts, 10, 64); err == nil {
		return time.Unix(0, msec*int64(time.Millisecond)), nil
	}
	
	return time.Time{}, fmt.Errorf("невозможно распарсить время: %s", ts)
}

func (p *CEFParser) mapCEFSeverityToGOST(severity string) string {
	sev, err := strconv.Atoi(severity)
	if err != nil {
		return models.SeverityInfo
	}
	
	switch {
	case sev >= 8:
		return models.SeverityCritical
	case sev >= 6:
		return models.SeverityHigh
	case sev >= 4:
		return models.SeverityMedium
	case sev >= 2:
		return models.SeverityLow
	default:
		return models.SeverityInfo
	}
}

func (p *CEFParser) categorizeCEFEvent(signatureID, name string, extensions map[string]string) string {
	nameLower := strings.ToLower(name)
	
	if strings.Contains(nameLower, "login") || strings.Contains(nameLower, "logon") ||
		strings.Contains(nameLower, "authentication") {
		return models.CategoryAuthentication
	}
	if strings.Contains(nameLower, "access") || strings.Contains(nameLower, "denied") ||
		strings.Contains(nameLower, "permission") {
		return models.CategoryAccess
	}
	if strings.Contains(nameLower, "modify") || strings.Contains(nameLower, "change") ||
		strings.Contains(nameLower, "update") || strings.Contains(nameLower, "delete") {
		return models.CategoryDataModification
	}
	if strings.Contains(nameLower, "network") || strings.Contains(nameLower, "connection") ||
		strings.Contains(nameLower, "firewall") {
		return models.CategoryNetworkEvent
	}
	if strings.Contains(nameLower, "security") || strings.Contains(nameLower, "threat") ||
		strings.Contains(nameLower, "attack") || strings.Contains(nameLower, "malware") {
		return models.CategorySecurityEvent
	}
	
	return models.CategorySystemEvent
}

func (p *CEFParser) determineResult(extensions map[string]string) string {
	if outcome, ok := extensions["outcome"]; ok {
		outcomeLower := strings.ToLower(outcome)
		if strings.Contains(outcomeLower, "success") {
			return models.ResultSuccess
		}
		if strings.Contains(outcomeLower, "fail") || strings.Contains(outcomeLower, "deny") {
			return models.ResultFailure
		}
	}
	
	if act, ok := extensions["act"]; ok {
		actLower := strings.ToLower(act)
		if strings.Contains(actLower, "allow") || strings.Contains(actLower, "permit") {
			return models.ResultSuccess
		}
		if strings.Contains(actLower, "block") || strings.Contains(actLower, "deny") {
			return models.ResultFailure
		}
	}
	
	return models.ResultUnknown
}

func (p *CEFParser) getExtensionValue(extensions map[string]string, keys ...string) string {
	for _, key := range keys {
		if val, ok := extensions[key]; ok && val != "" {
			return val
		}
	}
	return ""
}
