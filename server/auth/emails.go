package auth

import (
	"fmt"
	"net/mail"

	"github.com/soprasteria/docktor/model/types"
	"github.com/soprasteria/docktor/server/email"
)

// SendWelcomeEmail sends a welcome email after a user's registration
func SendWelcomeEmail(user types.User) {
	err := email.Send(email.SendOptions{
		To: []mail.Address{
			{Name: user.DisplayName, Address: user.Email},
		},
		Subject: "Welcome to Docktor",
		Body:    "Your account has been created !",
	})

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to send welcome mail to %q : %q", user.Username, err.Error()))
	}
}

// SendResetPasswordEmail sends an email to a user who forgot his password
func SendResetPasswordEmail(user types.User, url string) {
	err := email.Send(email.SendOptions{
		To: []mail.Address{
			{Name: user.DisplayName, Address: user.Email},
		},
		Subject: "Forgot your password ?",
		Body:    fmt.Sprintf("Your password has been reset and you can change it by following the link %s ", url),
	})

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to send reset password mail to %q : %q", user.Username, err.Error()))
	}
}
