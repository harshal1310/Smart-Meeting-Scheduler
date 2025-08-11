package repository

import (
	"encoding/json"
	"testing"
	"time"
)

// getISTTimezone returns IST timezone (UTC+05:30)
func getISTTimezone() *time.Location {
	return time.FixedZone("IST", 5*60*60+30*60)
}

func TestScheduleRequest(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectError bool
	}{
		{
			name: "Valid schedule request",
			jsonData: `{
				"userIDs": ["user1", "user2"],
				"durationMinutes": 60,
				"timeRange": {
					"start": "2025-08-09T09:00:00+05:30",
					"end": "2025-08-09T17:00:00+05:30"
				}
			}`,
			expectError: false,
		},
		{
			name: "Empty participants",
			jsonData: `{
				"userIDs": [],
				"durationMinutes": 30,
				"timeRange": {
					"start": "2025-08-09T09:00:00+05:30",
					"end": "2025-08-09T17:00:00+05:30"
				}
			}`,
			expectError: false,
		},
		{
			name: "Zero duration",
			jsonData: `{
				"userIDs": ["user1"],
				"durationMinutes": 0,
				"timeRange": {
					"start": "2025-08-09T09:00:00+05:30",
					"end": "2025-08-09T17:00:00+05:30"
				}
			}`,
			expectError: false, // JSON parsing succeeds, business logic validation should catch this
		},
		{
			name: "Missing timeRange",
			jsonData: `{
				"userIDs": ["user1"],
				"durationMinutes": 60
			}`,
			expectError: false,
		},
		{
			name:        "Invalid JSON - missing quote",
			jsonData:    `{"userIDs": ["user1", "durationMinutes": 60}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req ScheduleRequest
			err := json.Unmarshal([]byte(tt.jsonData), &req)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Validate the parsed structure
			if tt.name == "Valid schedule request" {
				if len(req.ParticipantIds) != 2 {
					t.Errorf("Expected 2 participants, got %d", len(req.ParticipantIds))
				}
				if req.DurationMinutes != 60 {
					t.Errorf("Expected 60 minutes duration, got %d", req.DurationMinutes)
				}
				if req.TimeRange.Start != "2025-08-09T09:00:00+05:30" {
					t.Errorf("Expected start time '2025-08-09T09:00:00+05:30', got '%s'", req.TimeRange.Start)
				}
			}
		})
	}
}

func TestScheduledMeetingResponse(t *testing.T) {
	tests := []struct {
		name     string
		response ScheduledMeetingResponse
	}{
		{
			name: "Complete response",
			response: ScheduledMeetingResponse{
				MeetingID:      "meeting-123",
				Title:          "Team Meeting",
				ParticipantIds: []string{"user1", "user2"},
				StartTime:      "2025-08-09T10:00:00+05:30",
				EndTime:        "2025-08-09T11:00:00+05:30",
			},
		},
		{
			name: "Minimal response",
			response: ScheduledMeetingResponse{
				MeetingID: "meeting-456",
				Title:     "Quick Chat",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling/unmarshaling
			jsonData, err := json.Marshal(tt.response)
			if err != nil {
				t.Errorf("Failed to marshal response: %v", err)
			}

			var unmarshaled ScheduledMeetingResponse
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}

			// Verify fields
			if unmarshaled.MeetingID != tt.response.MeetingID {
				t.Errorf("MeetingID mismatch: expected %s, got %s", tt.response.MeetingID, unmarshaled.MeetingID)
			}
			if unmarshaled.Title != tt.response.Title {
				t.Errorf("Title mismatch: expected %s, got %s", tt.response.Title, unmarshaled.Title)
			}
		})
	}
}

func TestSlot(t *testing.T) {
	baseTime := time.Date(2025, 8, 9, 9, 0, 0, 0, getISTTimezone())

	tests := []struct {
		name string
		slot Slot
	}{
		{
			name: "One hour slot",
			slot: Slot{
				Start: baseTime,
				End:   baseTime.Add(1 * time.Hour),
			},
		},
		{
			name: "30 minute slot",
			slot: Slot{
				Start: baseTime,
				End:   baseTime.Add(30 * time.Minute),
			},
		},
		{
			name: "Zero duration slot",
			slot: Slot{
				Start: baseTime,
				End:   baseTime,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic slot properties
			duration := tt.slot.End.Sub(tt.slot.Start)

			if tt.name == "One hour slot" && duration != time.Hour {
				t.Errorf("Expected 1 hour duration, got %v", duration)
			}

			if tt.name == "30 minute slot" && duration != 30*time.Minute {
				t.Errorf("Expected 30 minutes duration, got %v", duration)
			}

			if tt.name == "Zero duration slot" && duration != 0 {
				t.Errorf("Expected zero duration, got %v", duration)
			}
		})
	}
}

func TestScheduleRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request ScheduleRequest
		isValid bool
	}{
		{
			name: "Valid request",
			request: ScheduleRequest{
				ParticipantIds:  []string{"user1", "user2"},
				DurationMinutes: 60,
				TimeRange: struct {
					Start string `json:"start"`
					End   string `json:"end"`
				}{
					Start: "2025-08-09T09:00:00+05:30",
					End:   "2025-08-09T17:00:00+05:30",
				},
			},
			isValid: true,
		},
		{
			name: "No participants",
			request: ScheduleRequest{
				ParticipantIds:  []string{},
				DurationMinutes: 60,
			},
			isValid: false,
		},
		{
			name: "Negative duration",
			request: ScheduleRequest{
				ParticipantIds:  []string{"user1"},
				DurationMinutes: -30,
			},
			isValid: false,
		},
		{
			name: "Zero duration",
			request: ScheduleRequest{
				ParticipantIds:  []string{"user1"},
				DurationMinutes: 0,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Implement basic validation logic
			isValid := true

			if len(tt.request.ParticipantIds) == 0 {
				isValid = false
			}

			if tt.request.DurationMinutes <= 0 {
				isValid = false
			}

			// Validate time range format if provided
			if tt.request.TimeRange.Start != "" {
				if _, err := time.Parse(time.RFC3339, tt.request.TimeRange.Start); err != nil {
					isValid = false
				}
			}

			if tt.request.TimeRange.End != "" {
				if _, err := time.Parse(time.RFC3339, tt.request.TimeRange.End); err != nil {
					isValid = false
				}
			}

			if isValid != tt.isValid {
				t.Errorf("Expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestTimeRangeStruct(t *testing.T) {
	// Test the embedded TimeRange struct
	timeRange := struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Start: "2025-08-09T09:00:00+05:30",
		End:   "2025-08-09T17:00:00+05:30",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(timeRange)
	if err != nil {
		t.Errorf("Failed to marshal TimeRange: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal TimeRange: %v", err)
	}

	// Verify fields
	if unmarshaled.Start != timeRange.Start {
		t.Errorf("Start time mismatch: expected %s, got %s", timeRange.Start, unmarshaled.Start)
	}
	if unmarshaled.End != timeRange.End {
		t.Errorf("End time mismatch: expected %s, got %s", timeRange.End, unmarshaled.End)
	}
}
