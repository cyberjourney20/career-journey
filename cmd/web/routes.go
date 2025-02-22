package main

import (
	"net/http"

	"github.com/cyberjourney20/career-journey/internal/config"
	"github.com/cyberjourney20/career-journey/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Route("/user", func(mux chi.Router) {
		mux.Get("/login", handlers.Repo.Login)
		mux.Post("/login", handlers.Repo.PostShowLogin)
		mux.Get("/logout", handlers.Repo.UserLogout)
		mux.Get("/register", handlers.Repo.Register)
		mux.Post("/register", handlers.Repo.PostRegister)
	})

	mux.Route("/my", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.Dashboard)
		mux.Get("/resumes", handlers.Repo.ViewResumes)
		mux.Get("/applications", handlers.Repo.ViewApplications)
		mux.Get("/profile", handlers.Repo.UserEditProfile)
		mux.Get("/skills", handlers.Repo.SkillTracker)
		mux.Get("/certifications", handlers.Repo.CertTracker)
	})

	mux.Route("/jobs", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.JobListingAll)
		mux.Get("/search-manager", handlers.Repo.JobSearchManager)
		mux.Get("/listings", handlers.Repo.JobListingAll)
		mux.Get("/new", handlers.Repo.JobListingEdit)
		mux.Post("/new", handlers.Repo.JobListingEditPost)
		mux.Get("/view/{id}", handlers.Repo.JobListingViewByID)
		mux.Post("/view/{id}", handlers.Repo.JobListingViewByIDPost)
		mux.Get("/edit/{id}", handlers.Repo.JobListingEdit)
		// mux.Post("/edit/{id}", handlers.Repo.JobListingEditPost)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.AdminDashboard)
	})

	mux.Route("/contacts", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.ContactsAll)
		mux.Get("/new", handlers.Repo.ContactsEdit)
		mux.Post("/new", handlers.Repo.ContactsEditPost)
		mux.Get("/view/{id}", handlers.Repo.ContactViewByID)
		mux.Post("/view/{id}", handlers.Repo.ContactViewByIDPost)
		mux.Get("/edit/{id}", handlers.Repo.ContactsEdit)
		mux.Post("/edit/{id}", handlers.Repo.ContactsEditPost)
		// mux.Get("/edit/{src}/{id}", handlers.Repo.ContactEditByID)
		// mux.Post("/edit/{src}/{id}", handlers.Repo.ContactEditByIDPost)
	})
	mux.Route("/resources", func(mux chi.Router) {
		mux.Get("/", handlers.Repo.ResourcesJob)
		mux.Get("/job-search", handlers.Repo.ResourcesJob)
		mux.Get("/interview-prep", handlers.Repo.ResourcesInterview)
		mux.Get("/resume", handlers.Repo.ResourcesResume)
	})

	webFileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", webFileServer))

	userFileServer := http.FileServer(http.Dir("./users_files/"))
	mux.Handle("/users_files/*", http.StripPrefix("/users_files", userFileServer))

	return mux
}
