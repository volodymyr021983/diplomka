package auth

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
)

func smtp_init() emaildelivery.SMTPSettings {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpHost := os.Getenv("SMTP_HOST")
	appName := os.Getenv("APP_NAME")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	smtpSettings := emaildelivery.SMTPSettings{
		Host: smtpHost,
		From: emaildelivery.SMTPFrom{
			Name:  appName,
			Email: smtpEmail, //SMTP EMAIL MUST BE WITH VERIFICATED DOMAIN
		},
		Port:     587,
		Username: &smtpUsername, // this is optional. In case not given, from.email will be used
		Password: smtpPassword,
		Secure:   false,

		// this is optional. TLS config is used if Secure is set to true, or server supports STARTTLS
		// if not provided, the SDK will use a default config
		/*
			TLSConfig: &tls.Config{
				ServerName: "smtp.mailersend.net",
			},
		*/
	}
	return smtpSettings
}
