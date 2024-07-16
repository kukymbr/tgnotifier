package types

import (
	"sync"
	"time"
)

// TimeInterval is a from-to interval.
type TimeInterval struct {
	From KitchenTime
	To   KitchenTime
}

// Has checks if given time is in the interval.
func (t TimeInterval) Has(now time.Time) bool {
	from := t.From.Time()
	to := t.To.Time()

	if to.Before(from) {
		to = to.AddDate(0, 0, 1)
	}

	return isInInterval(from, to, now) || isInInterval(from, to, now.AddDate(0, 0, 1))
}

func isInInterval(from, to, now time.Time) bool {
	return (from.Before(now) || from.Equal(now)) && (to.After(now) || to.Equal(now))
}

// TimeSchedule is a list of intervals.
type TimeSchedule struct {
	mu        sync.RWMutex
	intervals []TimeInterval
}

// AddInterval adds TimeInterval into the schedule.
func (s *TimeSchedule) AddInterval(interval TimeInterval) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.intervals = append(s.intervals, interval)
}

// Has checks if given time is in one of the scheduled intervals.
func (s *TimeSchedule) Has(t time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, interval := range s.intervals {
		if interval.Has(t) {
			return true
		}
	}

	return false
}
