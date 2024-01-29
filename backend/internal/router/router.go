package router

import (
	"github.com/go-chi/chi"
	//"github.com/go-chi/cors"
	"github.com/kiloMIA/mediscan/backend/internal/handlers"
)

func NewRouter(userh *handlers.UserHandler, chatHandler *handlers.ChatHandler) *chi.Mux {
	r := chi.NewRouter()

	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins: []string{"http://localhost:3000"},
	// 	AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// }))

	// User routes
	r.Post("/register", userh.Register)
	r.Post("/login", userh.Login)
	r.Post("/logout", userh.Logout)
	r.Post("/book-checkup", userh.BookCheckUp)
	r.With(userh.RequireAuth).Get("/user", userh.GetUser)
	r.HandleFunc("/ws", chatHandler.HandleChat)

	
	return r
}