package main

import (
	"bufio"
	"fmt"
	"go/scanner"
	custom_errors "mechanus-compiler/error"
	"os"
	"strings"
)

// Constants for token types
const (
	T_PROGRAMA      = 1
	T_FIM           = 2
	T_VARIAVEIS     = 3
	T_VIRGULA       = 4
	T_PONTO_VIRGULA = 5
	T_SE            = 6
	T_SENAO         = 7
	T_FIM_SE        = 8
	T_ENQUANTO      = 9
	T_FIM_ENQUANTO  = 10
	T_PARA          = 11
	T_SETA          = 12
	T_ATE           = 13
	T_FIM_PARA      = 14
	T_LER           = 15
	T_ABRE_PAR      = 16
	T_FECHA_PAR     = 17
	T_ESCREVER      = 18
	T_MAIOR         = 19
	T_MENOR         = 20
	T_MAIOR_IGUAL   = 21
	T_MENOR_IGUAL   = 22
	T_IGUAL         = 23
	T_DIFERENTE     = 24
	T_MAIS          = 25
	T_MENOS         = 26
	T_VEZES         = 27
	T_DIVIDIDO      = 28
	T_RESTO         = 29
	T_ELEVADO       = 30
	T_NUMERO        = 31
	T_ID            = 32
	T_FIM_FONTE     = 90
	T_ERRO_LEX      = 98
	T_NULO          = 99
	FIM_ARQUIVO     = 26
)

// Lexical struct to hold lexical analyzer state
type Lexical struct {
	inputFile           *os.File
	rdFonte             *bufio.Reader
	outputFile          *os.File
	lookAhead           rune
	token               int
	lexema              string
	ponteiro            int
	linhaFonte          string
	linhaAtual          int
	colunaAtual         int
	mensagemDeErro      string
	tokensIdentificados strings.Builder
}

func (lex *Lexical) GetToken() (string, error) {
	err := lex.open()

	if err != nil {
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		return "", err
	}
	defer lex.close()

	lex.linhaAtual = 0
	lex.colunaAtual = 0
	lex.ponteiro = 0
	lex.linhaFonte = ""
	lex.token = T_NULO
	lex.mensagemDeErro = ""

	lex.movelookAhead()

	for lex.token != T_FIM_FONTE && lex.token != T_ERRO_LEX {
		lex.buscaProximoToken()
		lex.mostraToken()
	}

	if lex.token == T_ERRO_LEX {
		fmt.Println("Erro Léxico:", lex.mensagemDeErro)
	} else {
		fmt.Println("Análise Léxica terminada sem erros léxicos")
	}

	lex.exibeTokens()
	lex.gravaSaida()
}

func (lex *Lexical) movelookAhead() error {
	if lex.ponteiro+1 > len(lex.linhaFonte) {
		lex.linhaAtual++
		lex.ponteiro = 0
		line, err := lex.rdFonte.ReadString('\n')
		if err != nil {
			lex.lookAhead = FIM_ARQUIVO
			return err
		}
		lex.linhaFonte = line
		lex.lookAhead = rune(lex.linhaFonte[lex.ponteiro])
	} else {
		lex.lookAhead = rune(lex.linhaFonte[lex.ponteiro])
	}
	if lex.lookAhead >= 'a' && lex.lookAhead <= 'z' {
		lex.lookAhead = lex.lookAhead - 'a' + 'A'
	}
	lex.ponteiro++
	lex.colunaAtual = lex.ponteiro + 1
	return nil
}

func (lex *Lexical) buscaProximoToken() error {
	var sbLexema strings.Builder

	for lex.lookAhead == ' ' || lex.lookAhead == '\t' || lex.lookAhead == '\n' || lex.lookAhead == '\r' {
		lex.movelookAhead()
	}

	if lex.lookAhead >= 'A' && lex.lookAhead <= 'Z' {
		sbLexema.WriteRune(lex.lookAhead)
		lex.movelookAhead()

		for (lex.lookAhead >= 'A' && lex.lookAhead <= 'Z') || (lex.lookAhead >= '0' && lex.lookAhead <= '9') || lex.lookAhead == '_' {
			sbLexema.WriteRune(lex.lookAhead)
			lex.movelookAhead()
		}

		lex.lexema = sbLexema.String()

		switch lex.lexema {
		case "PROGRAMA":
			lex.token = T_PROGRAMA
		case "FIM":
			lex.token = T_FIM
		case "VARIAVEIS":
			lex.token = T_VARIAVEIS
		case "SE":
			lex.token = T_SE
		case "SENAO":
			lex.token = T_SENAO
		case "FIM_SE":
			lex.token = T_FIM_SE
		case "ENQUANTO":
			lex.token = T_ENQUANTO
		case "FIM_ENQUANTO":
			lex.token = T_FIM_ENQUANTO
		case "PARA":
			lex.token = T_PARA
		case "ATE":
			lex.token = T_ATE
		case "FIM_PARA":
			lex.token = T_FIM_PARA
		case "LER":
			lex.token = T_LER
		case "ESCREVER":
			lex.token = T_ESCREVER
		default:
			lex.token = T_ID
		}
	} else if lex.lookAhead >= '0' && lex.lookAhead <= '9' {
		sbLexema.WriteRune(lex.lookAhead)
		lex.movelookAhead()
		for lex.lookAhead >= '0' && lex.lookAhead <= '9' {
			sbLexema.WriteRune(lex.lookAhead)
			lex.movelookAhead()
		}
		lex.token = T_NUMERO
	} else {
		switch lex.lookAhead {
		case '(':
			lex.token = T_ABRE_PAR
		case ')':
			lex.token = T_FECHA_PAR
		case ';':
			lex.token = T_PONTO_VIRGULA
		case ',':
			lex.token = T_VIRGULA
		case '+':
			lex.token = T_MAIS
		case '-':
			lex.token = T_MENOS
		case '*':
			lex.token = T_VEZES
		case '/':
			lex.token = T_DIVIDIDO
		case '%':
			lex.token = T_RESTO
		case '<':
			lex.token = T_MENOR
		case '>':
			lex.token = T_MAIOR
		case '=':
			lex.token = T_IGUAL
		case FIM_ARQUIVO:
			lex.token = T_FIM_FONTE
		default:
			lex.token = T_ERRO_LEX
			lex.mensagemDeErro = fmt.Sprintf("Erro Léxico na linha: %d\nReconhecido ao atingir a coluna: %d\nLinha do Erro: <%s>\nToken desconhecido: %c", lex.linhaAtual, lex.colunaAtual, lex.linhaFonte, lex.lookAhead)
		}
		sbLexema.WriteRune(lex.lookAhead)
		lex.movelookAhead()
	}

	lex.lexema = sbLexema.String()
	return nil
}

