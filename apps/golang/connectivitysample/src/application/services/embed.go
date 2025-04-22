package services

import "embed"

//go:embed templates/*.json
var templateFS embed.FS

//templateFS is a package variable that is used to embed json template files.
