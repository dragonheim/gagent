package gstructs

import (
	"fmt"
)

type AgentStatus []string

var AgentStatuses = AgentStatus{
	"ERROR",
	"INIT",
	"SENDING",
	"RECEIVING",
	"ROUTING",
	"PROCESSING",
	"COMPLETED",
	"RETURNING",
	"ERROR",
}

func (a AgentStatus) GetByIndex(index int) (string, error) {
	if index < 0 || index >= len(a) {
		return "", fmt.Errorf("invalid index: %d", index)
	}
	return a[index], nil
}

func (a AgentStatus) GetByName(name string) (byte, error) {
	for i, status := range a {
		if status == name {
			return byte(i), nil
		}
	}
	return 0, fmt.Errorf("value not found: %s", name)
}
