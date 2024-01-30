package handler

import (
	"github.com/rs/zerolog/log"
	"net/smtp"
)

import (
	"encoding/json"
)

type EmailNotification struct {
	ToEmail   string
	FromEmail string
	Message   string
	Subject   string
}

type EmailHandler struct {
	token      string
	smtpServer string
}

func NewEmailHandler(token string, smtpServer string) *EmailHandler {
	return &EmailHandler{
		token:      token,
		smtpServer: smtpServer,
	}
}

func (s *EmailHandler) ParseData(rawData []byte) (*EmailNotification, *Error) {
	var emailData EmailNotification
	err := json.Unmarshal(rawData, &emailData)
	if err != nil {
		return nil, &Error{
			ErrorType: NON_RECOVERABLE,
		}
	}

	return &emailData, nil
}

func (s *EmailHandler) SendNotification(notification *EmailNotification) *Error {
	auth := smtp.PlainAuth("SumUp", notification.FromEmail, s.token, s.smtpServer)
	message := getEmailFromData(notification)

	err := smtp.SendMail(s.smtpServer+":587", auth, notification.FromEmail, []string{notification.ToEmail}, message)
	if err != nil {
		log.Error().Msgf("an error occurred while sending an email, retry will be attempted... Error: %s",
			err.Error())
		return &Error{
			ErrorType: RECOVERABLE,
		}
	}

	log.Info().Msg("Email sent successfully")

	return nil
}

func getEmailFromData(data *EmailNotification) []byte {
	return []byte("To: " + data.ToEmail + "\r\n" +
		"Subject: " + data.Subject + "\r\n" +
		"\r\n" +
		data.Message + "\r\n")
}
