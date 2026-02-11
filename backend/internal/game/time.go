package game

import (
	"stardew-agent/internal/models"
)

// TimeSystem manages game time
type TimeSystem struct {
	// Time scale: 1 real second = X game minutes
	timeScale int
}

// NewTimeSystem creates a new time system
func NewTimeSystem() *TimeSystem {
	return &TimeSystem{
		timeScale: 10,
	}
}

// AdvanceTime advances game time by one tick
func (ts *TimeSystem) AdvanceTime(time *models.TimeState) {
	time.Minute++

	// Handle minute overflow
	if time.Minute >= 60 {
		time.Minute = 0
		time.Hour++

		// Handle hour overflow (day ends at 2 AM, starts at 6 AM)
		if time.Hour >= 26 { // 2 AM = hour 26 (since day starts at 6 AM)
			ts.AdvanceDay(time)
		}
	}
}

// AdvanceDay advances to the next day
func (ts *TimeSystem) AdvanceDay(time *models.TimeState) {
	time.Day++
	time.Hour = 6
	time.Minute = 0

	// Handle season overflow (28 days per season)
	if time.Day > 28 {
		time.Day = 1
		ts.AdvanceSeason(time)
	}
}

// AdvanceSeason advances to the next season
func (ts *TimeSystem) AdvanceSeason(time *models.TimeState) {
	switch time.Season {
	case models.SeasonSpring:
		time.Season = models.SeasonSummer
	case models.SeasonSummer:
		time.Season = models.SeasonFall
	case models.SeasonFall:
		time.Season = models.SeasonWinter
	case models.SeasonWinter:
		time.Season = models.SeasonSpring
		time.Year++
	}
}

// GetFormattedTime returns a formatted time string
func (ts *TimeSystem) GetFormattedTime(time models.TimeState) string {
	hour := time.Hour
	period := "AM"

	// Convert from 24-hour format (starting at 6 AM = hour 6)
	if hour >= 12 && hour < 24 {
		period = "PM"
		if hour > 12 {
			hour -= 12
		}
	} else if hour >= 24 {
		// After midnight (hour 24 = 12 AM, hour 25 = 1 AM, hour 26 = 2 AM)
		period = "AM"
		hour -= 24
		if hour == 0 {
			hour = 12
		}
	} else if hour == 0 {
		hour = 12
	}

	return string(time.Season) + " " + itoa(time.Day) + ", Year " + itoa(time.Year) + " - " + itoa(hour) + ":" + padZero(time.Minute) + " " + period
}

// GetSeasonDisplayName returns a nice season name
func (ts *TimeSystem) GetSeasonDisplayName(season models.Season) string {
	switch season {
	case models.SeasonSpring:
		return "Spring"
	case models.SeasonSummer:
		return "Summer"
	case models.SeasonFall:
		return "Fall"
	case models.SeasonWinter:
		return "Winter"
	default:
		return string(season)
	}
}

// IsNight checks if it's currently night time
func (ts *TimeSystem) IsNight(time models.TimeState) bool {
	// Night is from 6 PM (hour 18) to 6 AM (hour 6)
	return time.Hour >= 18 || time.Hour < 6
}

// IsMorning checks if it's currently morning
func (ts *TimeSystem) IsMorning(time models.TimeState) bool {
	return time.Hour >= 6 && time.Hour < 12
}

// IsAfternoon checks if it's currently afternoon
func (ts *TimeSystem) IsAfternoon(time models.TimeState) bool {
	return time.Hour >= 12 && time.Hour < 18
}

// GetDayPhase returns the current phase of day
func (ts *TimeSystem) GetDayPhase(time models.TimeState) string {
	if time.Hour >= 6 && time.Hour < 12 {
		return "morning"
	} else if time.Hour >= 12 && time.Hour < 18 {
		return "afternoon"
	} else if time.Hour >= 18 && time.Hour < 24 {
		return "evening"
	} else {
		return "night"
	}
}

// GetTimeUntil checks how many game minutes until a specific hour
func (ts *TimeSystem) GetTimeUntil(time models.TimeState, targetHour int) int {
	currentMinutes := time.Hour*60 + time.Minute
	targetMinutes := targetHour * 60

	if targetMinutes < currentMinutes {
		// Target is tomorrow
		targetMinutes += 24 * 60
	}

	return targetMinutes - currentMinutes
}

// GetNextSeasonChange returns days until season change
func (ts *TimeSystem) GetNextSeasonChange(time models.TimeState) int {
	return 28 - time.Day
}

// DaysInSeason returns the number of days in a season
func (ts *TimeSystem) DaysInSeason() int {
	return 28
}

// Helper functions
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var negative bool
	if n < 0 {
		negative = true
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}

func padZero(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}
