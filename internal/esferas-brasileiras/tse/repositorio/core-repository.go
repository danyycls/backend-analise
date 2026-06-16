package repositorio

import (
	"github.com/danyele/laceu/internal/shared/database"
)

type Repositorio struct {
	db database.DB
}

func Novo(db database.DB) *Repositorio {
	return &Repositorio{db: db}
}
