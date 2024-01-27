package main

import (
	"projetgoreservation/config"
	"projetgoreservation/internal/handlers"
	"fmt"
	"net/http"
)

const port = ":8080"

func main() {
	var appConfig config.Config
	// Initialisation de la configuration de l'application et du cache de template
	templateCache, err := handlers.CreateTemplateCache()

	if err != nil {
		panic(err)
	}

	// on défini templateCahe et le port
	appConfig.TemplateCache = templateCache
	appConfig.Port = ":8080"

	handlers.CreateTemplates(&appConfig)

	http.HandleFunc("/signup", handlers.SignUp)
	http.HandleFunc("/signin", handlers.SignIn)

	fmt.Println("(http://localhost:8080) - Server started on port", port)
	http.ListenAndServe(appConfig.Port, nil)
}
