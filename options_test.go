package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGenerateOptionsParse(t *testing.T) {
	optionsArgsStr := "gen -b today -d 3day -d 12hour"
	optionsArgs := strings.Split(optionsArgsStr, " ")

	var options GenerateOptions
	remainingArgs, err := options.Parse(optionsArgs...)
	require.NoError(t, err)

	assert.Equal(t, "today", options.base)
	assert.Equal(t, []string{"3day", "12hour"}, options.delta)
	assert.Len(t, remainingArgs, 0)
}

func TestGenerateOptionsParseHelp(t *testing.T) {
	optionsArgsStr := "gen help"
	optionsArgs := strings.Split(optionsArgsStr, " ")

	var options GenerateOptions
	remainingArgs, err := options.Parse(optionsArgs...)
	require.NoError(t, err)
	assert.Equal(t, []string{"help"}, remainingArgs)
}

func TestTruncateOption(t *testing.T) {
	tests := []struct {
		input     string
		expected  TruncateOption
		shouldErr bool
	}{
		{"", TruncateOptionNone, false},
		{"second", TruncateOptionSecond, false},
		{"minute", TruncateOptionMinute, false},
		{"hour", TruncateOptionHour, false},
		{"day", TruncateOptionDay, false},
		{"foo", TruncateOptionNone, true},
	}

	for _, test := range tests {
		var opt TruncateOption
		err := opt.Set(test.input, nil)
		if test.shouldErr {
			assert.Errorf(t, err, "expected error for input %q", test.input)
		} else {
			assert.NoErrorf(t, err, "unexpected no error for input %q", test.input)
			assert.Equalf(t, test.expected, opt, "unexpected value for input %q", test.input)
		}
	}
}
