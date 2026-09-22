package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"nodistractions-online/backend/internal/db"
	"nodistractions-online/backend/internal/domain"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var input domain.Registration
	if !decodeJSON(w, r, &input) {
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if !strings.Contains(input.Email, "@") || len(input.Password) < 12 {
		errorResponse(w, http.StatusBadRequest, "a valid email and 12-character password are required")
		return
	}
	if input.Name == "" {
		errorResponse(w, http.StatusBadRequest, "a name is required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not secure password")
		return
	}

	user := domain.User{
		ID:           uuid.NewString(),
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: string(hash),
	}
	if err := s.users.Create(r.Context(), &user); err != nil {
		if errors.Is(err, db.ErrEmailTaken) {
			errorResponse(w, http.StatusConflict, "email already registered")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "database error")
		return
	}

	if err := s.sessions.Set(w, user.ID); err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not create session")
		return
	}
	respond(w, http.StatusCreated, map[string]domain.User{"user": user})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var input domain.Login
	if !decodeJSON(w, r, &input) {
		return
	}

	user, err := s.users.FindByEmail(r.Context(), strings.ToLower(strings.TrimSpace(input.Email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		errorResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err := s.sessions.Set(w, user.ID); err != nil {
		errorResponse(w, http.StatusInternalServerError, "could not create session")
		return
	}
	respond(w, http.StatusOK, map[string]domain.User{"user": user})
}

func (s *Server) logout(w http.ResponseWriter, _ *http.Request) {
	s.sessions.Clear(w)
	w.WriteHeader(http.StatusNoContent)
}
