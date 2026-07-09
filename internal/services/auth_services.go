package services

import (
	"github.com/testgithubanjali/url-shortener/internal/repository"
)

func FindUserByEmail(email string) error {

	_, err := repository.GetUserByEmail(email)

	return err
}
