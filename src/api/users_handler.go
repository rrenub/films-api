package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"films-api.rdelgado.es/src/internals/models"
	"films-api.rdelgado.es/src/internals/validator"
)

type userRequest struct {
	Name                string `json:"name"`
	Password            string `json:"password"`
	validator.Validator `json:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	var req userRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	req.CheckField(validator.NoBlank(req.Name), "name", "This field must no be blank")
	req.CheckField(validator.NoBlank(req.Password), "password", "This field must no be blank")

	if !req.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(req.FieldErrors)
		return
	}

	id, err := app.users.Authenticate(req.Name, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.clientError(w, http.StatusUnauthorized, err)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	token, err := app.tokens.CreateToken(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {

	var req userRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	req.CheckField(validator.NoBlank(req.Name), "name", "This field must no be blank")
	req.CheckField(validator.NoBlank(req.Password), "password", "This field must no be blank")
	req.CheckField(validator.MinChars(req.Password, 8), "password", "Password must be 8 characters long")
	req.CheckField(validator.MaxChars(req.Password, 24), "password", "Password must be less than 24 characters long")
	req.CheckField(validator.Matches(req.Name, validator.UsernameRX), "name", "This field must start with a letter")
	req.CheckField(validator.IsStrongPassword(req.Password), "password", "Password must contain all characters")

	if !req.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(req.FieldErrors)
		return
	}

	err = app.users.Insert(req.Name, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedEntry) {
			app.clientError(w, http.StatusConflict, err)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}
