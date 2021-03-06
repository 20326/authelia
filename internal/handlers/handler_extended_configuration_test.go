package handlers

import (
	"testing"

	"github.com/authelia/authelia/internal/authorization"
	"github.com/authelia/authelia/internal/configuration/schema"
	"github.com/authelia/authelia/internal/configuration/validator"
	"github.com/authelia/authelia/internal/mocks"

	"github.com/stretchr/testify/suite"
)

type SecondFactorAvailableMethodsFixture struct {
	suite.Suite
	mock *mocks.MockAutheliaCtx
}

func (s *SecondFactorAvailableMethodsFixture) SetupTest() {
	s.mock = mocks.NewMockAutheliaCtx(s.T())
	s.mock.Ctx.Providers.Authorizer = authorization.NewAuthorizer(schema.AccessControlConfiguration{
		DefaultPolicy: "deny",
		Rules:         []schema.ACLRule{},
	})
}

func (s *SecondFactorAvailableMethodsFixture) TearDownTest() {
	s.mock.Close()
}

func (s *SecondFactorAvailableMethodsFixture) TestShouldServeDefaultMethods() {
	s.mock.Ctx.Configuration = schema.Configuration{
		TOTP: &schema.TOTPConfiguration{
			Period: validator.DefaultTOTPPeriod,
		},
	}
	expectedBody := ExtendedConfigurationBody{
		AvailableMethods:    []string{"totp", "u2f"},
		SecondFactorEnabled: false,
		TOTPPeriod:          validator.DefaultTOTPPeriod,
	}
	ExtendedConfigurationGet(s.mock.Ctx)
	s.mock.Assert200OK(s.T(), expectedBody)
}

func (s *SecondFactorAvailableMethodsFixture) TestShouldServeDefaultMethodsAndMobilePush() {
	s.mock.Ctx.Configuration = schema.Configuration{
		DuoAPI: &schema.DuoAPIConfiguration{},
		TOTP: &schema.TOTPConfiguration{
			Period: validator.DefaultTOTPPeriod,
		},
	}
	expectedBody := ExtendedConfigurationBody{
		AvailableMethods:    []string{"totp", "u2f", "mobile_push"},
		SecondFactorEnabled: false,
		TOTPPeriod:          validator.DefaultTOTPPeriod,
	}
	ExtendedConfigurationGet(s.mock.Ctx)
	s.mock.Assert200OK(s.T(), expectedBody)
}

func (s *SecondFactorAvailableMethodsFixture) TestShouldCheckSecondFactorIsDisabledWhenNoRuleIsSetToTwoFactor() {
	s.mock.Ctx.Configuration = schema.Configuration{
		TOTP: &schema.TOTPConfiguration{
			Period: validator.DefaultTOTPPeriod,
		},
	}
	s.mock.Ctx.Providers.Authorizer = authorization.NewAuthorizer(schema.AccessControlConfiguration{
		DefaultPolicy: "bypass",
		Rules: []schema.ACLRule{
			schema.ACLRule{
				Domain: "example.com",
				Policy: "deny",
			},
			schema.ACLRule{
				Domain: "abc.example.com",
				Policy: "single_factor",
			},
			schema.ACLRule{
				Domain: "def.example.com",
				Policy: "bypass",
			},
		},
	})
	ExtendedConfigurationGet(s.mock.Ctx)
	s.mock.Assert200OK(s.T(), ExtendedConfigurationBody{
		AvailableMethods:    []string{"totp", "u2f"},
		SecondFactorEnabled: false,
		TOTPPeriod:          validator.DefaultTOTPPeriod,
	})
}

func (s *SecondFactorAvailableMethodsFixture) TestShouldCheckSecondFactorIsEnabledWhenDefaultPolicySetToTwoFactor() {
	s.mock.Ctx.Configuration = schema.Configuration{
		TOTP: &schema.TOTPConfiguration{
			Period: validator.DefaultTOTPPeriod,
		},
	}
	s.mock.Ctx.Providers.Authorizer = authorization.NewAuthorizer(schema.AccessControlConfiguration{
		DefaultPolicy: "two_factor",
		Rules: []schema.ACLRule{
			schema.ACLRule{
				Domain: "example.com",
				Policy: "deny",
			},
			schema.ACLRule{
				Domain: "abc.example.com",
				Policy: "single_factor",
			},
			schema.ACLRule{
				Domain: "def.example.com",
				Policy: "bypass",
			},
		},
	})
	ExtendedConfigurationGet(s.mock.Ctx)
	s.mock.Assert200OK(s.T(), ExtendedConfigurationBody{
		AvailableMethods:    []string{"totp", "u2f"},
		SecondFactorEnabled: true,
		TOTPPeriod:          validator.DefaultTOTPPeriod,
	})
}

func (s *SecondFactorAvailableMethodsFixture) TestShouldCheckSecondFactorIsEnabledWhenSomePolicySetToTwoFactor() {
	s.mock.Ctx.Configuration = schema.Configuration{
		TOTP: &schema.TOTPConfiguration{
			Period: validator.DefaultTOTPPeriod,
		},
	}
	s.mock.Ctx.Providers.Authorizer = authorization.NewAuthorizer(schema.AccessControlConfiguration{
		DefaultPolicy: "bypass",
		Rules: []schema.ACLRule{
			schema.ACLRule{
				Domain: "example.com",
				Policy: "deny",
			},
			schema.ACLRule{
				Domain: "abc.example.com",
				Policy: "two_factor",
			},
			schema.ACLRule{
				Domain: "def.example.com",
				Policy: "bypass",
			},
		},
	})
	ExtendedConfigurationGet(s.mock.Ctx)
	s.mock.Assert200OK(s.T(), ExtendedConfigurationBody{
		AvailableMethods:    []string{"totp", "u2f"},
		SecondFactorEnabled: true,
		TOTPPeriod:          validator.DefaultTOTPPeriod,
	})
}

func TestRunSuite(t *testing.T) {
	s := new(SecondFactorAvailableMethodsFixture)
	suite.Run(t, s)
}
