package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// HairSalon représente un salon de coiffure.
type HairSalon struct {
	IDHairSalon int
	Name        string
	Address     string
	Email       string
}

func main() {
	// Connexion à la base de données
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Définir un gestionnaire de fichiers statiques pour le répertoire assets
	fs := http.FileServer(http.Dir("assets"))

	// Utiliser le préfixe "/assets/" pour servir les fichiers statiques
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Récupération des salons depuis la base de données
	salons, err := getHairSalons(db)
	if err != nil {
		log.Fatal(err)
	}

	// Définition du gestionnaire d'URL pour afficher les salons
	http.HandleFunc("/salons", func(w http.ResponseWriter, r *http.Request) {
		// Chargement du template HTML
		tmpl, err := template.ParseFiles("templates/salons.html")
		if err != nil {
			log.Println("Erreur lors du chargement du template:", err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		// Exécution du template avec les données des salons et écriture du résultat dans la réponse HTTP
		err = tmpl.Execute(w, salons)
		if err != nil {
			log.Println("Erreur lors de l'exécution du template:", err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}
	})

	// Démarrage du serveur HTTP sur le port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Fonction pour récupérer les salons depuis la base de données
func getHairSalons(db *sql.DB) ([]HairSalon, error) {
	rows, err := db.Query("SELECT id_hair_salon, name, address, email FROM hair_salons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salons []HairSalon
	for rows.Next() {
		var salon HairSalon
		err := rows.Scan(&salon.IDHairSalon, &salon.Name, &salon.Address, &salon.Email)
		if err != nil {
			return nil, err
		}
		salons = append(salons, salon)
	}

	return salons, nil
}
