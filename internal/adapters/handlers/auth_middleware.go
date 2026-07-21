package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware intercepta la petición y valida el token JWT
func AuthMiddleware() gin.HandlerFunc {
	// Leemos la clave secreta directamente de tu archivo .env
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "MiClaveSecretaAltamenteSeguraYCompleja2026$" // Fallback por si acaso
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta el token de autorización"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido (se esperaba Bearer <token>)"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Extraemos el user_id (Buscamos 'user_id' o 'id' por compatibilidad con tu utils)
			var userIDStr string
			if val, exists := claims["user_id"]; exists {
				userIDStr = fmt.Sprintf("%v", val)
			} else if val, exists := claims["id"]; exists {
				userIDStr = fmt.Sprintf("%v", val)
			} else if val, exists := claims["sub"]; exists {
				userIDStr = fmt.Sprintf("%v", val)
			}

			if userIDStr == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "El token no contiene un identificador de usuario válido"})
				c.Abort()
				return
			}

			// Guardamos el ID del usuario como String en el contexto de Gin
			c.Set("user_id", userIDStr)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudieron leer las claims del token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
