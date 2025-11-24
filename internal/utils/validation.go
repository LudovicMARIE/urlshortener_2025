package utils

import (
	"bufio"
	"strings"
)

// Permet de lire et nettoyer la saisie d'un utilisateur depuis un reader
func ReaderLine(reader *bufio.Reader) (string, error) {
	readerValue, _ := reader.ReadString('\n')
	readerValue = strings.TrimSpace(readerValue)
	return readerValue, nil
}