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

	dbData, err := s.repo.GetMoneyFlow(ctx, userID, startDate, groupBy, req.AccountID, req.BankID)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]domain.TimeSeriesData)
	for _, d := range dbData {
		dataMap[d.Label] = d
	}

	finalData := []domain.TimeSeriesData{}
	for i := 0; i < daysToFill; i++ {
		var dateLabel string
		if req.Range == "1y" {
			dateLabel = startDate.AddDate(0, i, 0).Format("Jan 2006")
		} else {
			dateLabel = startDate.AddDate(0, 0, i).Format("02/01")
		}

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
        startDate = time.Now().AddDate(0, 0, -6)
        groupBy = "DD/MM" 
    case "1m":
        now := time.Now()
        startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
        daysToFill = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
        groupBy = "DD/MM"
case "1y":
        daysToFill = 12
        
        // 1. COMENTA ESTA LÍNEA TEMPORALMENTE (Para que no use la fecha real)
        // now := time.Now() 
        
        // 2. AGREGA ESTA LÍNEA (Tu máquina del tiempo al 31 de Julio de 2026)
        now := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.Local)
        
        // El resto del código que ya habías arreglado
        startDate = time.Date(now.Year()-1, now.Month(), 1, 0, 0, 0, 0, time.Local)
        groupBy = "Mon YYYY"
    default:
        return time.Time{}, "", 0, fmt.Errorf("rango no soportado")
    }
    return startDate, groupBy, daysToFill, nil
}
