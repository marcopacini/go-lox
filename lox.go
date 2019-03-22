//  MIT License
//
//  Copyright (c) 2019 Marco Pacini
//
//  Permission is hereby granted, free of charge, to any person obtaining a copy
//  of this software and associated documentation files (the "Software"), to deal
//  in the Software without restriction, including without limitation the rights
//  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//  copies of the Software, and to permit persons to whom the Software is
//  furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included in all
//  copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//  SOFTWARE.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"golox/ast"
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
