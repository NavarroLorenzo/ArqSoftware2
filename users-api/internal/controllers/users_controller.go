package controllers

import (
	"net/http"
	"strconv"

	"users-api/internal/middleware"
	"users-api/internal/repository"
)

type UsersController struct{ repo repository.UserRepository }

func NewUsersController(r repository.UserRepository) *UsersController {
	return &UsersController{repo: r}
}

func (c *UsersController) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user")
		return
	}
	u, err := c.repo.FindByID(uid)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}
	if u == nil {
		middleware.WriteError(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	u.Password = ""
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"user": u})
}

func (c *UsersController) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "bad_id", "invalid id")
		return
	}
	u, err := c.repo.FindByID(uint(id64))
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}
	if u == nil {
		middleware.WriteError(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	u.Password = ""
	middleware.WriteJSON(w, http.StatusOK, map[string]any{"user": u})
}
