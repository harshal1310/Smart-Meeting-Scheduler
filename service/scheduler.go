package service

import (
	"errors"
	"log"
	"smart-scheduler/db"
	database "smart-scheduler/db"
	"smart-scheduler/model"
	"smart-scheduler/repository"
	"sort"
	"time"

	"gorm.io/gorm"
)

func ScheduleEvent(req repository.ScheduleRequest) (*repository.ScheduledMeetingResponse, error) {
	startTime, _ := time.Parse(time.RFC3339, req.TimeRange.Start)
	endTime, _ := time.Parse(time.RFC3339, req.TimeRange.End)
	slotDuration := time.Duration(req.DurationMinutes) * time.Minute

	eventMap := make(map[string][]repository.Slot)
	for _, userId := range req.ParticipantIds {
		var events []model.Event
		// Fix: Query for events that overlap with the time range
		// An event overlaps if: event_start < range_end AND event_end > range_start
		database.DB.Where("user_id = ? AND start_time < ? AND end_time > ?", userId, endTime, startTime).Find(&events)

		log.Printf("User %s has %d existing events in range %v to %v", userId, len(events), startTime, endTime)
		for _, e := range events {
			log.Printf("  - Event: %s (%v to %v)", e.Title, e.StartTime, e.EndTime)
			eventMap[userId] = append(eventMap[userId], repository.Slot{Start: e.StartTime, End: e.EndTime})
		}
	}

	// generate potential slots in 30-min steps within range
	var candidateSlots []repository.Slot
	log.Printf("Generating candidate slots from %v to %v with duration %v", startTime, endTime, slotDuration)

	for t := startTime; t.Add(slotDuration).Before(endTime) || t.Add(slotDuration).Equal(endTime); t = t.Add(30 * time.Minute) {
		end := t.Add(slotDuration)
		valid := true
		log.Printf("Checking slot %v to %v", t, end)

		for userId, events := range eventMap {
			for _, ev := range events {
				if ev.Start.Before(end) && ev.End.After(t) {
					log.Printf("  Slot conflicts with %s's event: %v to %v", userId, ev.Start, ev.End)
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid {
			log.Printf("  Slot %v to %v is VALID", t, end)
			candidateSlots = append(candidateSlots, repository.Slot{Start: t, End: end})
		}
	}

	log.Printf("Found %d candidate slots", len(candidateSlots))

	if len(candidateSlots) == 0 {
		return nil, errors.New("no available time slot found for all participants")
	}

	type scoredSlot struct {
		repository.Slot
		score int
	}
	var slots []scoredSlot
	for _, slot := range candidateSlots {
		s := 0
		for _, events := range eventMap {
			s += ScoreSlot(slot, events)
		}
		slots = append(slots, scoredSlot{slot, s})
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].score < slots[j].score
	})

	chosen := slots[0].Slot
	meetingCode := "meeting-" + time.Now().Format("20060102150405")

	// Use provided title or default to "New Meeting"
	meetingTitle := req.Title
	if meetingTitle == "" {
		meetingTitle = "New Meeting"
	}

	for _, userId := range req.ParticipantIds {
		database.DB.Create(&model.Event{
			EventCode: meetingCode + "-" + userId,
			UserID:    userId,
			Title:     meetingTitle,
			StartTime: chosen.Start,
			EndTime:   chosen.End,
		})
	}

	return &repository.ScheduledMeetingResponse{
		MeetingID:      meetingCode,
		Title:          meetingTitle,
		ParticipantIds: req.ParticipantIds,
		StartTime:      chosen.Start.Format(time.RFC3339),
		EndTime:        chosen.End.Format(time.RFC3339),
	}, nil
}

func GetCalendarEvents(userId, start, end string) ([]model.Event, error) {

	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil && start != "" {
		return nil, errors.New("invalid start time format")
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil && end != "" {
		return nil, errors.New("invalid end time format")
	}

	// Debug logging
	log.Printf("GET /calendar/%s - start: %s, end: %s", userId, start, end)
	log.Printf("Parsed times - start: %v, end: %v", startTime, endTime)

	// Validate time range - if start time is after end time, return error
	if !startTime.IsZero() && !endTime.IsZero() && startTime.After(endTime) {
		log.Printf("Invalid time range: start (%v) is after end (%v)", startTime, endTime)
		return nil, errors.New("invalid time range: start time cannot be after end time")
	}

	var events []model.Event
	// Find events that overlap with the requested time range
	// An event overlaps if: event_start < range_end AND event_end > range_start
	var result *gorm.DB
	if start == "" || end == "" {
		result = db.DB.Where("user_id = ?", userId).Find(&events)

	} else {
		result = db.DB.Where("user_id = ? AND start_time < ? AND end_time > ?", userId, endTime, startTime).Find(&events)
	}

	log.Printf("Query result - found %d events, error: %v", len(events), result.Error)
	for _, event := range events {
		log.Printf("Event: %s (%s to %s)", event.Title, event.StartTime.Format(time.RFC3339), event.EndTime.Format(time.RFC3339))
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}
