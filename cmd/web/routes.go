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

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login-post", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)
	mux.Get("/user/Register", handlers.Repo.Register)

	mux.Get("/dashboard", handlers.Repo.Dashboard)
	mux.Get("/skill-tracker", handlers.Repo.SkillTracker)
	mux.Get("/cert-tracker", handlers.Repo.CertTracker)
	mux.Get("/application-tracker", handlers.Repo.SkillTracker)
	mux.Get("/job-search-manager", handlers.Repo.JobSearchManager)
	mux.Get("/resume-manager", handlers.Repo.ResumeManager)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
	})
	mux.Route("/contacts", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.ContactsAll)
		mux.Get("/new/{src}", handlers.Repo.ContactsNew)
		mux.Post("/new", handlers.Repo.ContactsNewPost)
		mux.Get("/view/{src}/{id}", handlers.Repo.ContactViewByID)
		mux.Post("/view/{src}/{id}", handlers.Repo.ContactViewByIDPost)
		mux.Get("/edit/{src}/{id}", handlers.Repo.ContactsNew)
		mux.Post("/edit/{src}/{id}", handlers.Repo.ContactsNewPost)
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
