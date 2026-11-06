package timeUtils

import (
	"math"
	"time"
)

const (
	brandedDateFormat     = "2 January 2006"
	brandedTimeFormat     = "15:04"
	brandedDateTimeFormat = brandedDateFormat + " at " + brandedTimeFormat
)

var loc *time.Location
var emptyTime *time.Time

func FromUnixEpochTimestamp(ts int64) *time.Time {
	result := time.Unix(ts, 0)
	return &result
}

func ToUnixEpochTimestamp(dt *time.Time) int64 {
	return dt.Unix()
}

func LoadLocation(name string) (*time.Location, error) {
	var err error
	loc, err = time.LoadLocation(name)

	return loc, err
}

func GetLocation() *time.Location {
	return loc
}

func DateTimeToUTC(datetime time.Time) time.Time {
	_, offset := datetime.Zone()

	return datetime.UTC().Add(time.Second * time.Duration(offset))
}

func IsEmpty(datetime time.Time) bool {
	return datetime.Before(time.Date(2000, time.January, 1, 1, 1, 1, 1, GetLocation()))
}

func EmptyLocalisedDateTime() *time.Time {
	if emptyTime == nil {
		emptyTime = &time.Time{}
		localisedTime := emptyTime.In(GetLocation())
		emptyTime = &localisedTime
	}

	return emptyTime
}

func CurrentLocalisedDateTime() *time.Time {
	currentDt := time.Now()
	return LocaliseDateTime(&currentDt)
}

func LocaliseDateTime(dt *time.Time) *time.Time {
	if dt == nil {
		return dt
	}
	localisedTime := dt.In(GetLocation())
	return &localisedTime
}

// Ignores the time. Result is always positive
func DifferenceInDays(dt1 time.Time, dt2 time.Time) int {
	d1 := StartOfDay(dt1)
	d2 := StartOfDay(dt2)
	return int(math.Abs(float64(d1.Sub(d2).Hours() / 24)))
}

// Removes the time from the given dt
func StartOfDay(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())
}

func EndOfDay(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 23, 59, 59, 999999999, dt.Location())
}

func ToBrandedDateFormat(dt time.Time) string {
	return dt.Format(brandedDateFormat)
}

func ToBrandedTimeFormat(dt time.Time) string {
	return dt.Format(brandedTimeFormat)
}

func ToBrandedDateTimeFormat(dt time.Time) string {
	return dt.Format(brandedDateTimeFormat)
}

func Max() time.Time {
	return time.Unix(1<<63-62135596801, 999999999)
}
