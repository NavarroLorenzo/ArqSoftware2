package controllers

import (
	"net/http"
	"strings"
	"time"

	"users-api/internal/middleware"
	"users-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct{ svc services.AuthService }

func NewAuthController(s services.AuthService) *AuthController { return &AuthController{svc: s} }

type registerReq struct {
	Name     string `json:"name" binding:"required,min=2,max=80"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4,max=64"`
}

func (c *AuthController) RegisterNormal(ctx *gin.Context) {
	var req registerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(ctx, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	u, err := c.svc.RegisterNormal(req.Name, email, req.Password)
	if err != nil {
		if err.Error() == "email_already_in_use" {
			middleware.JSONError(ctx, http.StatusConflict, "email_in_use", "email already in use")
			return
		}
		middleware.JSONError(ctx, http.StatusInternalServerError, "register_failed", err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": u})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(ctx, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	token, exp, u, err := c.svc.Login(strings.ToLower(strings.TrimSpace(req.Email)), req.Password)
	if err != nil {
		if err.Error() == "invalid_credentials" {
			middleware.JSONError(ctx, http.StatusUnauthorized, "invalid_credentials", "email or password incorrect")
			return
		}
		middleware.JSONError(ctx, http.StatusInternalServerError, "login_failed", err.Error())
		return
	}
	// normalizamos expires_at en ISO
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"expires_at":   exp.UTC().Format(time.RFC3339),
		"user":         u,
	})
}
