package parsers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"track-my-money/internal/core/domain/ports"
)

type pdfExtractor struct{}

func NewPDFExtractor() ports.DocumentExtractor {
	return &pdfExtractor{}
}

func (e *pdfExtractor) Extract(file io.Reader, password string) (string, error) {

	if err := e.ensureDependencies(); err != nil {
		return "", err
	}
	originalPath, err := e.saveToTemp(file)
	if err != nil {
		return "", err
	}
	defer os.Remove(originalPath)

	decryptedPath, err := e.createTempPath()
	if err != nil {
		return "", err
	}
	defer os.Remove(decryptedPath)

	if err := e.decryptPDF(originalPath, password, decryptedPath); err != nil {
		return "", err
	}

	return e.runExtractionFlow(decryptedPath)
}

func (e *pdfExtractor) saveToTemp(file io.Reader) (string, error) {
	tmpFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return "", fmt.Errorf("error creando temporal de entrada: %w", err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		return "", fmt.Errorf("error escribiendo bytes al disco: %w", err)
	}
	return tmpFile.Name(), nil
}

func (e *pdfExtractor) createTempPath() (string, error) {
	tmpFile, err := os.CreateTemp("", "decrypted-*.pdf")
	if err != nil {
		return "", err
	}
	path := tmpFile.Name()
	tmpFile.Close()
	return path, nil
}

func (e *pdfExtractor) decryptPDF(input, password, output string) error {
	args := []string{"--password=" + password, "--decrypt", input, output}
	cmd := exec.Command("qpdf", args...)
	var errOut bytes.Buffer
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("qpdf error: %v - %s", err, errOut.String())
	}
	return nil
}

func (e *pdfExtractor) runExtractionFlow(path string) (string, error) {
	// Intentar con layout
	text, err := e.pdftotext(path, true)
	if err != nil || text == "" {
		// Reintentar sin layout si falla
		text, err = e.pdftotext(path, false)
	}

	if err != nil {
		return "", err
	}
	if text == "" {
		return "", fmt.Errorf("el PDF no contiene texto extraíble")
	}
	return text, nil
}

func (e *pdfExtractor) pdftotext(path string, useLayout bool) (string, error) {
	args := []string{}
	if useLayout {
		args = append(args, "-layout")
	}
	args = append(args, path, "-") // El "-" indica que imprima a Stdout

	cmd := exec.Command("pdftotext", args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftotext error: %v", errOut.String())
	}
	return out.String(), nil
}

func (e *pdfExtractor) ensureDependencies() error {
	if _, err := exec.LookPath("qpdf"); err != nil {
		return fmt.Errorf("qpdf no instalado")
	}
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return fmt.Errorf("pdftotext no instalado")
	}
	return nil
}
