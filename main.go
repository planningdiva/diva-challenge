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

func decodeJSONBody(req *http.Request, data interface{}) error {
	defer req.Body.Close()
	return json.NewDecoder(req.Body).Decode(data)
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
		r.Get("/alive", func(resp http.ResponseWriter, req *http.Request) {
			_, _, _ = jwtauth.FromContext(req.Context())
			resp.Write([]byte("alive"))
		})
		r.Post("/slack", func(resp http.ResponseWriter, req *http.Request) {
			var body struct {
				Text string `json:"text"`
			}
			err := decodeJSONBody(req, &body)
			if err != nil {
				logrus.Errorf("%v", err)
				resp.WriteHeader(http.StatusBadRequest)
				return
			}
			err = postWebhook(body.Text)
			if err != nil {
				logrus.Errorf("%v", err)
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	})

	r.Group(func(r chi.Router) {
		r.Get("/login", func(resp http.ResponseWriter, req *http.Request) {
			var token string
			expireTime := time.Now().Add(15 * time.Second)
			token = makeToken(expireTime)
			http.SetCookie(resp, &http.Cookie{
				HttpOnly: true,
				Expires:  expireTime,
				SameSite: http.SameSiteLaxMode,
				Name:     "jwt",
				Value:    token,
			})
			jwtBody := response{JWT: token}
			resp.WriteHeader(http.StatusOK)
			json.NewEncoder(resp).Encode(jwtBody)
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
