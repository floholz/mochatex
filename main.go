package main

import (
	"github.com/floholz/mochatex/cmd/mochatex"
	"log"
	"os"
	"os/exec"
)

func main() {
	errLog := log.New(os.Stderr, "ERROR: ", log.Lshortfile|log.LstdFlags)
	infoLog := log.New(os.Stdout, "INFO: ", log.Lshortfile|log.LstdFlags)

	// Check for pdfLaTeX (pdfTex will do in a pinch)
	cmd := "pdflatex"
	if _, err := exec.LookPath(cmd); err != nil {
		errLog.Printf("error while searching checking pdflatex binary: %v\n\tchecking for pdftex binary", err)
		if _, err := exec.LookPath("pdftex"); err != nil {
			errLog.Fatal("neither pdflatex nor pdftex binary found in your $PATH")
		}
		infoLog.Printf("found pdftex binary; falling back to using pdftex instead of pdflatex")
		cmd = "pdftex"
	}

	if len(os.Args) < 1 {
		errLog.Fatal("no arguments provided")
	}

	mochatex.Cli(errLog, infoLog)
}
