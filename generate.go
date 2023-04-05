package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
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

func generate(w io.Writer, o GenerateOptions) error {
	var now time.Time

	switch o.base {
	case "now", "today", "":
		now = time.Now()
	case "yesterday":
		now = time.Now().AddDate(0, 0, -1)
	case "tomorrow":
		now = time.Now().AddDate(0, 0, 1)
	default:
		return fmt.Errorf("unknown base: %s", o.base)
	}

	if now.IsZero() {
		now = time.Now()
	}

	switch o.truncate {
	case "day":
		now = now.Truncate(time.Hour * 24)
	case "hour":
		now = now.Truncate(time.Hour)
	case "minute":
		now = now.Truncate(time.Minute)
	case "second":
		now = now.Truncate(time.Second)
	case "":
		// do nothing
	default:
		return fmt.Errorf("unknown truncate: %s", o.truncate)
	}

	for _, delta := range o.delta {
		var err error
		now, err = applyDelta(now, delta)
		if err != nil {
			return err
		}
	}

	now, err := transform(now, o.options)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "%d\n", now.Unix()); err != nil {
		return err
	}

	return nil
}
