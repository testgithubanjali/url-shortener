package services

import (
	"github.com/testgithubanjali/url-shortener/internal/models"
	"github.com/testgithubanjali/url-shortener/internal/repository"
)

func SaveShortURL(url *models.URL) error {
	return repository.SaveURL(url)
}
