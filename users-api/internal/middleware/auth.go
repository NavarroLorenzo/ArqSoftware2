package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ctxUserID = "uid"
	ctxRole   = "role"
)

type AuthMiddleware struct{ secret []byte }

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: []byte(secret)}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		const p = "Bearer "
		if len(h) <= len(p) || h[:len(p)] != p {
			JSONError(c, http.StatusUnauthorized, "unauthorized", "missing bearer token")
			return
		}
		tokenStr := h[len(p):]

		tk, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return m.secret, nil
		})
		if err != nil || !tk.Valid {
			JSONError(c, http.StatusUnauthorized, "unauthorized", "invalid token")
			return
		}
		claims, ok := tk.Claims.(jwt.MapClaims)
		if !ok {
			JSONError(c, http.StatusUnauthorized, "unauthorized", "invalid claims")
			return
		}
		subF, ok := claims["sub"].(float64)
		if !ok {
			JSONError(c, http.StatusUnauthorized, "unauthorized", "missing sub")
			return
		}
		role, _ := claims["role"].(string)

		c.Set(ctxUserID, uint(subF))
		c.Set(ctxRole, role)
		c.Next()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v, ok := c.Get(ctxRole); !ok || v.(string) != "admin" {
			JSONError(c, http.StatusForbidden, "forbidden", "admin required")
			return
		}
		c.Next()
	}
}

// Helpers que usan los controllers
func UserIDFromCtx(c *gin.Context) (uint, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return 0, false
	}
	return v.(uint), true
}
