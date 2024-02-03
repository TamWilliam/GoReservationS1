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
