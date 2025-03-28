package main

import (
	"bufio"
	"fmt"
	custom_errors "mechanus-compiler/error"
	"os"
	"strings"
)

// Constants for token types
const (
	TProgram    = 1
	TEnd        = 2
	TVariable   = 3
	TComma      = 4
	TSemiColon  = 5
	TIf         = 6
	TElse       = 7
	TEndIf      = 8
	TWhile      = 9
	TEndWhile   = 10
	TBreak      = 11
	TSeta       = 12
	TAte        = 13
	TFimPara    = 14
	TLer        = 15
	TAbrePar    = 16
	TFechaPar   = 17
	TEscrever   = 18
	TMaior      = 19
	TMenor      = 20
	TMaiorIgual = 21
	TMenorIgual = 22
	TIgual      = 23
	TDiferente  = 24
	TMais       = 25
	TMenos      = 26
	TVezes      = 27
	TDividido   = 28
	TResto      = 29
	TElevado    = 30
	TNumero     = 31
	TId         = 32
	TFimFonte   = 90
	TErroLex    = 98
	TNulo       = 99
	FimArquivo  = 26
)

// Lexical struct to hold lexical analyzer state
type Lexical struct {
	inputFile        *os.File
	rdFonte          *bufio.Reader
	outputFile       *os.File
	lookAhead        rune
	token            int
	lexeme           string
	pointer          int
	inputLine        string
	currentLine      int
	currentColumn    int
	errorMessage     string
	identifiedTokens strings.Builder
}

func (lex *Lexical) GetToken(inputFile string) (string, error) {
	err := lex.open(inputFile)

	if err != nil {
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		return "", err
	}
	defer lex.close("input")

	lex.currentLine = 0
	lex.currentColumn = 0
	lex.pointer = 0
	lex.inputLine = ""
	lex.token = TNulo
	lex.errorMessage = ""

	lex.movelookAhead()

	for lex.token != TFimFonte && lex.token != TErroLex {
		lex.buscaProximoToken()
		lex.mostraToken()
	}

	if lex.token == TErroLex {
		fmt.Println("Erro Léxico:", lex.errorMessage)
	} else {
		fmt.Println("Análise Léxica terminada sem erros léxicos")
	}

	lex.exibeTokens()
	lex.gravaSaida()
}

func (lex *Lexical) movelookAhead() error {
	if lex.pointer+1 > len(lex.inputLine) {
		lex.currentLine++
		lex.pointer = 0
		line, err := lex.rdFonte.ReadString('\n')
		if err != nil {
			lex.lookAhead = FimArquivo
			return err
		}
		lex.inputLine = line
		lex.lookAhead = rune(lex.inputLine[lex.pointer])
	} else {
		lex.lookAhead = rune(lex.inputLine[lex.pointer])
	}
	if lex.lookAhead >= 'a' && lex.lookAhead <= 'z' {
		lex.lookAhead = lex.lookAhead - 'a' + 'A'
	}
	lex.pointer++
	lex.currentColumn = lex.pointer + 1
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

		lex.lexeme = sbLexema.String()

		switch lex.lexeme {
		case "PROGRAMA":
			lex.token = TProgram
		case "FIM":
			lex.token = TEnd
		case "VARIAVEIS":
			lex.token = TVariable
		case "SE":
			lex.token = TIf
		case "SENAO":
			lex.token = TElse
		case "FIM_SE":
			lex.token = TEndIf
		case "ENQUANTO":
			lex.token = TWhile
		case "FIM_ENQUANTO":
			lex.token = TEndWhile
		case "PARA":
			lex.token = TBreak
		case "ATE":
			lex.token = TAte
		case "FIM_PARA":
			lex.token = TFimPara
		case "LER":
			lex.token = TLer
		case "ESCREVER":
			lex.token = TEscrever
		default:
			lex.token = TId
		}
	} else if lex.lookAhead >= '0' && lex.lookAhead <= '9' {
		sbLexema.WriteRune(lex.lookAhead)
		lex.movelookAhead()
		for lex.lookAhead >= '0' && lex.lookAhead <= '9' {
			sbLexema.WriteRune(lex.lookAhead)
			lex.movelookAhead()
		}
		lex.token = TNumero
	} else {
		switch lex.lookAhead {
		case '(':
			lex.token = TAbrePar
		case ')':
			lex.token = TFechaPar
		case ';':
			lex.token = TSemiColon
		case ',':
			lex.token = TComma
		case '+':
			lex.token = TMais
		case '-':
			lex.token = TMenos
		case '*':
			lex.token = TVezes
		case '/':
			lex.token = TDividido
		case '%':
			lex.token = TResto
		case '<':
			lex.token = TMenor
		case '>':
			lex.token = TMaior
		case '=':
			lex.token = TIgual
		case FimArquivo:
			lex.token = TFimFonte
		default:
			lex.token = TErroLex
			lex.errorMessage = fmt.Sprintf("Erro Léxico na linha: %d\nReconhecido ao atingir a coluna: %d\nLinha do Erro: <%s>\nToken desconhecido: %c", lex.currentLine, lex.currentColumn, lex.inputLine, lex.lookAhead)
		}
		sbLexema.WriteRune(lex.lookAhead)
		lex.movelookAhead()
	}

	lex.lexeme = sbLexema.String()
	return nil
}

