package rules

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func IsArch() validation.Rule {
	return &archRule{}
}

type archRule struct{}

// Validate implements validation.Rule.
func (r *archRule) Validate(value any) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	str, err := validation.EnsureString(value)
	if err != nil {
		return err
	}

	return validation.In(
		"386",
		"amd64",
		"arm",
		"arm64",
		"mips",
		"mipsle",
		"mips64",
		"mips64le",
		"ppc64",
		"ppc64le",
	).Error(
		fmt.Sprintf(
			"Must be one of the following values: %s",
			strings.Join([]string{
				"386",
				"amd64",
				"arm",
				"arm64",
				"mips",
				"mipsle",
				"mips64",
				"mips64le",
				"ppc64",
				"ppc64le",
			}, ", "),
		),
	).Validate(str)
}
