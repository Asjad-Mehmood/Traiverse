package main

import (
	_ "embed"
	"encoding/base64"
	"html"
	"strings"
)

// The start screen: shown in the window the moment the program opens, until the app itself
// is ready to show, so opening it is never a wait on an empty window.
//
//go:embed splash.html
var splashTemplate string

//go:embed icon.ico
var splashIcon []byte

// splashPage is the start screen for the program called name
func splashPage(name string) string {
	icon := "data:image/x-icon;base64," + base64.StdEncoding.EncodeToString(splashIcon)
	return strings.NewReplacer("{{NAME}}", html.EscapeString(name), "{{ICON}}", icon).Replace(splashTemplate)
}
