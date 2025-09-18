# qparams - boilerplate code for search api

`qparams` is a lightweight library for parsing and validating complex search requests provided via HTTP query parameters.  
It provides a configurable middleware that decodes JSON-based search payloads, validates operators and fields, and injects a structured `SearchRequest` into the request context.

## Project goal

The purpose of qparams is not to generate SQL queries or any ready-to-use database query language.
Instead, it provides a structured and validated SearchRequest object inside the HTTP request context.
This object can then be consumed by your application logic (e.g., ORM, query builder, or custom repository layer) to build queries in a way that best fits your use case.

## Features

- Parse structured search queries from a query parameter (default: `q`)  
- Enforce allowed logical and relational operators  
- Restrict filterable and sortable fields  
- Customizable defaults and per-handler overrides  
- Pluggable error handler  

---

## Installation

```bash
go get github.com/paccolamano/qparams
```

## Usage

The default handler has preset values and is ready to use out of the box.
The only mandatory values to set are the fields allowed for searching and sorting.

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/paccolamano/qparams"
)

func main() {
	// Create a handler that extracts the SearchRequest
	searchHandler := qparams.NewSearchHandler(
		qparams.WithFilterFields("id", "created_at", "updated_at", "name"), // Allow records to be filtered by (only) these fields
		qparams.WithOrderFields("name"), // Allow records ordering only by name field
		qparams.WithLimit(10) // Allow pagination between 0 and 10. if null, negative or >10 an error is automatically returned
	)

	usersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrive search request
		search := qparams.GetSearchRequest(r)
		if search == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Missing or invalid search request")
			return
		}

		// use search request in your business logic

		fmt.Fprintf(w, "Parsed search request: %+v\n", search)
	})

	// Wrap your handler with search handler
	http.Handle("/users", searchHandler(usersHandler))

	http.ListenAndServe(":8080", nil)
}
```

For other examples see the _examples_ folder.
