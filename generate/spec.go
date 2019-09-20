package generate

import (
	"encoding/json"
	"fmt"
	"github.com/VanMoof/gopenapi/interpret"
	"github.com/VanMoof/gopenapi/models"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Sink interface {
	Write(interface{}) error
}

type JSONSink struct {
	W io.WriteCloser
}

func (j *JSONSink) Write(i interface{}) error {
	defer j.W.Close()
	encoder := json.NewEncoder(j.W)
	encoder.SetIndent("", "  ")
	return encoder.Encode(i)
}

type YAMLSink struct {
	W io.WriteCloser
}

func (y *YAMLSink) Write(i interface{}) error {
	defer y.W.Close()
	encoder := yaml.NewEncoder(y.W)
	encoder.SetIndent(2)
	return encoder.Encode(i)
}

type FileVisitor interface {
	VisitFiles(func(filePath string, info os.FileInfo, err error) error) error
}

type GoFileVisitor struct {
	BasePath string
}

func (g GoFileVisitor) VisitFiles(f func(filePath string, info os.FileInfo, err error) error) error {
	return filepath.Walk(g.BasePath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(filePath, ".go") {
			return f(filePath, info, err)
		}
		return nil
	})
}

func Generate(f FileVisitor, i interpret.Interpreter, s Sink) error {
	root := models.Root{OpenAPI: "3.0.2", Components: &models.Components{}}

	err := f.VisitFiles(func(filePath string, info os.FileInfo, err error) error {
		file, _ := os.Open(filePath)
		defer file.Close()
		return i.InterpretFile(file, &root)
	})

	if err != nil {
		return fmt.Errorf("failed to read files: %w", err)
	}

	err = s.Write(&root)
	if err != nil {
		return err
	}
	return nil
}
