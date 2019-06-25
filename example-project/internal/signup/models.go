package signup

import (
	"geeks-accelerator/oss/saas-starter-kit/example-project/internal/account"
	"geeks-accelerator/oss/saas-starter-kit/example-project/internal/user"
)

// SignupRequest contains information needed perform signup.
type SignupRequest struct {
	Account struct {
		Name     string  `json:"name" validate:"required,unique" example:"Company {RANDOM_UUID}"`
		Address1 string  `json:"address1" validate:"required" example:"221 Tatitlek Ave"`
		Address2 string  `json:"address2" validate:"omitempty" example:"Box #1832"`
		City     string  `json:"city" validate:"required" example:"Valdez"`
		Region   string  `json:"region" validate:"required" example:"AK"`
		Country  string  `json:"country" validate:"required" example:"USA"`
		Zipcode  string  `json:"zipcode" validate:"required" example:"99686"`
		Timezone *string `json:"timezone" validate:"omitempty" example:"America/Anchorage"`
	} `json:"account" validate:"required"` // Account details.
	User struct {
		Name            string `json:"name" validate:"required" example:"Gabi May"`
		Email           string `json:"email" validate:"required,email,unique" example:"{RANDOM_EMAIL}"`
		Password        string `json:"password" validate:"required" example:"SecretString"`
		PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password" example:"SecretString"`
	} `json:"user" validate:"required"` // User details.
}

// SignupResponse contains information needed perform signup.
type SignupResponse struct {
	Account *account.Account `json:"account"`
	User    *user.User       `json:"user"`
}