package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

type DynamicParser struct {
	rule domain.BankParsingRule
}

func NewDynamicParser(rule domain.BankParsingRule) *DynamicParser {
	return &DynamicParser{rule: rule}
}

func (p *DynamicParser) Parse(text string) ([]domain.Movement, string, *domain.BankParsingRule, error) {
	if strings.TrimSpace(text) == "" {
		return nil, "", nil, fmt.Errorf("texto vacío")
	}

	// 1. Compilamos las reglas generadas por la IA
	reRow, err := regexp.Compile(p.rule.RegexRow)
	if err != nil {
		return nil, "", nil, fmt.Errorf("regex_row inválido en base de datos: %v", err)
	}

	var movements []domain.Movement
	lines := strings.Split(text, "\n")
	year := time.Now().Year()

	// 2. Escaneamos línea por línea
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if reRow.MatchString(line) {
			parts := regexp.MustCompile(`\s{2,}`).Split(line, -1)
			if len(parts) < 2 {
				continue
			}

			desc := parts[0]
			amountStr := parts[len(parts)-1]
			amountStr = strings.ReplaceAll(amountStr, "$", "")
			amountStr = strings.ReplaceAll(amountStr, ".", "")
			amountStr = strings.ReplaceAll(amountStr, ",", ".")
			amount, _ := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)

			if amount == 0 {
				continue
			}

			// --- NUEVA LÓGICA INTELIGENTE ---
			movType := domain.TypeIncome
			if amount < 0 {
				movType = domain.TypeExpense
			}

			// Si la BD tiene palabras clave de ingresos (Significa que es Tarjeta de Crédito)
			if p.rule.IncomeKeywords != "" {
				movType = domain.TypeExpense // En TDC, todo es gasto por defecto
				lowerDesc := strings.ToLower(desc)
				keywords := strings.Split(p.rule.IncomeKeywords, ",")

				for _, kw := range keywords {
					kw = strings.TrimSpace(strings.ToLower(kw))
					if kw != "" && strings.Contains(lowerDesc, kw) {
						movType = domain.TypeIncome // Si dice "pago" o "devolución", es ingreso
						break
					}
				}

				// En TDC, si es gasto y viene positivo, lo volvemos negativo para la BD
				if movType == domain.TypeExpense && amount > 0 {
					amount = -amount
				}
			}
			movements = append(movements, domain.Movement{
				ID:          uuid.Nil,
				Date:        time.Date(year, time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
				Description: desc,
				Amount:      amount,
				Type:        movType,
				CreatedAt:   time.Now(),
			})
		}
	}

	if len(movements) == 0 {
		return nil, "", nil, fmt.Errorf("el parser dinámico no encontró movimientos")
	}

	// Devuelve (movements, periodMonth, ruleSugerida=nil, error=nil)
	return movements, "2026-08", nil, nil
}
