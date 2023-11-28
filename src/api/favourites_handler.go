package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"films-api.rdelgado.es/src/internals/models"
	"films-api.rdelgado.es/src/internals/validator"
	"github.com/julienschmidt/httprouter"
)

type favouriteMovieRequest struct {
	MovieID             int `json:"movie_id"`
	validator.Validator `json:"-"`
}

func (app *application) getFavMovies(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdContextKey).(int)

	favMovies, err := app.favs.GetAll(userId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	if favMovies == nil {
		app.NotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(favMovies)
}

func (app *application) deleteMovieFromFav(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdContextKey).(int)

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil && id < 1 {
		app.NotFound(w)
		return
	}

	err = app.favs.Remove(id, userId)
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
}

func (app *application) addMovieToFav(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIdContextKey).(int)

	var req favouriteMovieRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	req.CheckField(validator.IsPositiveNumber(req.MovieID), "movie", "This field must be movie ID")

	if !req.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(req.FieldErrors)
		return
	}

	id, err := app.favs.Insert(userId, req.MovieID)
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
