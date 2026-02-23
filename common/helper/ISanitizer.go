package helper

import "html/template"

// ISanitizer adalah interface strategy
type ISanitizer interface {
	Sanitize(input string) template.HTML
}
