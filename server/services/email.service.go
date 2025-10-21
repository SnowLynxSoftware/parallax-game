package services

import (
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type IEmailService interface {
	SendEmail(options *EmailSendOptions) bool
	GetTemplates() IEmailTemplates
}

type EmailService struct {
	client    *mailjet.Client
	templates IEmailTemplates
}

func NewEmailService(apiKeyPublic string, apiKeyPrivate string, templates IEmailTemplates) IEmailService {
	emailClient := mailjet.NewMailjetClient(apiKeyPublic, apiKeyPrivate)
	return &EmailService{
		client:    emailClient,
		templates: templates,
	}
}

type EmailSendOptions struct {
	FromEmail   string
	ToEmail     string
	Subject     string
	HTMLContent string
}

func (s *EmailService) SendEmail(options *EmailSendOptions) bool {

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: options.FromEmail,
				Name:  options.FromEmail,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: options.ToEmail,
					Name:  options.ToEmail,
				},
			},
			Subject:  options.Subject,
			TextPart: "Parallax does not support plain text emails. Please enable HTML.",
			HTMLPart: options.HTMLContent,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := s.client.SendMailV31(&messages)
	if err != nil {
		util.LogErrorWithStackTrace(err)
		return false
	}
	return true
}

func (s *EmailService) GetTemplates() IEmailTemplates {
	return s.templates
}
