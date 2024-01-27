package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"projetgoreservation/config"
	"projetgoreservation/models"

	_ "github.com/lib/pq" // driver PostgreSQL
	"golang.org/x/crypto/bcrypt"
)

var appConfig *config.Config
var db *sql.DB

func CreateTemplates(app *config.Config) {
	appConfig = app
}

func init() {
	// Initialisez la connexion à la base de données ici
	const (
		host     = "localhost"
		port     = 5433
		user     = "postgres"
		password = "jeanne"
		dbname   = "projetgoreservation"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Impossible de se connecter à la base de données:", err)
	}

	fmt.Println("Connexion à la base de données réussie!")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur lors du parsing du formulaire", http.StatusInternalServerError)
			return
		}

		firstName := r.Form.Get("firstName")
		lastName := r.Form.Get("lastName")
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO customers (firstname, lastname, email, password) VALUES ($1, $2, $3, $4)", firstName, lastName, email, string(hashedPassword))
		if err != nil {
			http.Error(w, "Erreur lors de l'inscription de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// redirection
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	}

}

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur lors du parsing du formulaire", http.StatusInternalServerError)
			return
		}

		email := r.Form.Get("email")
		password := r.Form.Get("password")

		var hashedPassword string
		err = db.QueryRow("SELECT password FROM customers WHERE email = $1", email).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Erreur lors de la recherche de l'utilisateur", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Gérer la session ou rediriger l'utilisateur
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	// redirection
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// on récupère les fichiers page.tmpl
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// template du cache
	t, ok := appConfig.TemplateCache[tmpl]
	if !ok {
		http.Error(w, "Le template n'existe pas", http.StatusInternalServerError)
		return
	}

	// exrcute le template avec les données fourni
	err := t.Execute(w, td)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
}
