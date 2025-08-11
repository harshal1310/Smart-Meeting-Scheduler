package model

import (
	"encoding/json"
	"testing"
	"time"
)

// getISTTimezone returns IST timezone (UTC+05:30)
func getISTTimezone() *time.Location {
	return time.FixedZone("IST", 5*60*60+30*60)
}

func TestUser(t *testing.T) {
	tests := []struct {
		name string
		user User
	}{
		{
			name: "Basic user",
			user: User{
				ID:       1,
				UserCode: "user1",
				Name:     "John Doe",
			},
		},
		{
			name: "User with events",
			user: User{
				ID:       2,
				UserCode: "user2",
				Name:     "Jane Smith",
				Events: []Event{
					{
						ID:        1,
						EventCode: "event1",
						UserID:    "user2",
						Title:     "Meeting",
						StartTime: time.Date(2025, 8, 9, 9, 0, 0, 0, getISTTimezone()),
						EndTime:   time.Date(2025, 8, 9, 10, 0, 0, 0, getISTTimezone()),
					},
				},
			},
		},
		{
			name: "User with empty UserCode",
			user: User{
				ID:       3,
				UserCode: "",
				Name:     "Empty Code User",
			},
		},
		{
			name: "User with special characters",
			user: User{
				ID:       4,
				UserCode: "user-with_special.chars@domain",
				Name:     "Special User",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.user)
			if err != nil {
				t.Errorf("Failed to marshal user: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaled User
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("Failed to unmarshal user: %v", err)
			}

			// Verify basic fields
			if unmarshaled.ID != tt.user.ID {
				t.Errorf("ID mismatch: expected %d, got %d", tt.user.ID, unmarshaled.ID)
			}
			if unmarshaled.UserCode != tt.user.UserCode {
				t.Errorf("UserCode mismatch: expected %s, got %s", tt.user.UserCode, unmarshaled.UserCode)
			}
			if unmarshaled.Name != tt.user.Name {
				t.Errorf("Name mismatch: expected %s, got %s", tt.user.Name, unmarshaled.Name)
			}

			// Verify events count
			if len(unmarshaled.Events) != len(tt.user.Events) {
				t.Errorf("Events count mismatch: expected %d, got %d", len(tt.user.Events), len(unmarshaled.Events))
			}
		})
	}
}

func TestEvent(t *testing.T) {
	baseTime := time.Date(2025, 8, 9, 9, 0, 0, 0, getISTTimezone())

	tests := []struct {
		name  string
		event Event
	}{
		{
			name: "Basic event",
			event: Event{
				ID:        1,
				EventCode: "event1",
				UserID:    "user1",
				Title:     "Team Meeting",
				StartTime: baseTime,
				EndTime:   baseTime.Add(1 * time.Hour),
			},
		},
		{
			name: "All-day event",
			event: Event{
				ID:        2,
				EventCode: "event2",
				UserID:    "user1",
				Title:     "Conference",
				StartTime: time.Date(2025, 8, 9, 0, 0, 0, 0, getISTTimezone()),
				EndTime:   time.Date(2025, 8, 9, 23, 59, 59, 0, getISTTimezone()),
			},
		},
		{
			name: "Short event",
			event: Event{
				ID:        3,
				EventCode: "event3",
				UserID:    "user2",
				Title:     "Quick Call",
				StartTime: baseTime,
				EndTime:   baseTime.Add(15 * time.Minute),
			},
		},
		{
			name: "Event with special characters in title",
			event: Event{
				ID:        4,
				EventCode: "event4",
				UserID:    "user1",
				Title:     "Meeting: Q&A Session (Remote)",
				StartTime: baseTime,
				EndTime:   baseTime.Add(30 * time.Minute),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.event)
			if err != nil {
				t.Errorf("Failed to marshal event: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaled Event
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("Failed to unmarshal event: %v", err)
			}

			// Verify fields
			if unmarshaled.ID != tt.event.ID {
				t.Errorf("ID mismatch: expected %d, got %d", tt.event.ID, unmarshaled.ID)
			}
			if unmarshaled.EventCode != tt.event.EventCode {
				t.Errorf("EventCode mismatch: expected %s, got %s", tt.event.EventCode, unmarshaled.EventCode)
			}
			if unmarshaled.UserID != tt.event.UserID {
				t.Errorf("UserID mismatch: expected %s, got %s", tt.event.UserID, unmarshaled.UserID)
			}
			if unmarshaled.Title != tt.event.Title {
				t.Errorf("Title mismatch: expected %s, got %s", tt.event.Title, unmarshaled.Title)
			}

			// Verify times (JSON marshaling/unmarshaling might affect precision)
			if !unmarshaled.StartTime.Equal(tt.event.StartTime) {
				t.Errorf("StartTime mismatch: expected %v, got %v", tt.event.StartTime, unmarshaled.StartTime)
			}
			if !unmarshaled.EndTime.Equal(tt.event.EndTime) {
				t.Errorf("EndTime mismatch: expected %v, got %v", tt.event.EndTime, unmarshaled.EndTime)
			}
		})
	}
}

