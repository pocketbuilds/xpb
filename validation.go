package xpb

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
)

type PreValidator interface {
	PreValidate(app core.App) error
}

func PreValidatePlugins(app core.App) error {
	var errs []error

	for _, plugin := range plugins {
		var err error
		prevalidator, ok := plugin.(PreValidator)
		if !ok {
			continue
		}
		func() {
			defer func() {
				switch r := recover().(type) {
				case error:
					err = r
				case nil:
					return
				default:
					err = fmt.Errorf("%v", r)
				}
			}()
			err = prevalidator.PreValidate(app)
		}()
		if err != nil {
			errs = append(errs, fmt.Errorf(
				"error from plugin %s: %w",
				plugin.Name(),
				err,
			))
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("error(s) during load: %w", errors.Join(errs...))
	}

	return nil
}

func ValidatePlugins(app core.App) error {
	var errs []error

	for _, plugin := range plugins {
		var err error
		validator, ok := plugin.(validation.Validatable)
		if !ok {
			continue
		}
		func() {
			defer func() {
				switch r := recover().(type) {
				case error:
					err = r
				case nil:
					return
				default:
					err = fmt.Errorf("%v", r)
				}
			}()
			err = validator.Validate()
		}()
		if err != nil {
			errs = append(errs, fmt.Errorf(
				"error from plugin %s: %w",
				plugin.Name(),
				err,
			))
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("error(s) during load: %w", errors.Join(errs...))
	}

	return nil
}
