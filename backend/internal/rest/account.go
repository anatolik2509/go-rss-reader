package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"rss-reader-backend/internal/service/account"
	"rss-reader-backend/internal/session"

	"github.com/go-chi/chi/v5"
)

type SignUpRequest struct {
	Login    string
	Password string
}

type SignInRequest struct {
	Login    string
	Password string
}

type HttpHandler struct {
	accountManager account.Manager
	sessionManager session.Manager
}

func (h *HttpHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Request parsing error", http.StatusUnprocessableEntity)
		return
	}
	var accountDetails = account.Request{Login: payload.Login, Password: payload.Password}
	id, err := h.accountManager.AddAccount(r.Context(), accountDetails)
	if err != nil {
		if errors.Is(err, account.ErrSuchLoginAlreadyExists) {
			http.Error(w, account.ErrSuchLoginAlreadyExists.Error(), http.StatusConflict)
			return
		}
		slog.ErrorContext(r.Context(), "Processing SignUp request failed", "error", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	var userData = session.UserData{UserId: id}
	err = h.sessionManager.CreateNewSession(userData, r, w)
	if err != nil {
		slog.ErrorContext(r.Context(), "Processing SignUp request failed", "error", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = http.NoBody.WriteTo(w)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error writing response", "error", err)
		return
	}
}

func (h *HttpHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload SignInRequest
	json.NewDecoder(r.Body).Decode(&payload)
	var accountDetails account.Request = account.Request{payload.Login, payload.Password}
	id, ok, err := h.accountManager.VerifyAccount(r.Context(), accountDetails)
	if err != nil {
		slog.ErrorContext(r.Context(), "Processing SignIn request failed", "error", err)
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}
	if !ok {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}
	var session = session.UserData{UserId: id}
	err = h.sessionManager.CreateNewSession(session, r, w)
	if err != nil {
		slog.ErrorContext(r.Context(), "Processing SignIn request failed", "error", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = http.NoBody.WriteTo(w)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error writing response", "error", err)
		return
	}
}

func NewHttpHandler(accountManager account.Manager, sessionManager session.Manager) *HttpHandler {
	return &HttpHandler{
		accountManager,
		sessionManager,
	}
}

func MustConfigureRouter(accountManager account.Manager, sessionManager session.Manager) chi.Router {
	httpHandler := NewHttpHandler(accountManager, sessionManager)
	router := chi.NewRouter()
	router.Post("/signUp", httpHandler.SignUp)
	router.Post("/signIn", httpHandler.SignIn)
	return router
}
