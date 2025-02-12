package repository

import (
	"github.com/cyberjourney20/career-journey/internal/models"
)

type DatabaseRepo interface {
	Authenticate(email, testPassword string) (string, string, error)
	GetAllContacts() ([]models.Contact, error)
	GetFavoriteContacts() ([]models.Contact, error)
	GetContactByID(id int, user_id string) (models.Contact, error)
	AddNewContact(c models.Contact, user_id string) (int, error)
	CompanyExists(cmp models.Company) (int, error)
	GetAllCompanies() ([]models.Company, error)
	AddNewCompany(cmp models.Company) (int, error)

	// May not need this here if it is never called from outside of postgres.go
	//or moveit to helpers.go ? it dosnt interact with the DB, Just builds a JSON for it.
	//TimeLineBuilderJSON(t time.Time, s string) string
}
