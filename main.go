package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Customer struct {
	Id        int    `json:"id_customer"`
	Role      int    `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=postgres dbname=projetgoreservation password=Coucou75! sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func getCustomersJSON(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM customers")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	customers := make([]Customer, 0)
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.Id, &c.Role, &c.Firstname, &c.Lastname, &c.Email, &c.Password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, ok := vars["id_customer"]
	if !ok {
		http.Error(w, "ID du client manquant", http.StatusBadRequest)
		return
	}

	log.Printf("Tentative de suppression du client avec l'ID : %s\n", customerID)

	result, err := db.Exec("DELETE FROM customers WHERE id_customer = $1", customerID)
	if err != nil {
		log.Printf("Erreur lors de la suppression du client : %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Nombre de lignes affectées : %d\n", rowsAffected)

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Client avec l'ID %s supprimé avec succès\n", customerID)
}

func getCustomersHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/customer_list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT * FROM customers")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	customers := make([]Customer, 0)
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.Id, &c.Role, &c.Firstname, &c.Lastname, &c.Email, &c.Password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, c)
	}

	tmpl.Execute(w, customers)
}

func fetchCustomers() ([]Customer, error) {
	// remplacez par l'URL de l'api
	resp, err := http.Get("http://localhost:8000/customers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var customers []Customer
	err = json.NewDecoder(resp.Body).Decode(&customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func getCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, ok := vars["id_customer"]
	if !ok {
		http.Error(w, "ID du client manquant", http.StatusBadRequest)
		return
	}

	var c Customer
	err := db.QueryRow("SELECT * FROM customers WHERE id_customer = $1", customerID).Scan(
		&c.Id, &c.Role, &c.Firstname, &c.Lastname, &c.Email, &c.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("templates/customer_details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, c)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	customers, err := fetchCustomers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, customers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/customers", getCustomersJSON).Methods("GET")
	r.HandleFunc("/api/customers/{id}", deleteCustomer).Methods("DELETE") // Nouvelle route pour la suppression
	r.HandleFunc("/customers", getCustomersHTML).Methods("GET")
	r.HandleFunc("/compte-customers/{id}", getCustomerByID).Methods("GET")
	r.HandleFunc("/view-customers", serveTemplate)

	log.Fatal(http.ListenAndServe(":8000", r))
}
