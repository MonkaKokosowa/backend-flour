package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	env "github.com/MonkaKokosowa/backend-flour/internal/env"
	"github.com/MonkaKokosowa/backend-flour/internal/mail"
	"github.com/MonkaKokosowa/backend-flour/internal/proxy"
	"github.com/rs/cors"
	"gopkg.in/gomail.v2"
)

type App struct {
	Env    *env.Environment
	Dialer *gomail.Dialer
}

type MailWebRequest struct {
	Name    string `json:"name"`
	Mail    string `json:"email"`
	Message string `json:"message"`
}

func (a *App) blogProxy(w http.ResponseWriter, r *http.Request) {
	proxy := proxy.NewProxy(a.Env.Blog.FlatnotesURL)
	proxy.ServeHTTP(w, r)
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

	mail.SendMail(a.Dialer, mail.Message{
		From:    mailRequest.Mail,
		To:      a.Env.Dialer.To,
		Subject: fmt.Sprintf("New message from: %s", mailRequest.Name),
		Body:    mailRequest.Message,
	})

}
func StartWeb(app App) {

	mux := http.NewServeMux()
	mux.HandleFunc("/mail", app.mailHandler)
	mux.HandleFunc("/api/notes/", app.blogProxy)
	mux.HandleFunc("/api/attachments/", app.blogProxy)
	mux.HandleFunc("/api/search", app.blogProxy)
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
