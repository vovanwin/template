package httpserver

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	address           string           `option:"mandatory"`
	readHeaderTimeout time.Duration    `option:"mandatory"`
	middlewareSetup   func(r *chi.Mux) // кастомные мидлваре
}

func NewServer(opts Options) (*chi.Mux, *http.Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate options error: %w", err)
	}
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "endpoint not found"}`))
	})

	opts.middlewareSetup(r)

	httpServer := &http.Server{
		Addr:              opts.address,
		Handler:           r,
		ReadHeaderTimeout: opts.readHeaderTimeout,
	}

	return r, httpServer, nil
}
