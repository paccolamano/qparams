// Package qparams provides utilities to configure and handle
// search requests passed as query parameters in HTTP requests.
// It supports default global configuration, per-handler overrides,
// and validation of search expressions.
package qparams

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// contextKey is a custom type used to avoid collisions when
// storing values in request contexts.
type contextKey string

// searchKey is the context key under which parsed SearchRequest
// objects are stored.
const searchKey = contextKey("search")

// ErrorHandler defines the signature of a function responsible
// for handling request errors. It receives the HTTP response writer,
// the request, and the encountered error.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

var (
	// defaultQueryParam holds the default query parameter name used to
	// retrieve the search payload.
	defaultQueryParam string = "q"

	// defaultSearchMandatory indicates whether a search parameter is
	// required by default.
	defaultSearchMandatory bool = true

	// defaultLogicalOperators defines the default set of logical
	// operators allowed in filters.
	defaultLogicalOperators map[LogicalOperator]struct{} = logicalOperators

	// defaultRelationalOperators defines the default set of relational
	// operators allowed in filters.
	defaultRelationalOperators map[RelationalOperator]struct{} = relationalOperators

	// defaultLimit defines the default maximum limit applied to
	// search requests. Nil means "no limit".
	defaultLimit *int = nil

	// defaultErrorHandler is the fallback handler used when no custom
	// error handler is configured. It writes a error response with
	// HTTP 400 status code.
	defaultErrorHandler ErrorHandler = func(w http.ResponseWriter, r *http.Request, _ error) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		if err != nil {
			slog.Default().ErrorContext(r.Context(), "failed to send response", slog.String("err", err.Error()))
		}
	}

	// defaultFilterFields defines the default set of fields
	// allowed in filters.
	defaultFilterFields map[string]struct{} = map[string]struct{}{}

	// defaultOrderFields defines the default set of fields
	// allowed in order by.
	defaultOrderFields map[string]struct{} = map[string]struct{}{}
)

// SetDefaultQueryParam sets the default query parameter name
// used to extract search payloads.
func SetDefaultQueryParam(value string) {
	defaultQueryParam = value
}

// SetDefaultSearchMandatory sets the global default for requiring
// the search query parameter.
func SetDefaultSearchMandatory(value bool) {
	defaultSearchMandatory = value
}

// SetDefaultLogicalOperators replaces the default set of allowed
// logical operators with the provided ones.
func SetDefaultLogicalOperators(values ...LogicalOperator) {
	clear(defaultLogicalOperators)
	for _, v := range values {
		defaultLogicalOperators[v] = struct{}{}
	}
}

// SetDefaultRelationalOperators replaces the default set of allowed
// relational operators with the provided ones.
func SetDefaultRelationalOperators(values ...RelationalOperator) {
	clear(defaultRelationalOperators)
	for _, v := range values {
		defaultRelationalOperators[v] = struct{}{}
	}
}

// SetDefaultLimit sets the global default limit to apply to
// search requests. Negative means "no limit".
func SetDefaultLimit(value int) {
	if value < 0 {
		defaultLimit = nil
	} else {
		defaultLimit = ptr(value)
	}
}

// SetDefaultErrorHandler replaces the default global error handler
// with the provided function.
func SetDefaultErrorHandler(value ErrorHandler) {
	defaultErrorHandler = value
}

// SetDefaultFilterFields replaces the default set of allowed
// filter fields with the provided ones.
func SetDefaultFilterFields(values ...string) {
	clear(defaultFilterFields)
	for _, v := range values {
		defaultFilterFields[v] = struct{}{}
	}
}

// SetDefaultOrderFields replaces the default set of allowed
// order fields with the provided ones.
func SetDefaultOrderFields(values ...string) {
	clear(defaultOrderFields)
	for _, v := range values {
		defaultOrderFields[v] = struct{}{}
	}
}

// Options stores the configuration for a search handler,
// including query parameter names, validation rules,
// allowed operators, limits, and error handling.
type Options struct {
	queryParam                 string
	isSearchMandatory          bool
	allowedLogicalOperators    map[LogicalOperator]struct{}
	allowedRelationalOperators map[RelationalOperator]struct{}
	limit                      *int
	errorHandler               ErrorHandler
	allowedFilterFields        map[string]struct{}
	allowedOrderFields         map[string]struct{}
}

// Option is a functional option type used to configure Options
// when creating a new search handler.
type Option func(*Options)

// WithQueryParam sets a custom query parameter name for extracting
// search payloads.
func WithQueryParam(value string) Option {
	return func(o *Options) {
		o.queryParam = value
	}
}

// WithSearchMandatory configures whether the query parameter
// containing the search payload is required.
func WithSearchMandatory(value bool) Option {
	return func(o *Options) {
		o.isSearchMandatory = value
	}
}

// WithLogicalOperators restricts the set of logical operators
// allowed in filters.
func WithLogicalOperators(values ...LogicalOperator) Option {
	return func(o *Options) {
		clear(o.allowedLogicalOperators)
		for _, v := range values {
			o.allowedLogicalOperators[v] = struct{}{}
		}
	}
}

// WithRelationalOperators restricts the set of relational operators
// allowed in filters.
func WithRelationalOperators(values ...RelationalOperator) Option {
	return func(o *Options) {
		clear(o.allowedRelationalOperators)
		for _, v := range values {
			o.allowedRelationalOperators[v] = struct{}{}
		}
	}
}

