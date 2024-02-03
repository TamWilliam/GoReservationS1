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

// Data structures
type Customer struct {
	IDCustomer int
	Email      string
}

type HairSalon struct {
	IDHairSalon int
	Name        string
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
	tmpl = template.Must(template.ParseFiles("templates/reservation.html"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD")))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/reservation", func(w http.ResponseWriter, r *http.Request) {
		// Vérifier que la requête est une requête POST
		if r.Method == http.MethodPost {
			// Récupérer les données du formulaire
			utilisateurID := r.FormValue("utilisateur")
			salonID := r.FormValue("salon")
			coiffeurID := r.FormValue("coiffeur")
			date := r.FormValue("date")

			// Appeler la fonction saveReservation avec les données du formulaire
			err := saveReservation(db, utilisateurID, salonID, coiffeurID, date)
			if err != nil {
				http.Error(w, "Erreur lors de la réservation", http.StatusInternalServerError)
				return
			}

			// Rediriger ou afficher un message de confirmation, etc.
			http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
			return
		}

		// Si la méthode de la requête n'est pas POST, afficher simplement le formulaire
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
	})

	log.Fatal(http.ListenAndServe(":6060", nil))
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
	rows, err := db.Query("SELECT id_hair_salon, name FROM hair_salons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var salons []HairSalon
	for rows.Next() {
		var salon HairSalon
		err := rows.Scan(&salon.IDHairSalon, &salon.Name)
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
	// Convertir les identifiants en entiers
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

	// TODO: Mettez en œuvre la logique d'enregistrement de la réservation dans la base de données
	// Utilisez les paramètres (idCustomer, idSalon, idCoiffeur, date) pour insérer les données dans la table des réservations

	// Exemple (veuillez remplacer cela par la logique de votre application) :
	_, err = db.Exec("INSERT INTO reservations (id_customer, id_hair_salon, id_hair_dresser, reservation_date) VALUES ($1, $2, $3, $4)",
		idCustomer, idSalon, idCoiffeur, date)
	if err != nil {
		return err
	}

	return nil
}
