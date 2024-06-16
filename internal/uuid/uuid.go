package uuid

import "github.com/google/uuid"

type (
	GoogleUUID struct{}
)

func (u GoogleUUID) Generate() string {
	return uuid.NewString()
}

func New() *GoogleUUID {
	return &GoogleUUID{}
}
