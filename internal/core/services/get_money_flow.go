package services

import (
	"context"
	"fmt"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

func (s *service) GetMoneyFlow(ctx context.Context, userID uuid.UUID, req domain.MoneyFlowRequest) ([]domain.TimeSeriesData, error) {

	startDate, groupBy, daysToFill, err := getRangeTime(req.Range)
	if err != nil {
		return []domain.TimeSeriesData{}, err
	}

	// 2. Llamar al repositorio (Lo que hicimos en el paso anterior)
	dbData, err := s.analysisRepo.GetMoneyFlow(ctx, userID, startDate, groupBy, req.AccountID, req.BankID)
	if err != nil {
		return nil, err
	}

	// 3. LA MAGIA: Rellenar huecos (Data Filling)
	// Creamos un mapa para búsqueda rápida de lo que sí trajo la DB
	dataMap := make(map[string]domain.TimeSeriesData)
	for _, d := range dbData {
		dataMap[d.Label] = d
	}

	// 4. Generar la serie completa (Calendario)
	finalData := []domain.TimeSeriesData{}
	for i := 0; i < daysToFill; i++ {
		var dateLabel string
		if req.Range == "1y" {
			dateLabel = startDate.AddDate(0, i, 0).Format("Jan 2006")
		} else {
			dateLabel = startDate.AddDate(0, 0, i).Format("02/01")
		}

		// Si la DB tiene el dato, lo usamos; si no, ponemos 0
		if val, ok := dataMap[dateLabel]; ok {
			finalData = append(finalData, val)
		} else {
			finalData = append(finalData, domain.TimeSeriesData{
				Label:   dateLabel,
				Income:  0,
				Expense: 0,
			})
		}
	}

	return finalData, nil
}

func getRangeTime(tp string) (time.Time, string, int, error) {
	var startDate time.Time
	var groupBy string
	var daysToFill int

	switch tp {
	case "1w":
		daysToFill = 7
		startDate = time.Now().AddDate(0, 0, -7)
		groupBy = "DD/MM" // "30/03"
	case "1m":
		// En lugar de restar un mes relativo, vamos al inicio del mes actual
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)

		// Calculamos cuántos días tiene el mes actual para que el bucle sea exacto
		// (Ir al mes siguiente día 0 nos da el último día del mes actual)
		daysToFill = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()

		groupBy = "DD/MM"
	case "1y":
		daysToFill = 12
		startDate = time.Now().AddDate(-1, 0, 0)
		groupBy = "Mon YYYY" // "Mar 2026"
	default:
		return time.Time{}, "", 0, fmt.Errorf("rango no soportado")
	}
	return startDate, groupBy, daysToFill, nil
}
