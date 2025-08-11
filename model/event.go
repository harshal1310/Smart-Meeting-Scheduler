package model

import (
	"time"
)

type Event struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EventCode string    `gorm:"unique;not null" json:"eventCode"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
