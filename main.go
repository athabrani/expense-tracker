package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"expense-tracker/handlers" // pastikan path ini benar
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Buat connection string dari environment variables
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 3. Buat koneksi pool ke database TERLEBIH DAHULU
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close() // Pastikan koneksi ditutup saat program berakhir
	log.Println("Successfully connected to the database!")

	// 4. SETELAH 'dbpool' ada, BARU inisialisasi handler dengan koneksi tersebut
	h := &handlers.Handler{
		DB: dbpool,
		JWTSecret: os.Getenv("JWT_SECRET"),
		CookieDomain: os.Getenv("COOKIE_DOMAIN"),
	}

	// 5. Setup router dan daftarkan rute
	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Rute publik
	router.GET("/login", h.LoginPage)
	router.POST("/login", h.Login)
	router.GET("/register", h.RegisterPage)
	router.POST("/register", h.Register)
    router.GET("/logout", h.Logout)

	// Grup rute yang dilindungi middleware
	protected := router.Group("/")
	protected.Use(h.AuthMiddleware())
	{
		protected.GET("/", h.RenderIndexPage)
		protected.POST("/expenses", h.AddExpense)
	}

log.Println("Starting server on https://localhost:8000")
err = router.Run(":8000")
if err != nil {
    log.Fatalf("Failed to run server with TLS: %v", err)
}
}
