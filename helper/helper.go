package helper

import (
	"strings"

	"github.com/google/uuid"
)

// CreateID generate Id
func CreateID(prefix string) string {
	hash, _ := uuid.NewUUID()
	return prefix + strings.Replace(hash.String(), "-", ".", -1)
}
