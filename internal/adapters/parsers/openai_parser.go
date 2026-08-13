package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"track-my-money/internal/core/domain"

	"github.com/google/uuid"
)

type OpenAIParser struct {
	apiKey string
}

func NewOpenAIParser() *OpenAIParser {
	return &OpenAIParser{
		apiKey: os.Getenv("OPENAI_API_KEY"),
	}
}

// Estructura de la respuesta esperada de OpenAI
type aiResponseData struct {
	Movements []struct {
		Date        string  `json:"date"`
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
	} `json:"movements"`
	PeriodMonth    string `json:"period_month"`
	SuggestedRules struct {
		RegexRow       string `json:"regex_row"`
		RegexDate      string `json:"regex_date"`
		RegexAmount    string `json:"regex_amount"`
		IncomeKeywords string `json:"income_keywords"` // ← AGREGAR ESTA LÍNEA
	} `json:"suggested_rules"`
}

// Parse implementa la interfaz MovementParser
func (p *OpenAIParser) Parse(text string) ([]domain.Movement, string, *domain.BankParsingRule, error) {
	if p.apiKey == "" {
		return nil, "", nil, fmt.Errorf("OPENAI_API_KEY no está configurada")
	}

	reqBody := p.buildRequest(text)
	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error conectando con OpenAI: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, "", nil, fmt.Errorf("error de OpenAI (status %d): %s", res.StatusCode, string(bodyBytes))
	}

	var openAIRes struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(res.Body).Decode(&openAIRes); err != nil {
		return nil, "", nil, fmt.Errorf("error leyendo respuesta de IA: %w", err)
	}

	if len(openAIRes.Choices) == 0 {
		return nil, "", nil, fmt.Errorf("la IA no devolvió resultados")
	}

	var aiData aiResponseData
	if err := json.Unmarshal([]byte(openAIRes.Choices[0].Message.Content), &aiData); err != nil {
		return nil, "", nil, fmt.Errorf("error parseando JSON de la IA: %w", err)
	}

	var movements []domain.Movement
	for _, m := range aiData.Movements {
		parsedDate, err := time.Parse("2006-01-02", m.Date)
		if err != nil {
			parsedDate = time.Now()
		}

		movType := domain.TypeIncome
		finalAmount := m.Amount // Por defecto el monto viene positivo

		// SI LA IA DIJO QUE ES UN GASTO, LO VOLVEMOS NEGATIVO
		if m.Type == "expense" {
			movType = domain.TypeExpense
			finalAmount = -m.Amount // <--- ¡Esta es la magia que faltaba!
		}

		movements = append(movements, domain.Movement{
			ID:          uuid.Nil,
			Date:        parsedDate,
			Description: m.Description,
			Amount:      finalAmount, // Usamos la variable con el signo correcto
			Type:        movType,
			CreatedAt:   time.Now(),
		})
	}

	newRule := &domain.BankParsingRule{
		ID:             uuid.New(), // <-- ¡ESTA ES LA LÍNEA MÁGICA!
		RegexRow:       aiData.SuggestedRules.RegexRow,
		RegexDate:      aiData.SuggestedRules.RegexDate,
		RegexAmount:    aiData.SuggestedRules.RegexAmount,
		IncomeKeywords: aiData.SuggestedRules.IncomeKeywords,
	}

	return movements, aiData.PeriodMonth, newRule, nil
}

func (p *OpenAIParser) buildRequest(text string) map[string]interface{} {
	return map[string]interface{}{
		"model":       "gpt-4o-mini",
		"temperature": 0.0,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Eres un experto financiero y programador. Tu misión es doble: 1. Extraer movimientos del texto crudo (RAW) ignorando saldos iniciales/finales. REGLA VITAL: Si es cuenta de ahorros, valores negativos son 'expense' y positivos 'income'. Si es Tarjeta de Crédito, invierte la lógica: las compras son 'expense' y los pagos/devoluciones son 'income'. 2. Generar 3 expresiones regulares (Regex en Go RE2) para extraer esto en el futuro (regex_row, regex_date, regex_amount) y extraer palabras clave que identifiquen ingresos si es tarjeta de crédito.",
			},
			{
				"role":    "user",
				"content": text,
			},
		},
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name":   "statement_extraction",
				"strict": true,
				"schema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"movements": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"date":        map[string]string{"type": "string", "description": "YYYY-MM-DD"},
									"description": map[string]string{"type": "string"},
									"amount":      map[string]string{"type": "number"},
									"type":        map[string]interface{}{"type": "string", "enum": []string{"income", "expense"}},
								},
								"required":             []string{"date", "description", "amount", "type"},
								"additionalProperties": false,
							},
						},
						"period_month": map[string]string{"type": "string", "description": "YYYY-MM"},
						"suggested_rules": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"regex_row":       map[string]string{"type": "string"},
								"regex_date":      map[string]string{"type": "string"},
								"regex_amount":    map[string]string{"type": "string"},
								"income_keywords": map[string]string{"type": "string", "description": "Palabras separadas por comas que indican que es un ingreso o pago (ej: 'gracias por tu pago, devolucion'). Vacío si es cuenta de ahorros normal."},
							},
							"required":             []string{"regex_row", "regex_date", "regex_amount", "income_keywords"},
							"additionalProperties": false,
						},
					},
					"required":             []string{"movements", "period_month", "suggested_rules"},
					"additionalProperties": false,
				},
			},
		},
	}
}
