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
	qparams.SetDefaultSearchMandatory(true) // The api consumer must send query param for the search
	qparams.SetDefaultLimit(10)             // By setting a limit, the consumer must send it and must be between 0 and given number (included)

	// Wrap your handler (in this case usersHandler) with the search handler.
	// Without defining the filter and order fields, the api consumer can only paginate the data.
	mux.Handle("GET /api/v1/users", qparams.NewSearchHandler()(&usersHandler{}))

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
