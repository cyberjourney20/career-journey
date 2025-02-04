package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cyberjourney20/career-journey/internal/config"
	"github.com/cyberjourney20/career-journey/internal/driver"
	"github.com/cyberjourney20/career-journey/internal/handlers"
	"github.com/cyberjourney20/career-journey/internal/helpers"
	"github.com/cyberjourney20/career-journey/internal/models"
	"github.com/cyberjourney20/career-journey/internal/render"
	"github.com/joho/godotenv"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	//Close db connectiion when program closes
	defer db.SQL.Close()

	fmt.Printf("Starting application on port: %s\n", portNumber)

	// Starts HTTP Server
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// What goes in the session
	gob.Register(models.User{})
	gob.Register(models.Contact{})
	gob.Register(models.Application{})
	gob.Register(models.JobListing{})
	gob.Register(models.Company{})
	gob.Register(models.Location{})

	//change to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//godotenv to load secrets from .env file
	err := godotenv.Load(os.ExpandEnv("./.env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbString := os.Getenv("DBSTRING")

	//connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(dbString)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect to database! Dying... %v\n", err))
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
