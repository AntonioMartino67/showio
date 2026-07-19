package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOTPEmail(to, code string) error {
	from := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	if from == "" || password == "" {
		return fmt.Errorf("GMAIL_ADDRESS o GMAIL_APP_PASSWORD non impostate")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Il tuo codice di verifica Showio"
	body := fmt.Sprintf(`
Verifica il tuo account Showio

Il tuo codice di verifica è: %s

Il codice scade tra 10 minuti.
`, code)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}

func SendNewSeasonEmail(to, username, showTitle string, seasonNumber int) error {
	from := os.Getenv("GMAIL_ADDRESS")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	if from == "" || password == "" {
		return fmt.Errorf("GMAIL_ADDRESS o GMAIL_APP_PASSWORD non impostate")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := fmt.Sprintf("Nuova stagione disponibile: %s", showTitle)
	body := fmt.Sprintf(`
Ciao %s,

È uscita la stagione %d di "%s" che stai seguendo su Showio!

Vai su Showio per segnare i nuovi episodi.
`, username, seasonNumber, showTitle)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}