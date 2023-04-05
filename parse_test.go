package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	require.Error(t, parse(nil, []string{}, Options{}))
	require.Error(t, parse(os.Stdout, []string{}, Options{}))

	require.NoError(t, os.Setenv("TZ", "America/Toronto")) // UTC-4
	tests := []struct {
		entry   string
		want    string
		options Options
	}{
		{"1680704033", "2023-04-05 10:13:53 -0400 EDT", Options{}},
		{"1680704033", "2023-04-05 14:13:53 +0000 UTC", Options{utc: true}},
		{"1588059756238", "2020-04-28 03:42:36.238 -0400 EDT", Options{precision: "millisecond"}},
		{"1588059756238", "2020-04-28 07:42:36.238 +0000 UTC", Options{precision: "millisecond", utc: true}},
		{"1588059756238", "2020-04-28 16:42:36.238 +0900 JST", Options{precision: "millisecond", offset: "Asia/Tokyo"}},
		{"1588059756238", "2020-04-28 16:42:36.238 +0900 JST", Options{precision: "millisecond", offset: "JST"}},
		{"1588059756238", "2020-04-28 16:42:36.238 +0900 (+900)", Options{precision: "millisecond", offset: "9:00"}},
		{"1588059756238", "2020-04-28 16:42:36.238 +0900 (+900)", Options{precision: "millisecond", offset: "900"}},
		{"1588059756238", "2020-04-28 16:42:36.238 +0900 (+900)", Options{precision: "millisecond", offset: "09:00"}},

		{"1588059756238", "2020-04-28 15:42:36.238 +0800 (+800)", Options{precision: "millisecond", offset: "8:00"}},
		{"1588059756238", "2020-04-28 15:42:36.238 +0800 (+800)", Options{precision: "millisecond", offset: "08:00"}},

		{"1588059756238", "2020-04-28 04:42:36.238 -0300 (-300)", Options{precision: "millisecond", offset: "-3:00"}},
		{"1588059756238", "2020-04-28 04:42:36.238 -0300 (-300)", Options{precision: "millisecond", offset: "-300"}},
		{"1588059756238", "2020-04-28 04:42:36.238 -0300 (-300)", Options{precision: "millisecond", offset: "-0300"}},
		{"1588059756238", "2020-04-28 04:42:36.238 -0300 (-300)", Options{precision: "millisecond", offset: "-03:00"}},

		{"1680704033", "04/05/2023", Options{format: "%m/%d/%Y"}},
		{"1680704033", "05/04/2023", Options{format: "%d/%m/%Y"}},
		{"1680704033", "2023-04-05 10:13", Options{format: "%Y-%m-%d %H:%M"}},
		{"1588059756238", "2020-04-28 03:42:36.238000000", Options{precision: "millisecond", format: "%Y-%m-%d %H:%M:%S.%f"}},
	}

	for _, tt := range tests {
		var buf strings.Builder
		assert.NoError(t, parse(&buf, []string{tt.entry}, tt.options))

		result := buf.String()
		result = strings.Trim(result, "\n")
		assert.Equal(t, tt.want, result)
	}
}

func TestTime(t *testing.T) {
	myT := time.Date(2022, 4, 28, 14, 0, 0, 0, time.UTC)
	for _, d := range []struct {
		name     string
		expected string
	}{
		{"UTC", "2022-04-28 14:00:00 +0000 UTC"},
		{"America/Los_Angeles", "2022-04-28 07:00:00 -0700 PDT"},
		{"Asia/Tokyo", "2022-04-28 23:00:00 +0900 JST"},
		{"Asia/Taipei", "2022-04-28 22:00:00 +0800 CST"},
		{"Asia/Hong_Kong", "2022-04-28 22:00:00 +0800 HKT"},
		{"Asia/Shanghai", "2022-04-28 22:00:00 +0800 CST"},
	} {
		loc, _ := time.LoadLocation(d.name)
		if val := fmt.Sprintf("%s", myT.In(loc)); val != d.expected {
			fmt.Println(val)
			t.FailNow()
		}
	}
}
