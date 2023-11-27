package main

import (
	"log/slog"
	"net/http"
	"os"

	"films-api.rdelgado.es/src/internals/authentication"
	"films-api.rdelgado.es/src/internals/models"
)

func main() {

	// enviroment variables
	db_hostname := os.Getenv("MYSQL_HOSTNAME")
	db_name := os.Getenv("MYSQL_DATABASE")
	db_user := os.Getenv("MYSQL_USER")
	db_password := os.Getenv("MYSQL_PASSWORD")
	jwt_secret := os.Getenv("JWT_SECRET")
	server_port := os.Getenv("API_PORT")

	// init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// init database conn
	db, err := InitDB(db_hostname, db_name, db_user, db_password)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// migrate schemas db
	err = db.AutoMigrate(&models.User{}, &models.Movie{}, &models.Favourite{})
	if err != nil && db.Migrator().HasTable(&models.Movie{}) {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// create app struct (models, etc)
	app := &application{
		logger: logger,
		movies: &models.MovieModel{DB: db},
		users:  &models.UserModel{DB: db},
		favs:   &models.FavouriteModel{DB: db},
		tokens: &authentication.JwtToken{SecretJwt: []byte(jwt_secret)},
	}

	// seed db (if db is empty)
	app.seedDB(db)

	// init http server
	addr := ":" + server_port

	server := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("stating movies api server", slog.String("port", addr))

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
