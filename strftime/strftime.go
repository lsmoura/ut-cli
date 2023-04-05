package strftime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func StrTimeTokens(format string) []string {
	var pieces []string
	var block string

	for i := 0; i < len(format); i++ {
		c := format[i]

		if c == '%' {
			if len(block) > 0 {
				pieces = append(pieces, block)
				block = ""
			}
			if i == len(format)-1 {
				// malformed format string
				break
			}
			pieces = append(pieces, format[i:i+2])
			i++
		} else {
			block += string(c)
		}
	}

	if len(block) > 0 {
		pieces = append(pieces, block)
	}

	return pieces
}

// Strftime formats a time.Time according to the c's Strftime format
func Strftime(t time.Time, format string) string {
	pieces := StrTimeTokens(format)

	var output []string
	for _, piece := range pieces {
		switch piece {
		case "%Y":
			output = append(output, strconv.Itoa(t.Year()))
		case "%m":
			output = append(output, fmt.Sprintf("%02d", t.Month()))
		case "%d":
			output = append(output, fmt.Sprintf("%02d", t.Day()))
		case "%H":
			output = append(output, fmt.Sprintf("%02d", t.Hour()))
		case "%M":
			output = append(output, fmt.Sprintf("%02d", t.Minute()))
		case "%S":
			output = append(output, fmt.Sprintf("%02d", t.Second()))
		case "%s":
			output = append(output, strconv.Itoa(int(t.Unix())))
		case "%f":
			output = append(output, fmt.Sprintf("%06d", t.Nanosecond()))
		case "%z":
			_, offset := t.Zone()
			output = append(output, fmt.Sprintf("%+03d%02d", offset/3600, (offset%3600)/60))
		case "%Z":
			_, offset := t.Zone()
			output = append(output, fmt.Sprintf("%+03d:%02d", offset/3600, (offset%3600)/60))
		case "%I":
			hour := t.Hour()
			if hour == 0 {
				hour = 12
			}
			if hour > 12 {
				hour -= 12
			}
			output = append(output, fmt.Sprintf("%02d", hour))
		case "%p":
			if t.Hour() < 12 {
				output = append(output, "AM")
			} else {
				output = append(output, "PM")
			}
		case "%%":
			output = append(output, "%")
		default:
			output = append(output, piece)
		}
	}

	return strings.Join(output, "")
}
