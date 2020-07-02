package shell

import (
	"bufio"
	"context"
	"github.com/HakanSunay/gohil/env"
	"github.com/HakanSunay/gohil/eval"
	"io"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/logger"
	"github.com/HakanSunay/gohil/parser"
)

const prompt = "GOHIL=> "

func Start(ctx context.Context, reader io.Reader, writer io.Writer) {
	log := logger.GetFromContext(ctx)

	scanner := bufio.NewScanner(reader)
	environment := env.NewEnvironment()

	for {
		_, err := io.WriteString(writer, prompt)
		if err != nil {
			log.Errorf("Unable to redirect prompt output")
			continue
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		if len(p.GetErrors()) > 0 {
			log.Errorf("Parse encountered the following erros: %v", p.GetErrors())
			for _, errorMsg := range p.GetErrors() {
				_, err := io.WriteString(writer, errorMsg + "\n")
				if err != nil {
					log.Errorf("Unable to redirect error output")
				}
			}
			continue
		}

		result := eval.Eval(program, environment)
		if result == nil {
			log.Errorf("Unsupported evaluation type")
			continue
		}

		_, err = io.WriteString(writer, result.Inspect() + "\n")
		if err != nil {
			log.Errorf("Unable to redirect result output")
		}
	}
}
