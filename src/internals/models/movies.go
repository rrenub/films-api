package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MovieModel struct {
	DB *gorm.DB
}

type Cast []string

type Movie struct {
	gorm.Model

	Title       string    `gorm:"unique; not null"`
	Director    string    `gorm:"not null"`
	ReleaseDate time.Time `gorm:"not null"`
	Cast        Cast      `gorm:"serializer:json"`
	Genre       string    `gorm:"not null"`
	Synopsis    string    `gorm:"not null"`

	UserID uint
}

type MovieAndAuthor struct {
	Movie     `json:"movie" gorm:"embedded"`
	CreatedBy `json:"created_by"`
}

type CreatedBy struct {
	Name   string `json:"name"`
	UserId uint   `json:"userId"`
}

func (m *MovieModel) GetAll(filters map[string]interface{}) ([]Movie, error) {
	var movies []Movie

	query := m.DB
	title, exists := filters["title"]
	if exists {
		query = query.Where("title LIKE ?", "%"+title.(string)+"%")
	}

	genre, exists := filters["genre"]
	if exists {
		query = query.Where("genre = ?", genre)
		fmt.Print(genre)
	}

	year, exists := filters["year"]
	if exists {
		startDate := time.Date(year.(int), 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year.(int)+1, 1, 1, 0, 0, 0, 0, time.UTC)
		query = query.Where("release_date >= ? AND release_date < ?", startDate, endDate)
	}

	result := query.Model(&Movie{}).Find(&movies)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoRecord
		} else {
			return nil, result.Error
		}
	}

	return movies, nil
}

func (m *MovieModel) GetMovieAndAuthor(id int) (MovieAndAuthor, error) {
	var movie MovieAndAuthor

	result := m.DB.Model(&Movie{}).
		Select("movies.*", "users.id as user_id", "users.name").
		Joins("INNER JOIN users ON movies.user_id = users.id").
		Where("movies.id = ?", id).
		Scan(&movie)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return MovieAndAuthor{}, ErrNoRecord
		} else {
			return MovieAndAuthor{}, result.Error
		}
	}

	return movie, nil
}

func (m *MovieModel) Get(id int) (Movie, error) {
	var movie Movie

	result := m.DB.First(&movie, id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Movie{}, ErrNoRecord
		} else {
			return Movie{}, result.Error
		}
	}

	return movie, nil
}

func (m *MovieModel) Update(movie Movie) error {
	result := m.DB.Save(&movie)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (m *MovieModel) Insert(title, director, genre, synopsis string, releaseDate time.Time, cast []string, userId int) (int, error) {
	movie := &Movie{
		Title:       title,
		Director:    director,
		ReleaseDate: releaseDate,
		Cast:        cast,
		Genre:       genre,
		Synopsis:    synopsis,
		UserID:      uint(userId),
	}

	result := m.DB.Create(movie)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, ErrDuplicatedEntry
		} else {
			return 0, err
		}
	}

	return int(movie.ID), nil
}

func (m *MovieModel) Delete(id, userId int) error {

	// Retrieve movie to check user who created the movie
	movie, err := m.Get(id)
	if err != nil {
		return err
	}

	// If the user did not create the movie => not allowed to delete
	if movie.UserID != uint(userId) {
		return ErrNotAuthorized
	}

	result := m.DB.Unscoped().Delete(&Movie{}, id)
	if err := result.Error; err != nil {
		return err
	} else {
		return nil
	}
}
