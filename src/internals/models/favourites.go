package models

import (
	"errors"

	"gorm.io/gorm"
)

type FavouriteModel struct {
	DB *gorm.DB
}

type Favourite struct {
	gorm.Model
	UserID  uint `gorm:"uniqueIndex:idx_userid_movieid"`
	MovieID uint `gorm:"uniqueIndex:idx_userid_movieid"`
}

type GetFavouriteInfo struct {
	FavouriteID uint  `gorm:"column:fav_id"`
	Movie       Movie `gorm:"embedded"`
}

func (m *FavouriteModel) Remove(favId, userId int) error {

	result := m.DB.Unscoped().
		Where("id = ?", favId).
		Where("user_id = ?", userId).
		Delete(&Favourite{})

	if err := result.Error; err != nil {
		return err
	}

	// return Not Found error if no record has been deleted
	if result.RowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *FavouriteModel) GetAll(userId int) ([]GetFavouriteInfo, error) {
	var movieDetails []GetFavouriteInfo

	result := m.DB.Model(&Favourite{}).Select("favourites.id AS fav_id", "movies.*").
		Joins("LEFT JOIN users ON favourites.user_id = users.id").
		Joins("LEFT JOIN movies ON favourites.movie_id = movies.id").
		Where("favourites.user_id = ?", userId).
		Scan(&movieDetails)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []GetFavouriteInfo{}, ErrNoRecord
		} else {
			return []GetFavouriteInfo{}, result.Error
		}
	}

	return movieDetails, nil
}

func (m *FavouriteModel) Insert(userId, movieId int) (int, error) {

	favorite := Favourite{
		UserID:  uint(userId),
		MovieID: uint(movieId),
	}

	result := m.DB.Create(&favorite)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, ErrDuplicatedEntry
		}

		return 0, err
	}

	return int(favorite.ID), nil
}
