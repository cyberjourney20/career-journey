package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cyberjourney20/career-journey/internal/config"
	"github.com/cyberjourney20/career-journey/internal/driver"
	"github.com/cyberjourney20/career-journey/internal/forms"
	"github.com/cyberjourney20/career-journey/internal/helpers"
	"github.com/cyberjourney20/career-journey/internal/models"
	"github.com/cyberjourney20/career-journey/internal/render"
	"github.com/cyberjourney20/career-journey/internal/repository"
	"github.com/cyberjourney20/career-journey/internal/repository/dbrepo"
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

// Home is the home page handler
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	if !m.IsLoggedIn(w, r) {
		m.App.Session.Put(r.Context(), "error", "You Must Login to Access This Page")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
	contacts, err := m.DB.GetFavoriteContacts()
	//fmt.Println(contacts.FirstName, "Printed in Contacts Handler")
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
func (m *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	if !m.IsLoggedIn(w, r) {
		m.App.Session.Put(r.Context(), "error", "You Must Login to Access This Page")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
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
	if !m.IsLoggedIn(w, r) {
		m.App.Session.Put(r.Context(), "error", "You Must Login to Access This Page")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}
	contact := models.Contact{}
	m.App.Session.Put(r.Context(), "contact", contact)
	data := map[string]interface{}{
		"contact": contact,
	}

	render.Template(w, r, "contacts-new.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) ContactsNewPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("func ContactsNewPost is running")
	if !m.IsLoggedIn(w, r) {
		// fmt.Println("func ContactsNewPost IsLoggedIn is running")
		m.App.Session.Put(r.Context(), "error", "You Must Login to Complete Thie Action")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	contact, ok := m.App.Session.Get(r.Context(), "contact").(models.Contact)
	// fmt.Println("func ContactsNewPost Get Context is running")
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get contact data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	// fmt.Printf("Content-Length: %d\n", r.ContentLength)
	// fmt.Printf("Headers: %+v\n", r.Header)
	err := r.ParseForm()
	// fmt.Println("func ContactsNewPost parse form is running")
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	contact.FirstName = r.Form.Get("first_name")
	contact.LastName = r.Form.Get("last_name")
	contact.JobTitle = r.Form.Get("job_title")
	contact.Email = r.Form.Get("email")
	contact.Company.CompanyName = r.Form.Get("company")
	contact.Objective = r.Form.Get("objective")
	contact.MobilePhone = r.Form.Get("mobile_phone")
	contact.WorkPhone = r.Form.Get("work_phone")
	contact.Linkedin = r.Form.Get("linkedin")
	contact.Github = r.Form.Get("github")
	contact.Description = r.Form.Get("description")

	if r.FormValue("favorite") == "true" {
		contact.Favorite = true
	}

	// contact.Notes = r.Form.Get("notes")
	// if r.Form.Get("favorite") == "favorite" {
	// 	contact.Favorite = true
	// } else {
	// 	contact.Favorite = false
	// }

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
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

func (m *Repository) ContactsNewJSON(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts-new-JSON.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ContactsView(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contacts-view.page.tmpl", &models.TemplateData{})
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
