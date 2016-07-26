package main

import (
	"flag"
	"fmt"

	"github.com/v-yarotsky/ghprj/github"
)

func main() {
	var username, password, otp string
	flag.StringVar(&username, "username", "", "GitHub username")
	flag.StringVar(&password, "password", "", "GitHub password")
	flag.StringVar(&otp, "otp", "", "2FA OTP")
	flag.Parse()

	github.NewAuthenticator(func(c *github.Credentials, twoFactor bool) error {
		c.Username = username
		c.Password = password

		if twoFactor {
			c.TwoFactorToken = otp
		}
		return nil
	}, true).AccessToken()

	fmt.Println("Login succeeded")
}
