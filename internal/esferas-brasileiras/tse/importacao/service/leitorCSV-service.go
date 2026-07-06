package service

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	parse "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/parse"
	tipos "github.com/danyele/podp/internal/esferas-brasileiras/tse/importacao/types"
)

func (s *LeitorCSVService) NovoProcessador() *parse.ProcessadorLeitorCSV {
	return parse.NovoProcessadorLeitorCSV(s.tamanhoLote)
}

type LeitorCSVServiceInterface interface {
	ListarArquivos() ([]tipos.ArquivoImportacao, error)
	LerArquivo(ctx context.Context, processador *parse.ProcessadorLeitorCSV, arquivo tipos.ArquivoImportacao) (*tipos.ArquivoProcessado, error)
	NovoProcessador() *parse.ProcessadorLeitorCSV
}

type LeitorCSVService struct {
	diretorioCSV string
	tamanhoLote  int
}

func NovoLeitorCSVService(diretorioCSV string) *LeitorCSVService {
	return &LeitorCSVService{
		diretorioCSV: diretorioCSV,
		tamanhoLote:  tipos.TamanhoLotePadraoImportacao,
	}
}

func (s *LeitorCSVService) ListarArquivos() ([]tipos.ArquivoImportacao, error) {
	arquivos, err := s.localizarArquivos()
	if err != nil {
		return nil, err
	}
	if len(arquivos) == 0 {
		return nil, fmt.Errorf("nenhum arquivo CSV suportado encontrado em %s", s.diretorioCSV)
	}
	return arquivos, nil
}

func (s *LeitorCSVService) LerArquivo(ctx context.Context, processador *parse.ProcessadorLeitorCSV, arquivo tipos.ArquivoImportacao) (*tipos.ArquivoProcessado, error) {
	if processador == nil {
		processador = parse.NovoProcessadorLeitorCSV(s.tamanhoLote)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	n, err := parse.ProcessarArquivo(ctx, processador, arquivo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", arquivo.Nome, err)
	}

	return &tipos.ArquivoProcessado{
		NomeArquivo: arquivo.Nome,
		Tipo:        arquivo.Tipo,
		Registros:   n,
	}, nil
}

func (s *LeitorCSVService) localizarArquivos() ([]tipos.ArquivoImportacao, error) {
	arquivos := make([]tipos.ArquivoImportacao, 0)

	err := filepath.WalkDir(s.diretorioCSV, func(caminho string, entrada fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entrada.IsDir() {
			return nil
		}
		if strings.ToLower(filepath.Ext(entrada.Name())) != ".csv" {
			return nil
		}

		tipo, suportado := parse.IdentificarTipoArquivo(strings.ToLower(entrada.Name()))
		if !suportado {
			return nil
		}

		caminhoRelativo, _ := filepath.Rel(s.diretorioCSV, caminho)
		parts := strings.SplitN(caminhoRelativo, string(filepath.Separator), 2)
		if len(parts) == 2 {
			if _, err := strconv.Atoi(parts[0]); err == nil {
				caminhoRelativo = parts[1]
			}
		}
		diretorio := filepath.Dir(caminhoRelativo)
		if diretorio == "." {
			diretorio = ""
		}

		arquivos = append(arquivos, tipos.ArquivoImportacao{
			Caminho:         caminho,
			CaminhoRelativo: caminhoRelativo,
			Diretorio:       diretorio,
			DiretorioLower:  strings.ToLower(diretorio),
			Nome:            entrada.Name(),
			NomeLower:       strings.ToLower(entrada.Name()),
			Tipo:            tipo,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(arquivos, func(i, j int) bool {
		ordemI := parse.PrioridadeTipoArquivo(arquivos[i].Tipo)
		ordemJ := parse.PrioridadeTipoArquivo(arquivos[j].Tipo)
		if ordemI == ordemJ {
			return arquivos[i].Nome < arquivos[j].Nome
		}
		return ordemI < ordemJ
	})

	return arquivos, nil
}
