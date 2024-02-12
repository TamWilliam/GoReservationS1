package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Customer struct {
	IDCustomer int    `json:"id_customer"`
	Role       int    `json:"role"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type Hairdresser struct {
	IDHairDresser int    `json:"id_hair_dresser"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	IDHairSalon   int    `json:"id_hair_salon"`
}

type Hairsalon struct {
	IDHairSalon int    `json:"id_hair_salon"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func main() {
	router := gin.Default()

	// Configuration du store de sessions
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Charger les templates
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		customerName, exists := session.Get("customerName").(string) // Assurez-vous que la session contient bien une chaîne

		if !exists {
			// Si l'utilisateur n'est pas connecté ou si le nom n'est pas dans la session, passez une chaîne vide ou un message par défaut
			customerName = ""
		}

		// Passer le nom de l'utilisateur au template, qui sera vide si non connecté
		c.HTML(http.StatusOK, "accueil.html", gin.H{
			"CustomerName": customerName,
		})
	})

	// Définir la route pour afficher le formulaire de connexion
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// Définir la route pour traiter les données du formulaire
	router.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		customers, err := getCustomers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Problème lors de la récupération des clients"})
			return
		}

		session := sessions.Default(c) // Récupérer la session actuelle

		for _, customer := range customers {
			if customer.Email == email && customer.Password == password {
				session.Set("customerID", customer.IDCustomer)
				session.Set("customerName", customer.FirstName) // Stocker aussi le prénom
				if err := session.Save(); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de sauvegarder la session"})
					return
				}
				c.Redirect(http.StatusFound, "/")
				return
			}
		}

		// Si aucun utilisateur correspondant n'a été trouvé
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
	})

	router.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear() // Efface toutes les données de la session
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de sauvegarder la session"})
			return
		}
		c.Redirect(http.StatusFound, "/") // Redirige vers la page d'accueil après la déconnexion
	})

	router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	router.POST("/signup", func(c *gin.Context) {
		// Créer une instance de Customer à partir des données de formulaire
		customer := Customer{
			FirstName: c.PostForm("first_name"),
			LastName:  c.PostForm("last_name"),
			Email:     c.PostForm("email"),
			Password:  c.PostForm("password"),
			// Assurez-vous que le modèle Customer et l'API supportent tous les champs nécessaires
		}

		// Convertir l'instance Customer en JSON
		jsonData, err := json.Marshal(customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la conversion des données en JSON"})
			return
		}

		// Envoyer une requête POST à l'API pour créer le customer
		apiURL := "http://localhost:6060/customer" // Remplacez par l'URL réelle de votre API
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la requête"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'envoi de la requête à l'API"})
			return
		}
		defer response.Body.Close()

		// Vérifier la réponse de l'API
		if response.StatusCode == http.StatusCreated {
			// Rediriger l'utilisateur vers la page de connexion ou une page de confirmation après la création réussie du compte
			c.Redirect(http.StatusFound, "/login")
		} else {
			// Si l'API renvoie une erreur, lire le corps de la réponse pour plus de détails
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture de la réponse de l'API"})
				return
			}
			c.JSON(response.StatusCode, gin.H{"error": string(responseBody)})
		}
	})

	router.GET("/customers", func(c *gin.Context) {
		customers, err := getCustomers() // Utilisez votre fonction existante pour récupérer les clients
		if err != nil {
			// Gérer l'erreur
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Problème lors de la récupération des clients"})
			return
		}

		// Passer les données des clients au template
		c.HTML(http.StatusOK, "list_customers.html", customers)
	})

	router.GET("/edit-customer/:id_customer", func(c *gin.Context) {
		idCustomer := c.Param("id_customer")

		// Ici, ajoutez la logique pour récupérer les informations du client
		// en utilisant l'ID. Ceci est juste un exemple :
		customer, err := getCustomerByID(idCustomer)
		if err != nil {
			// Gérer l'erreur
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Client non trouvé"})
			return
		}

		c.HTML(http.StatusOK, "edit_customer.html", customer)
	})

	router.GET("/add-customer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_customer.html", nil)
	})

	router.POST("/create-customer", func(c *gin.Context) {
		// Créer une instance de Customer à partir des données de formulaire
		customer := Customer{
			FirstName: c.PostForm("first_name"),
			LastName:  c.PostForm("last_name"),
			Email:     c.PostForm("email"),
			Password:  c.PostForm("password"),
			// N'oubliez pas d'ajouter le rôle ici si nécessaire
		}

		// Convertir l'instance Customer en JSON
		jsonData, err := json.Marshal(customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la conversion des données en JSON"})
			return
		}

		// Envoyer une requête POST à l'API pour créer le customer
		apiURL := "http://localhost:6060/customer" // URL de votre API
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la requête"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'envoi de la requête à l'API"})
			return
		}
		defer response.Body.Close()

		// Vérifier la réponse de l'API
		if response.StatusCode == http.StatusCreated {
			c.Redirect(http.StatusFound, "/customers")
		} else {
			// Si l'API renvoie une erreur, lire le corps de la réponse pour plus de détails
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture de la réponse de l'API"})
				return
			}
			c.JSON(response.StatusCode, gin.H{"error": string(responseBody)})
		}
	})

	router.POST("/update-customer", func(c *gin.Context) {
		// Récupérer les données du formulaire
		idCustomer := c.PostForm("id_customer")
		firstName := c.PostForm("first_name")
		lastName := c.PostForm("last_name")
		email := c.PostForm("email")
		password := c.PostForm("password") // Soyez prudent avec la gestion des mots de passe
		// Assurez-vous que le champ "role" est également récupéré si nécessaire

		// Préparer l'instance Customer avec les données récupérées
		customer := Customer{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
			// Role:       role, // Ajoutez cette ligne si vous récupérez le rôle du formulaire
		}

		// Appeler la fonction d'API pour mettre à jour le client
		err := updateCustomer(idCustomer, customer)
		if err != nil {
			// En cas d'erreur, afficher l'erreur et peut-être rediriger vers une page d'erreur
			log.Println("Erreur lors de la mise à jour du client:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du client"})
			return
		}

		// Si tout va bien, rediriger vers la liste des clients
		c.Redirect(http.StatusFound, "/customers")
	})

	router.POST("/delete-customer/:id", func(c *gin.Context) {
		id := c.Param("id") // Récupérer l'ID du client à supprimer

		// Construire l'URL pour l'API
		apiUrl := fmt.Sprintf("http://localhost:6060/customer/%s", id)

		// Créer une requête DELETE
		req, err := http.NewRequest("DELETE", apiUrl, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la requête de suppression"})
			return
		}

		// Envoyer la requête à l'API
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'envoi de la requête de suppression à l'API"})
			return
		}
		defer response.Body.Close()

		// Vérifier la réponse de l'API
		if response.StatusCode == http.StatusOK {
			// Rediriger vers la liste des clients après la suppression
			c.Redirect(http.StatusFound, "/customers")
		} else {
			// Si l'API renvoie une erreur, lire le corps de la réponse pour plus de détails
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture de la réponse de l'API"})
				return
			}
			c.JSON(response.StatusCode, gin.H{"error": string(responseBody)})
		}
	})

	router.GET("/hairdressers", func(c *gin.Context) {
		hairdressers, err := getHairdressers() // Utilisez votre fonction existante pour récupérer les clients
		if err != nil {
			// Gérer l'erreur
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Problème lors de la récupération des clients"})
			return
		}

		// Passer les données des clients au template
		c.HTML(http.StatusOK, "list_hairdresser.html", hairdressers)
	})

	router.GET("/edit-hairdresser/:id_hair_dresser", func(c *gin.Context) {
		IDHairDresser := c.Param("id_hair_dresser")

		hairdresser, err := getHairdresserByID(IDHairDresser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Hairdresser non trouvé"})
			return
		}

		fmt.Print(hairdresser)

		c.HTML(http.StatusOK, "edit_hairdresser.html", hairdresser)
	})

	router.POST("/update-hairdresser", func(c *gin.Context) {
		// Récupérer les données du formulaire
		idHairDresserStr := c.PostForm("id_hair_dresser")
		firstName := c.PostForm("first_name")
		lastName := c.PostForm("last_name")
		idHairSalonStr := c.PostForm("idHairSalon")

		// Convertir idHairDresser et idHairSalon de string à int
		idHairDresser, err := strconv.Atoi(idHairDresserStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID du coiffeur invalide"})
			return
		}

		idHairSalon, err := strconv.Atoi(idHairSalonStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID du salon de coiffure invalide"})
			return
		}

		// Préparer l'instance Hairdresser avec les données récupérées
		hairdresser := Hairdresser{
			IDHairDresser: idHairDresser,
			FirstName:     firstName,
			LastName:      lastName,
			IDHairSalon:   idHairSalon,
		}

		// Appeler la fonction pour mettre à jour le hairdresser
		err = updateHairdresser(idHairDresserStr, hairdresser) // Assurez-vous que updateHairdresser peut gérer idHairDresser comme string
		if err != nil {
			log.Println("Erreur lors de la mise à jour du coiffeur:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du coiffeur"})
			return
		}

		// Si tout va bien, rediriger vers la liste des coiffeurs
		c.Redirect(http.StatusFound, "/hairdressers")
	})

	router.GET("/add-hairdresser", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_hairdresser.html", nil)
	})

	router.POST("/create-hairdresser", func(c *gin.Context) {
		var form Hairdresser
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := createHairdresser(form)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la création du coiffeur"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/hairdressers")
	})

	router.POST("/delete-hairdresser/:id", func(c *gin.Context) {
		id := c.Param("id") // Récupérer l'ID du client à supprimer

		// Construire l'URL pour l'API
		apiUrl := fmt.Sprintf("http://localhost:6060/hairdresser/%s", id)

		// Créer une requête DELETE
		req, err := http.NewRequest("DELETE", apiUrl, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la requête de suppression"})
			return
		}

		// Envoyer la requête à l'API
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'envoi de la requête de suppression à l'API"})
			return
		}
		defer response.Body.Close()

		// Vérifier la réponse de l'API
		if response.StatusCode == http.StatusOK {
			// Rediriger vers la liste des clients après la suppression
			c.Redirect(http.StatusFound, "/hairdressers")
		} else {
			// Si l'API renvoie une erreur, lire le corps de la réponse pour plus de détails
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture de la réponse de l'API"})
				return
			}
			c.JSON(response.StatusCode, gin.H{"error": string(responseBody)})
		}
	})

	router.GET("/hairsalons", func(c *gin.Context) {
		hairsalons, err := getHairsalons()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Problème lors de la récupération des salons de coiffure"})
			return
		}
		c.HTML(http.StatusOK, "list_hairsalons.html", hairsalons)
	})

	router.GET("/edit-hairsalon/:id_hairsalon", func(c *gin.Context) {
		idHairsalon := c.Param("id_hairsalon")
		hairsalon, err := getHairsalonByID(idHairsalon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Salon de coiffure non trouvé"})
			return
		}
		c.HTML(http.StatusOK, "edit_hairsalon.html", hairsalon)
	})

	router.GET("/add-hairsalon", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_hairsalon.html", nil)
	})

	router.POST("/create-hairsalon", func(c *gin.Context) {
		var form Hairsalon
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := createHairsalon(form)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la création du salon de coiffure"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/hairsalons")
	})

	router.POST("/update-hairsalon", func(c *gin.Context) {
		var form Hairsalon
		// Assurez-vous d'inclure un champ caché dans votre formulaire HTML pour l'ID du Hairsalon
		id := c.PostForm("id_hair_salon")

		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// La fonction updateHairsalon doit être ajustée pour prendre l'ID et l'objet form
		err := updateHairsalon(id, form)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la mise à jour du salon de coiffure"})
			return
		}

		c.Redirect(http.StatusSeeOther, "/hairsalons")
	})

	router.POST("/delete-hairsalon/:id_hairsalon", func(c *gin.Context) {
		idHairsalon := c.Param("id_hairsalon")
		err := deleteHairsalon(idHairsalon)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la suppression du salon de coiffure"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/hairsalons")
	})

	router.Run(":8080")
}

func getCustomers() ([]Customer, error) {
	resp, err := http.Get("http://localhost:6060/customers")
	if err != nil {
		log.Println("Erreur lors de la récupération des clients:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return nil, err
	}

	var customers []Customer
	err = json.Unmarshal(body, &customers)
	if err != nil {
		log.Println("Erreur lors du décodage des clients:", err)
		return nil, err
	}

	return customers, nil
}

func getCustomerByID(id string) (Customer, error) {
	var customer Customer

	// Formatter correctement l'URL avec l'ID du client
	url := fmt.Sprintf("http://localhost:6060/customer/%s", id)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erreur lors de la récupération du client:", err)
		return customer, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return customer, err
	}

	// Déserialiser la réponse JSON dans l'objet Customer
	err = json.Unmarshal(body, &customer)
	if err != nil {
		log.Println("Erreur lors du décodage du client:", err)
		return customer, err
	}

	return customer, nil
}

// updateCustomerAPI envoie une requête PUT pour mettre à jour un client
func updateCustomer(id string, customer Customer) error {
	url := fmt.Sprintf("http://localhost:6060/customer/%s", id)

	// Convertir le client en JSON
	jsonData, err := json.Marshal(customer)
	if err != nil {
		return err
	}

	// Créer la requête PUT
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier la réponse
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Erreur lors de la mise à jour du client: %s", string(body))
	}

	return nil
}

func createCustomerPost(c *gin.Context) {
	var form Customer
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertissez form en JSON pour l'envoi à l'API
	jsonData, err := json.Marshal(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode data to JSON"})
		return
	}

	// Définissez l'en-tête Content-Type de la requête sur application/json
	apiUrl := "http://localhost:6060/customer"
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to API"})
		return
	}
	defer response.Body.Close()

	// Vérifiez le statut de la réponse
	if response.StatusCode != http.StatusCreated {
		// Lire le corps de la réponse pour plus d'informations
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from API"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API did not create the customer", "response": string(responseBody)})
		return
	}

	// Redirection vers la liste des clients après la création
	c.Redirect(http.StatusSeeOther, "/customers")
}

func getHairdressers() ([]Hairdresser, error) {
	var hairdressers []Hairdresser
	resp, err := http.Get("http://localhost:6060/hairdressers")
	if err != nil {
		log.Println("Erreur lors de la récupération des hairdressers:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return nil, err
	}

	err = json.Unmarshal(body, &hairdressers)
	if err != nil {
		log.Println("Erreur lors du décodage des hairdressers:", err)
		return nil, err
	}

	return hairdressers, nil
}

func getHairdresserByID(id string) (Hairdresser, error) {
	var hairdresser Hairdresser

	// Formatter correctement l'URL avec l'ID du hairdresser
	url := fmt.Sprintf("http://localhost:6060/hairdresser/%s", id)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erreur lors de la récupération du hairdresser:", err)
		return hairdresser, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return hairdresser, err
	}

	// Déserialiser la réponse JSON dans l'objet Hairdresser
	err = json.Unmarshal(body, &hairdresser)
	if err != nil {
		log.Println("Erreur lors du décodage du hairdresser:", err)
		return hairdresser, err
	}

	return hairdresser, nil
}

func updateHairdresser(id string, hairdresser Hairdresser) error {
	url := fmt.Sprintf("http://localhost:6060/hairdresser/%s", id)

	jsonData, err := json.Marshal(hairdresser)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Erreur lors de la mise à jour du hairdresser: %s", string(body))
	}

	return nil
}

func createHairdresser(hairdresser Hairdresser) error {
	apiUrl := "http://localhost:6060/hairdresser" // URL de votre API
	jsonData, err := json.Marshal(hairdresser)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Erreur lors de la création du coiffeur : %s", string(responseBody))
	}

	return nil
}

func getHairsalons() ([]Hairsalon, error) {
	resp, err := http.Get("http://localhost:6060/hairsalons")
	if err != nil {
		log.Println("Erreur lors de la récupération des salons de coiffure:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return nil, err
	}

	var hairsalons []Hairsalon
	err = json.Unmarshal(body, &hairsalons)
	if err != nil {
		log.Println("Erreur lors du décodage des salons de coiffure:", err)
		return nil, err
	}

	return hairsalons, nil
}

func getHairsalonByID(id string) (Hairsalon, error) {
	var hairsalon Hairsalon

	// Formatter correctement l'URL avec l'ID du salon de coiffure
	url := fmt.Sprintf("http://localhost:6060/hairsalon/%s", id)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erreur lors de la récupération du salon de coiffure:", err)
		return hairsalon, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return hairsalon, err
	}

	// Déserialiser la réponse JSON dans l'objet Hairsalon
	err = json.Unmarshal(body, &hairsalon)
	if err != nil {
		log.Println("Erreur lors du décodage du salon de coiffure:", err)
		return hairsalon, err
	}

	return hairsalon, nil
}

func createHairsalon(hairsalon Hairsalon) error {
	apiUrl := "http://localhost:6060/hairsalon" // URL de votre API
	jsonData, err := json.Marshal(hairsalon)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Erreur lors de la création du salon de coiffure : %s", string(responseBody))
	}

	return nil
}

func updateHairsalon(id string, hairsalon Hairsalon) error {
	apiUrl := fmt.Sprintf("http://localhost:6060/hairsalon/%s", id)

	jsonData, err := json.Marshal(hairsalon)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Erreur lors de la mise à jour du salon de coiffure: %s", string(body))
	}

	return nil
}

func deleteHairsalon(id string) error {
	apiUrl := fmt.Sprintf("http://localhost:6060/hairsalon/%s", id) // URL de votre API pour supprimer un salon de coiffure

	// Créer la requête DELETE
	req, err := http.NewRequest("DELETE", apiUrl, nil)
	if err != nil {
		return err
	}

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier la réponse
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Erreur lors de la suppression du salon de coiffure: %s", string(body))
	}

	return nil
}