// WithLimit sets a maximum number of results for search requests.
// Negative values mean "no limit".
func WithLimit(value int) Option {
	return func(o *Options) {
		if value < 0 {
			o.limit = nil
		} else {
			o.limit = ptr(value)
		}
	}
}

// WithErrorHandler overrides the error handler used by the search handler.
func WithErrorHandler(h ErrorHandler) Option {
	return func(o *Options) {
		o.errorHandler = h
	}
}

// WithFilterFields restricts the fields that can be used
// in filter conditions. It replace the fields set by
// SetDefaultFilterFields
func WithFilterFields(fields ...string) Option {
	return func(o *Options) {
		clear(o.allowedFilterFields)
		for _, v := range fields {
			o.allowedFilterFields[v] = struct{}{}
		}
	}
}

// WithExtraFilterFields restricts the fields that can be used
// in filter conditions.
func WithExtraFilterFields(fields ...string) Option {
	return func(o *Options) {
		for _, v := range fields {
			o.allowedFilterFields[v] = struct{}{}
		}
	}
}

// WithOrderFields restricts the fields that can be used
// in order by clauses. It replace the fields set by
// SetDefaultOrderFields
func WithOrderFields(fields ...string) Option {
	return func(o *Options) {
		clear(o.allowedOrderFields)
		for _, v := range fields {
			o.allowedOrderFields[v] = struct{}{}
		}
	}
}

// WithOrderFields restricts the fields that can be used
// in order by clauses.
func WithExtraOrderFields(fields ...string) Option {
	return func(o *Options) {
		for _, v := range fields {
			o.allowedOrderFields[v] = struct{}{}
		}
	}
}

// NewSearchHandler creates a middleware that parses, validates,
// and injects a SearchRequest into the request context.
// It can be customized via Option functions, falling back to
// global defaults when not provided.
func NewSearchHandler(opts ...Option) func(http.Handler) http.Handler {
	options := &Options{
		queryParam:                 defaultQueryParam,
		isSearchMandatory:          defaultSearchMandatory,
		allowedLogicalOperators:    defaultLogicalOperators,
		allowedRelationalOperators: defaultRelationalOperators,
		limit:                      defaultLimit,
		errorHandler:               defaultErrorHandler,
		allowedFilterFields:        defaultFilterFields,
		allowedOrderFields:         defaultOrderFields,
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := r.URL.Query().Get(options.queryParam)
			if s == "" {
				if !options.isSearchMandatory {
					next.ServeHTTP(w, r)
					return
				}

				options.errorHandler(w, r, fmt.Errorf("missing %q query parameter", options.queryParam))
				return
			}

			decoder := json.NewDecoder(strings.NewReader(s))
			decoder.DisallowUnknownFields()

			var search SearchRequest
			if err := decoder.Decode(&search); err != nil {
				options.errorHandler(w, r, err)
				return
			}

			if err := validateSearchRequest(&search, options); err != nil {
				options.errorHandler(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), searchKey, &search)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateSearchRequest(s *SearchRequest, opts *Options) error {
	// even though it is optional, if it is less than zero, it returns an error
	if s.Limit != nil && *s.Limit < 0 {
		return errors.New("limit must be null or >= 0")
	}

	if opts.limit != nil {
		if s.Limit == nil {
			return errors.New("limit is mandatory")
		}
		if *s.Limit > *opts.limit {
			return fmt.Errorf("limit must be between 0 and %d", *opts.limit)
		}
	}

	// even though it is optional, if it is less than zero, it returns an error
	if s.Offset != nil && *s.Offset < 0 {
		return errors.New("offset must be null or >= 0")
	}

	for _, o := range s.OrderBy {
		if _, ok := opts.allowedOrderFields[o.Field]; !ok {
			return fmt.Errorf("field %q not allowed in order by", o.Field)
		}
	}

	var validateGroup func(g *FilterGroup) error
	validateGroup = func(g *FilterGroup) error {
		if g == nil {
			return nil
		}

		if _, ok := opts.allowedLogicalOperators[g.Op]; !ok {
			return fmt.Errorf("logical operator %q not allowed", g.Op)
		}

		for _, f := range g.Filters {
			if _, ok := opts.allowedFilterFields[f.Field]; !ok {
				return fmt.Errorf("field %q not allowed in filters", f.Field)
			}

			if _, ok := opts.allowedRelationalOperators[f.Op]; !ok {
				return fmt.Errorf("relational operator %q not allowed for field %q", f.Op, f.Field)
			}
		}

		for _, sg := range g.Groups {
			if err := validateGroup(&sg); err != nil {
				return err
			}
		}

		return nil
	}

	if err := validateGroup(s.Groups); err != nil {
		return err
	}

	return nil
}

// GetSearchRequest retrieves the parsed SearchRequest stored in the
// request context by NewSearchHandler. If no request is stored, it
// returns nil.
func GetSearchRequest(r *http.Request) *SearchRequest {
	v := r.Context().Value(searchKey)
	if v == nil {
		return nil
	}

	s, ok := v.(*SearchRequest)
	if !ok {
		return nil
	}

	return s
}
