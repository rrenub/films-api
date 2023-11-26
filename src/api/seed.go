package main

import (
	"errors"
	"time"

	"films-api.rdelgado.es/src/internals/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (app *application) seedData(db *gorm.DB) {
	if db.Migrator().HasTable(&models.Movie{}) {

		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {

			app.logger.Info("populating users database...")

			// Hash password for test
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Test.1234"), 12)
			if err != nil {
				return
			}

			users := []models.User{
				{Name: "test1", Password: string(hashedPassword)},
				{Name: "test2", Password: string(hashedPassword)},
				{Name: "test3", Password: string(hashedPassword)},
			}

			db.Create(&users)
		}

		if err := db.First(&models.Movie{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {

			app.logger.Info("populating movies database...")

			movies := [5]models.Movie{
				{
					Title:       "Inception",
					Director:    "Christopher Nolan",
					ReleaseDate: time.Date(2010, time.July, 16, 0, 0, 0, 0, time.UTC),
					Cast:        models.Cast{"Leonardo DiCaprio", "Joseph Gordon-Levitt", "Ellen Page"},
					Genre:       "Science Fiction",
					Synopsis:    "A thief who enters the dreams of others to steal their secrets.",
					UserID:      1,
				},
				{
					Title:       "Interstellar",
					Director:    "Christopher Nolan",
					ReleaseDate: time.Date(2014, time.October, 26, 0, 0, 0, 0, time.UTC),
					Cast:        models.Cast{"Matthew McConaughey", "Anne Hathaway", "Jessica Chastain"},
					Genre:       "Science Fiction",
					Synopsis:    "A group of explorers travels through a wormhole in space in an attempt to ensure humanity's survival.",
					UserID:      1,
				},
				{
					Title:       "The Dark Knight",
					Director:    "Christopher Nolan",
					ReleaseDate: time.Date(2008, time.July, 18, 0, 0, 0, 0, time.UTC),
					Cast:        models.Cast{"Christian Bale", "Heath Ledger", "Aaron Eckhart"},
					Genre:       "Action",
					Synopsis:    "Batman faces the Joker in a battle for Gotham City.",
					UserID:      1,
				},
				{
					Title:       "The Matrix",
					Director:    "Lana Wachowski, Lilly Wachowski",
					ReleaseDate: time.Date(1999, time.March, 31, 0, 0, 0, 0, time.UTC),
					Cast:        models.Cast{"Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss"},
					Genre:       "Science Fiction",
					Synopsis:    "A computer hacker learns about the true nature of his reality.",
					UserID:      1,
				},
				{
					Title:       "Pulp Fiction",
					Director:    "Quentin Tarantino",
					ReleaseDate: time.Date(1994, time.May, 21, 0, 0, 0, 0, time.UTC),
					Cast:        models.Cast{"John Travolta", "Samuel L. Jackson", "Uma Thurman"},
					Genre:       "Crime",
					Synopsis:    "Various interconnected stories of crime in Los Angeles.",
					UserID:      1,
				},
			}

			db.Create(&movies)
		}
	}
}
