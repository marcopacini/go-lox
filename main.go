package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/marcopacini/go-lox/ast"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()

	if len(flag.Args()) > 1 {
		println("usage: lox [script]")
		os.Exit(64)
	}

	if len(flag.Args()) == 1 {
		runFile(flag.Arg(0))
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := run(string(b)); err != nil {
		fmt.Println(err)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		if b, err := reader.ReadString('\n'); err == nil {
			if err := run(string(b)); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func run(source string) error {
	s := ast.Scanner{source}

	tokens, err := s.Scan()
	if err != nil {
		return err
	}

	//	for _, token := range tokens {
	//		fmt.Println(token)
	//	}

	p := ast.Parser{Tokens: tokens}

	stmts, err := p.Parse()
	if err != nil {
		return err
	}

	i := ast.Interpreter{}

	if err = i.Run(stmts); err != nil {
		return err
	}

	/*
		if err := parser.Resolve(ast); err != nil {
			return err
		}

		for _, stmt := range ast {
			if _, err := stmt.Evaluate(env); err != nil {
				return err
			}
		}
	*/

	return nil
}