func (lex *Lexical) mostraToken() {
	var tokenLexema string
	switch lex.token {
	case TProgram:
		tokenLexema = "T_PROGRAMA"
	case TEnd:
		tokenLexema = "T_FIM"
	case TVariable:
		tokenLexema = "T_VARIAVEIS"
	case TComma:
		tokenLexema = "T_VIRGULA"
	case TSemiColon:
		tokenLexema = "T_PONTO_VIRGULA"
	case TIf:
		tokenLexema = "T_SE"
	case TElse:
		tokenLexema = "T_SENAO"
	case TEndIf:
		tokenLexema = "T_FIM_SE"
	case TWhile:
		tokenLexema = "T_ENQUANTO"
	case TEndWhile:
		tokenLexema = "T_FIM_ENQUANTO"
	case TBreak:
		tokenLexema = "T_PARA"
	case TSeta:
		tokenLexema = "T_SETA"
	case TAte:
		tokenLexema = "T_ATE"
	case TFimPara:
		tokenLexema = "T_FIM_PARA"
	case TLer:
		tokenLexema = "T_LER"
	case TAbrePar:
		tokenLexema = "T_ABRE_PAR"
	case TFechaPar:
		tokenLexema = "T_FECHA_PAR"
	case TEscrever:
		tokenLexema = "T_ESCREVER"
	case TMaior:
		tokenLexema = "T_MAIOR"
	case TMenor:
		tokenLexema = "T_MENOR"
	case TMaiorIgual:
		tokenLexema = "T_MAIOR_IGUAL"
	case TMenorIgual:
		tokenLexema = "T_MENOR_IGUAL"
	case TIgual:
		tokenLexema = "T_IGUAL"
	case TDiferente:
		tokenLexema = "T_DIFERENTE"
	case TMais:
		tokenLexema = "T_MAIS"
	case TMenos:
		tokenLexema = "T_MENOS"
	case TVezes:
		tokenLexema = "T_VEZES"
	case TDividido:
		tokenLexema = "T_DIVIDIDO"
	case TResto:
		tokenLexema = "T_RESTO"
	case TElevado:
		tokenLexema = "T_ELEVADO"
	case TNumero:
		tokenLexema = "T_NUMERO"
	case TId:
		tokenLexema = "T_ID"
	case TFimFonte:
		tokenLexema = "T_FIM_FONTE"
	case TErroLex:
		tokenLexema = "T_ERRO_LEX"
	case TNulo:
		tokenLexema = "T_NULO"
	default:
		tokenLexema = "N/A"
	}
	fmt.Println(tokenLexema + " ( " + lex.lexeme + " )")
	lex.acumulaToken(tokenLexema + " ( " + lex.lexeme + " )")
}

func (lex *Lexical) open(fileName string) error {
	custom_errors.Log(fmt.Sprintf("opening '%s' file", fileName), nil, custom_errors.InfoLevel)
	file, err := os.Open(fileName)

	if err != nil {
		custom_errors.Log(custom_errors.FileOpenError, &err, custom_errors.ErrorLevel)
		return err
	}

	custom_errors.Log(custom_errors.FileOpenError, nil, custom_errors.ErrorLevel)
	lex.inputFile = file
	lex.rdFonte = bufio.NewReader(file)
	return nil
}

func (lex *Lexical) close(file string) error {
	custom_errors.Log(fmt.Sprintf("closing %s file", file), nil, custom_errors.InfoLevel)

	switch file {
	case "input":
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
	return nil
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

	_, err = file.WriteString(lex.identifiedTokens.String())
	if err != nil {
		return err
	}
	fmt.Println("Arquivo Salvo: output.txt")
	return nil
}

func (lex *Lexical) exibeTokens() {
	fmt.Println("Tokens Identificados (token/lexeme):")
	fmt.Println(lex.identifiedTokens.String())
}

func (lex *Lexical) acumulaToken(tokenIdentificado string) {
	lex.identifiedTokens.WriteString(tokenIdentificado)
	lex.identifiedTokens.WriteString("\n")
}
