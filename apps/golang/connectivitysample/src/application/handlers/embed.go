package handlers

import "embed"

//go:embed templates/*.html
var templateFS embed.FS

//templateFS is a package variable that is used to embed html template files.
