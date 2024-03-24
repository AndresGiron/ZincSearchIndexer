package main

import (
	"Indexation/server"
)

func main() {

	// Crear el index en ZincSearch
	//zinc.CreateIndex()
	// Subir los correos a zinc search
	//zinc.PushMails()
	//zinc.PushMailsQuick()
	//zinc.PushMailsQuickAP()
	//Encender el API Rest
	mux := server.Routes()
	server := server.NewServer(mux)
	server.Run()
}
