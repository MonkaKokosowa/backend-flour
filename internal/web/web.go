package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/rs/zerolog/log"

	env "github.com/MonkaKokosowa/backend-flour/internal/env"
	"github.com/MonkaKokosowa/backend-flour/internal/mail"
	"github.com/rs/cors"
	gomail "github.com/wneessen/go-mail"
)

type App struct {
	Env            *env.Environment
	Client         *gomail.Client
	FlatnotesProxy *httputil.ReverseProxy
}

type MailWebRequest struct {
	Name    string `json:"name"`
	Mail    string `json:"email"`
	Message string `json:"message"`
}

func (a *App) blogProxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Error().Msg("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	a.FlatnotesProxy.ServeHTTP(w, r)
}

func (a *App) mailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var mailRequest MailWebRequest

	err := json.NewDecoder(r.Body).Decode(&mailRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received user: %+v\n", mailRequest)

	mail.SendMail(a.Client, mail.Message{
		From: mail.Address{
			Name:  mail.LimitCharacters(mailRequest.Name, 80),
			Email: mail.LimitCharacters(mailRequest.Mail, 80),
		},
		To:      a.Env.Dialer.To,
		Subject: fmt.Sprintf("New message from: %s", mail.LimitCharacters(mailRequest.Name, 50)),
		Body:    mail.LimitCharacters(mailRequest.Message, 500),
		User:    a.Env.Dialer.Username,
	})

}
func StartWeb(app App) {

	mux := http.NewServeMux()
	mux.HandleFunc("/mail", app.mailHandler)
	mux.HandleFunc("/api/notes/", app.blogProxyHandler)
	mux.HandleFunc("/api/attachments/", app.blogProxyHandler)
	mux.HandleFunc("/api/search", app.blogProxyHandler)
	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(app.Env.WebServer.AllowedOrigins, ","), // Adjust this to your frontend's origin
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Info().Msg(fmt.Sprintf("Starting server on port %d...", app.Env.WebServer.Port))
	http.ListenAndServe(fmt.Sprintf(":%d", app.Env.WebServer.Port), handler)

}
