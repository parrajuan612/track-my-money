package postgres

import (
	"context"
	"fmt"
	"track-my-money/internal/core/domain"
)

func (r *postgresRepository) GetMovements(ctx context.Context, filters domain.MovementFilters) ([]domain.MovementTable, int64, error) {
	var movements []domain.MovementTable
	var total int64

	// 1. SELECT con nombres legibles. Usamos el esquema movements. para cada tabla.
	query := r.db.WithContext(ctx).Table("movements.movements m").
		Select(`
            m.id, 
            m.date, 
            m.description, 
            m.amount, 
            m.type, 
            c.name as category_name, 
            a.name as account_name, 
            b.name as bank_name
        `).
		// 2. JOINS con el esquema movements. explícito
		Joins("JOIN movements.categories c ON m.category_id = c.id").
		Joins("JOIN movements.accounts a ON m.account_id = a.id").
		Joins("JOIN movements.banks b ON a.bank_id = b.id").
		Where("m.user_id = ? AND m.is_active = true", filters.UserID)

	// 3. Filtros dinámicos (Usando el alias 'm' para evitar ambigüedad)
	if filters.AccountID != nil {
		query = query.Where("m.account_id = ?", *filters.AccountID)
	}
	if filters.CategoryID != nil {
		query = query.Where("m.category_id = ?", *filters.CategoryID)
	}
	if filters.StartDate != nil {
		query = query.Where("m.date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("m.date <= ?", *filters.EndDate)
	}
	if filters.Type != nil {
		query = query.Where("m.type = ?", *filters.Type)
	}
	if filters.Query != nil && *filters.Query != "" {
		query = query.Where("m.description ILIKE ?", "%"+*filters.Query+"%")
	}

	// 4. Conteo y Paginación
	query.Count(&total)

	offset := (filters.Page - 1) * filters.PageSize
	order := fmt.Sprintf("m.%s %s", filters.SortBy, filters.SortOrder)

	// Usamos Scan para que GORM mapee los nombres de las columnas (category_name, etc.)
	// a los campos del struct Movement que definimos antes.
	err := query.Order(order).Limit(filters.PageSize).Offset(offset).Scan(&movements).Error

	return movements, total, err
}
