package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"films-api.rdelgado.es/src/internals/models"
	"films-api.rdelgado.es/src/internals/validator"
	"github.com/julienschmidt/httprouter"
)

type movieRequest struct {
	Title               string
	Director            string
	ReleaseDate         string
	Cast                []string `json:"stringArray"`
	Genre               string
	Synopsis            string
	validator.Validator `json:"-"`
}

func (app *application) updateMovie(w http.ResponseWriter, r *http.Request) {

	// Get ID of movie to update
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil && id < 1 {
		app.NotFound(w)
		return
	}

	// Parse request fields to be updated
	var req movieRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	// Get movie to be update given the id
	movieToUpdate, err := app.movies.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Check that the user is the creator of this film
	userId := r.Context().Value(userIdContextKey).(int)

	app.logger.Info("userid", "id", userId)

	if movieToUpdate.UserID != uint(userId) {
		app.clientError(w, http.StatusForbidden, nil)
		return
	}

	// Validate each field from request to pass it to the model
	if validator.NoBlank(req.Title) {
		movieToUpdate.Title = req.Title
	}

	if validator.NoBlank(req.Director) {
		movieToUpdate.Director = req.Director
	}

	if validator.NoBlank(req.ReleaseDate) {
		parseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			app.clientError(w, http.StatusBadRequest, err)
			return
		}
		movieToUpdate.ReleaseDate = parseDate
	}

	if validator.NoEmptyTextSlice(req.Cast) {
		movieToUpdate.Cast = req.Cast
	}

	if validator.NoBlank(req.Genre) {
		movieToUpdate.Genre = req.Genre
	}

	if validator.NoBlank(req.Synopsis) {
		movieToUpdate.Synopsis = req.Synopsis
	}

	// Get movie to be updated
	err = app.movies.Update(movieToUpdate)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {

	// Get query params for optional filtering movies
	filters := make(map[string]interface{})
	if title := r.URL.Query().Get("title"); title != "" {
		filters["title"] = title
	}
	if genre := r.URL.Query().Get("genre"); genre != "" {
		filters["genre"] = genre
	}
	if year := r.URL.Query().Get("year"); year != "" {

		yearNumber, err := strconv.Atoi(year) // Check if year is a valid int (2019, 2022, etc)
		if err != nil && yearNumber < 1 {
			app.clientError(w, http.StatusBadRequest, err)
			return
		}

		filters["year"] = yearNumber
	}

	// Query movies from database (using filters if any)
	movies, err := app.movies.GetAll(filters)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(userIdContextKey).(int)

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil && id < 1 {
		app.NotFound(w)
		return
	}

	err = app.movies.Delete(id, userId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else if errors.Is(err, models.ErrNotAuthorized) {
			app.clientError(w, http.StatusForbidden, err)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil && id < 1 {
		app.NotFound(w)
		return
	}

	movie, err := app.movies.GetMovieAndAuthor(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(movie)
}

func (app *application) addMovie(w http.ResponseWriter, r *http.Request) {

	var req movieRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	// Check and parse release date
	parsedReleaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	req.CheckField(validator.NoBlank(req.Title), "title", "This field must no be blank")
	req.CheckField(validator.NoBlank(req.Director), "director", "This field must no be blank")
	req.CheckField(validator.NoBlank(req.Genre), "genre", "This field must no be blank")
	req.CheckField(validator.NoBlank(req.Synopsis), "synopsis", "This field must no be blank")
	req.CheckField(validator.NoEmptyTextSlice(req.Cast), "cast", "This field must not be empty")

	if !req.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(req.FieldErrors)
		return
	}

	// Get user_id for insertion
	userId := r.Context().Value(userIdContextKey).(int)

	id, err := app.movies.Insert(req.Title, req.Director, req.Genre, req.Synopsis, parsedReleaseDate, req.Cast, userId)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedEntry) {
			app.clientError(w, http.StatusConflict, err)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
}
