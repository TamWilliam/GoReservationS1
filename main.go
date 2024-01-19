package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "image/png"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "domapi92"
	dbname   = "projetgoreservation"
)

type Hairsalon struct {
	IDHairsalon int
	Name        string
	Address     string
	Email       string
}

type HairsalonPageData struct {
	PageTitle  string
	Hairsalons []Hairsalon
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

	/* gestion des assets statiques */
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	/* template pour affichage des salons de coiffure */
	tmplHairSalons := template.Must(template.ParseFiles("templates/salon-de-coiffure.html"))
	http.HandleFunc("/salon-de-coiffure", func(w http.ResponseWriter, r *http.Request) {
		var hairsalons []Hairsalon

		if db != nil {
			rows, err := db.Query("SELECT * FROM hairsalon")
			if err != nil {
				log.Printf("Erreur lors de l'exécution de la requête: %v", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var h Hairsalon
				if err := rows.Scan(&h.IDHairsalon, &h.Name, &h.Address, &h.Email); err != nil {
					log.Printf("Erreur lors de la lecture des lignes: %v", err)
					continue
				}
				hairsalons = append(hairsalons, h)
			}

			if err := rows.Err(); err != nil {
				log.Printf("Erreur lors de la récupération des lignes: %v", err)
			}
		}

		data := HairsalonPageData{
			PageTitle:  "Liste des salons de coiffure",
			Hairsalons: hairsalons,
		}
		tmplHairSalons.Execute(w, data)
	})

	/* listen and serve */
	log.Println("Le serveur démarre sur le port :80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur: ", err)
	}
}
