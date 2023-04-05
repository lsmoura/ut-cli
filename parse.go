package main

import (
	"bufio"
	"fmt"
	"github.com/lsmoura/ut-cli/strftime"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "time/tzdata"
)

var offsetToTZ = map[string]string{
	"JST": "Asia/Tokyo",
}

var smallOffsetMatch = regexp.MustCompile(`^(\d{1}):(\d{2})$`)
var timeOffsetMatch = regexp.MustCompile(`^([+-]?)(\d{1,2}):(\d{2})$`)
var timeHundredOffsetMatch = regexp.MustCompile(`^([+-]?)(\d{3,4})$`)

func sanitizeOffset(offset string) string {
	if offset == "" {
		return offset
	}

	var negative bool
	if timeOffsetMatch.MatchString(offset) {
		if offset[0] == '+' || offset[0] == '-' {
			negative = offset[0] == '-'
			offset = offset[1:]
		}
		if smallOffsetMatch.MatchString(offset) {
			offset = "0" + offset
		}
		if negative {
			offset = "-" + offset
		}
	}

	return offset
}

func transform(t time.Time, o Options) (time.Time, error) {
	if utc, _ := o.UTC(); utc {
		t = t.UTC()
	}
	if offset, _ := o.Offset(); offset != "" {
		offset = sanitizeOffset(offset)

		var offsetLocation *time.Location

		if matches := timeOffsetMatch.FindStringSubmatch(offset); len(matches) > 0 {
			hours, err := strconv.Atoi(matches[2])
			if err != nil {
				return t, err
			}
			minutes, err := strconv.Atoi(matches[3])
			if err != nil {
				return t, err
			}

			if matches[1] == "-" {
				hours = -hours
				minutes = -minutes
			}

			seconds := hours*3600 + minutes*60

			offsetName := fmt.Sprintf("%d", hours*100+minutes)
			if hours > 0 {
				offsetName = "+" + offsetName
			}
			offsetName = "(" + offsetName + ")"

			offsetLocation = time.FixedZone(offsetName, seconds)
		}

		if matches := timeHundredOffsetMatch.FindStringSubmatch(offset); len(matches) > 0 {
			value, err := strconv.Atoi(matches[2])
			if err != nil {
				return t, err
			}

			minutes := value % 100
			hours := value / 100

			if matches[1] == "-" {
				hours = -hours
				minutes = -minutes
			}

			seconds := hours*3600 + minutes*60

			offsetName := fmt.Sprintf("%d", hours*100+minutes)
			if hours > 0 {
				offsetName = "+" + offsetName
			}
			offsetName = "(" + offsetName + ")"

			offsetLocation = time.FixedZone(offsetName, seconds)
		}

		if offsetLocation == nil {
			if mapped, ok := offsetToTZ[offset]; ok {
				offset = mapped
			}

			loc, err := time.LoadLocation(offset)
			if err != nil {
				return t, err
			}

			offsetLocation = loc
		}

		t = t.In(offsetLocation)
	}

	return t, nil
}

func format(t time.Time, format string) string {
	if format == "" {
		return fmt.Sprintf("%s", t)
	}

	if strings.Contains(format, "%") {
		return strftime.Strftime(t, format)
	}

	return t.Format(format)
}

func parse(w io.Writer, args []string, options Options) error {
	if w == nil {
		return fmt.Errorf("no writer")
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	var data io.Reader
	if len(args) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no value to parse")
		}
		data = bufio.NewReader(os.Stdin)
	} else if args[0] == "-" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no input")
		}
		data = bufio.NewReader(os.Stdin)
	} else {
		data = strings.NewReader(args[0])
	}

	argBytes, err := io.ReadAll(data)
	if err != nil {
		return err
	}

	arg := strings.Trim(string(argBytes), "\n \t")

	if len(arg) == 0 {
		return fmt.Errorf("no input")
	}

	timestamp, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return err
	}

	var now time.Time
	switch options.precision {
	case "millisecond", "milli", "ms":
		now = time.UnixMilli(timestamp)
	case "microsecond", "micro", "us", "μs", "µs": // includes U+03BC (Greek letter mu) and U+00B5 (micro symbol)
		now = time.UnixMicro(timestamp)
	case "", "second", "s":
		now = time.Unix(timestamp, 0)
	default:
		return fmt.Errorf("unknown precision: %s", options.precision)
	}

	now, err = transform(now, options)
	if err != nil {
		return err
	}

	strFormat, _ := options.Format()
	if _, err := fmt.Fprintf(w, "%s\n", format(now, strFormat)); err != nil {
		return err
	}

	return nil
}
