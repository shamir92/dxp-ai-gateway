package mail

import "context"

type MagicCodeMessage struct {
	To               string
	Code             string
	Purpose          string // login | register
	ExpiresInMinutes int
}

type LoginAlertMessage struct {
	To         string
	Device     string
	IP         string
	Location   string
	LoggedInAt string
}

type Mailer interface {
	SendMagicCode(ctx context.Context, msg MagicCodeMessage) error
	SendLoginAlert(ctx context.Context, msg LoginAlertMessage) error
}