func TestEventDuration(t *testing.T) {
	baseTime := time.Date(2025, 8, 9, 9, 0, 0, 0, getISTTimezone())

	tests := []struct {
		name             string
		event            Event
		expectedDuration time.Duration
	}{
		{
			name: "One hour event",
			event: Event{
				StartTime: baseTime,
				EndTime:   baseTime.Add(1 * time.Hour),
			},
			expectedDuration: 1 * time.Hour,
		},
		{
			name: "30 minute event",
			event: Event{
				StartTime: baseTime,
				EndTime:   baseTime.Add(30 * time.Minute),
			},
			expectedDuration: 30 * time.Minute,
		},
		{
			name: "Zero duration event",
			event: Event{
				StartTime: baseTime,
				EndTime:   baseTime,
			},
			expectedDuration: 0,
		},
		{
			name: "Multi-day event",
			event: Event{
				StartTime: baseTime,
				EndTime:   baseTime.Add(24 * time.Hour),
			},
			expectedDuration: 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := tt.event.EndTime.Sub(tt.event.StartTime)
			if duration != tt.expectedDuration {
				t.Errorf("Expected duration %v, got %v", tt.expectedDuration, duration)
			}
		})
	}
}

func TestEventValidation(t *testing.T) {
	baseTime := time.Date(2025, 8, 9, 9, 0, 0, 0, getISTTimezone())

	tests := []struct {
		name    string
		event   Event
		isValid bool
	}{
		{
			name: "Valid event",
			event: Event{
				ID:        1,
				EventCode: "event1",
				UserID:    "user1",
				Title:     "Valid Meeting",
				StartTime: baseTime,
				EndTime:   baseTime.Add(1 * time.Hour),
			},
			isValid: true,
		},
		{
			name: "Event with end before start",
			event: Event{
				ID:        2,
				EventCode: "event2",
				UserID:    "user1",
				Title:     "Invalid Time Event",
				StartTime: baseTime.Add(1 * time.Hour),
				EndTime:   baseTime,
			},
			isValid: false,
		},
		{
			name: "Event with empty EventCode",
			event: Event{
				ID:        3,
				EventCode: "",
				UserID:    "user1",
				Title:     "No Code Event",
				StartTime: baseTime,
				EndTime:   baseTime.Add(1 * time.Hour),
			},
			isValid: false,
		},
		{
			name: "Event with empty UserID",
			event: Event{
				ID:        4,
				EventCode: "event3",
				UserID:    "",
				Title:     "No User Event",
				StartTime: baseTime,
				EndTime:   baseTime.Add(1 * time.Hour),
			},
			isValid: false,
		},
		{
			name: "Event with empty Title",
			event: Event{
				ID:        5,
				EventCode: "event4",
				UserID:    "user1",
				Title:     "",
				StartTime: baseTime,
				EndTime:   baseTime.Add(1 * time.Hour),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Implement basic validation logic
			isValid := true

			if tt.event.EventCode == "" {
				isValid = false
			}
			if tt.event.UserID == "" {
				isValid = false
			}
			if tt.event.Title == "" {
				isValid = false
			}
			if tt.event.EndTime.Before(tt.event.StartTime) {
				isValid = false
			}

			if isValid != tt.isValid {
				t.Errorf("Expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestEventJSONFormatting(t *testing.T) {
	// Test that time formatting in JSON matches expected format
	baseTime := time.Date(2025, 8, 9, 9, 0, 0, 0, getISTTimezone())

	event := Event{
		ID:        1,
		EventCode: "event1",
		UserID:    "user1",
		Title:     "Test Event",
		StartTime: baseTime,
		EndTime:   baseTime.Add(1 * time.Hour),
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Parse JSON to check format
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Check that time fields are properly formatted
	startTimeStr, ok := jsonMap["startTime"].(string)
	if !ok {
		t.Error("startTime should be a string in JSON")
	} else {
		// Should be in RFC3339 format
		if _, err := time.Parse(time.RFC3339, startTimeStr); err != nil {
			t.Errorf("startTime not in RFC3339 format: %s", startTimeStr)
		}
	}

	endTimeStr, ok := jsonMap["endTime"].(string)
	if !ok {
		t.Error("endTime should be a string in JSON")
	} else {
		// Should be in RFC3339 format
		if _, err := time.Parse(time.RFC3339, endTimeStr); err != nil {
			t.Errorf("endTime not in RFC3339 format: %s", endTimeStr)
		}
	}
}
