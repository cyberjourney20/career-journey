package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cyberjourney20/career-journey/internal/config"
	"github.com/cyberjourney20/career-journey/internal/driver"
	"github.com/cyberjourney20/career-journey/internal/forms"
	"github.com/cyberjourney20/career-journey/internal/helpers"
	"github.com/cyberjourney20/career-journey/internal/models"
	"github.com/cyberjourney20/career-journey/internal/render"
	"github.com/cyberjourney20/career-journey/internal/repository"
	"github.com/cyberjourney20/career-journey/internal/repository/dbrepo"
	"github.com/go-chi/chi"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewRepo creates a new repository for unit testing without a DB connection
// func NewTestRepo(a *config.AppConfig) *Repository {
// 	return &Repository{
// 		App: a,
// 		DB:  dbrepo.NewTestingRepo(a),
// 	}
// }

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	id := m.App.Session.Get(r.Context(), "user_id")
	if id == nil {
		return false
	}
	return true
}

func (m *Repository) UserLogout(w http.ResponseWriter, r *http.Request) {
	m.App.Session.Put(r.Context(), "user_id", nil)
	m.App.Session.Put(r.Context(), "error", "You have been logged out")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Home renders the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Dashboard is the users home page handler
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	if !m.IsLoggedIn(w, r) {
		m.App.Session.Put(r.Context(), "error", "You Must Login to Access This Page")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
	contacts, err := m.DB.GetFavoriteContacts()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["contacts"] = contacts

	render.Template(w, r, "dashboard.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Contacts page handler
func (m *Repository) ContactsAll(w http.ResponseWriter, r *http.Request) {
	// if !m.IsLoggedIn(w, r) {
	// 	m.App.Session.Put(r.Context(), "error", "You Must Login to Access This Page")
	// 	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	// }

	contacts, err := m.DB.GetAllContacts()
	//fmt.Println(contacts.FirstName, "Printed in Contacts Handler")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["contacts"] = contacts

	render.Template(w, r, "contacts.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ContactsNew(w http.ResponseWriter, r *http.Request) {
	// if !m.IsLoggedIn(w, r) {
	// 	m.App.Session.Put(r.Context(), "error", "You Must Login to Access This Page")
	// 	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	// }
	fmt.Print("GPT Version, This page is updatated!")
	var ctID int
	var contact models.Contact
	var err error
	form := forms.New(nil)
	src := chi.URLParam(r, "src")
	id := chi.URLParam(r, "id")
	editMode := id != ""

	if id != "" {
		ctID, err = strconv.Atoi(id) // Convert to int
		if err != nil {
			http.Error(w, "Invalid contact ID", http.StatusBadRequest)
			return
		}
	}

	fmt.Println("Raw Path:", r.URL.Path)
	fmt.Println("Extracted Source:", src)
	fmt.Println("Extracted ID:", id)
	fmt.Println("Edit Mode:", editMode)

	if editMode {
		user_id, ok := m.App.Session.Get(r.Context(), "user_id").(string)
		if !ok || user_id == "" {
			helpers.ServerError(w, errors.New("user ID not found in session"))
			return
		}
		fmt.Println("calling GetContactByID with ctID, user_id", ctID, user_id)
		contact, err = m.DB.GetContactByID(ctID, user_id)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		editMode = true
		fmt.Println("contact.FirstName:", contact.FirstName)
	} else {
		contact = models.Contact{} // This ensures contact is initialized when editMode is false
	}

	stringMap := make(map[string]string)
	data := make(map[string]interface{})
	data["contact"] = contact
	data["editMode"] = editMode
	stringMap["src"] = src
	stringMap["id"] = id

	fmt.Println("Outside of if statement contact.FirstName:", contact.FirstName)
	fmt.Println("Final Contact Data:", contact)
	fmt.Println("Company Data:", contact.Company)

	render.Template(w, r, "contacts-new.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Form:      form,
		Data:      data,
	})
}

// ContactsNewPost
func (m *Repository) ContactsNewPost(w http.ResponseWriter, r *http.Request) {
	// if !m.IsLoggedIn(w, r) {
	// 	m.App.Session.Put(r.Context(), "error", "You Must Login to Complete Thie Action")
	// 	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	// }
	var contact models.Contact

	contact, ok := m.App.Session.Get(r.Context(), "contact").(models.Contact)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get contact data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	favorite := false // Default value is false if the checkbox isn't checked
	if r.FormValue("favorite") == "on" {
		favorite = true
	}
	contact.Favorite = favorite
	contact.FirstName = r.Form.Get("first_name")
	contact.LastName = r.Form.Get("last_name")
	contact.Email = r.Form.Get("email")
	contact.Objective = r.Form.Get("objective")
	contact.MobilePhone = r.Form.Get("mobile_phone")
	contact.WorkPhone = r.Form.Get("work_phone")
	contact.Linkedin = r.Form.Get("linkedin")
	contact.Github = r.Form.Get("github")
	contact.JobTitle = r.Form.Get("job_title")
	contact.Company.CompanyName = r.Form.Get("company")
	contact.Company.URL = r.Form.Get("url")
	contact.Company.Industry = r.Form.Get("industry")
	contact.Company.Size = r.Form.Get("size")
	contact.Description = r.Form.Get("description")
	contact.Notes = r.Form.Get("notes")
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "job_title", "company")
	form.LengthTest("first_name", 3, 100)
	form.IsEmail("email")

	//Validate form data
	if !form.Valid() {
		data := make(map[string]interface{})
		data["contact"] = contact
		http.Error(w, "Invalid Form Data", http.StatusSeeOther)
		render.Template(w, r, "contacts-new.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	//Write new contact to DataBase
	user_id := m.App.Session.Get(r.Context(), "user_id").(string)

	_, err = m.DB.AddNewContact(contact, user_id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert new contact into database")
		http.Redirect(w, r, "/contacts", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Added new contact into database")
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	//render.Template(w, r, "contacts.page.tmpl", &models.TemplateData{})
}

// func (m *Repository) ContactsNewJSON(w http.ResponseWriter, r *http.Request) {
// 	render.Template(w, r, "contacts-new-JSON.page.tmpl", &models.TemplateData{})
// }

func (m *Repository) ContactViewByID(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	if len(exploded) < 5 {
		helpers.ServerError(w, errors.New("invalid URL format"))
		return
	}

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]

	user_id, ok := m.App.Session.Get(r.Context(), "user_id").(string)
	if !ok || user_id == "" {
		helpers.ServerError(w, errors.New("user ID not found in session"))
		return
	}

	contact, err := m.DB.GetContactByID(id, user_id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := map[string]interface{}{
		"contact": contact,
	}
	// store data in session so the user can access it if they chose to edit the contact.
	m.App.Session.Put(r.Context(), "contact", contact)

	render.Template(w, r, "contact-view.page.tmpl", &models.TemplateData{
		StringMap: map[string]string{"src": src},
		Data:      data,
		Form:      forms.New(nil),
	})
}

// ContactViewByIDPost posts changes to notes and description
func (m *Repository) ContactViewByIDPost(w http.ResponseWriter, r *http.Request) {
	// n := r.Form.Get("notes")
	// d := r.Form.Get("description")

	render.Template(w, r, "contacts-view.page.tmpl", &models.TemplateData{})
}

// ContactEditByID handles the contacted edit form
func (m *Repository) ContactEditByID(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts-new.page.tmpl", &models.TemplateData{})
}

// ContactEditByIDPost posts changes to a contact
func (m *Repository) ContactEditByIDPost(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts-new.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ContactsUpdateJSON(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts-update-JSON.page.tmpl", &models.TemplateData{})
}

// Skill-tracker page handler
func (m *Repository) SkillTracker(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "skill-tracker.page.tmpl", &models.TemplateData{})
}

// Availability page handler
func (m *Repository) Applications(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "applications.page.tmpl", &models.TemplateData{})
}

// Availability page handler
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	if m.IsLoggedIn(w, r) {
		m.App.Session.Put(r.Context(), "flash", "You are already logged in")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "user-register.page.tmpl", &models.TemplateData{})
}

// PostShowLogin handles logging user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	form.LengthTest("password", 8, 64)

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	user_id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		log.Println("Login Redirect 1")
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	} else {
		if len(user_id) == 0 {
			log.Println("Login Redirect 2")
			m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
	}
	//log.Println("User Authenticated", user_id)
	m.App.Session.Put(r.Context(), "user_id", user_id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) CertTracker(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "cert-tracker.page.tmpl", &models.TemplateData{})
}
func (m *Repository) JobSearchManager(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "job-search-manager.page.tmpl", &models.TemplateData{})
}
func (m *Repository) ResumeManager(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "resume-manager.page.tmpl", &models.TemplateData{})
}
func (m *Repository) ResourcesJob(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "job-resources.page.tmpl", &models.TemplateData{})
}
func (m *Repository) ResourcesInterview(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "interview-prep.page.tmpl", &models.TemplateData{})
}
func (m *Repository) ResourcesResume(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "resume.page.tmpl", &models.TemplateData{})
}
