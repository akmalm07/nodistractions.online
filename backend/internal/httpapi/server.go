package httpapi

import (
	"net/http"

	"nodistractions-online/backend/internal/auth"
	"nodistractions-online/backend/internal/db"
)

type Server struct {
	users    *db.UserRepository
	sessions *auth.SessionManager
}

func NewServer(users *db.UserRepository, sessions *auth.SessionManager) *Server {
	return &Server{users: users, sessions: sessions}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", s.health)
	mux.HandleFunc("POST /api/v1/auth/register", s.register)
	mux.HandleFunc("POST /api/v1/auth/login", s.login)
	mux.HandleFunc("POST /api/v1/auth/logout", s.logout)
	mux.HandleFunc("GET /api/v1/users/me", s.me)
	mux.HandleFunc("DELETE /api/v1/users/me", s.deleteAccount)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	respond(w, http.StatusOK, map[string]string{"status": "ok"})
}
