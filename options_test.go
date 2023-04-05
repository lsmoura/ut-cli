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
