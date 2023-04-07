package mail_service

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/mohammadmahdi255/Cloud-job-service/global"
	"time"
)

type Mailgun struct{}

func NewMailgun() *Mailgun {
	return &Mailgun{}
}

func (m *Mailgun) SendSimpleMessage(text string, to ...string) (string, error) {
	mg := mailgun.NewMailgun(global.Domain, global.ApiPrivateKey)
	mess := mg.NewMessage(
		global.From,
		"Cloud-job-service",
		text,
		to...,
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mx, id, err := mg.Send(ctx, mess)
	fmt.Println(mx)
	return id, err
}
