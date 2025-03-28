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
	TModule                 = 1
	TVariable               = 3
	TComma                  = 4
	TSemiColon              = 5
	TIf                     = 6
	TElse                   = 7
	TElseIf                 = 8
	TFor                    = 9
	TBreak                  = 11
	TOpenParentheses        = 16
	TCloseParentheses       = 17
	TGreaterThan            = 19
	TLessThan               = 20
	TGreaterEqual           = 21
	TLessEqual              = 22
	TEqual                  = 23
	TNotEqual               = 24
	TAdditionOperator       = 25
	TSubtractionOperator    = 26
	TMultiplicationOperator = 27
	TDivisionOperator       = 28
	TModuleOperator         = 29
	TExponentiationOperator = 30
	TInteger                = 31
	TFLoat                  = 32
	TId                     = 33

	TInputEnd = 90
	TLexError = 98
	TNil      = 99

	EOF = 26
)

// Lexical struct to hold lexical analyzer state
type Lexical struct {
	inputFile        *os.File
	rdInput          *bufio.Reader
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
	lex.token = TNil
	lex.errorMessage = ""

	lex.movelookAhead()

	for lex.token != TInputEnd && lex.token != TLexError {
		lex.buscaProximoToken()
		lex.mostraToken()
	}

	if lex.token == TLexError {
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
		line, err := lex.rdInput.ReadString('\n')
		if err != nil {
			lex.lookAhead = EOF
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
			lex.token = TModule
		case "FIM":
			lex.token = TEnd
		case "VARIAVEIS":
			lex.token = TVariable
		case "SE":
			lex.token = TIf
		case "SENAO":
			lex.token = TElse
		case "ELSE_IF":
			lex.token = TElseIf
		case "FIM_SE":
			lex.token = TEndIf
		case "ENQUANTO":
			lex.token = TFor
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
		lex.token = TInteger
	} else {
		switch lex.lookAhead {
		case '(':
			lex.token = TOpenParentheses
		case ')':
			lex.token = TCloseParentheses
		case ';':
			lex.token = TSemiColon
		case ',':
			lex.token = TComma
		case '+':
			lex.token = TAdditionOperator
		case '-':
			lex.token = TSubtractionOperator
		case '*':
			lex.token = TMultiplicationOperator
		case '/':
			lex.token = TDivisionOperator
		case '%':
			lex.token = TModuleOperator
		case '<':
			lex.token = TLessThan
		case '>':
			lex.token = TGreaterThan
		case '=':
			lex.token = TEqual
		case EOF:
			lex.token = TInputEnd
		default:
			lex.token = TLexError
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
	case TModule:
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
	case TFor:
		tokenLexema = "T_ENQUANTO"
	case TEndWhile:
		tokenLexema = "T_FIM_ENQUANTO"
	case TBreak:
		tokenLexema = "T_PARA"
	case TArrow:
		tokenLexema = "T_SETA"
	case TAte:
		tokenLexema = "T_ATE"
	case TFimPara:
		tokenLexema = "T_FIM_PARA"
	case TLer:
		tokenLexema = "T_LER"
	case TOpenParentheses:
		tokenLexema = "T_ABRE_PAR"
	case TCloseParentheses:
		tokenLexema = "T_FECHA_PAR"
	case TEscrever:
		tokenLexema = "T_ESCREVER"
	case TGreaterThan:
		tokenLexema = "T_MAIOR"
	case TLessThan:
		tokenLexema = "T_MENOR"
	case TGreaterEqual:
		tokenLexema = "T_MAIOR_IGUAL"
	case TLessEqual:
		tokenLexema = "T_MENOR_IGUAL"
	case TEqual:
		tokenLexema = "T_IGUAL"
	case TNotEqual:
		tokenLexema = "T_DIFERENTE"
	case TAdditionOperator:
		tokenLexema = "T_MAIS"
	case TSubtractionOperator:
		tokenLexema = "T_MENOS"
	case TMultiplicationOperator:
		tokenLexema = "T_VEZES"
	case TDivisionOperator:
		tokenLexema = "T_DIVIDIDO"
	case TModuleOperator:
		tokenLexema = "T_RESTO"
	case TExponentiationOperator:
		tokenLexema = "T_ELEVADO"
	case TInteger:
		tokenLexema = "T_NUMERO"
	case TId:
		tokenLexema = "T_ID"
	case TInputEnd:
		tokenLexema = "T_FIM_FONTE"
	case TLexError:
		tokenLexema = "T_ERRO_LEX"
	case TNil:
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
	lex.rdInput = bufio.NewReader(file)
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
