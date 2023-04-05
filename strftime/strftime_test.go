package strftime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStrftime(t *testing.T) {
	tests := []struct {
		format   string
		time     time.Time
		expected string
	}{
		{"%Y-%m-%d %H:%M:%S", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "2017-01-02 03:04:05"},
		{"%Y-%m-%d %H:%M:%S", time.Date(2017, 1, 2, 3, 4, 5, 0, time.FixedZone("JST", 9*60*60)), "2017-01-02 03:04:05"},
		{"Now it's %I:%M%p.", time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "Now it's 03:04AM."},
	}

	for _, test := range tests {
		actual := Strftime(test.time, test.format)
		assert.Equal(t, test.expected, actual)
	}
}
