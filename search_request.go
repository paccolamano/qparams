package qparams

// SearchRequest represents a structured query definition parsed from request parameters.
// It combines filtering (via FilterGroups), ordering, and pagination options.
//
// This object is typically extracted from query params (e.g. `?q=<json>`)
// and then used to build database queries or other filtering logic.
//
// Example:
//
//	{
//	  "groups": {
//	    "op": "and",
//	    "filters": [
//	      { "field": "status", "op": "eq", "value": "active" }
//	    ],
//	    "groups": [
//	      {
//	        "op": "or",
//	        "filters": [
//	          { "field": "role", "op": "eq", "value": "admin" },
//	          { "field": "role", "op": "eq", "value": "editor" }
//	        ]
//	      }
//	    ]
//	  },
//	  "order_by": [
//	    { "field": "created_at", "direction": "desc" }
//	  ],
//	  "limit": 20,
//	  "offset": 0
//	}
type SearchRequest struct {
	// Groups represents the root filter group, which can contain
	// multiple filters and nested groups combined with logical operators.
	Groups *FilterGroup `json:"groups,omitempty"`

	// OrderBy defines the sorting rules to apply to the result set.
	OrderBy []OrderClause `json:"order_by,omitempty"`

	// Limit restricts the maximum number of items returned.
	// If nil, no explicit limit is applied.
	Limit *int `json:"limit,omitempty"`

	// Offset specifies how many items to skip before starting to return results.
	// Useful for pagination in combination with Limit.
	Offset *int `json:"offset,omitempty"`
}
