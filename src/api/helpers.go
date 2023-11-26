package main

import (
	"net/http"
	"runtime/debug"

	"films-api.rdelgado.es/src/internals/models"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int, err error) {
	if err != nil {
		app.logger.Error(err.Error())
	}

	http.Error(w, http.StatusText(status), status)
}

func (app *application) NotFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound, models.ErrNoRecord)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
