package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"users-api/internal/middleware"
	"users-api/internal/services"
)

type AuthController struct{ svc services.AuthService }

func NewAuthController(s services.AuthService) *AuthController { return &AuthController{svc: s} }

type registerReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) RegisterNormal(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" || req.Name == "" {
		middleware.WriteError(w, http.StatusBadRequest, "validation_error", "name, email, password required")
		return
	}
	u, err := c.svc.RegisterNormal(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email_already_in_use" {
			middleware.WriteError(w, http.StatusConflict, "email_in_use", "email already in use")
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, "register_failed", err.Error())
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, map[string]any{"user": u})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
		return
	}
	token, exp, u, err := c.svc.Login(strings.ToLower(strings.TrimSpace(req.Email)), req.Password)
	if err != nil {
		if err.Error() == "invalid_credentials" {
			middleware.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "email or password incorrect")
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, "login_failed", err.Error())
		return
	}
	middleware.WriteJSON(w, http.StatusOK, map[string]any{
		"access_token": token,
		"expires_at":   exp.UTC(),
		"user":         u,
	})
}
