package main

import (
	"strings"

	"github.com/google/uuid"
)

func createID(prefix string) string {
	hash, _ := uuid.NewUUID()
	return prefix + strings.Replace(hash.String(), "-", ".", -1)
}
