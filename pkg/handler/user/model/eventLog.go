package model

import (
	"log"
	"strings"
	"time"
	"user-service/pkg/config"
)

// EventLog tabloyu temsil eder
type EventLog struct {
	ID          int       `json:"id"`
	UserID      *int      `json:"user_id,omitempty"`
	Email       string    `json:"email,omitempty"`
	SessionID   string    `json:"session_id,omitempty"`
	EventType   string    `json:"event_type"` // LOGIN, LOGOUT, PROFILE_REQUEST
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent,omitempty"`
	Status      string    `json:"status"` // SUCCESS, FAILED
	Reason      string    `json:"reason,omitempty"`
	RequestPath string    `json:"request_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Save logu DB'ye yazar
func (e *EventLog) Save() error {
	//  Status & EventType formatlama
	e.Status = strings.ToUpper(strings.TrimSpace(e.Status))
	e.EventType = strings.ToUpper(strings.TrimSpace(e.EventType))

	// ENUM validasyon
	if e.Status != "SUCCESS" && e.Status != "FAILED" {
		log.Printf("WARN: Invalid status '%s', fallback to 'FAILED'", e.Status)
		e.Status = "FAILED"
	}

	//  Debug log
	log.Printf("DEBUG: Saving EventLog - EventType=%s, Status=%s, Email=%s", e.EventType, e.Status, e.Email)

	query := `
		INSERT INTO event_logs (user_id, email, session_id, event_type, ip, user_agent, status, reason, request_path, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query,
		e.UserID,
		e.Email,
		e.SessionID,
		e.EventType,
		e.IP,
		e.UserAgent,
		e.Status,
		e.Reason,
		e.RequestPath,
		time.Now(),
	)

	if err != nil {
		log.Printf("event log save error: %v", err)
		return err
	}

	log.Printf("[AUDIT] %s event logged for email=%s status=%s", e.EventType, e.Email, e.Status)
	return nil
}

// Helper: Yeni EventLog oluştur
func NewEventLog(userID *int, email, sessionID, eventType, status, reason, ip, userAgent, requestPath string) *EventLog {
	return &EventLog{
		UserID:      userID,
		Email:       email,
		SessionID:   sessionID,
		EventType:   eventType,
		IP:          ip,
		UserAgent:   userAgent,
		Status:      status,
		Reason:      reason,
		RequestPath: requestPath,
		CreatedAt:   time.Now(),
	}
}
