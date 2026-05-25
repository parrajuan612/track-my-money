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
func (s *service) Register(ctx context.Context, name, email, password string) (string, *domain.User, error) {
	// 1. Validar si el correo ya existe en la base de datos
	existingUser, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("error al verificar el correo en la base de datos")
	}
	if existingUser != nil {
		return "", nil, errors.New("el correo electrónico ya está registrado")
	}

	// 2. Generar el Hash de la contraseña con bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, errors.New("error al encriptar la contraseña")
	}
	hashString := string(hashedPassword)

	// 3. Crear el objeto del nuevo usuario
	newUser := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: &hashString, // Guardamos el hash, NUNCA el texto plano
		AuthProvider: "local",
		IsActive:     true, // Lo activamos por defecto
	}

	// 4. Guardar en Postgres (GORM se encargará de generar el UUID automáticamente)
	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return "", nil, errors.New("error al registrar el usuario")
	}

	// 5. Generar el Access Token para auto-loguear al usuario
	token, err := utils.GenerateAccessToken(newUser.ID)
	if err != nil {
		return "", nil, errors.New("usuario creado, pero hubo un error al generar la sesión")
	}

	return token, newUser, nil
}
