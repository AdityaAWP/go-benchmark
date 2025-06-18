package main

import (
	"database/sql"
	"log"
	"os"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Struct untuk City - standardized field names
type City struct {
	ID          int    `json:"ID"`          // Capital ID to match Express
	Name        string `json:"Name"`        // Capital Name to match Express
	CountryCode string `json:"CountryCode"` // Capital CountryCode to match Express
	District    string `json:"District"`    // Capital District to match Express
	Population  int    `json:"Population"`  // Capital Population to match Express
}

// Struct Country - standardized field names, removed extra fields
type Country struct {
	Code       string `json:"Code"`       // Capital Code to match Express
	Name       string `json:"Name"`       // Capital Name to match Express
	Continent  string `json:"Continent"`  // Capital Continent to match Express
	Region     string `json:"Region"`     // Capital Region to match Express
	Population int    `json:"Population"` // Capital Population to match Express
}

// Database instance
var db *sql.DB

// Initialize database connection
func initDatabase() {
	var err error
	dsn := "root:root@tcp(localhost:3306)/world?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Database 'world' connected successfully")
}

// Get all cities - standardized to match Express
func getCities(c *fiber.Ctx) error {
	log.Println("Fetching all cities from database")
	
	// Same query as Express with ORDER BY ID
	query := "SELECT ID, Name, CountryCode, District, Population FROM city ORDER BY ID"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to fetch cities: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch cities",
		})
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var city City
		err := rows.Scan(
			&city.ID,
			&city.Name,
			&city.CountryCode,
			&city.District,
			&city.Population,
		)
		if err != nil {
			log.Printf("Failed to scan city data: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to scan city data",
			})
		}
		cities = append(cities, city)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over cities: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Error iterating over cities",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(cities),
		"data":    cities,
	})
}

// Get all countries - standardized to match Express
func getCountries(c *fiber.Ctx) error {
	log.Println("Fetching all countries from database")
	
	// Same query as Express: only basic fields, ORDER BY Population DESC
	query := "SELECT Code, Name, Continent, Region, Population FROM country ORDER BY Population DESC"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to fetch countries: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch countries",
		})
	}
	defer rows.Close()

	var countries []Country
	for rows.Next() {
		var country Country
		err := rows.Scan(
			&country.Code,
			&country.Name,
			&country.Continent,
			&country.Region,
			&country.Population,
		)
		if err != nil {
			log.Printf("Failed to scan country data: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to scan country data",
			})
		}
		countries = append(countries, country)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over countries: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Error iterating over countries",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(countries),
		"data":    countries,
	})
}

// Health check endpoint
func getHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "OK",
		"message": "World API is running",
	})
}

// Setup routes
func setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Get("/countries", getCountries)
	api.Get("/cities", getCities)
}

func main() {
	initDatabase()
	defer db.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", getHealth)

	// Setup API routes
	setupRoutes(app)

	// Get port from environment or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s...", port)
	log.Printf("World API is running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}