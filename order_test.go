package qparams

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestOrderDirectionSymbol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		direction OrderDirection
		expected  string
	}{
		{
			name:      `Symbol() should return "asc"`,
			direction: OrderAsc,
			expected:  "asc",
		},
		{
			name:      `Symbol() should return "desc"`,
			direction: OrderDesc,
			expected:  "desc",
		},
		{
			name:      `Given wrong direction, Symbol() should return "asc"`,
			direction: OrderDirection("foo"),
			expected:  "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.direction.Symbol(), tt.expected)
		})
	}
}
