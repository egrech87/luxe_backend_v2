package durationUtils

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// value must be in the format hh:mm
func Duration_FromHHMM(value string) (time.Duration, error) {
	timeParts := strings.Split(value, ":")
	if len(timeParts) != 2 {
		return time.Duration(0), errors.Errorf("Incorrect duration format")
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return time.Duration(0), err
	}
	if hour < 0 || hour > 23 {
		return time.Duration(0), errors.Errorf("Incorrect duration format")
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return time.Duration(0), err
	}
	if minute < 0 || minute > 59 {
		return time.Duration(0), errors.Errorf("Incorrect duration format")
	}

	return time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute, nil
}
