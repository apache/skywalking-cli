package table

import (
	"testing"

	"github.com/apache/skywalking-cli/graphql/schema"
)

func TestTableDisplay(t *testing.T) {
	var result []schema.Service
	display(t, result)
	result = make([]schema.Service, 0)
	display(t, result)
	result = append(result, schema.Service{
		ID:   "1",
		Name: "table",
	})
	display(t, result)
}

func display(t *testing.T, result []schema.Service) {
	if err := Display(result); err != nil {
		t.Error(err)
	}
}
