package middleware

import (
	"net/http"
	// "strings" // Ya no es necesario, lo podemos quitar

	"megabaseGo/internal/app/services"
	"megabaseGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware maneja la autenticación JWT
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware crea una nueva instancia del middleware de autenticación
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		authService: services.NewAuthService(),
	}
}

// RequireAuth middleware que requiere autenticación (MODIFICADO)
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ---- INICIO DEL CAMBIO ----
		// Ahora leemos el token desde la cookie "access_token"
		tokenString, err := c.Cookie("access_token")

		// Si hay un error (ej. la cookie no existe), devolvemos un error de no autorizado
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			c.Abort()
			return
		}
		// ---- FIN DEL CAMBIO ----

		// El resto de la lógica para validar el token y guardar los datos en el contexto
		// se mantiene exactamente igual.
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.UserName)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)
		c.Set("role_name", claims.RoleName)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole middleware que requiere un rol específico (SIN CAMBIOS)
// Funciona correctamente porque llama internamente a la nueva versión de RequireAuth.
func (m *AuthMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		userRoleName, exists := c.Get("role_name")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Role information not found",
			})
			c.Abort()
			return
		}

		if userRoleName != roleName {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware que requiere uno de varios roles (SIN CAMBIOS)
func (m *AuthMiddleware) RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		userRoleName, exists := c.Get("role_name")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Role information not found",
			})
			c.Abort()
			return
		}

		roleMatches := false
		for _, roleName := range roleNames {
			if userRoleName == roleName {
				roleMatches = true
				break
			}
		}

		if !roleMatches {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware que permite autenticación opcional (MODIFICADO)
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ---- INICIO DEL CAMBIO ----
		// Intentamos leer la cookie, pero no devolvemos un error si no existe.
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			// No hay token, pero la autenticación es opcional, así que continuamos.
			c.Next()
			return
		}
		// ---- FIN DEL CAMBIO ----

		// Si encontramos un token, intentamos validarlo.
		if claims, err := m.authService.ValidateToken(tokenString); err == nil {
			// Token válido, guardamos la información en el contexto.
			c.Set("user_id", claims.UserID)
			c.Set("user_name", claims.UserName)
			c.Set("email", claims.Email)
			c.Set("role_id", claims.RoleID)
			c.Set("role_name", claims.RoleName)
			c.Set("claims", claims)
			c.Set("authenticated", true)
		}

		c.Next()
	}
}

// GetCurrentUserID obtiene el ID del usuario actual del contexto (SIN CAMBIOS)
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUserClaims obtiene los claims del usuario actual (SIN CAMBIOS)
func GetCurrentUserClaims(c *gin.Context) (*utils.JWTClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	return claims.(*utils.JWTClaims), true
}

// IsAuthenticated verifica si el usuario está autenticado (SIN CAMBIOS)
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}