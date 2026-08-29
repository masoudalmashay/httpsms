package emails

import (
	"context"
	"errors"
)

// NewMultiMailer sends each notification through every supplied mailer. It
// attempts every channel even if another channel fails.
func NewMultiMailer(mailers ...Mailer) Mailer {
	return multiMailer{mailers: mailers}
}

type multiMailer struct {
	mailers []Mailer
}

func (mailer multiMailer) Send(ctx context.Context, email *Email) error {
	var sendErrors []error
	for _, channel := range mailer.mailers {
		if err := channel.Send(ctx, email); err != nil {
			sendErrors = append(sendErrors, err)
		}
	}
	return errors.Join(sendErrors...)
}
