package account

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"rss-reader-backend/internal/session"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
	accountManager AccountManager
	sessionManager session.Manager
	logger         *slog.Logger
}

func (h *HttpHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SignUpRequest
	json.NewDecoder(r.Body).Decode(&payload)
	var accountDetails AccountDetails = AccountDetails{payload.Login, payload.Password}
	id, err := h.accountManager.AddAccount(r.Context(), accountDetails)
	if err != nil {
		if errors.Is(err, ErrSuchLoginAlreadyExists) {
			http.Error(w, ErrSuchLoginAlreadyExists.Error(), http.StatusConflict)
			return
		}
		h.logger.ErrorContext(r.Context(), "Processing SignUp request failed", "error", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	var session session.UserData = session.UserData{UserId: id}
	err = h.sessionManager.CreateNewSession(session, r, w)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Processing SignUp request failed", "error", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	http.NoBody.WriteTo(w)
}

func (h *HttpHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload SignInRequest
	json.NewDecoder(r.Body).Decode(&payload)
	var accountDetails AccountDetails = AccountDetails{payload.Login, payload.Password}
	id, ok, err := h.accountManager.VerifyAccount(r.Context(), accountDetails)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Processing SignIn request failed", "error", err)
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}
	if !ok {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}
	var session session.UserData = session.UserData{UserId: id}
	err = h.sessionManager.CreateNewSession(session, r, w)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Processing SignIn request failed", "error", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	http.NoBody.WriteTo(w)
}

func NewHttpHandler(accountManager AccountManager, sessionManager session.Manager, logger *slog.Logger) *HttpHandler {
	return &HttpHandler{
		accountManager,
		sessionManager,
		logger,
	}
}

func MustConfigureRouter(pool *pgxpool.Pool, redisClient *redis.Client, accountManager AccountManager, sessionManager session.Manager, logger *slog.Logger) chi.Router {
	httpHandler := NewHttpHandler(accountManager, sessionManager, logger)
	router := chi.NewRouter()
	router.Post("/signUp", httpHandler.SignUp)
	router.Post("/signIn", httpHandler.SignIn)
	return router
}
