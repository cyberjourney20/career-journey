# Career Journey 
Career Journey is a Golang based web app to help users organize thier job hunt. 

## Build version 0.02
The current build is pre-alpha. It has limited functionality for data interaction beyond prepopulated psudo-data to ensure functionality with the databse. I upgraded teh version to 0.02 because there are som actual function sections now. 

## Workign Features
### Contacts
- Add, Edit, View Contacts
### Job Listings
- View Job Listings (they must be added to DB manually currently)
### Menus and Navigation 
- Navigation Menus are all working. Logout in user menu is working. 
- Nothing else in the top bar is currently functional.  

## Feature Developement
- Add & edit job listings 
- Add and edit applications
- Add and edit resumes and cover letters
- Add and edit skills and certifications
- Upload, Download, and export data
## AI Feature Developement
- Integrate local Deepseek-R1:8b & 14b
- Compare Resume to Job Listings Provide feedback and gap analysis
- Assist in resume revisions and coverletters 

## Requirements
- Uses PostgreSQL: This can be swapped out to a DB of your choice with slight modifications. 
- Database Migrations: Built using Soda [gobuffalo.io](https://gobuffalo.io/documentation/database/migrations/)
## Installation 
- Install Go
- Install PostgreSQL

- Clone the Repository 
```
gh repo clone https://github.com/cyberjourney20/career-journey
```
- Update .env with your local variables
- Install Soda CLI https://gobuffalo.io/documentation/database/soda/
```
go install github.com/gobuffalo/pop/v6/soda@latest
```
- Update database.yml then create the DB
```
soda create -a
soda create -e test
```
- Run the DB migrations. (remove any with test data if you want a fresh install)
```
soda migrate
```
- Build the program
```
go build -o career-journey.exe ./cmd/web 
```
- Run the exe and navigate to localhost:8080



## External Packages
- These are the exteral packages used int he program. 
- [SCS: HTTP Session Management for Go v2.8.0](github.com/alexedwards/scs)
- [chi: HTTP Router v1.5.5](https://github.com/go-chi/chi)
- [Go Validator](https://github.com/asaskevich/govalidator)
- [pgx: PostgreSQL Driver and Toolkit v5.7.2](github.com/jackc/pgx)
- [GoDotEnv: v1.5.1](github.com/joho/godotenv)
- [Nosurf: v1.1.1](github.com/justinas/nosurf)

# Photo Credit: 
- Sebastiaan Stam (plane over city)
- Josh Hild (Ny City with Reflections in rain)
- Nextvoyage (City Night)
s
