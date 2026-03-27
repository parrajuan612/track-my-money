package ports

import "io"

type DocumentExtractor interface {
	Extract(file io.Reader, password string) (string, error)
}
