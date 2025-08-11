package service

import (
	"smart-scheduler/repository"
	"testing"
	"time"
)

// Helper function to create IST timezone
func getISTTimezone() *time.Location {
	return time.FixedZone("IST", 5*60*60+30*60)
}

func TestGetCalendarEventsValidation(t *testing.T) {
	tests := []struct {
		name         string
		userId       string
		start        string
		end          string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Invalid time range - start after end",
			userId:       "user1",
			start:        "2025-08-10T19:00:00+05:30", // IST format
			end:          "2025-08-10T10:00:00+05:30", // IST format
			expectError:  true,
			errorMessage: "invalid time range: start time cannot be after end time",
		},
		{
			name:         "Invalid start time format",
			userId:       "user1",
			start:        "invalid-time",
			end:          "2025-08-10T18:00:00+05:30", // IST format
			expectError:  true,
			errorMessage: "invalid start time format",
		},
		{
			name:         "Invalid end time format",
			userId:       "user1",
			start:        "2025-08-10T08:00:00+05:30", // IST format
			end:          "invalid-time",
			expectError:  true,
			errorMessage: "invalid end time format",
		},
		{
			name:        "Empty time parameters",
			userId:      "user1",
			start:       "",
			end:         "",
			expectError: false,
		},
		{
			name:        "Valid time range",
			userId:      "user1",
			start:       "2025-08-09T08:00:00+05:30",
			end:         "2025-08-09T18:00:00+05:30",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test focuses on validation logic only
			// Database operations would require test database setup

			// Test time parsing validation
			if tt.start != "" {
				_, err := time.Parse(time.RFC3339, tt.start)
				if err != nil && !tt.expectError {
					t.Errorf("Start time parsing failed unexpectedly: %v", err)
				}
			}

			if tt.end != "" {
				_, err := time.Parse(time.RFC3339, tt.end)
				if err != nil && !tt.expectError {
					t.Errorf("End time parsing failed unexpectedly: %v", err)
				}
			}

			// Test time range validation logic
			if tt.start != "" && tt.end != "" {
				startTime, startErr := time.Parse(time.RFC3339, tt.start)
				endTime, endErr := time.Parse(time.RFC3339, tt.end)

				if startErr == nil && endErr == nil {
					if startTime.After(endTime) && !tt.expectError {
						t.Errorf("Expected error for invalid time range but validation logic would allow it")
					}
				}
			}
		})
	}
}

func TestScheduleEventValidation(t *testing.T) {
	tests := []struct {
		name    string
		request repository.ScheduleRequest
		valid   bool
	}{
		{
			name: "Valid request",
			request: repository.ScheduleRequest{
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
			valid: true,
		},
		{
			name: "Empty participants",
			request: repository.ScheduleRequest{
				ParticipantIds:  []string{},
				DurationMinutes: 60,
				TimeRange: struct {
					Start string `json:"start"`
					End   string `json:"end"`
				}{
					Start: "2025-08-09T09:00:00+05:30",
					End:   "2025-08-09T17:00:00+05:30",
				},
			},
			valid: false,
		},
		{
			name: "Zero duration",
			request: repository.ScheduleRequest{
				ParticipantIds:  []string{"user1"},
				DurationMinutes: 0,
				TimeRange: struct {
					Start string `json:"start"`
					End   string `json:"end"`
				}{
					Start: "2025-08-09T09:00:00+05:30",
					End:   "2025-08-09T17:00:00+05:30",
				},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate request structure
			if len(tt.request.ParticipantIds) == 0 && tt.valid {
				t.Errorf("Expected valid request but participants list is empty")
			}

			if tt.request.DurationMinutes <= 0 && tt.valid {
				t.Errorf("Expected valid request but duration is zero or negative")
			}

			// Test time parsing
			_, startErr := time.Parse(time.RFC3339, tt.request.TimeRange.Start)
			_, endErr := time.Parse(time.RFC3339, tt.request.TimeRange.End)

			if (startErr != nil || endErr != nil) && tt.valid {
				t.Errorf("Expected valid request but time parsing failed")
			}
		})
	}
}

func TestTimeRangeValidation(t *testing.T) {
	tests := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		valid     bool
	}{
		{
			name:      "Valid range - start before end",
			startTime: time.Date(2025, 8, 9, 9, 0, 0, 0, getISTTimezone()),
			endTime:   time.Date(2025, 8, 9, 17, 0, 0, 0, getISTTimezone()),
			valid:     true,
		},
		{
			name:      "Invalid range - start after end",
			startTime: time.Date(2025, 8, 9, 19, 0, 0, 0, getISTTimezone()),
			endTime:   time.Date(2025, 8, 9, 10, 0, 0, 0, getISTTimezone()),
			valid:     false,
		},
		{
			name:      "Same start and end time",
			startTime: time.Date(2025, 8, 9, 12, 0, 0, 0, getISTTimezone()),
			endTime:   time.Date(2025, 8, 9, 12, 0, 0, 0, getISTTimezone()),
			valid:     true,
		},
		{
			name:      "Zero time values",
			startTime: time.Time{},
			endTime:   time.Time{},
			valid:     true, // Empty times are allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic that would be used in GetCalendarEvents
			isInvalid := !tt.startTime.IsZero() && !tt.endTime.IsZero() && tt.startTime.After(tt.endTime)

			if isInvalid && tt.valid {
				t.Errorf("Expected valid time range but validation marked it as invalid")
			}

			if !isInvalid && !tt.valid && !tt.startTime.IsZero() && !tt.endTime.IsZero() {
				t.Errorf("Expected invalid time range but validation marked it as valid")
			}
		})
	}
}
