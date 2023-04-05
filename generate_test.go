package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestApplyDelta(t *testing.T) {
	tests := []struct {
		base     time.Time
		delta    string
		expected time.Time
	}{
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1year", time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "2year", time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "-2years", time.Date(2015, 1, 2, 3, 4, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "-17y", time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "3day", time.Date(2017, 1, 5, 3, 4, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "12hour", time.Date(2017, 1, 2, 15, 4, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1hour", time.Date(2017, 1, 2, 4, 4, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1min", time.Date(2017, 1, 2, 3, 5, 5, 0, time.UTC)},
		{time.Date(2017, 1, 2, 3, 4, 5, 0, time.UTC), "1s", time.Date(2017, 1, 2, 3, 4, 6, 0, time.UTC)},
	}

	for _, test := range tests {
		actual, err := applyDelta(test.base, test.delta)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, actual)
	}
}

func TestGenerate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		expected time.Time
		options  GenerateOptions
	}{
		{now.Truncate(time.Hour * 24), GenerateOptions{base: "now", truncate: "day"}},
		{now.Truncate(time.Hour*24).AddDate(0, 0, 1), GenerateOptions{base: "tomorrow", truncate: "day"}},
	}

	for _, test := range tests {
		var buf strings.Builder
		err := generate(&buf, test.options)

		actual := strings.Trim(buf.String(), "\n")

		assert.NoError(t, err)
		assert.Equal(t, strconv.Itoa(int(test.expected.Unix())), actual)
	}
}
