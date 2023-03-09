package services

import (
	"context"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/mohammadmahdi255/Cloud-job-service/global"
	"time"
)

type Mailgun struct {
	mg *mailgun.MailgunImpl
}

func NewMailgun() *Mailgun {
	mg := mailgun.NewMailgun(global.Domain, global.ApiPrivateKey)
	return &Mailgun{mg: mg}
}

func (m *Mailgun) SendSimpleMessage(text string, to ...string) (string, error) {
	mess := m.mg.NewMessage(
		global.From,
		"Cloud-job-service",
		text,
		to...,
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := m.mg.Send(ctx, mess)
	return id, err
}
