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
	password = "5712"
	dbname   = "projetgoreservation"
)

type Hairsalon struct {
	IDHairsalon int
	Name        string
	Address     string
	Email       string
}

type Customer struct {
	IDCustomer int
	role       int
	firstname  string
	lastname   string
	email      string
	password   string
}

type HairsalonPageData struct {
	PageTitle  string
	Hairsalons []Hairsalon
}

type AdminPageData struct {
	PageTitle  string
	Hairsalons []Hairsalon
	Customers  []Customer
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

	tmplHairSalons := template.Must(template.ParseFiles("templates/salon-de-coiffure.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	//page admin
	tmplAdmin := template.Must(template.ParseFiles("templates/page-admin.html"))
	http.HandleFunc("/page-admin", func(w http.ResponseWriter, r *http.Request) {
		//get all value from hairsalon
		var hairsalons []Hairsalon

		if db != nil {
			rows, err := db.Query("SELECT * FROM hairsalon ORDER BY id_hairsalon ASC ")
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

		//get all value from customer
		var customers []Customer

		if db != nil {
			rows, err := db.Query("SELECT * FROM customer")
			if err != nil {
				log.Printf("Erreur lors de l'exécution de la requête: %v", err)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var c Customer
				if err := rows.Scan(&c.IDCustomer, &c.role, &c.firstname, &c.lastname, &c.email, &c.password); err != nil {
					log.Printf("Erreur lors de la lecture des lignes: %v", err)
					continue
				}
				customers = append(customers, c)

			}

			if err := rows.Err(); err != nil {
				log.Printf("Erreur lors de la récupération des lignes: %v", err)
			}
		}

		data := AdminPageData{
			PageTitle:  "Admin",
			Hairsalons: hairsalons,
			Customers:  customers,
		}
		tmplAdmin.Execute(w, data)
	})

	//send data

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}

		// Access form values
		id := r.Form.Get("id")
		name := r.Form.Get("name")
		address := r.Form.Get("address")
		email := r.Form.Get("email")

		// Handle form data (e.g., insert into database)
		// Example:
		_, err = db.Exec("UPDATE hairsalon SET name=$2, address=$3, email=$4 WHERE id_hairsalon=$1", id, name, address, email)
		if err != nil {
			http.Error(w, "Failed to update data in database", http.StatusInternalServerError)
			return
		}

		// Send response
		fmt.Fprintf(w, "Data submitted successfully")
	})

	log.Println("Le serveur démarre sur le port :88")
	err = http.ListenAndServe(":88", nil)
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur: ", err)
	}
}
