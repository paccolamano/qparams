package qparams

// OrderDirection represents the direction of a sort clause in a query.
// Supported values are "asc" (ascending) and "desc" (descending).
type OrderDirection string

// Symbol returns the SQL keyword corresponding to the order direction.
// Defaults to "asc" if the direction is not explicitly "desc".
func (d OrderDirection) Symbol() string {
	switch d {
	case OrderDesc:
		return string(OrderDesc)
	default:
		return string(OrderAsc)
	}
}

const (
	// OrderAsc sorts results in ascending order (default).
	OrderAsc OrderDirection = "asc"

	// OrderDesc sorts results in descending order.
	OrderDesc OrderDirection = "desc"
)

// OrderClause represents a single ORDER BY clause in a query.
// It specifies the field to sort on and the direction of sorting.
//
// Example:
//
//	{ "field": "created_at", "direction": "desc" }
type OrderClause struct {
	// Field is the column or attribute to sort by.
	Field string `json:"field"`

	// Direction is the order direction ("asc" or "desc").
	Direction OrderDirection `json:"direction"`
}
