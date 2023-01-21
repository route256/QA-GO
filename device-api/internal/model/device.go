package model

import (
	"time"
)

// Device represents a device.
type Device struct {
	ID        uint64     `db:"id"         json:"id,omitempty"`
	Platform  string     `db:"platform"   json:"platform,omitempty"`
	UserID    uint64     `db:"user_id"    json:"user_id,omitempty"`
	EnteredAt *time.Time `db:"entered_at" json:"entered_at,omitempty"`
	Removed   bool       `db:"removed"    json:"removed,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// EventType represents an event type.
type EventType uint8

const (
	// Created is the event type for when a device is created.
	Created EventType = iota + 1
	// Updated is the event type for when a device is updated.
	Updated
)

// EventStatus represents a device event.
type EventStatus uint8

const (
	// Deferred means the event is not yet processed.
	Deferred EventStatus = iota + 1
	// Processed means the event has been processed.
	Processed
)

// DeviceEvent represents a device event.
type DeviceEvent struct {
	ID        uint64      `db:"id"`
	DeviceID  uint64      `db:"device_id"`
	Type      EventType   `db:"type"`
	Status    EventStatus `db:"status"`
	Device    *Device     `db:"payload"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}
