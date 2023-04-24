package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
)

const port = ":8080"

var (
	tokenAuth *jwtauth.JWTAuth //Persistent private object for creating valid JWTs
)

type response struct {
	JWT string `json:"jwt"`
}

func initJWT() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func makeToken(expireTime time.Time) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"exp": expireTime})
	return tokenString
}

func router() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/alive", func(w http.ResponseWriter, r *http.Request) {
			_, _, _ = jwtauth.FromContext(r.Context())
			w.Write([]byte("alive"))
		})
	})

	r.Group(func(r chi.Router) {
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			var token string
			expireTime := time.Now().Add(15 * time.Second)
			token = makeToken(expireTime)
			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				Expires:  expireTime,
				SameSite: http.SameSiteLaxMode,
				Name:     "jwt",
				Value:    token,
			})
			jwtBody := response{JWT: token}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(jwtBody)
		})
	})
	return r
}

func main() {
	logrus.SetReportCaller(true)
	initJWT()
	fmt.Printf("Starting server on %v\n", port)
	http.ListenAndServe(port, router())
}
