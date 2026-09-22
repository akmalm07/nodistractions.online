package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SessionManager struct {
	secret   []byte
	issuer   string
	audience string
}

func NewSessionManager(secret, issuer, audience string) *SessionManager {
	return &SessionManager{
		secret:   []byte(secret),
		issuer:   issuer,
		audience: audience,
	}
}

func (s *SessionManager) Set(w http.ResponseWriter, userID string) error {
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": s.issuer,
		"aud": s.audience,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   os.Getenv("NODE_ENV") == "production",
		MaxAge:   24 * 60 * 60,
	})
	return nil
}

func (s *SessionManager) UserID(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", false
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		cookie.Value,
		claims,
		func(_ *jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience(s.audience),
	)
	if err != nil || !token.Valid {
		return "", false
	}
	id, ok := claims["sub"].(string)
	return id, ok
}

func (s *SessionManager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
