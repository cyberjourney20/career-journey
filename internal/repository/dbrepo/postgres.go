package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cyberjourney20/career-journey/internal/models"
	"github.com/cyberjourney20/career-journey/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// func (m *postgresDBRepo) AllUsers() bool {
// 	return true
// }

// Authenticate authenticates a user by comparing password hash to stored password hash
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

// GetAllContacts gets all contacts associated with a user's UUID
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

// GetAllContacts gets all contacts identified as a favorite associated with a user's UUID
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

// GetContactByID gets a contact by thier ID
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
	where c.id=$1`

	// Get UserID in QueryContext call from session
	row := m.DB.QueryRowContext(ctx, query, id)
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

// CompanyExists searches for an existing company by Name and returns its ID
func (m *postgresDBRepo) CompanyExists(cmp models.Company) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select cmp.id, cmp.company_name
	from companies cmp
	where lower(cmp.company_name)=lower($1)`
	row := m.DB.QueryRowContext(ctx, query, cmp.CompanyName)
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

// AddNewCompany adds a new company to companies table
func (m *postgresDBRepo) AddNewCompany(cmp models.Company) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//check if company exists  with func CompanyExists in parent function

	var returnID int
	stmt := `insert into companies (company_name, url, address, industry, size, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $6) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		cmp.CompanyName,
		cmp.URL,
		cmp.Address,
		cmp.Industry,
		cmp.Size,
		time.Now(),
	).Scan(&returnID)

	if err != nil {
		fmt.Println("Error in AddNewContact insert company", err)
		return 0, err
	}
	log.Printf("Got New ID: %d", returnID)
	return returnID, nil

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
	tempID, err := m.CompanyExists(c.Company)
	var returnID int
	if err != nil {
		log.Println("error in CompanyExists", err)
	}
	if tempID > 0 {
		returnID = tempID
		log.Printf("Temp ID: %d", tempID)
	} else {
		returnID, err := m.AddNewCompany(c.Company)
		if err != nil {
			log.Println("error in AddNewCompany", err)
			return 0, err
		}
		return returnID, nil
	}
	log.Printf("Got New Company ID: %d", returnID)

	stmt := `
	insert into contacts (first_name, last_name, job_title, email, objective, mobile_phone, work_phone, 
	linkedin, github, description, notes, user_id, timeline, favorite, created_at, updated_at, company_id) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	c.ContactTimeLine = utils.TimeLineBuilderJSON(time.Now(), "Contact Created")

	_, err = m.DB.ExecContext(ctx, stmt,
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

// Get all Companies returns a slice of models.company with all companies
func (m *postgresDBRepo) GetAllCompanies() ([]models.Company, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var companies []models.Company

	query := `
	select company.id, company.company_name, company.url, company.address, 
	company.industry, company.size, created_at, updated_at
	from companies cmp
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return companies, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Company
		err := rows.Scan(
			&i.ID,
			&i.CompanyName,
			&i.URL,
			&i.Address,
			&i.Industry,
			&i.Size,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return companies, err
		}
		companies = append(companies, i)
	}
	if err = rows.Err(); err != nil {
		return companies, err
	}
	return companies, nil
}

// UdateContactByID updates a contact details by ID
func (m *postgresDBRepo) UpdateContactByID(c models.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update contacts (company_id =$1, first_name=$2, last_name=$3, job_title=$4, 
	email=$5, mobile_phone=$6, work_phone=$7, phone_3=$8, linkedin=$9, github=$10, 
	website=$11, notes=$12, description=$13, objective=$14, timeline=$15, favorite=$16, updated_at=$17) 
	where id = $18
	`
	//validate this work right
	c.ContactTimeLine = utils.TimeLineBuilderJSON(time.Now(), "Contact Updated")

	_, err := m.DB.ExecContext(ctx, stmt,
		c.CompanyID, // I need to think about this. I need to check if it changed, then add new company or update
		c.FirstName,
		c.LastName,
		c.JobTitle,
		c.Email,
		c.MobilePhone,
		c.WorkPhone,
		c.Phone3,
		c.Linkedin,
		c.Github,
		c.Website,
		c.Notes,
		c.Description,
		c.Objective,
		c.ContactTimeLine,
		c.Favorite,
		time.Now(),
		c.ID,
	)
	if err != nil {
		return err
	}
	return nil

}

// UserExists checks to see if a username is already registered in the DB
func (m *postgresDBRepo) UserExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var dbEmail string

	row := m.DB.QueryRowContext(ctx, "select email from users where lower(email)=lower($1)", email)
	err := row.Scan(&dbEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("UserExists sql.ErrNoRows")
			return false, nil // No user found, which means email doesn't exist
		}
		fmt.Println("UserExists Error 2")
		return false, err // Some other error occurred
	}
	fmt.Println("UserExists No Error")
	return true, nil // User found, email exists
}

// AddNewUser adds a new user to the database
func (m *postgresDBRepo) AddNewUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var err error
	bpass := []byte(u.Password)
	phash, err := bcrypt.GenerateFromPassword([]byte(bpass), 12)
	if err != nil {
		fmt.Println("AddNewUser in in bcrypt", err)
		return err
	}

	stmt := `
	insert into users (first_name, last_name, email, password, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6)
	`
	_, err = m.DB.ExecContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		phash,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		fmt.Println("Error in AddNewUser insert user", err)
		return err
	}
	fmt.Println("AddNewUser no error")
	return nil
}

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
