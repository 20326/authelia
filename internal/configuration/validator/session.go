package validator

import (
	"errors"
	"fmt"
	"github.com/authelia/authelia/internal/configuration/schema"
	"github.com/authelia/authelia/internal/utils"
)

// ValidateSession validates and update session configuration.
func ValidateSession(configuration *schema.SessionConfiguration, validator *schema.StructValidator) {
	if configuration.Name == "" {
		configuration.Name = schema.DefaultSessionConfiguration.Name
	}

	if configuration.Redis != nil && configuration.Secret == "" {
		validator.Push(errors.New("Set secret of the session object"))
	}

	// TODO(james-d-elliott): Convert to duration notation
	if configuration.Expiration == 0 {
		configuration.Expiration = schema.DefaultSessionConfiguration.Expiration // 1 hour
	} else if configuration.Expiration < 1 {
		validator.Push(errors.New("Set expiration of the session above 0"))
	}

	// TODO(james-d-elliott): Convert to duration notation
	if configuration.Inactivity < 0 {
		validator.Push(errors.New("Set inactivity of the session to 0 or above"))
	}

	if configuration.RememberMeDuration == "" {
		configuration.RememberMeDuration = schema.DefaultSessionConfiguration.RememberMeDuration
	} else {
		if _, err := utils.ParseDurationString(configuration.RememberMeDuration); err != nil {
			validator.Push(errors.New(fmt.Sprintf("Error occurred parsing remember_me_duration string: %s", err)))
		}
	}

	if configuration.Domain == "" {
		validator.Push(errors.New("Set domain of the session object"))
	}
}
