package main

/*
& goes in front of a variable when you want to get that variable's memory address.

`*` goes in front of a variable that holds a memory address and resolves it (it is therefore the counterpart to the & operator).

*/

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	// Not using anything in this package but need the `init()` function to run so
	// it can register itself with the `database/sql` package
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"snippetbox.dkimhw.com/internal/models"
)

// inject dependencies by using application struct
// useful for when all handlers are in the same package
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	godotenv.Load()

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", fmt.Sprintf("web:%s@/snippetbox?parseTime=true", os.Getenv("DB_PASSWORD")), "MySQL data source name")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)) // initialize a new structured logger

	// Pass openDB() the DSN from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close() // defer call - closes db before main() function closes

	templateCache, err := newTemplateCache() // initialize a new template cache
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db) // uses mysql to manage sessions
	sessionManager.Lifetime = 12 * time.Hour  // lifetime of 12 hours for each session

	// initialize a new instance of applicaiton struct containing dependencies
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db}, // Initialize a models.SnippetModel instance containing the connection pool
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Use the Info() method to log the starting server message at Info severity
	// (along with the listen address as an attribute).
	logger.Info("starting server", slog.String("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes()) // Pass the dereferenced addr pointer
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
