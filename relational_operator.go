package qparams

// RelationalOperator defines the set of supported comparison operators
// that can be used in query parameters to filter results.
type RelationalOperator string

// Symbol returns the SQL equivalent symbol for a given RelationalOperator.
// For example, "eq" maps to "=", "lt" maps to "<", "ilike" maps to "ilike", etc.
// If the operator is unknown, it defaults to "=".
func (o RelationalOperator) Symbol() string {
	switch o {
	case NotEqualsOperator:
		return "<>"
	case GreaterThanOperator:
		return ">"
	case GreaterThanEqualsOperator:
		return ">="
	case LowerThanOperator:
		return "<"
	case LowerThanEqualsOperator:
		return "<="
	case LikeOperator:
		return "like"
	case ILikeOperator:
		return "ilike"
	case InOperator:
		return "in"
	default:
		return "="
	}
}

const (
	// EqualsOperator represents equality comparison (=).
	EqualsOperator RelationalOperator = "eq"

	// NotEqualsOperator represents inequality comparison (<>).
	NotEqualsOperator RelationalOperator = "ne"

	// GreaterThanOperator represents greater-than comparison (>).
	GreaterThanOperator RelationalOperator = "gt"

	// GreaterThanEqualsOperator represents greater-than-or-equal comparison (>=).
	GreaterThanEqualsOperator RelationalOperator = "gte"

	// LowerThanOperator represents less-than comparison (<).
	LowerThanOperator RelationalOperator = "lt"

	// LowerThanEqualsOperator represents less-than-or-equal comparison (<=).
	LowerThanEqualsOperator RelationalOperator = "lte"

	// LikeOperator represents a case-sensitive pattern match (LIKE).
	LikeOperator RelationalOperator = "like"

	// ILikeOperator represents a case-insensitive pattern match (ILIKE).
	ILikeOperator RelationalOperator = "ilike"

	// InOperator represents an inclusion check (IN).
	InOperator RelationalOperator = "in"
)

var relationalOperators = map[RelationalOperator]struct{}{
	EqualsOperator:            {},
	NotEqualsOperator:         {},
	GreaterThanOperator:       {},
	GreaterThanEqualsOperator: {},
	LowerThanOperator:         {},
	LowerThanEqualsOperator:   {},
	LikeOperator:              {},
	ILikeOperator:             {},
	InOperator:                {},
}
