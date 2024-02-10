package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

type Hairdresserschedule struct {
	IDHairDresserSchedule int       `json:"id_hair_dresser_schedule"`
	IDHairDresser         int       `json:"id_hair_dresser"`
	Day                   int       `json:"day"`
	StartShift            time.Time `json:"start_shift"`
	EndShift              time.Time `json:"end_shift"`
}

type Hairsalon struct {
	IDHairSalon int    `json:"id_hair_salon"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type Openinghours struct {
	IDOpeningHours int       `json:"id_opening_hours"`
	IDHairSalon    int       `json:"id_hair_salon"`
	Day            int       `json:"day"`
	Opening        time.Time `json:"opening"`
	Closing        time.Time `json:"closing"`
}

type Reservation struct {
	IDReservation   int       `json:"id_reservation"`
	IDCustomer      int       `json:"id_customer"`
	IDHairSalon     int       `json:"id_hair_salon"`
	IDHairDresser   int       `json:"id_hair_dresser"`
	ReservationDate time.Time `json:"reservation_date"`
}

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=postgres dbname=projetGoReservation password=admin sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func getCustomers(c *gin.Context) {
	rows, err := db.Query("SELECT id_customer, role, first_name, last_name, email, password FROM customers ORDER BY id_customer ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.IDCustomer, &customer.Role, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, customers)
}

func getCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer Customer
	err := db.QueryRow("SELECT id_customer, role, first_name, last_name, email, password FROM customers WHERE id_customer = $1", id).
		Scan(&customer.IDCustomer, &customer.Role, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No customer found with the specified ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func createCustomer(c *gin.Context) {
	var customer Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO customers (role, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5)",
		customer.Role, customer.FirstName, customer.LastName, customer.Email, customer.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func updateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE customers SET role = $1, first_name = $2, last_name = $3, email = $4, password = $5 WHERE id_customer = $6",
		customer.Role, customer.FirstName, customer.LastName, customer.Email, customer.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM customers WHERE id_customer = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getHairdressers(c *gin.Context) {
	rows, err := db.Query("SELECT id_hair_dresser, first_name, last_name, id_hair_salon FROM hair_dressers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var hairdressers []Hairdresser
	for rows.Next() {
		var hd Hairdresser
		if err := rows.Scan(&hd.IDHairDresser, &hd.FirstName, &hd.LastName, &hd.IDHairSalon); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hairdressers = append(hairdressers, hd)
	}

	c.JSON(http.StatusOK, hairdressers)
}

func getHairdresser(c *gin.Context) {
	id := c.Param("id")
	var hd Hairdresser
	err := db.QueryRow("SELECT id_hair_dresser, first_name, last_name, id_hair_salon FROM hair_dressers WHERE id_hair_dresser = $1", id).
		Scan(&hd.IDHairDresser, &hd.FirstName, &hd.LastName, &hd.IDHairSalon)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No hairdresser found with the specified ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hd)
}

func createHairdresser(c *gin.Context) {
	var hd Hairdresser
	if err := c.ShouldBindJSON(&hd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO hair_dressers (first_name, last_name, id_hair_salon) VALUES ($1, $2, $3)", hd.FirstName, hd.LastName, hd.IDHairSalon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func updateHairdresser(c *gin.Context) {
	id := c.Param("id")
	var hd Hairdresser
	if err := c.ShouldBindJSON(&hd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE hair_dressers SET first_name = $1, last_name = $2, id_hair_salon = $3 WHERE id_hair_dresser = $4", hd.FirstName, hd.LastName, hd.IDHairSalon, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteHairdresser(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM hair_dressers WHERE id_hair_dresser = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getAllHairdresserSchedules(c *gin.Context) {
	rows, err := db.Query("SELECT id_hair_dresser_schedule, id_hair_dresser, day, start_shift, end_shift FROM hair_dresser_schedules")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var schedules []Hairdresserschedule
	for rows.Next() {
		var schedule Hairdresserschedule
		if err := rows.Scan(&schedule.IDHairDresserSchedule, &schedule.IDHairDresser, &schedule.Day, &schedule.StartShift, &schedule.EndShift); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, schedules)
}

func getHairdresserSchedule(c *gin.Context) {
	id := c.Param("id")
	var schedule Hairdresserschedule
	err := db.QueryRow("SELECT id_hair_dresser_schedule, id_hair_dresser, day, start_shift, end_shift FROM hair_dresser_schedules WHERE id_hair_dresser_schedule = $1", id).
		Scan(&schedule.IDHairDresserSchedule, &schedule.IDHairDresser, &schedule.Day, &schedule.StartShift, &schedule.EndShift)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No schedule found with the specified ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func createHairdresserSchedule(c *gin.Context) {
	var schedule Hairdresserschedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO hair_dresser_schedules (id_hair_dresser, day, start_shift, end_shift) VALUES ($1, $2, $3, $4)",
		schedule.IDHairDresser, schedule.Day, schedule.StartShift, schedule.EndShift)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func updateHairdresserSchedule(c *gin.Context) {
	id := c.Param("id")
	var schedule Hairdresserschedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE hair_dresser_schedules SET id_hair_dresser = $1, day = $2, start_shift = $3, end_shift = $4 WHERE id_hair_dresser_schedule = $5",
		schedule.IDHairDresser, schedule.Day, schedule.StartShift, schedule.EndShift, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteHairdresserSchedule(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM hair_dresser_schedules WHERE id_hair_dresser_schedule = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getHairsalons(c *gin.Context) {
	rows, err := db.Query("SELECT id_hair_salon, name, address, email, password FROM hair_salons")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var salons []Hairsalon
	for rows.Next() {
		var salon Hairsalon
		if err := rows.Scan(&salon.IDHairSalon, &salon.Name, &salon.Address, &salon.Email, &salon.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		salons = append(salons, salon)
	}

	c.JSON(http.StatusOK, salons)
}

func getHairsalon(c *gin.Context) {
	id := c.Param("id")
	var salon Hairsalon
	err := db.QueryRow("SELECT id_hair_salon, name, address, email, password FROM hair_salons WHERE id_hair_salon = $1", id).
		Scan(&salon.IDHairSalon, &salon.Name, &salon.Address, &salon.Email, &salon.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No salon found with the specified ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, salon)
}

func createHairsalon(c *gin.Context) {
	var salon Hairsalon
	if err := c.ShouldBindJSON(&salon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO hair_salons (name, address, email, password) VALUES ($1, $2, $3, $4)",
		salon.Name, salon.Address, salon.Email, salon.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func updateHairsalon(c *gin.Context) {
	id := c.Param("id")
	var salon Hairsalon
	if err := c.ShouldBindJSON(&salon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE hair_salons SET name = $1, address = $2, email = $3, password = $4 WHERE id_hair_salon = $5",
		salon.Name, salon.Address, salon.Email, salon.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteHairsalon(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM hair_salons WHERE id_hair_salon = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getAllOpeningHours(c *gin.Context) {
	rows, err := db.Query("SELECT id_opening_hours, id_hair_salon, day, opening, closing FROM opening_hours")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var hoursList []Openinghours
	for rows.Next() {
		var hours Openinghours
		if err := rows.Scan(&hours.IDOpeningHours, &hours.IDHairSalon, &hours.Day, &hours.Opening, &hours.Closing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		hoursList = append(hoursList, hours)
	}

	c.JSON(http.StatusOK, hoursList)
}

func getOpeningHours(c *gin.Context) {
	id := c.Param("id")
	var hours Openinghours
	err := db.QueryRow("SELECT id_opening_hours, id_hair_salon, day, opening, closing FROM opening_hours WHERE id_opening_hours = $1", id).
		Scan(&hours.IDOpeningHours, &hours.IDHairSalon, &hours.Day, &hours.Opening, &hours.Closing)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No opening hours found with the specified ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hours)
}

func createOpeningHours(c *gin.Context) {
	var hours Openinghours
	if err := c.ShouldBindJSON(&hours); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO opening_hours (id_hair_salon, day, opening, closing) VALUES ($1, $2, $3, $4)",
		hours.IDHairSalon, hours.Day, hours.Opening, hours.Closing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func updateOpeningHours(c *gin.Context) {
	id := c.Param("id")
	var hours Openinghours
	if err := c.ShouldBindJSON(&hours); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE opening_hours SET id_hair_salon = $1, day = $2, opening = $3, closing = $4 WHERE id_opening_hours = $5",
		hours.IDHairSalon, hours.Day, hours.Opening, hours.Closing, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteOpeningHours(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM opening_hours WHERE id_opening_hours = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func getAllReservations(c *gin.Context) {
	rows, err := db.Query("SELECT id_reservation, id_customer, id_hair_salon, id_hair_dresser, reservation_date FROM reservations")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var reservation Reservation
		if err := rows.Scan(&reservation.IDReservation, &reservation.IDCustomer, &reservation.IDHairSalon, &reservation.IDHairDresser, &reservation.ReservationDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reservations = append(reservations, reservation)
	}

	c.JSON(http.StatusOK, reservations)
}

func getReservation(c *gin.Context) {
	id := c.Param("id")
	var reservation Reservation
	err := db.QueryRow("SELECT id_reservation, id_customer, id_hair_salon, id_hair_dresser, reservation_date FROM reservations WHERE id_reservation = $1", id).
		Scan(&reservation.IDReservation, &reservation.IDCustomer, &reservation.IDHairSalon, &reservation.IDHairDresser, &reservation.ReservationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No reservation found with the specified ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func createReservation(c *gin.Context) {
	var reservation Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO reservations (id_customer, id_hair_salon, id_hair_dresser, reservation_date) VALUES ($1, $2, $3, $4)",
		reservation.IDCustomer, reservation.IDHairSalon, reservation.IDHairDresser, reservation.ReservationDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func updateReservation(c *gin.Context) {
	id := c.Param("id")
	var reservation Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE reservations SET id_customer = $1, id_hair_salon = $2, id_hair_dresser = $3, reservation_date = $4 WHERE id_reservation = $5",
		reservation.IDCustomer, reservation.IDHairSalon, reservation.IDHairDresser, reservation.ReservationDate, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteReservation(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM reservations WHERE id_reservation = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func main() {
	initDB()
	defer db.Close()

	router := gin.Default()

	router.POST("/customer", createCustomer)
	router.GET("/customers", getCustomers)
	router.GET("/customer/:id", getCustomer)
	router.PUT("/customer/:id", updateCustomer)
	router.DELETE("/customer/:id", deleteCustomer)

	router.POST("/hairdresser", createHairdresser)       // Créer un coiffeur
	router.GET("/hairdressers", getHairdressers)         // Lire tous les coiffeurs
	router.GET("/hairdresser/:id", getHairdresser)       // Lire un coiffeur par ID
	router.PUT("/hairdresser/:id", updateHairdresser)    // Mettre à jour un coiffeur
	router.DELETE("/hairdresser/:id", deleteHairdresser) // Supprimer un coiffeur

	router.POST("/hairdresserschedule", createHairdresserSchedule)
	router.GET("/hairdresserschedules", getAllHairdresserSchedules)
	router.GET("/hairdresserschedule/:id", getHairdresserSchedule)
	router.PUT("/hairdresserschedule/:id", updateHairdresserSchedule)
	router.DELETE("/hairdresserschedule/:id", deleteHairdresserSchedule)

	router.POST("/hairsalon", createHairsalon)       // Créer un salon de coiffure
	router.GET("/hairsalons", getHairsalons)         // Lire tous les salons de coiffure
	router.GET("/hairsalon/:id", getHairsalon)       // Lire un salon de coiffure par ID
	router.PUT("/hairsalon/:id", updateHairsalon)    // Mettre à jour un salon de coiffure
	router.DELETE("/hairsalon/:id", deleteHairsalon) // Supprimer un salon de coiffure

	router.POST("/openinghours", createOpeningHours)       // Créer des heures d'ouverture
	router.GET("/openinghours", getAllOpeningHours)        // Lire toutes les heures d'ouverture
	router.GET("/openinghours/:id", getOpeningHours)       // Lire des heures d'ouverture par ID
	router.PUT("/openinghours/:id", updateOpeningHours)    // Mettre à jour des heures d'ouverture
	router.DELETE("/openinghours/:id", deleteOpeningHours) // Supprimer des heures d'ouverture

	router.POST("/reservation", createReservation)       // Créer une réservation
	router.GET("/reservations", getAllReservations)      // Lire toutes les réservations
	router.GET("/reservation/:id", getReservation)       // Lire une réservation par ID
	router.PUT("/reservation/:id", updateReservation)    // Mettre à jour une réservation
	router.DELETE("/reservation/:id", deleteReservation) // Supprimer une réservation

	router.Run(":6060")
}
