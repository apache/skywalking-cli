package model

import (
	"fmt"
	"strings"

	event "skywalking/network/event/v3"
)

// EventTypeEnumValue defines the values domain of --type option.
type EventTypeEnumValue struct {
	Enum     []event.Type
	Default  event.Type
	Selected event.Type
}

// Set the --type value, from raw string to EventTypeEnumValue.
func (s *EventTypeEnumValue) Set(value string) error {
	for _, enum := range s.Enum {
		if strings.EqualFold(enum.String(), value) {
			s.Selected = enum
			return nil
		}
	}
	types := make([]string, len(event.Type_name))
	for index := range event.Type_name {
		types[index] = event.Type_name[index]
	}
	return fmt.Errorf("allowed types are %s", strings.Join(types, ", "))
}

// String representation of the event type.
func (s EventTypeEnumValue) String() string {
	return s.Selected.String()
}
