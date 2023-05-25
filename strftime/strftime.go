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
		case "%a": // Weekday as locale's abbreviated name
			output = append(output, t.Weekday().String()[:3])
		case "%A": // Weekday as locale's full name
			output = append(output, t.Weekday().String())
		case "%w": // Weekday as a decimal number, where 0 is Sunday and 6 is Saturday
			output = append(output, strconv.Itoa(int(t.Weekday())))
		case "%d": // Day of the month as a zero-padded decimal number
			output = append(output, fmt.Sprintf("%02d", t.Day()))
		case "%-d": // Day of the month as a decimal number
			output = append(output, strconv.Itoa(t.Day()))
		case "%b": // Month as locale's abbreviated name
			output = append(output, t.Month().String()[:3])
		case "%B": // Month as locale's full name
			output = append(output, t.Month().String())
		case "%m": // Month as a zero-padded decimal number
			output = append(output, fmt.Sprintf("%02d", t.Month()))
		case "%-m": // Month as a decimal number
			output = append(output, strconv.Itoa(int(t.Month())))
		case "%y": // Year without century as a zero-padded decimal number
			output = append(output, fmt.Sprintf("%02d", t.Year()%100))
		case "%Y": // Year with century as a decimal number
			output = append(output, strconv.Itoa(t.Year()))
		case "%H": // Hour (24-hour clock) as a zero-padded decimal number
			output = append(output, fmt.Sprintf("%02d", t.Hour()))
		case "%-H": // Hour (24-hour clock) as a decimal number
			output = append(output, strconv.Itoa(t.Hour()))
		case "%I": // Hour (12-hour clock) as a zero-padded decimal number
			hour := t.Hour() % 12
			if hour == 0 {
				hour = 12
			}
			output = append(output, fmt.Sprintf("%02d", hour))
		case "%-I": // Hour (12-hour clock) as a decimal number
			hour := t.Hour() % 12
			if hour == 0 {
				hour = 12
			}
			output = append(output, strconv.Itoa(hour))
		case "%p": // Locale's equivalent of either AM or PM
			if t.Hour() < 12 {
				output = append(output, "AM")
			} else {
				output = append(output, "PM")
			}
		case "%M": // Minute as a zero-padded decimal number
			output = append(output, fmt.Sprintf("%02d", t.Minute()))
		case "%-M": // Minute as a decimal number
			output = append(output, strconv.Itoa(t.Minute()))
		case "%S": // Second as a zero-padded decimal number
			output = append(output, fmt.Sprintf("%02d", t.Second()))
		case "%-S": // Second as a decimal number
			output = append(output, strconv.Itoa(t.Second()))
		case "%f": // Microsecond as a decimal number, zero-padded on the left
			output = append(output, fmt.Sprintf("%06d", t.Nanosecond()/1000))
		case "%z": // UTC offset in the form +HHMM or -HHMM
			_, offset := t.Zone()
			output = append(output, fmt.Sprintf("%+03d%02d", offset/3600, (offset%3600)/60))
		case "%Z": // Time zone name
			_, offset := t.Zone()
			output = append(output, fmt.Sprintf("%+03d:%02d", offset/3600, (offset%3600)/60))
		case "%j": // Day of the year as a zero-padded decimal number
			output = append(output, fmt.Sprintf("%03d", t.YearDay()))
		case "%-j": // Day of the year as a decimal number
			output = append(output, strconv.Itoa(t.YearDay()))
		case "%U": // Week number of the year (Sunday as the first day of the week)
			_, week := t.ISOWeek()
			output = append(output, fmt.Sprintf("%02d", week))
		case "%-U": // Week number of the year (Sunday as the first day of the week)
			_, week := t.ISOWeek()
			output = append(output, strconv.Itoa(week))
		//case "%W": // Week number of the year (Monday as the first day of the week)
		//case "%-W": // Week number of the year (Monday as the first day of the week)
		case "%x": // Locale's appropriate date representation
			output = append(output, t.Format("01/02/06"))
		case "%X": // Locale's appropriate time representation
			output = append(output, t.Format("15:04:05"))
		case "%%":
			output = append(output, "%")

		case "%s":
			output = append(output, strconv.Itoa(int(t.Unix())))

		default:
			output = append(output, piece)
		}
	}

	return strings.Join(output, "")
}
