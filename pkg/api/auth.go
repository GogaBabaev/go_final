package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type SignInRequest struct {
	Password string `json:"password"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {

	var req SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	password := os.Getenv("TODO_PASSWORD")

	if req.Password != password {

		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Неверный пароль",
		})
		return

	}

	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	claims := jwt.MapClaims{
		"password_hash": passwordHash,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(password))

	if err != nil {

		log.Printf("Ошибка создания JWT: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка сервера",
		})
		return

	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")

		if password == "" {
			next(w, r)
			return
		}

		var tokenString string

		cookie, err := r.Cookie("token")

		if err == nil {
			tokenString = cookie.Value
		}

		valid := false

		if tokenString != "" {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				return []byte(password), nil
			})

			if err == nil && token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)

				if ok {
					tokenHash, ok := claims["password_hash"].(string)

					hash := sha256.Sum256([]byte(password))
					currentHash := hex.EncodeToString(hash[:])

					if ok && tokenHash == currentHash {
						valid = true
					}
				}
			}
		}
		if !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
