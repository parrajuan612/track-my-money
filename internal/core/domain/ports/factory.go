package ports

type MovementParserFactory interface {
	GetParser(bankID string) (MovementParser, error)
}
