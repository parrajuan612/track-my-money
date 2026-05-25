package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"track-my-money/internal/core/domain"
	"track-my-money/internal/utils"

	"cloud.google.com/go/auth/credentials/idtoken"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) Login(ctx context.Context, email string, password string) (string, *domain.User, error) {
	// 1. Buscar al usuario usando el puerto del Repositorio

	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return "", nil, err
	}

	if user == nil {
		fmt.Println("esto entro a que es igual a nil")
		return "", nil, errors.New("credenciales inválidas")
	}

	// 3. Validar si está activo
	if !user.IsActive {
		return "", nil, errors.New("el usuario se encuentra inactivo")
	}

	// 4. Validar que tenga contraseña local
	if user.PasswordHash == nil {

		return "", nil, errors.New("credenciales inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {

		return "", nil, errors.New("credenciales inválidas")
	}

	// 6. CAMBIO: ¡Generamos el JWT Real usando el ID del usuario de la DB!
	token, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", nil, errors.New("error al generar el token de acceso")
	}

	return token, user, nil
}

func (s *service) LoginWithGoogle(ctx context.Context, idTokenStr string) (string, *domain.User, error) {
	// 1. Validar el token directamente con los servidores de Google
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	payload, err := idtoken.Validate(ctx, idTokenStr, clientID)
	if err != nil {
		return "", nil, errors.New("el token de Google es inválido o ha expirado")
	}

	// 2. Extraer la información básica del Payload de Google
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	externalID := payload.Subject // Identificador único e inmutable de Google

	// 3. Buscar si el usuario ya existe en nuestra base de datos
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err // Error de conexión a la BD
	}

	// 4. Lógica de "Upsert" (Actualizar o Insertar)
	if user == nil {
		// El usuario es nuevo: Lo registramos automáticamente sin contraseña
		newUser := &domain.User{
			Name:         name,
			Email:        email,
			AuthProvider: "google",
			ExternalID:   externalID,
			IsActive:     true,
		}

		// Guardamos en BD. GORM se encarga de rellenar el ID (UUID) tras el Create.
		err = s.repo.CreateUser(ctx, newUser)
		if err != nil {
			return "", nil, errors.New("error al registrar al nuevo usuario")
		}
		user = newUser // Asignamos el nuevo usuario para el resto del flujo
	} else {
		// El usuario ya existía: Validamos que no esté suspendido o inactivo
		if !user.IsActive {
			return "", nil, errors.New("el usuario se encuentra inactivo")
		}
		// Nota: Podrías querer actualizar su "ExternalID" si antes era solo 'local'
	}

	// 5. Generar NUESTRO token de sesión (Access Token)
	token, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", nil, errors.New("error al generar el token de acceso")
	}

	return token, user, nil
}
