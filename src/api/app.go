package main

import (
	"log/slog"

	"films-api.rdelgado.es/src/internals/authentication"
	"films-api.rdelgado.es/src/internals/models"
)

type application struct {
	logger *slog.Logger
	movies *models.MovieModel
	users  *models.UserModel
	favs   *models.FavouriteModel
	tokens *authentication.JwtToken
}

func NewApplication(logger *slog.Logger, movies *models.MovieModel, users *models.UserModel, favs *models.FavouriteModel, tokens *authentication.JwtToken) *application {
	return &application{
		logger: logger,
		movies: movies,
		users:  users,
		favs:   favs,
		tokens: tokens,
	}
}
