package qparams

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLogicalOperatorSymbol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		operator LogicalOperator
		expected string
	}{
		{
			name:     `Symbol() should return "and"`,
			operator: AndOperator,
			expected: "and",
		},
		{
			name:     `Symbol() should return "or"`,
			operator: OrOperator,
			expected: "or",
		},
		{
			name:     `Given wrong operator, Symbol() should return "and"`,
			operator: LogicalOperator("foo"),
			expected: "and",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.operator.Symbol(), tt.expected)
		})
	}
}
