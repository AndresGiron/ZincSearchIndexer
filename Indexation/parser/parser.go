package parser

import (
	"bytes"
	"fmt"
	"net/mail"
	"os"
	"regexp"
	"time"
)

// Estructura del Mail
type Email struct {
	MessageId string    `json:"messageId"`
	Date      time.Time `json:"date"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
}

func EmailFromFile(path string) (*Email, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	msg, err := mail.ReadMessage(file)
	if err != nil {
		return nil, err
	}

	emailObj := &Email{
		MessageId: msg.Header.Get("Message-ID"),
		Subject:   msg.Header.Get("Subject"),
	}

	date, err := msg.Header.Date()
	if err != nil {
		return nil, err
	}
	emailObj.Date = date

	fromHeader := msg.Header.Get("From")
	if fromHeader != "" {
		addres := extractEmailAddresses(fromHeader)
		if len(addres) > 0 {
			emailObj.From = addres[0]
		} else {
			emailObj.From = ""
			fmt.Println("Este from esta vacio")
		}
	} else {
		fmt.Println("El encabezado From está vacío o no existe")
	}

	toHeader := msg.Header.Get("To")
	if toHeader != "" {
		addresses := extractEmailAddresses(toHeader)
		emailObj.To = append(emailObj.To, addresses...)
	}
	// parsear el cuerpo del correo
	buf := new(bytes.Buffer)
	buf.ReadFrom(msg.Body)
	emailObj.Body = buf.String()

	return emailObj, nil
}

// Expresion regular para correos de todos los formatos
func extractEmailAddresses(header string) []string {
	re := regexp.MustCompile(`[\w.-]+@[\w.-]+\.\w+`)
	return re.FindAllString(header, -1)
}
