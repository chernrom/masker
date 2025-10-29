package main

import (
	"fmt"
	"os"

	"github.com/chernrom/masker/internal/fileio"
	"github.com/chernrom/masker/internal/service"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: app <input_path> [output_path]")
		os.Exit(2)
	}
	inPath := os.Args[1]
	outPath := ""
	if len(os.Args) >= 3 {
		outPath = os.Args[2]
	}

	prod := fileio.NewFileProducer(inPath)
	pres := fileio.NewFilePresenter(outPath)

	svc := service.NewService(prod, pres)
	if err := svc.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
