package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cyberjourney20/career-journey/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) GetAllContacts() ([]models.Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var contacts []models.Contact

	query := `
	select c.id, c.first_name, c.last_name, c.job_title, c.email, c.mobile_phone, 
	c.work_phone, c.phone_3, c.linkedin, c.github, c.website, c.notes, c.description, c.objective,
	c.timeline, c.favorite, c.created_at, c.updated_at, cp.company_name 
	from contacts c
	left join companies cp on c.company_id=cp.id
	where c.user_id='2041a735-715f-4fb1-a08c-1055fe1916e6'`

	// Get UserID in QueryContext call from session
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return contacts, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Contact
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.JobTitle,
			&i.Email,
			&i.MobilePhone,
			&i.WorkPhone,
			&i.Phone3,
			&i.Linkedin,
			&i.Github,
			&i.Website,
			&i.Notes,
			&i.Description,
			&i.Objective,
			&i.ContactTimeLine,
			&i.Favorite,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Company.CompanyName,
		)
		if err != nil {
			return contacts, err
		}
		// add an if statement to only check for future reservations.
		contacts = append(contacts, i)
	}
	if err = rows.Err(); err != nil {
		return contacts, err
	}
	return contacts, nil
}

func (m *postgresDBRepo) GetFavoriteContacts() ([]models.Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var contacts []models.Contact

	query := `
	select c.id, c.first_name, c.last_name, c.job_title, c.email, c.mobile_phone, 
	c.work_phone, c.phone_3, c.linkedin, c.github, c.website, c.notes, c.description, c.objective,
	c.timeline, c.favorite, c.created_at, c.updated_at, cp.company_name 
	from contacts c
	left join companies cp on c.company_id=cp.id
	where c.user_id='2041a735-715f-4fb1-a08c-1055fe1916e6' and favorite='true'`

	// Get UserID in QueryContext call from session
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return contacts, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Contact
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.JobTitle,
			&i.Email,
			&i.MobilePhone,
			&i.WorkPhone,
			&i.Phone3,
			&i.Linkedin,
			&i.Github,
			&i.Website,
			&i.Notes,
			&i.Description,
			&i.Objective,
			&i.ContactTimeLine,
			&i.Favorite,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Company.CompanyName,
		)
		if err != nil {
			return contacts, err
		}
		// add an if statement to only check for future reservations.
		contacts = append(contacts, i)
	}
	if err = rows.Err(); err != nil {
		return contacts, err
	}
	return contacts, nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id string
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select user_id, password from users where email =$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		log.Println("func Authenticate, error 1")
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("func Authenticate, error 2")
		return "", "", errors.New("incorrect password")
	} else if err != nil {
		log.Println("func Authenticate, error 3")
		return "", "", err
	}
	log.Println("func Authenticate, no error")
	return id, hashedPassword, nil
}

//CreateNewUser
//UpdateUser

//CreateContact
//UpdateContact
//DeleteContact

//CreateJobListing
//UpdateJobLinting
//DeleteJobListing

//CreateCert
//UpdateCert
//DeleteCert

//CreateSkill
//UpdateSkill
//DeleteSkil

//CreateCompany
//UpdateCompany
//DeleteCompany
