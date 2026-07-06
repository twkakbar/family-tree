// Package validation registers project-specific rules into Gin's validator
// engine. Register must be called once at startup (cmd/api).
package validation

import (
	"fmt"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Register installs all custom validation tags.
func Register() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("unexpected validator engine %T", binding.Validator.Engine())
	}
	if err := v.RegisterValidation("strong_password", strongPassword); err != nil {
		return fmt.Errorf("registering strong_password: %w", err)
	}
	return nil
}

// strongPassword enforces the project password policy (DEVELOPMENT_RULES §7):
// at least 8 characters, containing at least one letter and one digit.
// Deliberately no special-character requirement — length beats complexity
// theater, and the policy must not push users into predictable patterns.
func strongPassword(fl validator.FieldLevel) bool {
	pw := fl.Field().String()
	if len(pw) < 8 || len(pw) > 128 {
		return false
	}
	var hasLetter, hasDigit bool
	for _, r := range pw {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
