package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/danyele/podp/internal/shared/database"
	"github.com/danyele/podp/internal/shared/logger"
)

type ConvenioDTO struct {
	CNPJ      string `json:"cnpj"`
	UF        string `json:"uf"`
	Municipio string `json:"municipio"`
	NomeOrgao string `json:"nome_orgao"`
	Tipo      string `json:"tipo"`
}

type ListarConveniosResultado struct {
	Dados        []ConvenioDTO `json:"dados"`
	Total        int           `json:"total"`
	Pagina       int           `json:"pagina"`
	PorPagina    int           `json:"por_pagina"`
	TotalPaginas int           `json:"total_paginas"`
}

type ConsultaConvenioUseCase struct {
	db database.DB
}

func NovoConsultaConvenioUseCase(db database.DB) *ConsultaConvenioUseCase {
	return &ConsultaConvenioUseCase{db: db}
}

func (u *ConsultaConvenioUseCase) Listar(ctx context.Context, pagina, porPagina int, uf, municipio, tipo string) (*ListarConveniosResultado, error) {
	log := logger.New("Convenio: UseCase: Listar")

	if pagina < 1 {
		pagina = 1
	}
	if porPagina < 1 || porPagina > 100 {
		porPagina = 10
	}
	offset := (pagina - 1) * porPagina

	where := "WHERE deleted_at IS NULL"
	args := make([]interface{}, 0)
	argIdx := 1

	if uf != "" {
		where += fmt.Sprintf(" AND uf = $%d", argIdx)
		args = append(args, uf)
		argIdx++
	}
	if municipio != "" {
		where += fmt.Sprintf(" AND LOWER(nome_municipio) LIKE LOWER($%d)", argIdx)
		args = append(args, "%"+municipio+"%")
		argIdx++
	}
	if tipo != "" {
		where += fmt.Sprintf(" AND tipo_convenente = $%d", argIdx)
		args = append(args, tipo)
		argIdx++
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM convenio %s", where)
	var total int
	if err := u.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		log.Error("erro ao contar convenios", "erro", err)
		return nil, err
	}

	totalPaginas := int(math.Ceil(float64(total) / float64(porPagina)))
	if totalPaginas < 1 {
		totalPaginas = 1
	}

	query := fmt.Sprintf(`
		SELECT codigo_convenente, uf, COALESCE(nome_municipio,''), COALESCE(nome_convenente,''), COALESCE(tipo_convenente,'')
		FROM convenio
		%s
		ORDER BY nome_convenente
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, porPagina, offset)

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("erro ao consultar convenios", "erro", err)
		return nil, err
	}
	defer rows.Close()

	dados := make([]ConvenioDTO, 0)
	for rows.Next() {
		var d ConvenioDTO
		if err := rows.Scan(&d.CNPJ, &d.UF, &d.Municipio, &d.NomeOrgao, &d.Tipo); err != nil {
			log.Error("erro ao scanear convenio", "erro", err)
			continue
		}
		dados = append(dados, d)
	}

	return &ListarConveniosResultado{
		Dados:        dados,
		Total:        total,
		Pagina:       pagina,
		PorPagina:    porPagina,
		TotalPaginas: totalPaginas,
	}, nil
}
