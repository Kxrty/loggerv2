package models

import "time"

// GOSTEvent представляет событие в формате ГОСТ Р 59710-2021
type GOSTEvent struct {
	// Основные поля согласно ГОСТ Р 59710-2021
	EventID          string    `json:"event_id"`           // Идентификатор события
	Timestamp        time.Time `json:"timestamp"`          // Время события
	Source           Source    `json:"source"`             // Источник события
	Category         string    `json:"category"`           // Категория события
	Severity         string    `json:"severity"`           // Критичность события
	Description      string    `json:"description"`        // Описание события
	AdditionalData   map[string]interface{} `json:"additional_data,omitempty"` // Дополнительные данные
	SubjectAccount   *Account  `json:"subject_account,omitempty"`   // Субъект
	ObjectAccount    *Account  `json:"object_account,omitempty"`    // Объект
	Result           string    `json:"result"`             // Результат события (успех/неуспех)
	Action           string    `json:"action"`             // Действие
}

// Source содержит информацию об источнике события
type Source struct {
	Hostname    string `json:"hostname"`
	IPAddress   string `json:"ip_address,omitempty"`
	Application string `json:"application,omitempty"`
	Process     string `json:"process,omitempty"`
	ProcessID   int    `json:"process_id,omitempty"`
}

// Account представляет учетную запись
type Account struct {
	Username string `json:"username,omitempty"`
	Domain   string `json:"domain,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

// Severity levels согласно ГОСТ
const (
	SeverityCritical = "КРИТИЧЕСКИЙ"
	SeverityHigh     = "ВЫСОКИЙ"
	SeverityMedium   = "СРЕДНИЙ"
	SeverityLow      = "НИЗКИЙ"
	SeverityInfo     = "ИНФОРМАЦИОННЫЙ"
)

// Result values
const (
	ResultSuccess = "УСПЕХ"
	ResultFailure = "НЕУСПЕХ"
	ResultUnknown = "НЕИЗВЕСТНО"
)

// Category values согласно ГОСТ
const (
	CategoryAuthentication    = "АУТЕНТИФИКАЦИЯ"
	CategoryAuthorization     = "АВТОРИЗАЦИЯ"
	CategoryAccess            = "ДОСТУП"
	CategoryDataModification  = "ИЗМЕНЕНИЕ_ДАННЫХ"
	CategorySystemEvent       = "СИСТЕМНОЕ_СОБЫТИЕ"
	CategorySecurityEvent     = "СОБЫТИЕ_БЕЗОПАСНОСТИ"
	CategoryNetworkEvent      = "СЕТЕВОЕ_СОБЫТИЕ"
)
