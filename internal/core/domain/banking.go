package domain

import "time"

type Bank struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (Bank) TableName() string {
	return "movements.banks" // Especificamos el esquema 'movements'
}

type Account struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id"`
	BankID    int       `json:"bank_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	// Relación para que devuelva el nombre del banco
	Bank *Bank `json:"bank,omitempty" gorm:"foreignKey:BankID"`
}

func (Account) TableName() string {
	return "movements.accounts" // Especificamos el esquema 'movements'
}

type CreateAccountRequest struct {
	BankID int    `json:"bank_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}
