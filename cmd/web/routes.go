package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	// create file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// register the file server as the handler for all URL paths
	// that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippets/create", app.snippetCreate)
	mux.HandleFunc("POST /snippets/create", app.snippetCreatePost)

	// Create a standard reusable middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// use commonHeaders middleware
	return standard.Then(mux)
}
