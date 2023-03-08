package services

import (
	"context"
	"github.com/mailgun/mailgun-go/v3"
	"time"
)

type Mailgun struct {
	mg *mailgun.MailgunImpl
}

const (
	DOMAIN        = "sandboxd1bc33ec4ae34c949acd84663555b903.mailgun.org"
	ApiPrivateKey = "27f5708ed00814daa31005ed6ddf7aa0-7764770b-58ee5c33"
	FROM          = "nemati.mahdi255@gmail.com"
)

func NewMailgun() *Mailgun {
	mg := mailgun.NewMailgun(DOMAIN, ApiPrivateKey)
	return &Mailgun{mg: mg}
}

func (m *Mailgun) SendSimpleMessage(text string, to ...string) (string, error) {
	mess := m.mg.NewMessage(
		FROM,
		"Cloud-job-service",
		text,
		to...,
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := m.mg.Send(ctx, mess)
	return id, err
}
