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
	Id        int    `json:"id_customer"`
	Role      int    `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Hairdresser struct {
	Id           int    `json:"id_hairdresser"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Id_hairsalon int    `json:"id_hairsalon"`
}

type Hairdresserschedule struct {
	Id             int       `json:"id_hairdresserschedule"`
	Id_hairdresser int       `json:"id_hairdresser"`
	Day            int       `json:"day"`
	Startshift     time.Time `json:"startshift"`
	Endshift       time.Time `json:"endshift"`
}

type Hairsalon struct {
	Id       int    `json:"id_hairsalon"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Openinghours struct {
	Id           int       `json:"id_openinghours"`
	Id_hairsalon int       `json:"id_hairsalon"`
	Day          int       `json:"day"`
	Opening      time.Time `json:"opening"`
	Closing      time.Time `json:"closing"`
}

type Reservation struct {
	Id               int       `json:"id_reservation"`
	Id_customer      int       `json:"id_customer"`
	Id_hairsalon     int       `json:"id_hairsalon"`
	Id_hairdresser   int       `json:"id_hairdresser"`
	Reservation_date time.Time `json:"reservation_date"`
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
	rows, err := db.Query("SELECT id_customer, role, firstname, lastname, email, password FROM customer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	customers := make([]Customer, 0)
	for rows.Next() {
		var cust Customer
		if err := rows.Scan(&cust.Id, &cust.Role, &cust.Firstname, &cust.Lastname, &cust.Email, &cust.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		customers = append(customers, cust)
	}

	c.JSON(http.StatusOK, customers)
}

func getCustomer(c *gin.Context) {
	id := c.Param("id")
	var cust Customer
	err := db.QueryRow("SELECT * FROM customer WHERE id_customer = $1", id).Scan(&cust.Id, &cust.Role, &cust.Firstname, &cust.Lastname, &cust.Email, &cust.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aucun client trouvé avec l'ID spécifié"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cust)
}

func updateCustomer(c *gin.Context) {
	id := c.Param("id")
	var cust Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE customer SET role = $1, firstname = $2, lastname = $3, email = $4, password = $5 WHERE id_customer = $6", cust.Role, cust.Firstname, cust.Lastname, cust.Email, cust.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteCustomer(c *gin.Context) {
	id := c.Param("id")

	log.Printf(id)

	result, err := db.Exec("DELETE FROM customer WHERE id_customer = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucun client trouvé avec l'ID spécifié"})
		return
	}

	c.Status(http.StatusOK)
}

func createCustomer(c *gin.Context) {
	var customer Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO customer (role, firstname, lastname, email, password) VALUES ($1, $2, $3, $4, $5)", customer.Role, customer.Firstname, customer.Lastname, customer.Email, customer.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func getHairdressers(c *gin.Context) {
	rows, err := db.Query("SELECT id_hairdresser, firstname, lastname, id_hairsalon FROM hairdresser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	hairdressers := make([]Hairdresser, 0)
	for rows.Next() {
		var hd Hairdresser
		if err := rows.Scan(&hd.Id, &hd.Firstname, &hd.Lastname, &hd.Id_hairsalon); err != nil {
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
	err := db.QueryRow("SELECT id_hairdresser, firstname, lastname, id_hairsalon FROM hairdresser WHERE id_hairdresser = $1", id).Scan(&hd.Id, &hd.Firstname, &hd.Lastname, &hd.Id_hairsalon)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aucun coiffeur trouvé avec l'ID spécifié"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hd)
}

func updateHairdresser(c *gin.Context) {
	id := c.Param("id")
	var hd Hairdresser
	if err := c.ShouldBindJSON(&hd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE hairdresser SET firstname = $1, lastname = $2, id_hairsalon = $3 WHERE id_hairdresser = $4", hd.Firstname, hd.Lastname, hd.Id_hairsalon, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteHairdresser(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM hairdresser WHERE id_hairdresser = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucun coiffeur trouvé avec l'ID spécifié"})
		return
	}

	c.Status(http.StatusOK)
}

func createHairdresser(c *gin.Context) {
	var customer Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO customer (role, firstname, lastname, email, password) VALUES ($1, $2, $3, $4, $5)", customer.Role, customer.Firstname, customer.Lastname, customer.Email, customer.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func getAllHairdresserSchedules(c *gin.Context) {
	rows, err := db.Query("SELECT id_hairdresserschedule, id_hairdresser, day, startshift, endshift FROM hairdresserschedule")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var schedules []Hairdresserschedule
	for rows.Next() {
		var schedule Hairdresserschedule
		if err := rows.Scan(&schedule.Id, &schedule.Id_hairdresser, &schedule.Day, &schedule.Startshift, &schedule.Endshift); err != nil {
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
	err := db.QueryRow("SELECT id_hairdresserschedule, id_hairdresser, day, startshift, endshift FROM hairdresserschedule WHERE id_hairdresserschedule = $1", id).Scan(&schedule.Id, &schedule.Id_hairdresser, &schedule.Day, &schedule.Startshift, &schedule.Endshift)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aucun horaire trouvé avec l'ID spécifié"})
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

	_, err := db.Exec("INSERT INTO hairdresserschedule (id_hairdresser, day, startshift, endshift) VALUES ($1, $2, $3, $4)", schedule.Id_hairdresser, schedule.Day, schedule.Startshift, schedule.Endshift)
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

	_, err := db.Exec("UPDATE hairdresserschedule SET id_hairdresser = $1, day = $2, startshift = $3, endshift = $4 WHERE id_hairdresserschedule = $5", schedule.Id_hairdresser, schedule.Day, schedule.Startshift, schedule.Endshift, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteHairdresserSchedule(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM hairdresserschedule WHERE id_hairdresserschedule = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucun horaire trouvé avec l'ID spécifié"})
		return
	}

	c.Status(http.StatusOK)
}

func getHairsalons(c *gin.Context) {
	rows, err := db.Query("SELECT id_hairsalon, name, address, email, password FROM hairsalon")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	salons := make([]Hairsalon, 0)
	for rows.Next() {
		var salon Hairsalon
		if err := rows.Scan(&salon.Id, &salon.Name, &salon.Address, &salon.Email, &salon.Password); err != nil {
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
	err := db.QueryRow("SELECT id_hairsalon, name, address, email, password FROM hairsalon WHERE id_hairsalon = $1", id).Scan(&salon.Id, &salon.Name, &salon.Address, &salon.Email, &salon.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aucun salon trouvé avec l'ID spécifié"})
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

	_, err := db.Exec("INSERT INTO hairsalon (name, address, email, password) VALUES ($1, $2, $3, $4)", salon.Name, salon.Address, salon.Email, salon.Password)
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

	_, err := db.Exec("UPDATE hairsalon SET name = $1, address = $2, email = $3, password = $4 WHERE id_hairsalon = $5", salon.Name, salon.Address, salon.Email, salon.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteHairsalon(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM hairsalon WHERE id_hairsalon = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucun salon trouvé avec l'ID spécifié"})
		return
	}

	c.Status(http.StatusOK)
}

func getAllOpeningHours(c *gin.Context) {
	rows, err := db.Query("SELECT id_openinghours, id_hairsalon, day, opening, closing FROM openinghours")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var hoursList []Openinghours
	for rows.Next() {
		var hours Openinghours
		if err := rows.Scan(&hours.Id, &hours.Id_hairsalon, &hours.Day, &hours.Opening, &hours.Closing); err != nil {
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
	err := db.QueryRow("SELECT id_openinghours, id_hairsalon, day, opening, closing FROM openinghours WHERE id_openinghours = $1", id).Scan(&hours.Id, &hours.Id_hairsalon, &hours.Day, &hours.Opening, &hours.Closing)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aucune heure d'ouverture trouvée avec l'ID spécifié"})
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

	_, err := db.Exec("INSERT INTO openinghours (id_hairsalon, day, opening, closing) VALUES ($1, $2, $3, $4)", hours.Id_hairsalon, hours.Day, hours.Opening, hours.Closing)
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

	_, err := db.Exec("UPDATE openinghours SET id_hairsalon = $1, day = $2, opening = $3, closing = $4 WHERE id_openinghours = $5", hours.Id_hairsalon, hours.Day, hours.Opening, hours.Closing, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteOpeningHours(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM openinghours WHERE id_openinghours = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucune heure d'ouverture trouvée avec l'ID spécifié"})
		return
	}

	c.Status(http.StatusOK)
}

func getAllReservations(c *gin.Context) {
	rows, err := db.Query("SELECT id_reservation, id_customer, id_hairsalon, id_hairdresser, reservation_date FROM reservation")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var reservation Reservation
		if err := rows.Scan(&reservation.Id, &reservation.Id_customer, &reservation.Id_hairsalon, &reservation.Id_hairdresser, &reservation.Reservation_date); err != nil {
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
	err := db.QueryRow("SELECT id_reservation, id_customer, id_hairsalon, id_hairdresser, reservation_date FROM reservation WHERE id_reservation = $1", id).Scan(&reservation.Id, &reservation.Id_customer, &reservation.Id_hairsalon, &reservation.Id_hairdresser, &reservation.Reservation_date)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aucune réservation trouvée avec l'ID spécifié"})
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

	_, err := db.Exec("INSERT INTO reservation (id_customer, id_hairsalon, id_hairdresser, reservation_date) VALUES ($1, $2, $3, $4)", reservation.Id_customer, reservation.Id_hairsalon, reservation.Id_hairdresser, reservation.Reservation_date)
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

	_, err := db.Exec("UPDATE reservation SET id_customer = $1, id_hairsalon = $2, id_hairdresser = $3, reservation_date = $4 WHERE id_reservation = $5", reservation.Id_customer, reservation.Id_hairsalon, reservation.Id_hairdresser, reservation.Reservation_date, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func deleteReservation(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM reservation WHERE id_reservation = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucune réservation trouvée avec l'ID spécifié"})
		return
	}

	c.Status(http.StatusOK)
}

func main() {
	initDB()
	defer db.Close()

	router := gin.Default()

	router.GET("/customers", getCustomers)
	router.POST("/customer", createCustomer)
	router.GET("/customer/:id", getCustomer)
	router.PUT("/customer/:id", updateCustomer)
	router.DELETE("/customer/:id", deleteCustomer)

	router.GET("/hairdressers", getHairdressers)
	router.POST("/hairdresser", createHairdresser)
	router.GET("/hairdresser/:id", getHairdresser)
	router.PUT("/hairdresser/:id", updateHairdresser)
	router.DELETE("/hairdresser/:id", deleteHairdresser)

	router.POST("/hairdresserschedule", createHairdresserSchedule)
	router.GET("/hairdresserschedules", getAllHairdresserSchedules)
	router.GET("/hairdresserschedule/:id", getHairdresserSchedule)
	router.PUT("/hairdresserschedule/:id", updateHairdresserSchedule)
	router.DELETE("/hairdresserschedule/:id", deleteHairdresserSchedule)

	router.POST("/hairsalon", createHairsalon)
	router.GET("/hairsalons", getHairsalons)
	router.GET("/hairsalon/:id", getHairsalon)
	router.PUT("/hairsalon/:id", updateHairsalon)
	router.DELETE("/hairsalon/:id", deleteHairsalon)

	router.POST("/openinghours", createOpeningHours)
	router.GET("/openinghours", getAllOpeningHours)
	router.GET("/openinghours/:id", getOpeningHours)
	router.PUT("/openinghours/:id", updateOpeningHours)
	router.DELETE("/openinghours/:id", deleteOpeningHours)

	router.POST("/reservation", createReservation)
	router.GET("/reservations", getAllReservations)
	router.GET("/reservation/:id", getReservation)
	router.PUT("/reservation/:id", updateReservation)
	router.DELETE("/reservation/:id", deleteReservation)

	router.Run(":8000")
}
