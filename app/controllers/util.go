package controllers

func CleanSlashes(in string) string {
	return SlashRegexp.ReplaceAllString(in, "/")
}
