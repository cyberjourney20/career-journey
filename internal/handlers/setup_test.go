package handlers

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/cyberjourney20/career-journey/internal/config"
)

var pathToTemplates = "./../../templates"
var app config.AppConfig
var session *scs.SessionManager
var functions = template.FuncMap{}
