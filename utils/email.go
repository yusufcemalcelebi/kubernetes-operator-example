package utils

import "strings"

// IsValidEmail checks if the provided email address is in a valid format
func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
