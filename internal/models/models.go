package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Organization string    `json:"organization"`
	IsActive     bool      `json:"is_active"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	User         User   `json:"user"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents a refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FullName     string `json:"full_name,omitempty"`
	Organization string `json:"organization,omitempty"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// Request represents a ticket/request in the system
type Request struct {
	ID                   int       `json:"id"`
	UUID                 string    `json:"uuid"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Status               string    `json:"status"` // waiting, in_progress, done
	RefLink              string    `json:"ref_link"`
	GroupID              *int      `json:"group_id"`
	GroupName            string    `json:"group_name"`
	PIC                  string    `json:"pic"`
	StartPIC             string    `json:"start_pic"`
	EndPIC               string    `json:"end_pic"`
	StartIP              string    `json:"start_ip"`
	EndIP                string    `json:"end_ip"`
	IsPinned             bool      `json:"is_pinned"`
	IsScheduled          bool      `json:"is_scheduled"`
	ScheduledDate        *string   `json:"scheduled_date"`
	URLToken             string    `json:"url_token"`
	PublicURL            string    `json:"public_url"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	StartedAt            *time.Time `json:"started_at"`
	FinishedAt           *time.Time `json:"finished_at"`
	DurationSeconds      *int      `json:"duration_seconds"`
	ResponseTimeSeconds  *int      `json:"response_time_seconds"`
}

// CreateRequestRequest represents a request creation payload
type CreateRequestRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	RefLink       string `json:"ref_link,omitempty"`
	GroupID       *int   `json:"group_id,omitempty"`
	PIC           string `json:"pic,omitempty"`
	IsPinned      bool   `json:"is_pinned"`
	IsScheduled   bool   `json:"is_scheduled"`
	ScheduledDate string `json:"scheduled_date,omitempty"`
}

// UpdateRequestRequest represents a request update payload
type UpdateRequestRequest struct {
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	RefLink       string `json:"ref_link,omitempty"`
	GroupID       *int   `json:"group_id,omitempty"`
	PIC           string `json:"pic,omitempty"`
	IsPinned      *bool  `json:"is_pinned,omitempty"`
	IsScheduled   *bool  `json:"is_scheduled,omitempty"`
	ScheduledDate string `json:"scheduled_date,omitempty"`
}

// AssignVendorRequest represents a vendor assignment payload
type AssignVendorRequest struct {
	GroupID *int   `json:"group_id,omitempty"`
	PIC     string `json:"pic,omitempty"`
}

// RequestStats represents statistics for requests
type RequestStats struct {
	Total      int `json:"total"`
	Waiting    int `json:"waiting"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
}

// RequestStatsPremium represents premium statistics
type RequestStatsPremium struct {
	RequestStats
	AvgResponseTimeMinutes float64            `json:"avg_response_time_minutes"`
	AvgDurationMinutes     float64            `json:"avg_duration_minutes"`
	DailyStats             []DailyStat        `json:"daily_stats"`
}

// DailyStat represents daily statistics
type DailyStat struct {
	Date        string `json:"date"`
	Total       int    `json:"total"`
	Completed   int    `json:"completed"`
	AvgDuration int    `json:"avg_duration"` // minutes
}

// VendorGroup represents a vendor group
type VendorGroup struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Phone     string   `json:"phone"`
	PICs      []string `json:"pics"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// CreateVendorGroupRequest represents a vendor group creation payload
type CreateVendorGroupRequest struct {
	Name  string   `json:"name"`
	Phone string   `json:"phone,omitempty"`
	PICs  []string `json:"pics"`
}

// UpdateVendorGroupRequest represents a vendor group update payload
type UpdateVendorGroupRequest struct {
	Name  string   `json:"name,omitempty"`
	Phone string   `json:"phone,omitempty"`
	PICs  []string `json:"pics,omitempty"`
}

// Note represents a user note
type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Color     string    `json:"color"`
	Reminder  *time.Time `json:"reminder"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateNoteRequest represents a note creation payload
type CreateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Color    string `json:"color,omitempty"`
	Reminder string `json:"reminder,omitempty"`
}

// UpdateNoteRequest represents a note update payload
type UpdateNoteRequest struct {
	Title    string `json:"title,omitempty"`
	Content  string `json:"content,omitempty"`
	Color    string `json:"color,omitempty"`
	Reminder string `json:"reminder,omitempty"`
}

// Notification represents a notification
type Notification struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

// PublicRequestAction represents a public request action (start/finish)
type PublicRequestAction struct {
	PIC   string `json:"pic,omitempty"`
	Notes string `json:"notes,omitempty"`
}

// PublicMonitoringResponse represents public monitoring data
type PublicMonitoringResponse struct {
	Username string    `json:"username"`
	Requests []Request `json:"requests"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
