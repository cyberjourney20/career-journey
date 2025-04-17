package handlers

import (
	"encoding/json"
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
	"github.com/cyberjourney20/career-journey/internal/repository/llmrepo"
	"github.com/go-chi/chi"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
	LLM repository.LLMRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
		LLM: llmrepo.NewOllamaRepo(a),
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

// Availability page handler
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	if m.IsLoggedIn(r) {
		m.App.Session.Put(r.Context(), "flash", "You are already logged in")
		http.Redirect(w, r, "/my/dashboard", http.StatusSeeOther)
	}
	render.Template(w, r, "user-login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// PostShowLogin handles the user login post request
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println("error parsing form", err)
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
		log.Println("error in authentication", err)
		// log.Println("Login Redirect 1")
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	} else {
		if len(user_id) == 0 {
			// log.Println("Login Redirect 2")
			m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}
	}
	//log.Println("User Authenticated", user_id)
	m.App.Session.Put(r.Context(), "user_id", user_id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/my/dashboard", http.StatusSeeOther)
}

// IsLoggedIn checks for a user_id to validate user is logged in.
func (m *Repository) IsLoggedIn(r *http.Request) bool {
	return m.App.Session.Exists(r.Context(), "user_id")
}

// Logout logs a user out
func (m *Repository) UserLogout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Register handles the Registration page
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {

	if m.IsLoggedIn(r) {
		m.App.Session.Put(r.Context(), "flash", "You are already logged in")
		http.Redirect(w, r, "/my/dashboard", http.StatusSeeOther)
	}

	render.Template(w, r, "user-register.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// PostRegister handles the user registration post request
func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {

	// log.Println("PostRegister Running")
	var err error
	var exists bool

	newUser := models.User{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Password:  r.Form.Get("password"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "password", "password_2")
	form.LengthTest("first_name", 3, 64)
	form.LengthTest("last_name", 3, 64)
	form.LengthTest("password", 12, 64)
	form.PasswordsMatch(form.Get("password"), form.Get("password_2"))
	form.IsEmail("email")

	exists, err = m.DB.UserExists(newUser.Email)
	if err != nil {
		log.Println("Error checking user existence:", err)
		m.App.Session.Put(r.Context(), "warning", "Database error, please try again.")
		render.Template(w, r, "user-register.page.tmpl", &models.TemplateData{
			Form: form,
			Data: map[string]interface{}{
				"newUser": newUser,
			},
		})
		return
	}

	if exists {
		form.Errors.Add("email", "This username is unavailable")
		render.Template(w, r, "user-register.page.tmpl", &models.TemplateData{
			Form: form,
			Data: map[string]interface{}{
				"newUser": newUser,
			},
		})
		return
	}

	if !form.Valid() {
		// Re-render registration page with form errors and user data
		log.Println("form != Valid")
		render.Template(w, r, "user-register.page.tmpl", &models.TemplateData{
			Form: form,
			Data: map[string]interface{}{
				"newUser": newUser, // Retain user inputs
			},
		})
		return
	}

	err = m.DB.AddNewUser(newUser)
	if err != nil {
		log.Println("error adding new user to database", err)
		m.App.Session.Put(r.Context(), "warning", "error adding new user to database")
		http.Redirect(w, r, "/user/register", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Your account was created successfully, Please login.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Home renders the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// ResourcesJob handels the job resources page
func (m *Repository) ResourcesJob(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "resources-job.page.tmpl", &models.TemplateData{})
}

// ResourcesInterview handels the interview resources page
func (m *Repository) ResourcesInterview(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "resources-interview.page.tmpl", &models.TemplateData{})
}

// ResourcesResume handels the resume resources page
func (m *Repository) ResourcesResume(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "resources-resume.page.tmpl", &models.TemplateData{})
}

// Dashboard is the users home page handler
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	if !m.IsLoggedIn(r) {
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

// ContactsAll handles the Contacs page
func (m *Repository) ContactsAll(w http.ResponseWriter, r *http.Request) {
	fmt.Print("ContactAll Running")
	contacts, err := m.DB.GetAllContacts()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["contacts"] = contacts

	render.Template(w, r, "contacts-all.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// ContactsEdit handles the new and edit contact page
func (m *Repository) ContactsEdit(w http.ResponseWriter, r *http.Request) {

	//fmt.Print("GPT Version, This page is updatated!")
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

	if editMode {
		user_id, ok := m.App.Session.Get(r.Context(), "user_id").(string)
		if !ok || user_id == "" {
			helpers.ServerError(w, errors.New("user ID not found in session"))
			return
		}
		// fmt.Println("calling GetContactByID with ctID, user_id", ctID, user_id)
		contact, err = m.DB.GetContactByID(ctID, user_id)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		editMode = true
		// fmt.Println("contact.FirstName:", contact.FirstName)
	} else {
		contact = models.Contact{}
	}

	stringMap := make(map[string]string)
	data := make(map[string]interface{})
	data["contact"] = contact
	data["editMode"] = editMode
	stringMap["src"] = src
	stringMap["id"] = id

	// fmt.Println("Outside of if statement contact.FirstName:", contact.FirstName)
	// fmt.Println("Final Contact Data:", contact)
	// fmt.Println("Company Data:", contact.Company)

	render.Template(w, r, "contacts-new.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Form:      form,
		Data:      data,
	})
}

// ContactsEditPost handels the new and edit contact post data
func (m *Repository) ContactsEditPost(w http.ResponseWriter, r *http.Request) {

	editMode := r.FormValue("edit_mode") == "true"
	var contact models.Contact
	if editMode {
		var ok bool
		contact, ok = m.App.Session.Get(r.Context(), "contact").(models.Contact)
		if !ok {
			m.App.Session.Put(r.Context(), "error", "can't get contact data from session")
			http.Redirect(w, r, "/my/dashboard", http.StatusTemporaryRedirect)
			return
		}
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/my/dashboard", http.StatusTemporaryRedirect)
		return
	}
	favorite := false // Default value is false if the checkbox isn't checked

	if r.FormValue("favorite") == "true" {
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
		m.App.Session.Put(r.Context(), "error", "Invalid form data")
		render.Template(w, r, "contacts-new.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	//Write new contact to DataBase
	user_id, ok := m.App.Session.Get(r.Context(), "user_id").(string)
	if !ok || user_id == "" {
		helpers.ServerError(w, errors.New("user ID not found in session"))
		return
	}

	if r.FormValue("edit_mode") == "true" {
		err = m.DB.UpdateContactByID(contact)
		if err != nil {
			fmt.Println("error returned from UpdateContactByID: ", err)
			m.App.Session.Put(r.Context(), "error", "can't update contact in database")
			http.Redirect(w, r, "/contacts", http.StatusTemporaryRedirect)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "Updatd contact in database")
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)

	} else {
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
}

// ContactViewByID handels data for one contact by thier ID
func (m *Repository) ContactViewByID(w http.ResponseWriter, r *http.Request) {

	var err error
	var id int

	cid := chi.URLParam(r, "id")
	if cid != "" {
		id, err = strconv.Atoi(cid) // Convert to int
		if err != nil {
			fmt.Println("strconv.Atoi error: ", err)
			http.Error(w, "Invalid job listing  ID", http.StatusBadRequest)
			return
		}
	}

	m.App.Session.Put(r.Context(), "return_to", r.Referer())
	returnPath := m.App.Session.PopString(r.Context(), "return_to")
	if returnPath == "" {
		returnPath = "/my/dashboard"
	}

	fmt.Println("Return Path: ", returnPath)

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
		StringMap: map[string]string{"returnPath": returnPath},
		Data:      data,
		Form:      forms.New(nil),
	})
}

// ContactVieByIDPost handels post data for ContactViewByID data
func (m *Repository) ContactViewByIDPost(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "#", &models.TemplateData{})
}

// JobSearchManager handels the Job Search Manager dashboard
func (m *Repository) JobSearchManager(w http.ResponseWriter, r *http.Request) {

	user_id, ok := m.App.Session.Get(r.Context(), "user_id").(string)
	if !ok || user_id == "" {
		helpers.ServerError(w, errors.New("user ID not found in session"))
		return
	}

	//fmt.Println("running JobSearchManager, uuid: ", id)
	listings, err := m.DB.GetAllJobListing(user_id)
	if err != nil {
		log.Println("error getting job listings from database", err)
		m.App.Session.Put(r.Context(), "warning", "error getting job listings from the database")
		http.Redirect(w, r, "/my/dashboard", http.StatusSeeOther)

	}

	data := make(map[string]interface{})
	data["listings"] = listings

	render.Template(w, r, "job-search-manager.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// UserEditProfile handels the user profile page
func (m *Repository) UserEditProfile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "my-profile.page.tmpl", &models.TemplateData{})
}

// ViewResumes handels the ViewResumes page
func (m *Repository) ViewResumes(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "my-resumes.page.tmpl", &models.TemplateData{})
}

// ViewApplications handels the ViewApplications page
func (m *Repository) ViewApplications(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "my-applications.page.tmpl", &models.TemplateData{})
}

// SkillTracker User skills tracker page handler
func (m *Repository) SkillTracker(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "my-skills.page.tmpl", &models.TemplateData{})
}

// CertTracker User certifications tracker page handler
func (m *Repository) CertTracker(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "my-certifications.page.tmpl", &models.TemplateData{})
}

// JobListingAll handels the job listings page
func (m *Repository) JobListingAll(w http.ResponseWriter, r *http.Request) {

	user_id, ok := m.App.Session.Get(r.Context(), "user_id").(string)
	if !ok || user_id == "" {
		helpers.ServerError(w, errors.New("user ID not found in session"))
		return
	}

	listings, err := m.DB.GetAllJobListing(user_id)
	if err != nil {
		log.Println("error getting job listings from database", err)
		m.App.Session.Put(r.Context(), "warning", "error getting job listings from the database")
		http.Redirect(w, r, "/my/dashboard", http.StatusSeeOther)

	}

	data := make(map[string]interface{})
	data["listings"] = listings

	render.Template(w, r, "job-listing-all.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// JobListingEdit handels the new and edit job listings pages
func (m *Repository) JobListingEdit(w http.ResponseWriter, r *http.Request) {

	var listID int
	var err error
	var listing models.JobListing
	form := forms.New(nil)
	id := chi.URLParam(r, "id")
	editMode := id != ""

	m.App.Session.Put(r.Context(), "return_to", r.Referer()) // need to handle bad return addressess like llm or just stor the first one?

	returnPath := m.App.Session.PopString(r.Context(), "return_to")
	if returnPath == "" {
		returnPath = "/my/dashboard"
	}

	if id != "" {
		listID, err = strconv.Atoi(id) // Convert to int
		if err != nil {
			http.Error(w, "Invalid contact ID", http.StatusBadRequest)
			return
		}
	}

	if editMode {
		sessionListing, ok := m.App.Session.Get(r.Context(), "listing").(models.JobListing)
		if !ok {
			listing = sessionListing
		} else {
			user_id, ok := m.App.Session.Get(r.Context(), "user_id").(string)
			if !ok {
				helpers.ServerError(w, errors.New("user ID not found in session"))
				return
			}
			listing, err = m.DB.GetJobListingByID(listID, user_id)
			if err != nil {
				log.Println("Error retrieving job listing:", err)
				http.Redirect(w, r, "/jobs", http.StatusSeeOther)
				return
			}
		}
	} else {
		listing = models.JobListing{}
	}

	stringMap := make(map[string]string)
	data := make(map[string]interface{})
	data["listing"] = listing
	data["editMode"] = editMode
	stringMap["id"] = id
	// set edit mode
	form.Required("company", "job_title")

	render.Template(w, r, "job-listing-edit.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Form:      form,
		Data:      data,
	})
}

// JobListingEditPost handels the new and edit job listings post data
func (m *Repository) JobListingEditPost(w http.ResponseWriter, r *http.Request) {
	//check if edit or new
	//call add func
	//call update func

	//return message
	//redirect

	render.Template(w, r, "job-listing-view.page.tmpl", &models.TemplateData{})
}

// JobListingViewByID displays one joblisting by ID
func (m *Repository) JobListingViewByID(w http.ResponseWriter, r *http.Request) {

	fmt.Println("JobListingViewByID Running...")
	var listID int
	var err error
	var listing models.JobListing
	lid := chi.URLParam(r, "id")
	userID, ok := m.App.Session.Get(r.Context(), "user_id").(string)
	if !ok || userID == "" {
		helpers.ServerError(w, errors.New("user ID not found in session"))
		return
	}
	m.App.Session.Put(r.Context(), "return_to", r.Referer())

	returnPath := m.App.Session.PopString(r.Context(), "return_to")
	if returnPath == "" {
		returnPath = "/my/dashboard"
	}
	//http.Redirect(w, r, returnPath, http.StatusSeeOther)
	log.Println("returnPath: ", returnPath)
	log.Println("Listing ID: ", lid)

	if lid != "" {
		listID, err = strconv.Atoi(lid) // Convert to int
		if err != nil {
			fmt.Println("strconv.Atoi error: ", err)
			http.Error(w, "Invalid job listing  ID", http.StatusBadRequest)
			return
		}
	}

	listing, err = m.DB.GetJobListingByID(listID, userID)
	if err != nil {
		log.Println("Error in GetJobListingByID: ", err)
		m.App.Session.Put(r.Context(), "error", "error getting job listing from DB")
		http.Redirect(w, r, returnPath, http.StatusSeeOther)
	}

	stringMap := make(map[string]string)
	data := make(map[string]interface{})
	data["listing"] = listing
	stringMap["returnPath"] = returnPath
	stringMap["id"] = lid

	render.Template(w, r, "job-listing-view.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}

// JobListingViewByIDPOST displays one joblisting by ID post data
func (m *Repository) JobListingViewByIDPost(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "job-listing-view.page.tmpl", &models.TemplateData{})
}

// AdminDashboard handels the admin dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "job-listing-edit.page.tmpl", &models.TemplateData{})
}

func (m *Repository) JobListingLLM(w http.ResponseWriter, r *http.Request) {
	fmt.Println("JobListingLLM Handler Started")

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	jobDescription := r.FormValue("paste_description")
	//fmt.Println("Received Job Description:", jobDescription)

	if jobDescription == "" {
		fmt.Println("Error: Job description is empty!")
		http.Error(w, "Job description is required", http.StatusBadRequest)
		return
	}

	systemPrompt := m.LLM.JobListingPrompt(jobDescription)

	// ---- Use the exact same working function from the standalone test ----
	response, err := m.LLM.OllamaGenerateResponse(systemPrompt, false)
	if err != nil {
		fmt.Println("Error calling Ollama:", err)
		http.Error(w, "Error processing AI request", http.StatusInternalServerError)
		return
	}

	jsonStart := strings.Index(response, "{")
	if jsonStart == -1 {
		fmt.Println("Error: No JSON found in LLM response")
		fmt.Println("Raw LLM Response:", response)
		return
	}

	// Extract only the JSON part

	cleanResponse := response[jsonStart:]
	trimResponse := strings.Trim(cleanResponse, "`")

	fmt.Println("Clean JSON Response:", trimResponse) // Debugging

	// Now try parsing only the valid JSON
	var newListing models.JobListing
	if err := json.Unmarshal([]byte(trimResponse), &newListing); err != nil {
		fmt.Println("Error while decoding the data:", err.Error())
		return
	}
	m.App.Session.Put(r.Context(), "listing", newListing)
	m.App.Session.Put(r.Context(), "flash", "Job Listing Populated. Please validate the data and select submit")
	http.Redirect(w, r, "/jobs/edit/0", http.StatusSeeOther)
}

// func (m *Repository) JobListingLLM(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Running JobListingLLM")
// 	ollama := driver.NewOllamaDriverNoStream()

// 	err := r.ParseForm()
// 	if err != nil {
// 		fmt.Println("Error parsing form:", err)
// 		http.Error(w, "Error parsing form", http.StatusBadRequest)
// 		return
// 	}

// 	jobDescription := r.Form.Get("paste_description")
// 	fmt.Println("Received Job Description:", jobDescription)

// 	if jobDescription == "" {
// 		fmt.Println("Error: Job description is empty!")
// 		http.Error(w, "Job description is required", http.StatusBadRequest)
// 		return
// 	}

// 	prompt := m.LLM.JobListingPrompt(jobDescription)

// 	response, err := ollama.OllamaGenerateResponse(prompt)
// 	if err != nil {
// 		log.Println("Error calling LLM:", err)
// 		http.Error(w, "Error processing AI request", http.StatusInternalServerError)
// 	} else {
// 		log.Println("LLM Response:", response)
// 	}

// jsonData, err := json.MarshalIndent(response, "", "  ")
// if err != nil {
// 	fmt.Println("Error formatting JSON:", err)
// 	return
// }
// fmt.Println(string(jsonData))

//render.Template(w, r, "job-listing-edit.page.tmpl", &models.TemplateData{})
//}

// LLM Usage

// ollama := NewOllamaDriver()
// prompt := fmt.Sprintf("%s\n\nUser Input:\n%s", llm.JobListingPrompt(jobDescription), userPrompt)
// response, err := ollama.GenerateResponse(prompt)
// if err != nil {
//     log.Println("Error calling LLM:", err)
// } else {
//     log.Println("LLM Response:", response)
// }
