package controllers

import (
	"github.com/revel/revel"
)

func CleanSlashes(in string) string {
	return SlashRegexp.ReplaceAllString(in, "/")
}

func GitBasePath() string {
	return revel.Config.StringDefault("git.baseDir", "/")
}