func (lex *Lexical) mostraToken() {
	var tokenLexema string
	switch lex.token {
	case T_PROGRAMA:
		tokenLexema = "T_PROGRAMA"
	case T_FIM:
		tokenLexema = "T_FIM"
	case T_VARIAVEIS:
		tokenLexema = "T_VARIAVEIS"
	case T_VIRGULA:
		tokenLexema = "T_VIRGULA"
	case T_PONTO_VIRGULA:
		tokenLexema = "T_PONTO_VIRGULA"
	case T_SE:
		tokenLexema = "T_SE"
	case T_SENAO:
		tokenLexema = "T_SENAO"
	case T_FIM_SE:
		tokenLexema = "T_FIM_SE"
	case T_ENQUANTO:
		tokenLexema = "T_ENQUANTO"
	case T_FIM_ENQUANTO:
		tokenLexema = "T_FIM_ENQUANTO"
	case T_PARA:
		tokenLexema = "T_PARA"
	case T_SETA:
		tokenLexema = "T_SETA"
	case T_ATE:
		tokenLexema = "T_ATE"
	case T_FIM_PARA:
		tokenLexema = "T_FIM_PARA"
	case T_LER:
		tokenLexema = "T_LER"
	case T_ABRE_PAR:
		tokenLexema = "T_ABRE_PAR"
	case T_FECHA_PAR:
		tokenLexema = "T_FECHA_PAR"
	case T_ESCREVER:
		tokenLexema = "T_ESCREVER"
	case T_MAIOR:
		tokenLexema = "T_MAIOR"
	case T_MENOR:
		tokenLexema = "T_MENOR"
	case T_MAIOR_IGUAL:
		tokenLexema = "T_MAIOR_IGUAL"
	case T_MENOR_IGUAL:
		tokenLexema = "T_MENOR_IGUAL"
	case T_IGUAL:
		tokenLexema = "T_IGUAL"
	case T_DIFERENTE:
		tokenLexema = "T_DIFERENTE"
	case T_MAIS:
		tokenLexema = "T_MAIS"
	case T_MENOS:
		tokenLexema = "T_MENOS"
	case T_VEZES:
		tokenLexema = "T_VEZES"
	case T_DIVIDIDO:
		tokenLexema = "T_DIVIDIDO"
	case T_RESTO:
		tokenLexema = "T_RESTO"
	case T_ELEVADO:
		tokenLexema = "T_ELEVADO"
	case T_NUMERO:
		tokenLexema = "T_NUMERO"
	case T_ID:
		tokenLexema = "T_ID"
	case T_FIM_FONTE:
		tokenLexema = "T_FIM_FONTE"
	case T_ERRO_LEX:
		tokenLexema = "T_ERRO_LEX"
	case T_NULO:
		tokenLexema = "T_NULO"
	default:
		tokenLexema = "N/A"
	}
	fmt.Println(tokenLexema + " ( " + lex.lexema + " )")
	lex.acumulaToken(tokenLexema + " ( " + lex.lexema + " )")
}

func (lex *Lexical) open() error {
	file, err := os.Open("input.grm")
	if err != nil {
		return err
	}
	lex.inputFile = file
	lex.rdFonte = bufio.NewReader(file)
	return nil
}

func (lex *Lexical) close(file string) error {
	custom_errors.Log(fmt.Sprintf("closing %s file", file), nil, custom_errors.InfoLevel)
	switch file {
	case "source":
		err := lex.inputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
			return err
		}
	case "output":
		err := lex.outputFile.Close()
		if err != nil {
			custom_errors.Log(custom_errors.FileCloseError, &err, custom_errors.ErrorLevel)
			return err
		}
	}
	custom_errors.Log(custom_errors.FileCloseSuccess, nil, custom_errors.InfoLevel)
}

func (lex *Lexical) gravaSaida() error {
	if lex.outputFile == nil {
		return fmt.Errorf("Nome de Arquivo Inválido")
	}
	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(lex.tokensIdentificados.String())
	if err != nil {
		return err
	}
	fmt.Println("Arquivo Salvo: output.txt")
	return nil
}

func (lex *Lexical) exibeTokens() {
	fmt.Println("Tokens Identificados (token/lexema):")
	fmt.Println(lex.tokensIdentificados.String())
}

func (lex *Lexical) acumulaToken(tokenIdentificado string) {
	lex.tokensIdentificados.WriteString(tokenIdentificado)
	lex.tokensIdentificados.WriteString("\n")
}
