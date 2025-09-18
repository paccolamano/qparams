package qparams

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestRelationalOperatorSymbol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		operator RelationalOperator
		expected string
	}{
		{
			name:     `Symbol() should return "="`,
			operator: EqualsOperator,
			expected: "=",
		},
		{
			name:     `Symbol() should return "<>"`,
			operator: NotEqualsOperator,
			expected: "<>",
		},
		{
			name:     `Symbol() should return ">"`,
			operator: GreaterThanOperator,
			expected: ">",
		},
		{
			name:     `Symbol() should return ">="`,
			operator: GreaterThanEqualsOperator,
			expected: ">=",
		},
		{
			name:     `Symbol() should return "<"`,
			operator: LowerThanOperator,
			expected: "<",
		},
		{
			name:     `Symbol() should return "<="`,
			operator: LowerThanEqualsOperator,
			expected: "<=",
		},
		{
			name:     `Symbol() should return "like"`,
			operator: LikeOperator,
			expected: "like",
		},
		{
			name:     `Symbol() should return "ilike"`,
			operator: ILikeOperator,
			expected: "ilike",
		},
		{
			name:     `Symbol() should return "in"`,
			operator: InOperator,
			expected: "in",
		},
		{
			name:     `Given wrong operator, Symbol() should return "="`,
			operator: RelationalOperator("foo"),
			expected: "=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.operator.Symbol(), tt.expected)
		})
	}
}
