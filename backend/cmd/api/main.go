package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/AntonioMartino67/showio/backend/internal/database"
	"github.com/AntonioMartino67/showio/backend/internal/handlers"
	"github.com/AntonioMartino67/showio/backend/internal/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nessun file .env trovato, uso variabili d'ambiente di sistema")
	}

	database.Connect()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "https://*.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Cron-Secret"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	r.Post("/register", handlers.RegisterHandler)
	r.Post("/login", handlers.LoginHandler)
	r.Post("/verify-otp", handlers.VerifyOTPHandler)
	r.Post("/resend-otp", handlers.ResendOTPHandler)
	r.Post("/auth/google", handlers.GoogleLoginHandler)
	r.Post("/cron/sync-all", handlers.SyncAllHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(auth.RequireAuth)
		protected.Get("/me", handlers.MeHandler)
		protected.Put("/me/avatar", handlers.UpdateAvatarHandler)
		protected.Put("/me/password", handlers.ChangePasswordHandler)
		protected.Post("/me/google", handlers.LinkGoogleHandler)
		protected.Delete("/me/google", handlers.UnlinkGoogleHandler)
		protected.Put("/me/notifications", handlers.UpdateNotificationsHandler)
		protected.Post("/progress", handlers.AddProgressHandler)
		protected.Get("/progress", handlers.ListProgressHandler)
		protected.Put("/progress/{mediaItemId}/episode", handlers.UpdateEpisodeHandler)
		protected.Delete("/progress/{mediaItemId}", handlers.RemoveProgressHandler)
		protected.Put("/progress/{mediaItemId}/rating", handlers.UpdateRatingHandler)
		protected.Get("/media/{mediaItemId}", handlers.MediaDetailHandler)
		protected.Get("/calendar", handlers.CalendarHandler)
		protected.Get("/stats", handlers.StatsHandler)
		protected.Get("/tags", handlers.ListTagsHandler)
		protected.Post("/tags", handlers.CreateTagHandler)
		protected.Delete("/tags/{tagId}", handlers.DeleteTagHandler)
		protected.Post("/progress/{mediaItemId}/tags", handlers.AssignTagHandler)
		protected.Delete("/progress/{mediaItemId}/tags/{tagId}", handlers.RemoveTagFromProgressHandler)
	})

	r.Get("/search", handlers.SearchHandler)
	r.Get("/trending", handlers.TrendingHandler)

	log.Printf("Server in ascolto sulla porta %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}