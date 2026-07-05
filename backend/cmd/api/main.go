package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/AntonioMartino67/showio/backend/internal/database"
	"github.com/AntonioMartino67/showio/backend/internal/handlers"
	"github.com/AntonioMartino67/showio/backend/internal/auth"
)

func main() {
	// Carica variabili d'ambiente dal file .env (se presente)
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, uso variabili d'ambiente di sistema")
	}

	database.Connect()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// Middleware di base: logging delle richieste e recovery da panic
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check: serve per verificare che il server sia vivo
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})
	
	r.Post("/register", handlers.RegisterHandler)
	r.Post("/login", handlers.LoginHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(auth.RequireAuth)
		protected.Get("/me", handlers.MeHandler)
	})

	r.Get("/search", handlers.SearchHandler)

	log.Printf("Server in ascolto sulla porta %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}