package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "projetgoreservation"
)

type Hairdresser struct {
	IDHairdresser int
	FirstName     string
	LastName      string
	IDHairSalon   int
}

type HairdresserPageData struct {
	PageTitle    string
	Hairdressers []Hairdresser
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	tmplHairSalons := template.Must(template.ParseFiles("templates/salon-de-coiffure.html"))
	http.HandleFunc("/salon-de-coiffure", func(w http.ResponseWriter, r *http.Request) {
		var hairdressers []Hairdresser

		if db != nil {
			rows, err := db.Query("SELECT * FROM hairdresser")
			if err != nil {
				log.Printf("Erreur lors de l'exécution de la requête: %v", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var h Hairdresser
				if err := rows.Scan(&h.IDHairdresser, &h.FirstName, &h.LastName, &h.IDHairSalon); err != nil {
					log.Printf("Erreur lors de la lecture des lignes: %v", err)
					continue
				}
				hairdressers = append(hairdressers, h)
			}

			if err := rows.Err(); err != nil {
				log.Printf("Erreur lors de la récupération des lignes: %v", err)
			}
		}

		data := HairdresserPageData{
			PageTitle:    "Liste des coiffeurs",
			Hairdressers: hairdressers,
		}
		tmplHairSalons.Execute(w, data)
	})

	tmplCustomerAccount := template.Must(template.ParseFiles("templates/compte-utilisateur.html"))
	http.HandleFunc("/compte-utilisateur", func(w http.ResponseWriter, r *http.Request) {
		var hairdressers []Hairdresser

		if db != nil {
			rows, err := db.Query("SELECT * FROM hairdresser")
			if err != nil {
				log.Printf("Erreur lors de l'exécution de la requête: %v", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var h Hairdresser
				if err := rows.Scan(&h.IDHairdresser, &h.FirstName, &h.LastName, &h.IDHairSalon); err != nil {
					log.Printf("Erreur lors de la lecture des lignes: %v", err)
					continue
				}
				hairdressers = append(hairdressers, h)
			}

			if err := rows.Err(); err != nil {
				log.Printf("Erreur lors de la récupération des lignes: %v", err)
			}
		}

		data := HairdresserPageData{
			PageTitle:    "Compte utilisateur",
			Hairdressers: hairdressers,
		}
		tmplCustomerAccount.Execute(w, data)
	})

	log.Println("Le serveur démarre sur le port :80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur: ", err)
	}
}
