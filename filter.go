package qparams

// Filter represents a single filtering condition in a query.
// It targets a specific field, applies a relational operator,
// and compares against the given value.
//
// Example:
//
//	{ "field": "name", "op": "eq", "value": "Alice" }
type Filter struct {
	// Field is the name of the column or attribute being filtered.
	Field string `json:"field"`

	// Op is the relational operator to apply (e.g., eq, lt, in).
	Op RelationalOperator `json:"op"`

	// Value is the comparison value used with the operator.
	Value string `json:"value"`
}

// FilterGroup represents a collection of filters combined together
// with a logical operator (AND/OR). FilterGroups can be nested,
// enabling the construction of complex, tree-like query conditions.
//
// Example:
//
//	{ "op": "and", "filters": [...], "groups": [...] }
type FilterGroup struct {
	// Op determines how Filters and Groups inside this group are combined.
	// Supported values: "and", "or".
	Op LogicalOperator `json:"op"`

	// Filters is the list of individual filtering conditions in this group.
	Filters []Filter `json:"filters,omitempty"`

	// Groups allows nesting of additional filter groups for more complex queries.
	Groups []FilterGroup `json:"groups,omitempty"`
}
