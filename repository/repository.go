package repository

import (
	"smart-scheduler/model"
	"time"

	"gorm.io/gorm"
)

// Request/Response types
type ScheduleRequest struct {
	Title           string   `json:"title"`
	ParticipantIds  []string `json:"userIDs"`
	DurationMinutes int      `json:"durationMinutes"`
	TimeRange       struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"timeRange"`
}

type ScheduledMeetingResponse struct {
	MeetingID      string   `json:"meetingId"`
	Title          string   `json:"title"`
	ParticipantIds []string `json:"participantIds"`
	StartTime      string   `json:"startTime"`
	EndTime        string   `json:"endTime"`
}

// Service types
type Slot struct {
	Start time.Time
	End   time.Time
}

// Dummy data for testing
func CreateDummyData(db *gorm.DB) error {
	// Clear existing test data first
	db.Where("event_code LIKE ?", "event%").Delete(&model.Event{})
	db.Where("user_code LIKE ?", "user%").Delete(&model.User{})

	// Create dummy users
	users := []model.User{
		{
			UserCode: "user1",
			Name:     "Alice Johnson",
		},
		{
			UserCode: "user2",
			Name:     "Bob Smith",
		},
		{
			UserCode: "user3",
			Name:     "Charlie Brown",
		},
		{
			UserCode: "user4",
			Name:     "Diana Prince",
		},
		{
			UserCode: "user5",
			Name:     "Eve Wilson",
		},
	}

	// Insert users
	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	// Create dummy events for testing with IST timezone
	istTimezone := time.FixedZone("IST", 5*60*60+30*60) // 5 hours 30 minutes offset

	// Base time: August 9, 2025, 9:00 AM IST
	baseTime := time.Date(2025, 8, 9, 9, 0, 0, 0, istTimezone) // 9:00 AM IST
	events := []model.Event{

		{
			EventCode: "event1",
			UserID:    "user1",
			Title:     "Team Standup",
			StartTime: baseTime,                       // 9:00 AM IST
			EndTime:   baseTime.Add(30 * time.Minute), // 9:30 AM IST
		},
		{
			EventCode: "event2",
			UserID:    "user1",
			Title:     "Project Review",
			StartTime: baseTime.Add(3 * time.Hour), // 12:00 PM IST
			EndTime:   baseTime.Add(4 * time.Hour), // 1:00 PM IST
		},

		{
			EventCode: "event3",
			UserID:    "user2",
			Title:     "Client Call",
			StartTime: baseTime.Add(2 * time.Hour),                // 11:00 AM IST
			EndTime:   baseTime.Add(2*time.Hour + 45*time.Minute), // 11:45 AM IST
		},
		{
			EventCode: "event4",
			UserID:    "user2",
			Title:     "Code Review",
			StartTime: baseTime.Add(5 * time.Hour),                // 2:00 PM IST
			EndTime:   baseTime.Add(5*time.Hour + 30*time.Minute), // 2:30 PM IST
		},

		{
			EventCode: "event5",
			UserID:    "user3",
			Title:     "Design Meeting",
			StartTime: baseTime.Add(1*time.Hour + 30*time.Minute), // 10:30 AM IST
			EndTime:   baseTime.Add(2*time.Hour + 30*time.Minute), // 11:30 AM IST
		},
		{
			EventCode: "event6",
			UserID:    "user3",
			Title:     "Sprint Planning",
			StartTime: baseTime.Add(6 * time.Hour), // 3:00 PM IST
			EndTime:   baseTime.Add(7 * time.Hour), // 4:00 PM IST
		},

		{
			EventCode: "event7",
			UserID:    "user4",
			Title:     "1:1 Meeting",
			StartTime: baseTime.Add(4 * time.Hour),                // 1:00 PM IST
			EndTime:   baseTime.Add(4*time.Hour + 30*time.Minute), // 1:30 PM IST
		},

		{
			EventCode: "event8",
			UserID:    "user5",
			Title:     "Training Session",
			StartTime: baseTime.Add(7 * time.Hour), // 4:00 PM IST
			EndTime:   baseTime.Add(8 * time.Hour), // 5:00 PM IST
		},
	}

	for _, event := range events {
		if err := db.Create(&event).Error; err != nil {
			return err
		}
	}

	return nil
}
