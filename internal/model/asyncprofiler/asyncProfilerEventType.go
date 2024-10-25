package asyncprofiler

import (
	"fmt"
	api "skywalking.apache.org/repo/goapi/query"
	"strings"
)

type AsyncProfilerEventTypeEnumValue struct {
	Enum     []api.AsyncProfilerEventType
	Default  []api.AsyncProfilerEventType
	Selected []api.AsyncProfilerEventType
}

func (e *AsyncProfilerEventTypeEnumValue) Set(value string) error {
	values := strings.Split(value, ",")
	types := make([]api.AsyncProfilerEventType, 0)
	for _, v := range values {
		for _, enum := range e.Enum {
			if strings.EqualFold(enum.String(), v) {
				types = append(types, enum)
				break
			}
		}
	}

	if len(types) != 0 {
		e.Selected = types
		return nil
	}

	orders := make([]string, len(api.AllAsyncProfilerEventType))
	for i, order := range api.AllAsyncProfilerEventType {
		orders[i] = order.String()
	}
	return fmt.Errorf("allowed analysis aggregate type are %s", strings.Join(orders, ", "))
}

func (e *AsyncProfilerEventTypeEnumValue) String() string {
	selected := make([]string, len(e.Selected))
	for i, item := range e.Selected {
		selected[i] = item.String()
	}
	return strings.Join(selected, ",")
}
