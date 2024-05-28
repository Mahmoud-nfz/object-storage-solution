package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func hmsToSeconds(hms string) (int, error) {
    parts := strings.Split(hms, ":")
    if len(parts) != 3 {
        return 0, fmt.Errorf("invalid time format")
    }
    hours, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, err
    }
    minutes, err := strconv.Atoi(parts[1])
    if err != nil {
        return 0, err
    }
    seconds, err := strconv.Atoi(parts[2])
    if err != nil {
        return 0, err
    }
    return hours*3600 + minutes*60 + seconds, nil
}

func secondsToHMS(seconds int) string {
    hours := seconds / 3600
    minutes := (seconds % 3600) / 60
    secs := seconds % 60
    return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
