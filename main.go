package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Customer représente la structure de vos données client.
type Customer struct {
	Id        int    `json:"id_customer"`
	Role      int    `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// fetchCustomers récupère les données des clients depuis votre API.
func fetchCustomers() ([]Customer, error) {
	// Remplacez par l'URL de votre API
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

// serveTemplate sert votre page web avec les données des clients.
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	customers, err := fetchCustomers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("templates/salon-de-coiffure.html")
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
	r := mux.NewRouter()
	r.HandleFunc("/", serveTemplate)

	log.Fatal(http.ListenAndServe(":80", r))
}
