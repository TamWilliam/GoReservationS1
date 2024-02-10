package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

func main() {
	router := gin.Default()

	// Charger les templates
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "accueil.html", nil)
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
			// Gérer l'erreur
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Problème lors de la récupération des clients"})
			return
		}

		for _, customer := range customers {
			println(customer.IDCustomer, customer.Email, customer.Password)
			if customer.Email == email && customer.Password == password {
				// Utilisateur trouvé, gérer la connexion
				c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie"})
				return
			}
		}

		// Si aucun utilisateur correspondant n'a été trouvé
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
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

	router.GET("/edit-hairdresser/:id_hairdresser", func(c *gin.Context) {
		IDHairDresser := c.Param("id_hairdresser")

		// Ici, ajoutez la logique pour récupérer les informations du client
		// en utilisant l'ID. Ceci est juste un exemple :
		hairdresser, err := getCustomerByID(IDHairDresser)
		if err != nil {
			// Gérer l'erreur
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Client non trouvé"})
			return
		}

		println("infos : ", hairdresser)

		c.HTML(http.StatusOK, "edit_hairdresser.html", hairdresser)
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

func newCustomer(c *gin.Context) {
	// Affiche le formulaire pour un nouveau client
	c.HTML(http.StatusOK, "new_customer.html", nil)
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
	resp, err := http.Get("http://localhost:6060/hairdressers")
	if err != nil {
		log.Println("Erreur lors de la récupération des sallons:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return nil, err
	}

	var hairdressers []Hairdresser
	err = json.Unmarshal(body, &hairdressers)
	if err != nil {
		log.Println("Erreur lors du décodage des sallons:", err)
		return nil, err
	}

	return hairdressers, nil
}

func getHairdresserByID(id string) (Hairdresser, error) {
	var hairdresser Hairdresser

	// Formatter correctement l'URL avec l'ID du client
	url := fmt.Sprintf("http://localhost:6060/hairdresser/%s", id)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Erreur lors de la récupération du client:", err)
		return hairdresser, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erreur lors de la lecture du corps de la réponse:", err)
		return hairdresser, err
	}

	// Déserialiser la réponse JSON dans l'objet Customer
	err = json.Unmarshal(body, &hairdresser)
	if err != nil {
		log.Println("Erreur lors du décodage du client:", err)
		return hairdresser, err
	}

	print(hairdresser)

	return hairdresser, nil
}
