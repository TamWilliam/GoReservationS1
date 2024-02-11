package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Customer struct {
	IDCustomer int
	Email      string
}

type HairSalon struct {
	IDHairSalon int
	Name        string
	Address     string
	Email       string
}

type HairDresser struct {
	IDHairDresser int
	FirstName     string
}

type PageVariables struct {
	PageTitle    string
	Utilisateurs []Customer
	Salons       []HairSalon
	Coiffeurs    []HairDresser
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("templates/reservation.html", "templates/confirmation_reservation.html"))
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

	http.HandleFunc("/reservation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			utilisateurID := r.FormValue("utilisateur")
			salonID := r.FormValue("salon")
			coiffeurID := r.FormValue("coiffeur")
			date := r.FormValue("date")

			err := saveReservation(db, utilisateurID, salonID, coiffeurID, date)
			if err != nil {
				http.Error(w, "Erreur lors de la réservation", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/confirmation_reservation", http.StatusSeeOther)
			return
		}

		utilisateurs, err := getCustomers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		salons, err := getHairSalons(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		coiffeurs, err := getHairDressers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageVariables := PageVariables{
			PageTitle:    "Réservation",
			Utilisateurs: utilisateurs,
			Salons:       salons,
			Coiffeurs:    coiffeurs,
		}

		tmpl.Execute(w, pageVariables)
	},
	)

	http.HandleFunc("/confirmation_reservation", func(w http.ResponseWriter, r *http.Request) {

		confirmationPageVariables := PageVariables{
			PageTitle: "Confirmation réservation",
		}

		tmpl.ExecuteTemplate(w, "confirmation_reservation.html", confirmationPageVariables)
	})

	// Définition du gestionnaire d'URL pour afficher les salons
	http.HandleFunc("/salons", func(w http.ResponseWriter, r *http.Request) {
		// Chargement du template HTML
		tmpl, err := template.ParseFiles("templates/salons.html")
		if err != nil {
			log.Println("Erreur lors du chargement du template:", err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		salons, err := getHairSalons(db)
		if err != nil {
			log.Fatal(err)
		}

		// Exécution du template avec les données des salons et écriture du résultat dans la réponse HTTP
		err = tmpl.Execute(w, salons)
		if err != nil {
			log.Println("Erreur lors de l'exécution du template:", err)
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query("SELECT id_customer, email FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.IDCustomer, &customer.Email)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

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

func getHairDressers(db *sql.DB) ([]HairDresser, error) {
	rows, err := db.Query("SELECT id_hair_dresser, first_name FROM hair_dressers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coiffeurs []HairDresser
	for rows.Next() {
		var coiffeur HairDresser
		err := rows.Scan(&coiffeur.IDHairDresser, &coiffeur.FirstName)
		if err != nil {
			return nil, err
		}
		coiffeurs = append(coiffeurs, coiffeur)
	}

	return coiffeurs, nil
}

func saveReservation(db *sql.DB, utilisateurID, salonID, coiffeurID, date string) error {
	idCustomer, err := strconv.Atoi(utilisateurID)
	if err != nil {
		return err
	}

	idSalon, err := strconv.Atoi(salonID)
	if err != nil {
		return err
	}

	idCoiffeur, err := strconv.Atoi(coiffeurID)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO reservations (id_customer, id_hair_salon, id_hair_dresser, reservation_date) VALUES ($1, $2, $3, $4)",
		idCustomer, idSalon, idCoiffeur, date)
	if err != nil {
		return err
	}

	return nil
}
