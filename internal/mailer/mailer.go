package mailer

import "embed"

const (
	FromName            = "GopherSocial"
	maxRetries          = 3
	UserWelcomeTemplate = "userInvitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(
		template, username, userEmail string,
		data any,
		isSandbox bool,
	) error
}
