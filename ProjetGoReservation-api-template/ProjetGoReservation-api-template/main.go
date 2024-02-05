package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	IDCustomer int    `json:"id_customer"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/customers", listCustomers)
	// Définissez d'autres routes pour new, edit, delete
	router.GET("/customer/new", newCustomer)
	router.POST("/customer/new", createCustomerPost)
	router.GET("/customer/edit/:id", editCustomer)
	router.POST("/customer/edit/:id", updateCustomerPost)
	// Utilisez GET pour la démo, mais DELETE est plus approprié pour les opérations réelles de suppression
	//router.GET("/customer/delete/:id", deleteCustomer)

	router.Run(":8080")
}

func listCustomers(c *gin.Context) {
	// URL de l'API pour récupérer les clients
	apiUrl := "http://localhost:6060/customers"

	// Effectuer la requête à l'API
	resp, err := http.Get(apiUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to request customers data"})
		return
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read response from customers API"})
		return
	}

	// Désérialiser les données JSON dans un slice de Customer
	var customers []Customer
	err = json.Unmarshal(body, &customers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal customers response"})
		return
	}

	// Passer les clients récupérés au template
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":     "Customers List",
		"Customers": customers,
	})
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
		responseBody, err := io.ReadAll(response.Body)
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

func editCustomer(c *gin.Context) {
	id := c.Param("id")

	// Récupérer le client par ID depuis l'API
	apiUrl := fmt.Sprintf("http://localhost:6060/customer/%s", id)
	client := &http.Client{}
	response, err := client.Get(apiUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer data from API"})
		return
	}
	defer response.Body.Close()

	// Vérifier le statut de la réponse
	if response.StatusCode != http.StatusOK {
		// Lire le corps de la réponse pour plus d'informations
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from API"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API did not return customer data", "response": string(responseBody)})
		return
	}

	/* décoder les données JSON de la réponse */
	var customer Customer
	if err := json.NewDecoder(response.Body).Decode(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode customer data"})
		return
	}

	// Passer les données au template edit_customer.html
	c.HTML(http.StatusOK, "edit_customer.html", gin.H{
		"title":    "Edit Customer",
		"customer": customer,
	})
}

func updateCustomerPost(c *gin.Context) {
	// Logique pour traiter les données du formulaire et mettre à jour le client
}
