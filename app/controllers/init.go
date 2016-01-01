package controllers

import (
	"regexp"
)

var SlashRegexp *regexp.Regexp

func init() {
	SlashRegexp = regexp.MustCompile("/{2,}")
}
