package repository

import (
	"github.com/testgithubanjali/url-shortener/internal/database"
	"github.com/testgithubanjali/url-shortener/internal/models"
)

func SaveURL(url *models.URL) error {
	return database.DB.Create(url).Error
}

func GetURLByShortCode(shortCode string) (*models.URL, error) {

	var url models.URL

	err := database.DB.
		Where("short_code = ?", shortCode).
		First(&url).Error

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func UpdateClickCount(shortCode string) error {

	return database.DB.
		Model(&models.URL{}).
		Where("short_code = ?", shortCode).
		Update("click_count", database.DB.Raw("click_count + 1")).Error
}
