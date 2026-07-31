package auth

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func currentCode(t *testing.T, secret string) string {
	t.Helper()

	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return code
}
