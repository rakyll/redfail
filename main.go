package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

var highlight = color.New(color.FgHiRed)

func main() {
	if len(os.Args) < 2 {
		usage("Missing target program.")
		os.Exit(1)
	}

	r, w := io.Pipe()

	var args = []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	cmd := exec.Command(os.Args[1], args...)
	cmd.Stderr = w
	cmd.Stdout = os.Stdout

	go consume(r)

	cmd.Env = os.Environ()
	cmd.Run()
}

func consume(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		_ = line
		highlight.Fprintf(os.Stderr, "%s\n", line)
	}
}

const usageText = `Usage: redfail <program> [args...]`

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, usageText)
}
