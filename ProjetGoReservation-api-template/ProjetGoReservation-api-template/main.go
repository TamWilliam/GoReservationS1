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
	Password   string `json:"password"`
}

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/customers", listCustomers)
	router.GET("/customer/new", newCustomer)
	router.POST("/customer/new", createCustomerPost)
	router.GET("/customer/edit/:id", editCustomer)
	router.POST("/customer/edit/:id", updateCustomerPost)
	router.POST("/customer/delete/:id", deleteCustomer)

	router.Run(":8080")
}

func listCustomers(c *gin.Context) {
	/* URL de l'API */
	apiUrl := "http://localhost:6060/customers"

	resp, err := http.Get(apiUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to request customers data"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read response from customers API"})
		return
	}

	/* gestion des données json */
	var customers []Customer
	err = json.Unmarshal(body, &customers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal customers response"})
		return
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":     "Customers List",
		"Customers": customers,
	})
}

/* nouveau form */
func newCustomer(c *gin.Context) {
	c.HTML(http.StatusOK, "new_customer.html", nil)
}

func createCustomerPost(c *gin.Context) {
	var form Customer
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/* conversion en json pour envoyer vers l'API */
	jsonData, err := json.Marshal(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode data to JSON"})
		return
	}

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

	if response.StatusCode != http.StatusCreated {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from API"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API did not create the customer", "response": string(responseBody)})
		return
	}

	c.Redirect(http.StatusSeeOther, "/customers")
}

func editCustomer(c *gin.Context) {

	id := c.Param("id")
	/* URL de l'API */
	apiUrl := fmt.Sprintf("http://localhost:6060/customer/%s", id)
	client := &http.Client{}
	response, err := client.Get(apiUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer data from API"})
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from API"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "API did not return customer data", "response": string(responseBody)})
		return
	}

	var customer Customer
	if err := json.NewDecoder(response.Body).Decode(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode customer data"})
		return
	}

	c.HTML(http.StatusOK, "edit_customer.html", gin.H{
		"title":    "Edit Customer",
		"customer": customer,
	})
}

func updateCustomerPost(c *gin.Context) {
	println(c.Request.GetBody())
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")

	/* URL de l'API */
	apiUrl := fmt.Sprintf("http://localhost:6060/customer/%s", id)

	req, err := http.NewRequest("DELETE", apiUrl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la requête de suppression"})
		return
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'envoi de la requête de suppression à l'API"})
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		c.Redirect(http.StatusFound, "/customers")
	} else {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la lecture de la réponse de l'API"})
			return
		}
		c.JSON(response.StatusCode, gin.H{"error": string(responseBody)})
	}
}
