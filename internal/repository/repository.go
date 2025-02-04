package repository

import "github.com/cyberjourney20/career-journey/internal/models"

type DatabaseRepo interface {
	GetAllContacts() ([]models.Contact, error)
	Authenticate(email, testPassword string) (string, string, error)
	GetFavoriteContacts() ([]models.Contact, error)
}
