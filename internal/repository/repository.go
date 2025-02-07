package repository

import "github.com/cyberjourney20/career-journey/internal/models"

type DatabaseRepo interface {
	GetAllContacts() ([]models.Contact, error)
	Authenticate(email, testPassword string) (string, string, error)
	GetFavoriteContacts() ([]models.Contact, error)
	AddNewContact(c models.Contact, user_id string) (int, error)
	CompanyExists(c models.Contact) (int, error)
}
