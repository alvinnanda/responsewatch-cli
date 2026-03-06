package models

import (
	"encoding/json"
)

// ==================== USER ====================

type User struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	FullName     *string `json:"full_name,omitempty"`
	Organization *string `json:"organization,omitempty"`
	IsActive     bool    `json:"is_active"`
	Role         string  `json:"role,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	AccessToken  string `json:"access_token"`  // Alias for compatibility
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UpdateProfileRequest struct {
	FullName     *string `json:"full_name,omitempty"`
	Organization *string `json:"organization,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ==================== VENDOR GROUP ====================

// PIC represents a Person In Charge with name and phone
type PIC struct {
	Name  string `json:"name"`
	Phone string `json:"phone,omitempty"`
}

type VendorGroup struct {
	ID          int64    `json:"id"`
	GroupName   string   `json:"group_name"`
	VendorPhone string   `json:"vendor_phone,omitempty"`
	PICs        []PIC    `json:"pics"`
	PICNames    []string `json:"pic_names,omitempty"` // For backward compatibility
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// VendorGroupListResponse wraps the list response
type VendorGroupListResponse struct {
	VendorGroups []*VendorGroup `json:"vendor_groups"`
	Pagination   Pagination     `json:"pagination"`
}

type CreateVendorGroupRequest struct {
	GroupName   string   `json:"group_name"`
	VendorPhone string   `json:"vendor_phone,omitempty"`
	PICs        []PIC    `json:"pics,omitempty"`
	PICNames    []string `json:"pic_names,omitempty"` // Deprecated
}

type UpdateVendorGroupRequest struct {
	GroupName   string   `json:"group_name,omitempty"`
	VendorPhone string   `json:"vendor_phone,omitempty"`
	PICs        []PIC    `json:"pics,omitempty"`
	PICNames    []string `json:"pic_names,omitempty"` // Deprecated
}

// ==================== REQUEST ====================

type Request struct {
	ID                    int64     `json:"id"`
	UUID                  string    `json:"uuid"`
	URLToken              string    `json:"url_token"`
	Title                 string    `json:"title"`
	Description           *string   `json:"description,omitempty"`
	FollowupLink          *string   `json:"followup_link,omitempty"`
	Status                string    `json:"status"`
	EmbeddedPICList       []string  `json:"embedded_pic_list"`
	StartPIC              *string   `json:"start_pic,omitempty"`
	EndPIC                *string   `json:"end_pic,omitempty"`
	StartIP               *string   `json:"start_ip,omitempty"`
	EndIP                 *string   `json:"end_ip,omitempty"`
	CreatedAt             string    `json:"created_at"`
	StartedAt             *string   `json:"started_at,omitempty"`
	FinishedAt            *string   `json:"finished_at,omitempty"`
	DurationSeconds       *int      `json:"duration_seconds,omitempty"`
	ResponseTimeSeconds   *int      `json:"response_time_seconds,omitempty"`
	PICIsPublic           *bool     `json:"pic_is_public,omitempty"`
	IsDescriptionSecure   bool      `json:"is_description_secure"`
	VendorGroupID         *int64    `json:"vendor_group_id,omitempty"`
	VendorName            *string   `json:"vendor_name,omitempty"`
	ScheduledTime         *string   `json:"scheduled_time,omitempty"`
	ReopenedAt            *string   `json:"reopened_at,omitempty"`
	ReopenCount           int       `json:"reopen_count"`
	CheckboxIssueMismatch bool      `json:"checkbox_issue_mismatch"`
	ResolutionNotes       *string   `json:"resolution_notes,omitempty"`
	CompletionPhotoURL    *string   `json:"completion_photo_url,omitempty"`
	CompletionPhotoURLs   []string  `json:"completion_photo_urls,omitempty"`
	// Joined fields
	VendorGroup *VendorGroup `json:"vendor_group,omitempty"`
}

// RequestListResponse wraps the list response
type RequestListResponse struct {
	Requests   []Request  `json:"requests"`
	Pagination Pagination `json:"pagination"`
}

type CreateRequestRequest struct {
	Title               string   `json:"title"`
	Description         *string  `json:"description,omitempty"`
	FollowupLink        *string  `json:"followup_link,omitempty"`
	EmbeddedPICList     []string `json:"embedded_pic_list,omitempty"`
	IsDescriptionSecure bool     `json:"is_description_secure"`
	DescriptionPIN      *string  `json:"description_pin,omitempty"`
	VendorGroupID       *int64   `json:"vendor_group_id,omitempty"`
	ScheduledTime       *string  `json:"scheduled_time,omitempty"`
}

type UpdateRequestRequest struct {
	Title               string   `json:"title,omitempty"`
	Description         *string  `json:"description,omitempty"`
	FollowupLink        *string  `json:"followup_link,omitempty"`
	EmbeddedPICList     []string `json:"embedded_pic_list,omitempty"`
	IsDescriptionSecure *bool    `json:"is_description_secure,omitempty"`
	VendorGroupID       *int64   `json:"vendor_group_id,omitempty"`
	ScheduledTime       *string  `json:"scheduled_time,omitempty"`
}

type AssignVendorRequest struct {
	VendorGroupID *int64 `json:"vendor_group_id,omitempty"`
	PIC           string `json:"pic,omitempty"`
}

type RequestStats struct {
	Total      int `json:"total"`
	Waiting    int `json:"waiting"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
}

type RequestStatsPremium struct {
	RequestStats
	AvgResponseTimeMinutes float64     `json:"avg_response_time_minutes"`
	AvgDurationMinutes     float64     `json:"avg_duration_minutes"`
	DailyStats             []DailyStat `json:"daily_stats"`
}

type DailyStat struct {
	Date        string  `json:"date"`
	Total       int     `json:"total"`
	Completed   int     `json:"completed"`
	AvgDuration int     `json:"avg_duration"`
}

type PublicRequestAction struct {
	PIC   string `json:"pic,omitempty"`
	Notes string `json:"notes,omitempty"`
}

type PublicMonitoringResponse struct {
	Username string    `json:"username"`
	Requests []Request `json:"requests"`
}

// ==================== NOTE ====================

// NoteLinkedRequest is a slim struct for note-request relation
type NoteLinkedRequest struct {
	UUID     string `json:"uuid"`
	Title    string `json:"title"`
	URLToken string `json:"url_token"`
}

type Note struct {
	ID              string            `json:"id"` // UUID as string
	UserID          int64             `json:"user_id"`
	Title           string            `json:"title"`
	Content         string            `json:"content"`
	RemindAt        *string           `json:"remind_at,omitempty"` // RFC3339 format
	IsReminder      bool              `json:"is_reminder"`
	ReminderChannel string            `json:"reminder_channel"`
	WebhookURL      *string           `json:"webhook_url,omitempty"`
	WebhookPayload  *string           `json:"webhook_payload,omitempty"`
	WhatsAppPhone   *string           `json:"whatsapp_phone,omitempty"`
	BackgroundColor string            `json:"background_color,omitempty"`
	Tagline         string            `json:"tagline,omitempty"`
	RequestUUID     *string           `json:"request_uuid,omitempty"`
	Request         *NoteLinkedRequest `json:"request,omitempty"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
}

// NoteListResponse wraps the list response
type NoteListResponse struct {
	Notes      []Note     `json:"notes"`
	Pagination Pagination `json:"pagination"`
}

type CreateNoteRequest struct {
	Title           string `json:"title"`
	Content         string `json:"content"`
	RemindAt        *string `json:"remind_at,omitempty"` // RFC3339 format
	IsReminder      bool   `json:"is_reminder"`
	ReminderChannel string `json:"reminder_channel,omitempty"`
	WebhookURL      *string `json:"webhook_url,omitempty"`
	WebhookPayload  *string `json:"webhook_payload,omitempty"`
	WhatsAppPhone   *string `json:"whatsapp_phone,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	Tagline         string `json:"tagline,omitempty"`
	RequestUUID     *string `json:"request_uuid,omitempty"`
}

type UpdateNoteRequest struct {
	Title           string  `json:"title,omitempty"`
	Content         string  `json:"content,omitempty"`
	RemindAt        *string `json:"remind_at,omitempty"`
	IsReminder      bool    `json:"is_reminder"`
	ReminderChannel string  `json:"reminder_channel,omitempty"`
	WebhookURL      *string `json:"webhook_url,omitempty"`
	WebhookPayload  *string `json:"webhook_payload,omitempty"`
	WhatsAppPhone   *string `json:"whatsapp_phone,omitempty"`
	BackgroundColor string  `json:"background_color,omitempty"`
	Tagline         string  `json:"tagline,omitempty"`
	RequestUUID     *string `json:"request_uuid,omitempty"`
}

// ==================== NOTIFICATION ====================

// NotificationMetadata for storing extra information
type NotificationMetadata struct {
	// Status change fields
	OldStatus    string `json:"old_status,omitempty"`
	NewStatus    string `json:"new_status,omitempty"`
	RequestTitle string `json:"request_title,omitempty"`
	RequestToken string `json:"request_token,omitempty"`

	// Reminder fields
	NoteID         string `json:"note_id,omitempty"`
	NoteTitle      string `json:"note_title,omitempty"`
	Channel        string `json:"channel,omitempty"`         // email, whatsapp, webhook
	Recipient      string `json:"recipient,omitempty"`       // email address, phone number, or webhook url
	DeliveryStatus string `json:"delivery_status,omitempty"` // sent, failed
	DeliveryError  string `json:"delivery_error,omitempty"`

	// Plan change fields
	Plan      string `json:"plan,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type Notification struct {
	ID        int64                `json:"id"`
	UserID    int64                `json:"user_id"`
	RequestID *int64               `json:"request_id,omitempty"`
	Type      string               `json:"type"`
	Title     string               `json:"title"`
	Message   string               `json:"message"`
	IsRead    bool                 `json:"is_read"`
	ReadAt    *string              `json:"read_at,omitempty"`
	Metadata  NotificationMetadata `json:"metadata"`
	CreatedAt string               `json:"created_at"`
}

// NotificationListResponse wraps the list response
type NotificationListResponse struct {
	Notifications []Notification `json:"notifications"`
	Pagination    Pagination     `json:"pagination"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

// ==================== COMMON ====================

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse for simple success messages
type SuccessResponse struct {
	Success bool `json:"success"`
}

// RawJSON for flexible JSON parsing
type RawJSON json.RawMessage
