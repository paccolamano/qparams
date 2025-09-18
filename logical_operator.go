package qparams

// LogicalOperator defines how multiple filters or filter groups
// are combined in a search query (e.g., with "AND" or "OR").
type LogicalOperator string

// Symbol returns the SQL equivalent keyword for a given LogicalOperator.
// Defaults to "and" if the operator is not recognized.
func (o LogicalOperator) Symbol() string {
	switch o {
	case OrOperator:
		return string(OrOperator)
	default:
		return string(AndOperator)
	}
}

const (
	// AndOperator represents a logical AND between filters or groups.
	AndOperator LogicalOperator = "and"

	// OrOperator represents a logical OR between filters or groups.
	OrOperator LogicalOperator = "or"
)

var logicalOperators = map[LogicalOperator]struct{}{
	AndOperator: {},
	OrOperator:  {},
}
