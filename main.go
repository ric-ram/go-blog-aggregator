package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ric-ram/go-blog-aggregator/internal/database"
)

type apiConfig struct {
	runPort string
	DB      *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	runPort := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error loading database")
	}

	dbQueries := database.New(db)

	apiConfig := apiConfig{
		runPort: runPort,
		DB:      dbQueries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	appRouter := chi.NewRouter()
	appRouter.Get("/users", apiConfig.handlerGetUserByApiKey)
	appRouter.Post("/users", apiConfig.handlerUsers)
	router.Mount("/v1", appRouter)

	server := &http.Server{
		Addr:    ":" + apiConfig.runPort,
		Handler: router,
	}
	log.Printf("Server starting on port %v", apiConfig.runPort)
	log.Fatal(server.ListenAndServe())

}
