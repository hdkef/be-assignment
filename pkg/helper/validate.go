package helper

import "regexp"

func IsValidEmail(email string) bool {
	// A simple regex for validating email addresses
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
