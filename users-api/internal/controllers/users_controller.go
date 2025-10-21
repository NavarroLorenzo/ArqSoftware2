package controllers

import (
	"net/http"
	"strconv"

	"users-api/internal/middleware"
	"users-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type UsersController struct{ repo repository.UserRepository }

func NewUsersController(r repository.UserRepository) *UsersController {
	return &UsersController{repo: r}
}

func (c *UsersController) Me(ctx *gin.Context) {
	uid, ok := middleware.UserIDFromCtx(ctx)
	if !ok {
		middleware.JSONError(ctx, http.StatusUnauthorized, "unauthorized", "missing user")
		return
	}
	u, err := c.repo.FindByID(uid)
	if err != nil {
		middleware.JSONError(ctx, http.StatusInternalServerError, "db_error", err.Error())
		return
	}
	if u == nil {
		middleware.JSONError(ctx, http.StatusNotFound, "not_found", "user not found")
		return
	}
	u.Password = ""
	ctx.JSON(http.StatusOK, gin.H{"user": u})
}

func (c *UsersController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		middleware.JSONError(ctx, http.StatusBadRequest, "bad_id", "invalid id")
		return
	}
	u, err := c.repo.FindByID(uint(id64))
	if err != nil {
		middleware.JSONError(ctx, http.StatusInternalServerError, "db_error", err.Error())
		return
	}
	if u == nil {
		middleware.JSONError(ctx, http.StatusNotFound, "not_found", "user not found")
		return
	}
	u.Password = ""
	ctx.JSON(http.StatusOK, gin.H{"user": u})
}
