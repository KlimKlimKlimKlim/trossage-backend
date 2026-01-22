package validate

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

var (
	loginRegex       = regexp.MustCompile(`^[a-z0-9_]{3,20}$`)
	passwordRegex    = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{}|;:,.<>?]{8,63}$`)
	displayNameRegex = regexp.MustCompile(`^[\p{L}\p{N}\s._-]{1,20}$`)
)

func Login(login string) bool {
	return loginRegex.MatchString(login)
}

func Password(password string) bool {
	if !passwordRegex.MatchString(password) {
		return false
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}

		if unicode.IsDigit(char) {
			hasDigit = true
		}

		if hasLetter && hasDigit {
			return true
		}
	}

	return false
}

func DisplayName(displayName string) bool {
	runeCount := utf8.RuneCountInString(displayName)
	if runeCount == 0 || runeCount > 20 {
		return false
	}

	return displayNameRegex.MatchString(displayName)
}
