package asyncprofiler

import (
	"fmt"
	api "skywalking.apache.org/repo/goapi/query"
	"strings"
)

type JFREventTypeEnumValue struct {
	Enum     []api.JFREventType
	Default  api.JFREventType
	Selected api.JFREventType
}

func (e *JFREventTypeEnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if strings.EqualFold(enum.String(), value) {
			e.Selected = enum
			return nil
		}
	}
	orders := make([]string, len(api.AllJFREventType))
	for i, order := range api.AllJFREventType {
		orders[i] = order.String()
	}
	return fmt.Errorf("allowed analysis aggregate type are %s", strings.Join(orders, ", "))
}

func (e *JFREventTypeEnumValue) String() string {
	return e.Selected.String()
}
