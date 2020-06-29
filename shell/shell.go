package shell

import (
	"bufio"
	"context"
	"io"

	"github.com/HakanSunay/gohil/lexer"
	"github.com/HakanSunay/gohil/logger"
	"github.com/HakanSunay/gohil/parser"
)

const prompt = "GOHIL=> "

func Start(ctx context.Context, reader io.Reader, writer io.Writer) {
	log := logger.GetFromContext(ctx)

	scanner := bufio.NewScanner(reader)
	for {
		io.WriteString(writer, prompt)
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
				io.WriteString(writer, errorMsg + "\n")
			}
			continue
		}


		io.WriteString(writer, program.String() + "\n")
		io.WriteString(writer, "\n")
	}
}