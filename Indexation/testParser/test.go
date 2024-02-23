package main

import (
	"Indexation/parser"
	"fmt"
	"log"
)

func main() {
	emailObj, err := parser.EmailFromFile("/home/gyron/Desktop/EmailIndex/Indexation/maildir/allen-p/all_documents/1.")
	if err != nil {
		log.Fatalf("Error al leer el archivo de correo: %v", err)
	}

	fmt.Printf("Mensaje ID: %s\n", emailObj.MessageId)
	fmt.Printf("Fecha: %s\n", emailObj.Date)
	fmt.Printf("De: %s\n", emailObj.From)
	fmt.Printf("Para: %s\n", emailObj.To)
	fmt.Printf("Asunto: %s\n", emailObj.Subject)
	fmt.Printf("Cuerpo: %s\n", emailObj.Body)
}
