package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
// func (app *application) clientError(w http.ResponseWriter, status int) {
// 	http.Error(w, http.StatusText(status), status)
// }

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page] // check the cache if the page exists
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)                     // initalize a new buffer
	err := ts.ExecuteTemplate(buf, "base", data) // write template to buffer
	if err != nil {
		// if there is an error writing to buffer, call serverError()
		app.serverError(w, r, err)
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		// add the flash message to the template data if one exists
		// Use the PopString() method to retrieve the value for the "flash" key.
		// PopString() also deletes the key and value from the session data, so it
		// acts like a one-time fetch. If there is no matching key in the session
		// data this will return the empty string.
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}

// dst is the target destination that we want to decode the form data into
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}

	return nil
}
