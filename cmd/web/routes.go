package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/")) // create file server

	// Register the file server as the handler for all URL paths that start with "/static/".
	// For matching paths, we strip the "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware specific to our
	// dynamic applicaiton routes. This middleware automatically loads and saves session data with every HTTP request and response.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf) // Unprotected application routes using the "dynamic" middleware chain.

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected application routes, using a new "protected" middleware chain which includes
	// requireAuthentication middleware
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create a standard reusable middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// use commonHeaders middleware
	return standard.Then(mux)
}
