package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/paccolamano/qparams"
)

func main() {
	mux := http.NewServeMux()

	// Set default for all search handler
	qparams.SetDefaultQueryParam("s")
	qparams.SetDefaultLimit(50)
	qparams.SetDefaultFilterFields("id")
	qparams.SetDefaultOrderFields("id")
	qparams.SetDefaultErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	})

	search := qparams.NewSearchHandler(
		// Set the allowed filter query parameters in this handler. It merge default filter fields (id) with provided ones (name, email)
		qparams.WithExtraFilterFields("name", "email"),
		// Set the allowed order query parameters in this handler. It replace default order fields with provided ones (created_at, updated_at)
		qparams.WithOrderFields("created_at", "updated_at"),
		// Set the default limit for this handler (override default limit)
		qparams.WithLimit(10),
	)
	// Wrap your handler (in this case usersHandler) with the search handler
	mux.Handle("GET /api/v1/users", search(&usersHandler{}))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

type usersHandler struct{}

func (h usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Retrieve your search and use it to filter, sort and paginate the requested data
	s := qparams.GetSearchRequest(r)
	if s == nil {
		log.Println("no search provided")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Printf("search request: %+v", *s)

	w.Header().Add("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		log.Println("failed to send response")
	}
}
