package dbrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

func (m *postgresDBRepo) GetContactByID(id int, user_id string) (models.Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var contact models.Contact

	query := `
	select c.id, c.first_name, c.last_name, c.job_title, c.email, c.mobile_phone, 
	c.work_phone, c.phone_3, c.linkedin, c.github, c.website, c.notes, c.description, c.objective,
	c.timeline, c.favorite, c.created_at, c.updated_at, cp.company_name 
	from contacts c
	left join companies cp on c.company_id=cp.id
	where c.user_id='$1' and c.id='$2'`

	// Get UserID in QueryContext call from session
	row := m.DB.QueryRowContext(ctx, query, user_id, id)
	err := row.Scan(
		&contact.ID,
		&contact.FirstName,
		&contact.LastName,
		&contact.JobTitle,
		&contact.Email,
		&contact.MobilePhone,
		&contact.WorkPhone,
		&contact.Phone3,
		&contact.Linkedin,
		&contact.Github,
		&contact.Website,
		&contact.Notes,
		&contact.Description,
		&contact.Objective,
		&contact.ContactTimeLine,
		&contact.Favorite,
		&contact.CreatedAt,
		&contact.UpdatedAt,
		&contact.Company.CompanyName,
	)
	if err != nil {
		return contact, err
	}
	return contact, nil
}

// SearchCompany foir use in new contacts form
func (m *postgresDBRepo) CompanyExists(c models.Contact) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select cmp.id, cmp.company_name
	from companies cmp
	where lower(cmp.company_name)=lower($1)`
	row := m.DB.QueryRowContext(ctx, query, c.Company.CompanyName)
	var cpny models.Company
	err := row.Scan(
		&cpny.ID,
		&cpny.CompanyName,
	)
	if err == sql.ErrNoRows {
		return 0, nil // No company found, return 0 but not an error
	} else {
		if err != nil {
			return 0, err
		}
	}
	return cpny.ID, nil

}

// func (m *postgresDBRepo) AddNewCompany(cpny_name string) int {
// 		stmt1 := `insert into companies (company_name, created_at, updated_at) values ($1, $2, $2) returning id`

// 		err := m.DB.QueryRowContext(ctx, stmt1, cpny_name, time.Now()).Scan(&newID)
// 		if err != nil {
// 			fmt.Println("Error in AddNewContact insert company", err)
// 			return 0, err
// 		}
// 		returnID = newID
// 		log.Printf("Got New ID: %d", newID)
// 	return cpnyID
// }

// AddNewContact Adds a new contact to the database
func (m *postgresDBRepo) AddNewContact(c models.Contact, user_id string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// check for no company listed
	// Check if company exists, if yes, get existing ID, if no, add it and get a new ID
	tempID, err := m.CompanyExists(c)
	var returnID int
	if err != nil {
		log.Println("error in CompanyExists", err)
	}

	if tempID > 0 {
		returnID = tempID
		log.Printf("Temp ID: %d", tempID)
	} else {
		//AddNewCompany(c.Company.CompanyName)
		stmt1 := `insert into companies (company_name, created_at, updated_at) values ($1, $2, $2) returning id`
		err := m.DB.QueryRowContext(ctx, stmt1, c.Company.CompanyName, time.Now()).Scan(&returnID)
		if err != nil {
			fmt.Println("Error in AddNewContact insert company", err)
			return 0, err
		}
		log.Printf("Got New ID: %d", returnID)
	}

	stmt2 := `
	insert into contacts (first_name, last_name, job_title, email, objective, mobile_phone, work_phone, 
	linkedin, github, description, notes, user_id, timeline, favorite, created_at, updated_at, company_id) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	c.ContactTimeLine = TimeLineBuilderJSON(time.Now(), "Contact Created")

	_, err = m.DB.ExecContext(ctx, stmt2,
		c.FirstName,
		c.LastName,
		c.JobTitle,
		c.Email,
		c.Objective,
		c.MobilePhone,
		c.WorkPhone,
		c.Linkedin,
		c.Github,
		c.Description,
		c.Notes,
		user_id,
		c.ContactTimeLine,
		c.Favorite,
		time.Now(),
		time.Now(),
		returnID,
	)
	if err != nil {
		fmt.Println("Error in AddNewContact insert contact", err)
		return 0, err
	}
	fmt.Println("AddNewContact insert contact Ran")
	return 1, nil
}

func TimeLineBuilderJSON(t time.Time, s string) string {
	timeline := map[string]string{
		"marker": t.Format(time.RFC3339), // Use a standard format
		"event":  s,
	}
	// Marshal the map into JSON bytes
	jsonBytes, err := json.Marshal(timeline)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return ""
	}

	// Return the JSON as a string
	return string(jsonBytes)

	// ts := t.String()
	// json := fmt.Sprintf("{'marker':'%s', 'event':'%s'}", ts, s)
	// return json

}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id string
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select user_id, password from users where lower(email)=lower($1)", email)
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
