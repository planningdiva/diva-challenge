package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	port                 = ":8080"
	webhookPlainText     = "WEBHOOK_URL"
	webhookdOldPlainText = "WEBHOOK_URL_OLD"
)

var (
	tokenAuth *jwtauth.JWTAuth //Persistent private object for creating valid JWTs
)

type response struct {
	JWT string `json:"jwt"`
}

type baseResponse struct {
	Message string `json:"message"`
}

func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

func decodeJSONBody(req *http.Request, data interface{}) error {
	defer req.Body.Close()
	return json.NewDecoder(req.Body).Decode(data)
}

// RespondJSON is the standard wrapper used in diva-server for creating and sending JSON endpoint responses
func respondJSON(resp http.ResponseWriter, status int, data interface{}) {
	resp.WriteHeader(status)
	if data != nil {
		EncodeJSONBody(resp, data)
	}
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
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte("alive"))
		})
		r.Post("/slack", func(resp http.ResponseWriter, req *http.Request) {
			var body struct {
				Text string `json:"text"`
			}
			err := decodeJSONBody(req, &body)
			if err != nil {
				logrus.Errorf("%v", err)
				respondJSON(resp, http.StatusBadRequest, nil)
				return
			}
			err = postWebhook(body.Text)
			if err != nil {
				logrus.Errorf("%v", err)
				respondJSON(resp, http.StatusInternalServerError, nil)
				return
			}
			respondJSON(resp, http.StatusOK, baseResponse{Message: "Success"})
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

	var envVars map[string]string
	envVars, err := godotenv.Read(".env")
	if err != nil {
		logrus.Warnf("No .env file.  Loading from environment: %v", err)
	}
	if envVars[webhookPlainText] == "" {
		envVars[webhookPlainText] = os.Getenv(webhookPlainText)
	}
	if envVars[webhookdOldPlainText] == "" {
		envVars[webhookdOldPlainText] = os.Getenv(webhookdOldPlainText)
	}
	if envVars[webhookPlainText] == "" { //Dont fail if the old webhook is missing
		logrus.Fatal("No environment var for slack webhook")
	}
	initSlack(envVars[webhookPlainText], envVars[webhookdOldPlainText])

	fmt.Printf("Starting server on %v\n", port)
	http.ListenAndServe(port, router())
}
