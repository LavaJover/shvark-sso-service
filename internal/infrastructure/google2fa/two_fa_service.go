package google2fa

import "github.com/pquerna/otp/totp"


func Generate2FASecret(login string) (secret string, qrURL string,err  error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer: "ShvarkPay",
		AccountName: login,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}