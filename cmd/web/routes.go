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

	// Create a new middleware chain containing the middleware specific to our
	// dynamic applicaiton routes. This middleware automatically loads and saves session data with every HTTP request and response.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippets/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippets/create", dynamic.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// Create a standard reusable middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// use commonHeaders middleware
	return standard.Then(mux)
}
