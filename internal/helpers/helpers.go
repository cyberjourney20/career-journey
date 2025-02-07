package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/cyberjourney20/career-journey/internal/config"
	"github.com/joho/godotenv"
)

var app *config.AppConfig

// Sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}

func LoadEnv(key string) string {
	err := godotenv.Load(os.ExpandEnv("./.env"))
	if err != nil {
		log.Println("Error loading .env file")
	}

	return os.Getenv(key)
}
