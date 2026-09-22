package httpapi

import "net/http"

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	id, ok := s.sessions.UserID(r)
	if !ok {
		errorResponse(w, http.StatusUnauthorized, "authentication required")
		return
	}

	user, err := s.users.FindByID(r.Context(), id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "user not found")
		return
	}
	respond(w, http.StatusOK, user)
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	id, ok := s.sessions.UserID(r)
	if !ok {
		errorResponse(w, http.StatusUnauthorized, "authentication required")
		return
	}

	if err := s.users.Delete(r.Context(), id); err != nil {
		errorResponse(w, http.StatusNotFound, "user not found")
		return
	}
	s.sessions.Clear(w)
	w.WriteHeader(http.StatusNoContent)
}
