package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	/*router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.clientError(w, r, "Not found")
	})*/

	// Movies endpoints (auth required)
	router.Handler(http.MethodGet, "/movies", app.requireAuthentication(http.HandlerFunc(app.getAllMovies)))
	router.Handler(http.MethodPost, "/movie", app.requireAuthentication(http.HandlerFunc(app.addMovie)))
	router.Handler(http.MethodGet, "/movie/:id", app.requireAuthentication(http.HandlerFunc(app.getMovie)))
	router.Handler(http.MethodDelete, "/movie/:id", app.requireAuthentication(http.HandlerFunc(app.deleteMovie)))
	router.Handler(http.MethodPut, "/movie/:id", app.requireAuthentication(http.HandlerFunc(app.updateMovie)))

	// Favourites movies endpoints (auth required)
	router.Handler(http.MethodPost, "/favourites", app.requireAuthentication(http.HandlerFunc(app.addMovieToFav)))
	router.Handler(http.MethodGet, "/favourites", app.requireAuthentication(http.HandlerFunc(app.getFavMovies)))
	router.Handler(http.MethodDelete, "/favourites/:id", app.requireAuthentication(http.HandlerFunc(app.deleteMovieFromFav)))

	// Authentication endpoints
	router.HandlerFunc(http.MethodPost, "/user/signup", app.userSignup)
	router.HandlerFunc(http.MethodPost, "/user/login", app.userLogin)

	return app.recoverPanic(app.logRequest(app.authenticate(app.logResponse(router))))
}
