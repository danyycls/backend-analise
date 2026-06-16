package parse

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/encoding/charmap"
)

var poolMapaCSV = sync.Pool{
	New: func() any {
		return make(map[string]string, 50)
	},
}

func lerArquivoCSV(caminho string, callback func(numeroLinha int, registro map[string]string) error) error {
	arquivo, err := os.Open(caminho)
	if err != nil {
		return fmt.Errorf("falha ao abrir CSV: %w", err)
	}
	defer arquivo.Close()

	leitorBuf := bufio.NewReaderSize(arquivo, 64*1024)
	decoder := charmap.ISO8859_1.NewDecoder().Reader(leitorBuf)
	leitor := csv.NewReader(decoder)
	leitor.Comma = ';'
	leitor.FieldsPerRecord = -1
	leitor.LazyQuotes = true

	cabecalho, err := leitor.Read()
	if err != nil {
		return err
	}

	for i := range cabecalho {
		cabecalho[i] = strings.TrimSpace(cabecalho[i])
	}

	numeroLinha := 1
	for {
		linha, err := leitor.Read()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("erro de leitura na linha %d: %w", numeroLinha+1, err)
		}

		numeroLinha++
		registro := poolMapaCSV.Get().(map[string]string)
		for i, coluna := range cabecalho {
			if i < len(linha) {
				registro[coluna] = strings.TrimSpace(linha[i])
				continue
			}
			registro[coluna] = ""
		}

		if err := callback(numeroLinha, registro); err != nil {
			return err
		}

		for k := range registro {
			delete(registro, k)
		}
		poolMapaCSV.Put(registro)
	}
}
