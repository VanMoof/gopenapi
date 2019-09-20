package cmd

import (
	"fmt"
	"github.com/VanMoof/gopenapi/generate"
	"github.com/VanMoof/gopenapi/interpret"
	"io"
	"os"
	"path/filepath"
)

func GenerateSpec(format string, output string, args []string) error {
	givenPath := ""
	if len(args) != 0 {
		givenPath = args[0]
	}
	normalizedPath, err := NormalizeInputPath(givenPath)
	if err != nil {
		return fmt.Errorf("failed to normalize working directory: %w", err)
	}

	out, err := ResolveOutputWriter(output)
	if err != nil {
		return err
	}
	s := ResolveOutputSink(format, out)
	return generate.Generate(generate.GoFileVisitor{BasePath: normalizedPath}, &interpret.ASTInterpreter{}, s)
}

func ResolveOutputSink(format string, out io.WriteCloser) generate.Sink {
	var s generate.Sink
	if format == "json" {
		s = &generate.JSONSink{W: out}
	} else if format == "yaml" {
		s = &generate.YAMLSink{W: out}
	}
	return s
}

func ResolveOutputWriter(output string) (io.WriteCloser, error) {
	if output == "-" {
		return os.Stdout, nil
	}
	out, outFileError := os.Create(output)
	if outFileError != nil {
		return nil, outFileError
	}
	return out, nil
}

func NormalizeInputPath(inputPath string) (string, error) {
	if filepath.IsAbs(inputPath) {
		return inputPath, nil
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(currentDir, inputPath), nil
}
