package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kuromii5/time-tracker/internal/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env loading error: %v", err)
	}

	r := gin.Default()

	r.GET("/info", func(c *gin.Context) {
		passportSerie := c.Query("passportSerie")
		passportNumber := c.Query("passportNumber")

		if passportSerie == "" || passportNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
			return
		}

		// response example
		c.JSON(http.StatusOK, models.People{
			Name:       "Иван",
			Surname:    "Иванов",
			Patronymic: "Иванович",
			Address:    "г. Москва, ул. Ленина, д. 5, кв. 1",
		})
	})

	port := os.Getenv("EXTERNAL_API_PORT")

	fmt.Printf("external server is running on port: %s\n", port)

	r.Run(fmt.Sprintf(":%s", port))
}
