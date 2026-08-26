package security

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type SignUpRequest struct {
	Login string
	Password string
}

type SignInRequest struct {
	Login string
	Password string
}

type HttpHandler struct{
	accountManager AccountManager
	sessionManager SessionManager
	logger *slog.Logger
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
	var session UserSession = UserSession{UserId: id}
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
	var session UserSession = UserSession{UserId: id}
	err = h.sessionManager.CreateNewSession(session, r, w)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "Processing SignIn request failed", "error", err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	http.NoBody.WriteTo(w)
}

func NewHttpHandler(accountManager AccountManager, sessionManager SessionManager, logger *slog.Logger) *HttpHandler {
	return &HttpHandler{
		accountManager,
		sessionManager,
		logger,
	}
}

