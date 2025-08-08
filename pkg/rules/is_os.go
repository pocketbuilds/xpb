package rules

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func IsOs() validation.Rule {
	return &osRule{}
}

type osRule struct{}

// Validate implements validation.Rule.
func (r *osRule) Validate(value any) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	str, err := validation.EnsureString(value)
	if err != nil {
		return err
	}

	return validation.In(
		"android",
		"darwin",
		"dragonfly",
		"freebsd",
		"linux",
		"netbsd",
		"openbsd",
		"plan9",
		"solaris",
		"windows",
	).Error(
		fmt.Sprintf(
			"Must be one of the following values: %s",
			strings.Join([]string{
				"android",
				"darwin",
				"dragonfly",
				"freebsd",
				"linux",
				"netbsd",
				"openbsd",
				"plan9",
				"solaris",
				"windows",
			}, ", "),
		),
	).Validate(str)
}
