package utils

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(durationStr string) (time.Duration, error) {
	// Split the duration string into number and unit
	var unit string
	var numberStr string
	for i, char := range durationStr {
		if char >= '0' && char <= '9' {
			numberStr = durationStr[:i+1]
			unit = durationStr[i+1:]
			break
		}
	}

	if unit == "" || numberStr == "" {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	// Convert the number part to an integer
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	// Map unit to time.Duration
	var duration time.Duration
	switch strings.ToLower(unit) {
	case "ms":
		duration = time.Duration(number) * time.Millisecond
	case "s":
		duration = time.Duration(number) * time.Second
	case "m":
		duration = time.Duration(number) * time.Minute
	case "h":
		duration = time.Duration(number) * time.Hour
	case "d":
		duration = time.Duration(number) * 24 * time.Hour
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unit)
	}

	return duration, nil
}
