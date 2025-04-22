package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// Handles all requests coming to the '/' (root) endpoint.
type HomeHandler struct {
}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// Renders 'home.html' when '/' (root) endpoint gets requested (commonly from MC)
func (hh *HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {

	// For root path we return the home page
	path := "templates/home.html"
	tmpl, err := template.ParseFS(templateFS, path)
	if err != nil {
		log.Println("Error parsing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
