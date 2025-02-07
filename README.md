# Career Journey 
Career Journey is a Golang based web app to help users organize thier job hunt. 

## Build version 0.01
The current build is pre-alpha. It has no actual functionality for data interaction beyond prepopulated psudo-data to ensure functionality with the databse. 

## Requirements
- Uses PostgreSQL: This can be swapped out to a DB of your choice with slight modifications. 
- Database Migrations: Built using Soda [gobuffalo.io](https://gobuffalo.io/documentation/database/migrations/)
## Installation 
- Install Go
- Install PostgreSQL
- Install Ollama
- Clone the Repository 
```
gh repo clone https://github.com/cyberjourney20/career-journey
```
- Update .env with your local variables
- Build the program
```
go build -o career-journey.exe ./cmd/web 
```
- Run the exe and navigate to localhost:8080

## Feature Developement
- Currently working to integrate local Deepseek-R1:8B for various tasks.

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

