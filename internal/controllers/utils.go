package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"app/internal/templates"
)

func isHtmxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func wrapComponentTemplate(component templ.Component, r *http.Request) templ.Component {
	if !isHtmxRequest(r) {
		return templates.Layout(component)
	}
	return component
}
