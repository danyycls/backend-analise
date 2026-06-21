package repositorio

import (
	"github.com/danyele/podp/internal/shared/database"
)

type Repositorio struct {
	db database.DB
}

func Novo(db database.DB) *Repositorio {
	return &Repositorio{db: db}
}
