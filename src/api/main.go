package main

import (
	"log/slog"
	"net/http"
	"os"

	"films-api.rdelgado.es/src/internals/authentication"
	"films-api.rdelgado.es/src/internals/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//Estas variables debería cogerlas del entorno
	dsn := "movies_user:movies789@tcp(127.0.0.1:3306)/moviesdb?charset=utf8mb4&parseTime=True&loc=Local"
	secretJwt := []byte("test_secret")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = db.AutoMigrate(&models.User{}, &models.Movie{}, &models.Favourite{})
	if err != nil && db.Migrator().HasTable(&models.Movie{}) {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger: logger,
		movies: &models.MovieModel{DB: db},
		users:  &models.UserModel{DB: db},
		favs:   &models.FavouriteModel{DB: db},
		tokens: &authentication.JwtToken{SecretJwt: secretJwt},
	}

	app.seedData(db)

	addr := ":4000"
	logger.Info("stating movies api server", slog.String("port", addr))

	server := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
