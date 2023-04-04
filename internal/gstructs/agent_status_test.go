package gstructs_test

import (
	"testing"

	"github.com/dragonheim/gagent/internal/gstructs"
	"github.com/stretchr/testify/assert"
)

func TestGetByIndex(t *testing.T) {
	agentStatuses := gstructs.AgentStatuses

	tests := []struct {
		index        int
		expected     string
		shouldReturn bool
	}{
		{0, "ERROR", true},
		{1, "INIT", true},
		{8, "ERROR", true},
		{9, "", false},
		{-1, "", false},
	}

	for _, test := range tests {
		res, err := agentStatuses.GetByIndex(test.index)
		if test.shouldReturn {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, res)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestGetByName(t *testing.T) {
	agentStatuses := gstructs.AgentStatuses

	tests := []struct {
		name         string
		expected     byte
		shouldReturn bool
	}{
		{"ERROR", 0, true},
		{"INIT", 1, true},
		{"COMPLETED", 6, true},
		{"INVALID", 0, false},
	}

	for _, test := range tests {
		res, err := agentStatuses.GetByName(test.name)
		if test.shouldReturn {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, res)
		} else {
			assert.Error(t, err)
		}
	}
}
