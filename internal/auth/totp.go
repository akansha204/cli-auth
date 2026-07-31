package auth

import (
	"github.com/pquerna/otp/totp"
)

const totpIssuer = "cli-auth"

func GenerateTOTPSecret(accountName string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      totpIssuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func ValidateTOTPCode(secret, code string) bool {
	if secret == "" || code == "" {
		return false
	}

	return totp.Validate(code, secret)
}
