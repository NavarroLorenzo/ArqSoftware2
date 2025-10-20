package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	ctxUserID ctxKey = "uid"
	ctxRole   ctxKey = "role"
)

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: []byte(secret)}
}

// Middleware: exige JWT válido y coloca uid/role en el contexto
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := bearer(r.Header.Get("Authorization"))
		if tokenStr == "" {
			WriteError(w, http.StatusUnauthorized, "unauthorized", "missing bearer token")
			return
		}
		tk, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return m.secret, nil
		})
		if err != nil || !tk.Valid {
			WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
			return
		}
		claims, ok := tk.Claims.(jwt.MapClaims)
		if !ok {
			WriteError(w, http.StatusUnauthorized, "unauthorized", "invalid claims")
			return
		}

		// sub puede venir como float64 desde JSON Web Token
		subF, ok := claims["sub"].(float64)
		if !ok {
			WriteError(w, http.StatusUnauthorized, "unauthorized", "missing sub")
			return
		}
		role, _ := claims["role"].(string)

		uid := uint(subF) // conversión explícita
		ctx := context.WithValue(r.Context(), ctxUserID, uid)
		ctx = context.WithValue(ctx, ctxRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware: requiere rol admin (encadena RequireAuth)
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role, _ := r.Context().Value(ctxRole).(string); role != "admin" {
			WriteError(w, http.StatusForbidden, "forbidden", "admin required")
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// -------- Helpers exportados (los que usa tu controller) --------

// UserIDFromCtx devuelve el userID seteado por RequireAuth
func UserIDFromCtx(ctx context.Context) (uint, bool) {
	uid, ok := ctx.Value(ctxUserID).(uint)
	return uid, ok
}

func RoleFromCtx(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ctxRole).(string)
	return role, ok
}

func bearer(h string) string {
	const p = "Bearer "
	if len(h) > len(p) && h[:len(p)] == p {
		return h[len(p):]
	}
	return ""
}
