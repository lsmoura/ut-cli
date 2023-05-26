package strftime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStrTimeTokens(t *testing.T) {
	tests := []struct {
		format   string
		expected []string
	}{
		{"%Y-%m-%d %H:%M:%S", []string{"%Y", "-", "%m", "-", "%d", " ", "%H", ":", "%M", ":", "%S"}},
		{"%A%B%C", []string{"%A", "%B", "%C"}},
		{"%A %B %C", []string{"%A", " ", "%B", " ", "%C"}},
		{"%Afoo%Bbar %C", []string{"%A", "foo", "%B", "bar ", "%C"}},
		{"%-Afoo%Bbar %C", []string{"%-A", "foo", "%B", "bar ", "%C"}},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, StrTimeTokens(test.format))
	}
}

func TestStrftime(t *testing.T) {
	tests := []struct {
		format   string
		time     time.Time
		expected string
	}{
		{"%Y-%m-%d %H:%M:%S", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "2017-01-02 03:04:05"},
		{"%Y-%m-%d %H:%M:%S", time.Date(2017, 1, 2, 3, 4, 5, 0, time.FixedZone("JST", 9*60*60)), "2017-01-02 03:04:05"},
		{"Now it's %I:%M%p.", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "Now it's 03:04AM."},
		{"%a", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "Mon"},
		{"%A", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "Monday"},
		{"%w", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1"},
		{"%d", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "02"},
		{"%-d", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "2"},
		{"%b", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "Jan"},
		{"%B", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "January"},
		{"%m", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "01"},
		{"%-m", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1"},
		{"%y", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "17"},
		{"%Y", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "2017"},
		{"%H", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "03"},
		{"%-H", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "3"},
		{"%I", time.Date(2017, 1, 2, 15, 4, 5, 0, time.UTC), "03"},
		{"%I", time.Date(2017, 1, 2, 0, 4, 5, 0, time.UTC), "12"},
		{"%-I", time.Date(2017, 1, 2, 15, 4, 5, 0, time.UTC), "3"},
		{"%p", time.Date(2017, 1, 2, 15, 4, 5, 0, time.UTC), "PM"},
		{"%p", time.Date(2017, 1, 2, 2, 4, 5, 0, time.UTC), "AM"},
		{"%M", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "04"},
		{"%-M", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "4"},
		{"%S", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "05"},
		{"%-S", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "5"},
		{"%f", time.Date(2017, 1, 2, 3, 4, 5, 678*1000*1000, time.UTC), "678000"},
		{"%z", time.Date(2017, 1, 2, 3, 4, 5, 0, time.FixedZone("JST", 9*60*60)), "+0900"},
		{"%Z", time.Date(2017, 1, 2, 3, 4, 5, 0, time.FixedZone("JST", 9*60*60)), "JST"},
		{"%Z", time.Date(2017, 1, 2, 3, 4, 5, 0, time.FixedZone("", 9*60*60)), "+09:00"},
		{"%Z", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "UTC"},
		{"%j", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "002"},
		{"%-j", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "2"},
		{"%U", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "01"},
		{"%-U", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1"},
		{"%x", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "01/02/17"},
		{"%X", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "03:04:05"},
		{"%%", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "%"},
		{"%s", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1483326245"},
	}

	for _, test := range tests {
		actual := Strftime(test.time, test.format)
		assert.Equalf(t, test.expected, actual, "error parsing %q", test.format)
	}
}
