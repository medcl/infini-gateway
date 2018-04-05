package timeago

import (
	"time"
)

import "testing"

func check(t *testing.T, d time.Duration, result string) {
	start := time.Now()
	end := time.Now().Add(d)
	got, error := TimeAgoWithTime(start, end)
	if error == nil {
		if got != result {
			t.Errorf("Wrong result: %s", got)
		}
	}
}

func TestThreeHoursAgo(t *testing.T) {
	d, error := time.ParseDuration("-3h")
	if error == nil {
		check(t, d, "3 hours ago")
	}
}

func TestAnHourAgo(t *testing.T) {
	d, error := time.ParseDuration("-1.5h")
	if error == nil {
		check(t, d, "An hour ago")
	}
}

func TestThreeMinutesAgo(t *testing.T) {
	d, error := time.ParseDuration("-3m")
	if error == nil {
		check(t, d, "3 minutes ago")
	}
}

func TestAMinuteAgo(t *testing.T) {
	d, error := time.ParseDuration("-1.2m")
	if error == nil {
		check(t, d, "A minute ago")
	}
}

func TestJustNow(t *testing.T) {
	d, error := time.ParseDuration("-1.2s")
	if error == nil {
		check(t, d, "Just now")
	}
}

func TestFromNow(t *testing.T) {
	d, error := time.ParseDuration("-1.2m")
	if error == nil {
		end := time.Now().Add(d)
		got, err := TimeAgoFromNowWithTime(end)
		if err == nil {
			if got != "A minute ago" {
				t.Errorf("Wrong result: %s", got)
			}
		}
	}
}

func TestFromNowWithString(t *testing.T) {
	d, error := time.ParseDuration("-1.2m")
	if error == nil {
		end := time.Now().Add(d)
		got, err := TimeAgoFromNowWithString(time.RFC3339, end.Format(time.RFC3339))
		if err == nil {
			if got != "A minute ago" {
				t.Errorf("Wrong result: %s", got)
			}
		}
	}
}
