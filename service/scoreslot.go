package service

import (
	"smart-scheduler/repository"
	"time"
)

func ScoreSlot(slot repository.Slot, userEvents []repository.Slot) int {
	score := 0
	startHour := slot.Start.Hour()

	// Prefer earlier slots (penalize late hours)
	if startHour >= 16 {
		score += 3
	} else if startHour >= 12 {
		score += 2
	} else if startHour >= 9 {
		score += 1
	} else {
		score += 4 // outside working hours
	}

	// Buffer check (penalize if meeting is adjacent)
	buffer := 15 * time.Minute
	for _, event := range userEvents {
		if absDuration(slot.Start.Sub(event.End)) < buffer || absDuration(event.Start.Sub(slot.End)) < buffer {
			score += 2 // no buffer before/after
		}
	}

	// Penalize small gaps before/after
	for _, event := range userEvents {
		gapBefore := slot.Start.Sub(event.End)
		gapAfter := event.Start.Sub(slot.End)
		if gapBefore > 0 && gapBefore < 30*time.Minute {
			score += 1
		}
		if gapAfter > 0 && gapAfter < 30*time.Minute {
			score += 1
		}
	}

	return score
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
