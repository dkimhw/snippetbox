package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// inject dependencies by using application struct
// useful for when all handlers are in the same package
type application struct {
	logger *slog.Logger
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)) // initialize a new structured logger

	app := &application{
		logger: logger,
	} // initialize a new instance of applicaiton struct containing dependencies

	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("starting server", slog.String("addr", *addr))

	// And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	mux := app.routes()
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
