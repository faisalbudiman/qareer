package locations

import (
	"qareer/pkg/validator"
	"time"
)

var (
	stringValidator = validator.NewStringValidator()
)

type Location struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (l Location) ValidateInsert() error {
	return validator.MergeError(
		validator.ValidateString(
			"name",
			l.Name,
			stringValidator.Required,
			stringValidator.StringMinLength(2),
		),
	)
}

func validateFilterName(s string) error {
	return validator.ValidateString(
		"name",
		s,
		stringValidator.StringMinLength(2),
	)
}
