package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
	"time"
)

type TruncateOption string

const (
	TruncateOptionNone   TruncateOption = ""
	TruncateOptionDay    TruncateOption = "day"
	TruncateOptionHour   TruncateOption = "hour"
	TruncateOptionMinute TruncateOption = "minute"
	TruncateOptionSecond TruncateOption = "second"
)

func (opt TruncateOption) Truncate(t time.Time) (time.Time, error) {
	switch opt {
	case TruncateOptionDay:
		return t.Truncate(time.Hour * 24), nil
	case TruncateOptionHour:
		return t.Truncate(time.Hour), nil
	case TruncateOptionMinute:
		return t.Truncate(time.Minute), nil
	case TruncateOptionSecond:
		return t.Truncate(time.Second), nil
	case TruncateOptionNone:
		return t, nil
	}

	return t, fmt.Errorf("unknown truncate option: %s", opt)
}

func (opt *TruncateOption) Set(value string, _ getopt.Option) error {
	switch v := TruncateOption(value); v {
	case TruncateOptionNone, TruncateOptionDay, TruncateOptionHour, TruncateOptionMinute, TruncateOptionSecond:
		*opt = v
	default:
		return fmt.Errorf("unknown truncate option: %s", value)
	}

	return nil
}

func (opt *TruncateOption) String() string {
	return string(*opt)
}

type Options struct {
	utc       bool
	utcOption getopt.Option
	help      *bool
	version   *bool

	format       string
	formatOption getopt.Option

	offset          string
	offsetOption    getopt.Option
	precision       string
	precisionOption getopt.Option

	flags *getopt.Set
}

const (
	offsetEnvVar    = "UT_OFFSET"
	precisionEnvVar = "UT_PRECISION"
	formatEnvVar    = "UT_DATETIME_FORMAT"
)

func (o *Options) Flags() *getopt.Set {
	if o.flags != nil {
		return o.flags
	}

	o.flags = getopt.New()

	o.utcOption = o.flags.FlagLong(&o.utc, "utc", 'u', "Use utc timezone")
	o.help = o.flags.BoolLong("help", 'h', "Prints help information")
	o.version = o.flags.BoolLong("version", 'V', "Prints version information")

	o.formatOption = o.flags.FlagLong(&o.format, "format", 'f', "", "Format output using given format (used for generate command)")
	o.offsetOption = o.flags.FlagLong(&o.offset, "offset", 'o', "", "Use given value as timezone offset")
	o.precisionOption = o.flags.FlagLong(&o.precision, "precision", 'p', "", "Use given value as precision")

	return o.flags
}

// Parse parses the command line arguments and returns the remaining arguments.
func (o *Options) Parse(args ...string) ([]string, error) {
	if err := o.Flags().Getopt(args, nil); err != nil {
		return nil, err
	}

	return o.Flags().Args(), nil
}

// UTC returns the value of the UTC flag and whether it was set.
func (o *Options) UTC() (bool, bool) {
	var seen bool
	if o.utcOption != nil {
		seen = o.utcOption.Seen()
	}
	return o.utc, seen
}

func (o *Options) Offset() (string, bool) {
	var seen bool
	if o.offsetOption != nil {
		seen = o.offsetOption.Seen()
	}

	if !seen {
		if os.Getenv(offsetEnvVar) != "" {
			return os.Getenv(offsetEnvVar), true
		}
	}

	return o.offset, seen
}

func (o *Options) Precision() (string, bool) {
	if !o.precisionOption.Seen() {
		if os.Getenv(precisionEnvVar) != "" {
			return os.Getenv(precisionEnvVar), true
		}
	}

	return o.precision, o.precisionOption.Seen()
}

func (o *Options) Format() (string, bool) {
	var seen bool

	if o.formatOption != nil {
		seen = o.formatOption.Seen()
	}

	if !seen {
		if os.Getenv(formatEnvVar) != "" {
			return os.Getenv(formatEnvVar), true
		}
	}

	return o.format, seen
}

type GenerateOptions struct {
	options Options

	base           string
	baseOption     getopt.Option
	delta          []string
	deltaOption    getopt.Option
	truncate       TruncateOption
	truncateOption getopt.Option

	flags *getopt.Set
}

func (o *GenerateOptions) Flags() *getopt.Set {
	if o.flags != nil {
		return o.flags
	}

	o.flags = getopt.New()

	o.baseOption = o.flags.FlagLong(&o.base, "base", 'b', "", "Use given value as base timestamp")
	o.deltaOption = o.flags.FlagLong(&o.delta, "delta", 'd', "", "Use given value as delta")
	o.truncateOption = o.flags.FlagLong(&o.truncate, "truncate", 't', "", "Truncate the timestamp to the given precision")

	return o.flags
}

func (o *GenerateOptions) Parse(args ...string) ([]string, error) {
	if err := o.Flags().Getopt(args, nil); err != nil {
		return nil, err
	}

	return o.Flags().Args(), nil
}
