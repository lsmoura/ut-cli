package main

import (
	"fmt"
	"github.com/lsmoura/ut-cli/strftime"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var deltaMatch = regexp.MustCompile(`^([+-]?)(\d+)(days?|d|years?|y|hours?|h|min|minutes?|seconds?|s)$`)

func applyDelta(t time.Time, delta string) (time.Time, error) {
	matches := deltaMatch.FindStringSubmatch(delta)
	if len(matches) == 0 {
		return t, fmt.Errorf("invalid delta: %s", delta)
	}

	value, err := strconv.Atoi(matches[2])
	if err != nil {
		return t, err
	}
	if matches[1] == "-" {
		value = -value
	}

	switch matches[3] {
	case "year", "years", "y":
		return t.AddDate(value, 0, 0), nil
	case "day", "days", "d":
		return t.AddDate(0, 0, value), nil
	case "hour", "hours", "h":
		return t.Add(time.Duration(value) * time.Hour), nil
	case "min", "minute", "minutes":
		return t.Add(time.Duration(value) * time.Minute), nil
	case "second", "seconds", "s":
		return t.Add(time.Duration(value) * time.Second), nil
	default:
		return t, fmt.Errorf("unknown delta: %s", delta)
	}
}

// timeFormatFromPercent converts a time format string with percent values to a go time format string
func timeFormatFromPercent(format string) string {
	pieces := strftime.StrTimeTokens(format)

	var goFormat string
	for _, piece := range pieces {
		switch piece {
		case "%Y":
			goFormat += "2006"
		case "%m":
			goFormat += "01"
		case "%d":
			goFormat += "02"
		case "%H":
			goFormat += "15"
		case "%M":
			goFormat += "04"
		case "%S":
			goFormat += "05"
		case "%z":
			goFormat += "Z07:00"
		case "%Z":
			goFormat += "MST"
		case "%%":
			goFormat += "%"
		default:
			goFormat += piece
		}
	}

	return goFormat
}

func generate(w io.Writer, o GenerateOptions) (err error) {
	var now time.Time

	switch o.base {
	case "now", "today", "":
		now = time.Now()
	case "yesterday":
		now = time.Now().AddDate(0, 0, -1)
	case "tomorrow":
		now = time.Now().AddDate(0, 0, 1)
	default:
		layout := o.options.format
		if layout == "" {
			layout = time.RFC3339
		}
		if strings.Contains(layout, "%") {
			layout = timeFormatFromPercent(layout)
		}
		t, err := time.Parse(layout, o.base)
		if err != nil {
			return err
		}
		now = t
	}

	if now.IsZero() {
		now = time.Now()
	}

	now, err = o.truncate.Truncate(now)
	if err != nil {
		return err
	}

	for _, delta := range o.delta {
		now, err = applyDelta(now, delta)
		if err != nil {
			return err
		}
	}

	now, err = transform(now, o.options)
	if err != nil {
		return err
	}

	var n int64
	switch o.options.precision {
	case "millisecond", "ms":
		n = now.UnixNano() / 1000000
	case "microsecond", "us":
		n = now.UnixNano() / 1000
	case "nanosecond", "ns":
		n = now.UnixNano()
	case "second", "s", "":
		n = now.Unix()
	default:
		return fmt.Errorf("unknown precision: %s", o.options.precision)
	}

	if _, err := fmt.Fprintf(w, "%d\n", n); err != nil {
		return err
	}

	return nil
}
